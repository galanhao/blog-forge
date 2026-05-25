package site

import (
	"os"
	"path/filepath"

	"github.com/galanhao/blog-builder/internal/theme"
)

// copyAssets copies theme assets, theme static, and site static to dist.
func copyAssets(t *theme.Theme, distDir string) error {
	if err := copyDirIfExists(t.AssetsDir(), filepath.Join(distDir, "assets")); err != nil {
		return err
	}
	if err := copyDirIfExists(t.StaticDir(), distDir); err != nil {
		return err
	}
	if err := copyDirIfExists("static", distDir); err != nil {
		return err
	}
	return nil
}

// copyDirIfExists copies a directory only if it exists.
func copyDirIfExists(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return nil
	}
	return copyDir(src, dst)
}
