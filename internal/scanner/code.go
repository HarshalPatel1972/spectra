package scanner

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/HarshalPatel1972/spectra/internal/detector"
)

// extensionToLanguage maps file extensions to the pattern language key used in
// the crypto_patterns.yaml registry. Extensions not in this map are still
// scanned with generic patterns.
var extensionToLanguage = map[string]string{
	".go":    "go",
	".py":    "python",
	".java":  "java",
	".kt":    "java",
	".scala": "java",
	".js":    "javascript",
	".ts":    "javascript",
	".tsx":   "javascript",
	".jsx":   "javascript",
	".rs":    "rust",
	".c":     "c_cpp",
	".cpp":   "c_cpp",
	".cc":    "c_cpp",
	".h":     "c_cpp",
	".cs":    "csharp",
	".rb":    "ruby",
	".php":   "php",
	".swift": "swift",
}

// defaultSkipDirs contains directory basenames that are always skipped during
// code scanning regardless of user-supplied excludes.
var defaultSkipDirs = map[string]bool{
	".git":        true,
	"node_modules": true,
	"vendor":      true,
	"dist":        true,
	"build":       true,
	"__pycache__": true,
	".idea":       true,
	".vscode":     true,
}

// keySizeRe matches common cryptographic key sizes in source code lines.
var keySizeRe = regexp.MustCompile(`\b(128|192|256|384|512|521|1024|2048|3072|4096)\b`)

// sensitivePatterns are substrings that indicate a line contains secret material
// and should be redacted in outputs.
var sensitivePatterns = []string{
	"PRIVATE KEY",
	"SECRET",
	"PASSWORD",
	"TOKEN",
}

// isSensitiveLine returns true if the line content suggests it contains
// secret material that should be redacted.
func isSensitiveLine(line string) bool {
	upper := strings.ToUpper(line)
	for _, pat := range sensitivePatterns {
		if strings.Contains(upper, pat) {
			return true
		}
	}
	return false
}

// isBinaryFile checks the first 512 bytes of the given file for NUL characters.
// Files containing NUL bytes are considered binary and should be skipped.
func isBinaryFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		// Empty file or read error — treat as non-binary.
		return false, nil
	}

	for _, b := range buf[:n] {
		if b == 0 {
			return true, nil
		}
	}
	return false, nil
}

// shouldSkipDir returns true if the directory should be excluded from scanning.
func shouldSkipDir(dirPath, dirName string, excludes []string) bool {
	if defaultSkipDirs[dirName] {
		return true
	}
	for _, pattern := range excludes {
		// Match the glob against the directory name.
		if matched, _ := filepath.Match(pattern, dirName); matched {
			return true
		}
		// Match the glob against the relative path.
		if matched, _ := filepath.Match(pattern, dirPath); matched {
			return true
		}
	}
	return false
}

// shouldExcludeFile returns true if the file should be excluded from scanning
// based on user-supplied glob patterns.
func shouldExcludeFile(filePath string, excludes []string) bool {
	baseName := filepath.Base(filePath)
	for _, pattern := range excludes {
		if matched, _ := filepath.Match(pattern, baseName); matched {
			return true
		}
		if matched, _ := filepath.Match(pattern, filePath); matched {
			return true
		}
	}
	return false
}

