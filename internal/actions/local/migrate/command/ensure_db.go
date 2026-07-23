package command

import (
	"github.com/Drafteame/draft/internal/pkg/exec"
)

func runInitSQL() error {
	cmd := `docker exec -i api-local-postgres psql -U root -d postgres -f /docker-entrypoint-initdb.d/init.sql`
	_, err := exec.Command(cmd)
	return err
}
