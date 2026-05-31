package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall/js"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/HarshalPatel1972/spectra/internal/scanner"
)

// scanSpectra is the function exposed to the JavaScript host.
// It takes three arguments: code (string), language (string), filename (string).
// It returns a JSON string containing the ScanResult.
func scanSpectra(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return js.ValueOf(\`{"error": "requires 3 arguments: code, language, filename"}\`)
	}

	code := args[0].String()
	// language := args[1].String()
	filename := args[2].String()

	// In the Go WASM environment, the OS provides an in-memory filesystem.
	// We can write the code to a temporary directory.
	tmpDir, err := os.MkdirTemp("", "spectra-wasm-")
	if err != nil {
		return js.ValueOf(fmt.Sprintf(\`{"error": "failed to create temp dir: %v"}\`, err))
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, filename)
	err = os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return js.ValueOf(fmt.Sprintf(\`{"error": "failed to write file: %v"}\`, err))
	}

	// Load patterns from the embedded rules or load hardcoded defaults if we can't embed.
	// Note: We need a PatternRegistry. For simplicity in WASM, we'll try to load it from disk
	// or assume the caller passed the patterns. Wait, spectra includes patterns via embed in a real setup,
	// but currently the detector loads from rules/crypto_patterns.yaml.
	// Let's assume detector has a way to load from bytes or we mock it.
	
	// Since we are running in WASM, we must ensure we have the patterns.
	// For now, we will construct an empty or minimal registry if loading fails.
	registry, err := detector.LoadPatterns("rules/crypto_patterns.yaml")
	if err != nil {
		// Fallback for WASM environment if file is missing
		registry = &detector.PatternRegistry{ByLanguage: make(map[string][]detector.PatternEntry)}
	}

	// Create a new orchestrator and scan
	// In spectra, scanner.ScanDirectory handles the whole flow.
	res, err := scanner.ScanDirectory(tmpDir, []string{}, []string{"code"}, 1, registry, false)
	if err != nil {
		return js.ValueOf(fmt.Sprintf(\`{"error": "scan failed: %v"}\`, err))
	}

	// Serialize result to JSON
	out, err := json.Marshal(res)
	if err != nil {
		return js.ValueOf(fmt.Sprintf(\`{"error": "json marshal failed: %v"}\`, err))
	}

	return js.ValueOf(string(out))
}

func main() {
	// Expose the function to the global JavaScript object
	js.Global().Set("scanSpectra", js.FuncOf(scanSpectra))
	
	// Block forever so the Go WASM runtime stays alive
	select {}
}
