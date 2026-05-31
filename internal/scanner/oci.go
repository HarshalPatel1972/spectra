package scanner

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/HarshalPatel1972/spectra/internal/detector"
)

// ScanImage pulls a container image using crane, extracts its layers one by one,
// scans each layer, and attributes findings to the specific layer digest.
func ScanImage(imageRef string, registry *detector.PatternRegistry, concurrency int) ([]Finding, error) {
	img, err := crane.Pull(imageRef)
	if err != nil {
		return nil, fmt.Errorf("pulling image %s: %w", imageRef, err)
	}

	layers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("getting image layers: %w", err)
	}

	var allFindings []Finding

	for _, layer := range layers {
		digest, err := layer.Digest()
		if err != nil {
			continue
		}
		
		layerStr := digest.Hex

		tmpDir, err := os.MkdirTemp("", "spectra-layer-"+layerStr+"-*")
		if err != nil {
			return nil, fmt.Errorf("creating temp dir for layer %s: %w", layerStr, err)
		}

		if err := extractLayer(layer, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return nil, fmt.Errorf("extracting layer %s: %w", layerStr, err)
		}

		// Scan the extracted layer
		codeFindings, _, _ := ScanCodeFiles(tmpDir, nil, registry, concurrency, false)
		certFindings, _ := ScanCertFiles(tmpDir, nil)
		depsFindings, _ := ScanDepsFiles(tmpDir, nil)
		configFindings, _ := ScanConfigFiles(tmpDir, nil, registry)

		var layerFindings []Finding
		layerFindings = append(layerFindings, codeFindings...)
		layerFindings = append(layerFindings, certFindings...)
		layerFindings = append(layerFindings, depsFindings...)
		layerFindings = append(layerFindings, configFindings...)

		for i := range layerFindings {
			layerFindings[i].Source = "CONTAINER"
			layerFindings[i].ContainerLayer = layerStr
		}

		allFindings = append(allFindings, layerFindings...)
		os.RemoveAll(tmpDir)
	}

	return allFindings, nil
}

func extractLayer(layer v1.Layer, dest string) error {
	rc, err := layer.Uncompressed()
	if err != nil {
		return err
	}
	defer rc.Close()

	tr := tar.NewReader(rc)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, filepath.Clean(header.Name))
		
		if !filepath.HasPrefix(target, filepath.Clean(dest)) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				continue
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				continue
			}
			f.Close()
		}
	}
	return nil
}

// ExtractImage is preserved for backward compatibility
func ExtractImage(imageRef string) (string, error) {
	img, err := crane.Pull(imageRef)
	if err != nil {
		return "", fmt.Errorf("pulling image %s: %w", imageRef, err)
	}

	tmpDir, err := os.MkdirTemp("", "spectra-image-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}

	layers, err := img.Layers()
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("getting image layers: %w", err)
	}

	for _, layer := range layers {
		if err := extractLayer(layer, tmpDir); err != nil {
			os.RemoveAll(tmpDir)
			return "", fmt.Errorf("extracting layer: %w", err)
		}
	}

	return tmpDir, nil
}
