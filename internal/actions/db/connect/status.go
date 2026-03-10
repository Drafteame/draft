package connect

import (
	"fmt"
	"sort"
	"strings"
)

// Status prints the current state of all tracked tunnels, verifying each PID is alive.
func Status() error {
	state, err := loadRuntimeState()
	if err != nil {
		return err
	}

	if len(state.Connections) == 0 {
		fmt.Println("No active connections tracked.")
		return nil
	}

	type row struct {
		dbType    string
		name      string
		status    string
		pid       int
		localPort int
		host      string
	}

	var rows []row
	stateChanged := false

	for key, entry := range state.Connections {
		alive := isProcessAlive(entry.PID)
		statusStr := "up"

		if !alive {
			statusStr = "down"
			delete(state.Connections, key)
			stateChanged = true
		}

		rows = append(rows, row{
			dbType:    entry.DBType,
			name:      entry.Name,
			status:    statusStr,
			pid:       entry.PID,
			localPort: entry.LocalPort,
			host:      entry.Host,
		})
	}

	if stateChanged {
		_ = saveRuntimeState(state) // best-effort cleanup of dead entries
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].dbType != rows[j].dbType {
			return rows[i].dbType < rows[j].dbType
		}

		return rows[i].name < rows[j].name
	})

	wType := len("TYPE")
	wName := len("NAME")
	wStatus := len("STATUS")
	wAddress := len("ADDRESS")

	for _, r := range rows {
		if len(r.dbType) > wType {
			wType = len(r.dbType)
		}

		if len(r.name) > wName {
			wName = len(r.name)
		}

		addr := fmt.Sprintf("localhost:%d", r.localPort)
		if len(addr) > wAddress {
			wAddress = len(addr)
		}
	}

	wType += 2
	wName += 2
	wStatus += 2
	wAddress += 2

	fmt.Printf("%-*s %-*s %-*s %-10s %s\n", wType, "TYPE", wName, "NAME", wStatus, "STATUS", "PID", "ADDRESS")
	fmt.Println(strings.Repeat("-", wType+wName+wStatus+10+wAddress+4))

	for _, r := range rows {
		pidStr := fmt.Sprintf("%d", r.pid)
		if r.status == "down" {
			pidStr = "-"
		}

		addr := fmt.Sprintf("localhost:%d", r.localPort)
		fmt.Printf("%-*s %-*s %-*s %-10s %s\n", wType, r.dbType, wName, r.name, wStatus, r.status, pidStr, addr)
	}

	return nil
}