// ScanCodeFiles walks the root directory, scanning each text file for
// cryptographic algorithm patterns defined in the PatternRegistry. It uses a
// semaphore-based worker pool limited by the concurrency parameter.
//
// Returns the collected findings, file statistics, and any walk error.
func ScanCodeFiles(root string, excludes []string, registry *detector.PatternRegistry, concurrency int, verbose bool) ([]Finding, FileStats, error) {
	if concurrency < 1 {
		concurrency = 4
	}

	var (
		mu       sync.Mutex
		findings []Finding
		stats    FileStats
	)

	// Collect files to scan.
	var filePaths []string
	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip inaccessible entries
		}

		relPath, _ := filepath.Rel(root, path)
		if relPath == "." {
			return nil
		}

		if d.IsDir() {
			if shouldSkipDir(relPath, d.Name(), excludes) {
				return filepath.SkipDir
			}
			return nil
		}

		mu.Lock()
		stats.TotalFiles++
		mu.Unlock()

		if shouldExcludeFile(relPath, excludes) {
			mu.Lock()
			stats.SkippedFiles++
			mu.Unlock()
			return nil
		}

		filePaths = append(filePaths, path)
		return nil
	})
	if walkErr != nil {
		return nil, stats, walkErr
	}

	// Process files with a bounded worker pool.
	g := new(errgroup.Group)
	g.SetLimit(concurrency)

	for _, fp := range filePaths {
		fp := fp // capture loop variable
		g.Go(func() error {
			relPath, _ := filepath.Rel(root, fp)

			// Skip binary files.
			binary, bErr := isBinaryFile(fp)
			if bErr != nil {
				mu.Lock()
				stats.SkippedFiles++
				mu.Unlock()
				return nil
			}
			if binary {
				mu.Lock()
				stats.SkippedFiles++
				mu.Unlock()
				return nil
			}

			// Determine language from extension.
			ext := strings.ToLower(filepath.Ext(fp))
			language := extensionToLanguage[ext]

			// Get applicable patterns.
			var patterns []detector.PatternEntry
			if language != "" && registry != nil {
				patterns = append(patterns, registry.ByLanguage[language]...)
			}
			// Always apply generic patterns.
			if registry != nil {
				patterns = append(patterns, registry.Generic...)
			}

			if len(patterns) == 0 {
				mu.Lock()
				stats.SkippedFiles++
				mu.Unlock()
				return nil
			}

			mu.Lock()
			stats.ScannedFiles++
			mu.Unlock()

			// Scan the file line by line.
			fileFindings, scanErr := scanFileLines(fp, relPath, language, patterns)
			if scanErr != nil {
				return nil // skip unreadable files
			}

			if len(fileFindings) > 0 {
				mu.Lock()
				findings = append(findings, fileFindings...)
				mu.Unlock()
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return findings, stats, err
	}

	return findings, stats, nil
}

// scanFileLines reads a file line by line and matches each line against the
// provided patterns. Returns findings for all matches.
func scanFileLines(absPath, relPath, language string, patterns []detector.PatternEntry) ([]Finding, error) {
	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var findings []Finding
	var lines []string

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 256*1024), 1024*1024) // up to 1MB lines
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Track consecutive duplicate findings for deduplication.
	type dedupKey struct {
		algorithm string
	}
	var lastKey dedupKey
	var lastFindingIdx int = -1
	consecutiveCount := 0

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			// Reset consecutive tracking on blank lines.
			lastKey = dedupKey{}
			consecutiveCount = 0
			continue
		}

		for _, entry := range patterns {
			for _, re := range entry.Patterns {
				if re.MatchString(line) {
					algo := entry.Algorithm
					canonical, found := detector.NormaliseAlgorithm(algo)
					if !found {
						continue
					}
					info, ok := detector.LookupAlgorithm(canonical)
					if !ok {
						continue
					}

					// Deduplication: collapse consecutive matches for the
					// same algorithm into a single finding with incremented
					// occurrence count.
					key := dedupKey{algorithm: canonical}
					if key == lastKey && consecutiveCount > 0 && lastFindingIdx >= 0 {
						findings[lastFindingIdx].OccurrenceCount++
						consecutiveCount++
						goto nextLine
					}

					// Build line content, redact if sensitive.
					lineContent := trimmed
					if isSensitiveLine(lineContent) {
						lineContent = "[REDACTED]"
					}

					// Extract key size from nearby lines (±3 lines).
					keySize := extractKeySize(lines, lineNum)

					findings = append(findings, Finding{
						ID:              uuid.New().String(),
						Algorithm:       canonical,
						AlgorithmInfo:   info,
						Source:          detector.SourceCode,
						FilePath:        relPath,
						LineNumber:      lineNum + 1,
						LineContent:     lineContent,
						Language:        language,
						KeySize:         keySize,
						OccurrenceCount: 1,
					})

					lastFindingIdx = len(findings) - 1
					lastKey = key
					consecutiveCount = 1
					goto nextLine
				}
			}
		}

		// No match on this line resets consecutive tracking.
		lastKey = dedupKey{}
		consecutiveCount = 0

	nextLine:
	}

	return findings, nil
}

// extractKeySize looks in a window of ±3 lines around the given line index for
// a recognisable cryptographic key size value.
func extractKeySize(lines []string, lineIdx int) int {
	start := lineIdx - 3
	if start < 0 {
		start = 0
	}
	end := lineIdx + 3
	if end >= len(lines) {
		end = len(lines) - 1
	}

	for i := start; i <= end; i++ {
		matches := keySizeRe.FindAllString(lines[i], -1)
		for _, m := range matches {
			var size int
			for _, c := range m {
				size = size*10 + int(c-'0')
			}
			// Only return plausible key sizes.
			switch size {
			case 128, 192, 256, 384, 512, 521, 1024, 2048, 3072, 4096:
				return size
			}
		}
	}
	return 0
}
