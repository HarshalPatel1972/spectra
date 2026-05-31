package scanner

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// BlameFinding adds Git blame information (Author, CommitHash, IntroducedAt)
// to a finding by shelling out to git blame.
func BlameFinding(scanRoot string, f *Finding) {
	if f.Source != "CODE" && f.Source != "CONFIG" && f.Source != "DEPS" {
		return
	}
	if f.LineNumber <= 0 {
		return
	}

	absPath := filepath.Join(scanRoot, f.FilePath)
	
	cmd := exec.Command("git", "blame", "-p", "-L", fmt.Sprintf("%d,%d", f.LineNumber, f.LineNumber), absPath)
	cmd.Dir = scanRoot

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		// Ignore errors (file not tracked, git not installed, etc.)
		return
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) == 0 {
		return
	}

	// First line format: <hash> <original_line> <final_line> <group_lines>
	parts := strings.Fields(lines[0])
	if len(parts) > 0 {
		f.CommitHash = parts[0]
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "author ") {
			f.Author = strings.TrimPrefix(line, "author ")
		} else if strings.HasPrefix(line, "author-time ") {
			tsStr := strings.TrimPrefix(line, "author-time ")
			if ts, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
				f.IntroducedAt = time.Unix(ts, 0)
			}
		}
	}
}
