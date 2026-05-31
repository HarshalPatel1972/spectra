// Package scanner implements the dependency file scanner.
package scanner

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HarshalPatel1972/spectra/internal/detector"
	"github.com/google/uuid"
)

type DependencyInfo struct {
	Algorithms []string
	Risk       string
}

// knownCryptoDeps from the spec §7.4
var knownCryptoDeps = map[string]DependencyInfo{
	// Go
	"golang.org/x/crypto":           {Algorithms: []string{"various"}, Risk: "check-usage"},
	"github.com/square/go-jose":     {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},
	"github.com/dvsekhvalnov/jose2go": {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},

	// npm
	"node-forge": {Algorithms: []string{"RSA", "DES", "3DES", "RC4", "MD5", "SHA1"}, Risk: "CRITICAL"},
	"jsencrypt":  {Algorithms: []string{"RSA"}, Risk: "CRITICAL"},
	"crypto-js":  {Algorithms: []string{"MD5", "SHA1", "DES", "3DES", "RC4"}, Risk: "HIGH"},
	"elliptic":   {Algorithms: []string{"ECDSA", "ECDH"}, Risk: "HIGH"},
	"bn.js":      {Algorithms: []string{"RSA", "DH"}, Risk: "HIGH"},
	"forge":      {Algorithms: []string{"RSA", "DES", "MD5"}, Risk: "CRITICAL"},
	"jose":       {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},

	// Python
	"pycryptodome": {Algorithms: []string{"RSA", "DES", "3DES", "RC4", "MD5", "SHA1"}, Risk: "HIGH"},
	"pycrypto":     {Algorithms: []string{"RSA", "DES", "RC4", "MD5", "SHA1"}, Risk: "CRITICAL"},
	"pyOpenSSL":    {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},
	"paramiko":     {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},
	"cryptography": {Algorithms: []string{"RSA", "ECC", "various"}, Risk: "check-usage"},
	"PyJWT":        {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},

	// Java (Maven artifact IDs)
	"bcprov-jdk15on":  {Algorithms: []string{"RSA", "ECDSA", "DES", "MD5", "SHA1"}, Risk: "HIGH"},
	"bcpkix-jdk15on":  {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},
	"nimbus-jose-jwt": {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},
	"java-jwt":        {Algorithms: []string{"RSA", "ECDSA"}, Risk: "HIGH"},
	"commons-codec":   {Algorithms: []string{"MD5", "SHA1"}, Risk: "HIGH"},

	// Rust (crate names)
	"rsa":  {Algorithms: []string{"RSA"}, Risk: "CRITICAL"},
	"p256": {Algorithms: []string{"ECDSA", "ECDH"}, Risk: "HIGH"},
	"p384": {Algorithms: []string{"ECDSA", "ECDH"}, Risk: "HIGH"},
	"k256": {Algorithms: []string{"ECDSA", "ECDH"}, Risk: "HIGH"},
	"sha1": {Algorithms: []string{"SHA1"}, Risk: "HIGH"},
	"md5":  {Algorithms: []string{"MD5"}, Risk: "HIGH"},
	"des":  {Algorithms: []string{"DES"}, Risk: "CRITICAL"},
}

// targetManifestFiles defines which manifest files we scan.
var targetManifestFiles = map[string]string{
	"go.mod":           "go",
	"go.sum":           "go",
	"package.json":     "npm",
	"yarn.lock":        "yarn",
	"requirements.txt": "pip",
	"Pipfile":          "pipenv",
	"pyproject.toml":   "poetry",
	"pom.xml":          "maven",
	"build.gradle":     "gradle",
	"build.gradle.kts": "gradle",
	"Cargo.toml":       "cargo",
	"composer.json":    "composer",
}

// ScanDepsFiles walks the directory and parses manifest files to detect known
// cryptographic dependencies.
func ScanDepsFiles(root string, excludes []string) ([]Finding, error) {
	var findings []Finding

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors accessing paths
		}

		// Skip directories.
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		baseName := filepath.Base(path)
		pkgMgr, ok := targetManifestFiles[baseName]
		if !ok {
			return nil
		}

		switch baseName {
		case "package.json":
			f, err := scanPackageJSON(path)
			if err == nil {
				findings = append(findings, f...)
			}
		default:
			// Fallback: simple line scanning for all other manifest files.
			f, err := scanManifestLines(path, pkgMgr)
			if err == nil {
				findings = append(findings, f...)
			}
		}

		return nil
	})

	return findings, err
}

func scanPackageJSON(path string) ([]Finding, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var findings []Finding
	checkDep := func(name, version string, isDev bool) {
		if info, ok := knownCryptoDeps[name]; ok {
			for _, algoName := range info.Algorithms {
				algoInfo, found := detector.LookupAlgorithm(algoName)
				if !found {
					// Fallback for generic/various algorithms
					algoInfo = detector.AlgorithmInfo{Name: algoName, Family: "generic"}
				}

				f := Finding{
					ID:              uuid.NewString(),
					Algorithm:       algoInfo.Name,
					AlgorithmInfo:   algoInfo,
					Source:          detector.SourceDeps,
					FilePath:        path,
					LineNumber:      0, // Line number not tracked in JSON parse
					LineContent:     name + "@" + version,
					Language:        "json",
					Context:         "Dependency: " + name + " (" + version + ")",
					OccurrenceCount: 1,
				}
				findings = append(findings, f)
			}
		}
	}

	for k, v := range pkg.Dependencies {
		checkDep(k, v, false)
	}
	for k, v := range pkg.DevDependencies {
		checkDep(k, v, true)
	}

	return findings, nil
}

func scanManifestLines(path, pkgMgr string) ([]Finding, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var findings []Finding
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		
		for depName, info := range knownCryptoDeps {
			// A simple substring check is usually sufficient for manifests,
			// though it might yield false positives if a dep name is a common word.
			// Using regex bounds for a safer match:
			match, _ := regexp.MatchString(`\b`+regexp.QuoteMeta(depName)+`\b`, line)
			if match {
				for _, algoName := range info.Algorithms {
					algoInfo, found := detector.LookupAlgorithm(algoName)
					if !found {
						algoInfo = detector.AlgorithmInfo{Name: algoName, Family: "generic"}
					}

					f := Finding{
						ID:              uuid.NewString(),
						Algorithm:       algoInfo.Name,
						AlgorithmInfo:   algoInfo,
						Source:          detector.SourceDeps,
						FilePath:        path,
						LineNumber:      lineNum,
						LineContent:     strings.TrimSpace(line),
						Language:        pkgMgr,
						Context:         "Dependency: " + depName,
						OccurrenceCount: 1,
					}
					findings = append(findings, f)
				}
			}
		}
		lineNum++
	}

	return findings, scanner.Err()
}
