package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/galanhao/blog-forge/internal/content"
)

// copyPostAssets copies per-post resource directories to the output directory.
// For a post at "content/posts/2026-05-24-hello.md", if a directory
// "content/posts/2026-05-24-hello/" exists, its contents are copied to
// the post's output directory (e.g. "dist/2026/05/24/hello/").
func copyPostAssets(distDir string, posts []*content.Post) error {
	for _, p := range posts {
		if p.FilePath == "" {
			continue
		}
		assetDir := trimExt(p.FilePath)
		info, err := os.Stat(assetDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("stat %s: %w", assetDir, err)
		}
		if !info.IsDir() {
			continue
		}
		outDir := filepath.Join(distDir, filepath.Clean(p.Permalink))
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", outDir, err)
		}
		if err := copyDir(assetDir, outDir); err != nil {
			return fmt.Errorf("copy post assets %s: %w", assetDir, err)
		}
	}
	return nil
}

// trimExt removes the file extension from a path.
// "content/posts/2026-05-24-hello.md" → "content/posts/2026-05-24-hello"
func trimExt(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}
