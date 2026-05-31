package scanner

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// ExtractImage pulls a container image using crane, extracts its layers into a temporary
// directory, and returns the path to that directory. The caller is responsible for
// removing the directory when done.
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
		
		// Basic path traversal protection
		if !filepath.HasPrefix(target, filepath.Clean(dest)) {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				// Ignore errors like permission denied on weird files
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
