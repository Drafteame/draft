package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Drafteame/draft/internal/pkg/files"
)

// shouldSkipForStage decides whether a service should be skipped for the given
// stage based on its .deployignore file. The semantics mirror
// scripts/deploy/check_deploy_ignore.sh in the api repo:
//
//   - File missing                       -> deploy (skip=false)
//   - File present but empty/whitespace  -> skip every stage
//   - File with content                  -> one stage per line; skip if the
//     current stage matches any non-empty trimmed line (case-insensitive).
//
// An empty stage is treated as "unknown" and never matches a listed stage; an
// empty file still skips everything in that case (preserves the legacy blunt
// behaviour for callers that do not yet pass a stage).
func shouldSkipForStage(absPath, stage string) (skip bool, reason string, err error) {
	ignorePath := filepath.Join(absPath, ".deployignore")
	if !files.Exists(ignorePath) {
		return false, "", nil
	}

	content, err := files.Read(ignorePath)
	if err != nil {
		return false, "", fmt.Errorf("failed to read %s: %w", ignorePath, err)
	}

	stages := make([]string, 0)
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			stages = append(stages, trimmed)
		}
	}

	if len(stages) == 0 {
		return true, ".deployignore is empty (skipping all stages)", nil
	}

	stageLC := strings.ToLower(stage)
	for _, s := range stages {
		if strings.ToLower(s) == stageLC {
			return true, fmt.Sprintf("stage %q listed in .deployignore", stage), nil
		}
	}

	return false, "", nil
}
