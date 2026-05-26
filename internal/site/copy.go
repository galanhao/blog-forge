package site

import (
	"os"
	"path/filepath"

	"github.com/galanhao/blog-forge/internal/theme"
)

// copyAssets copies theme assets, theme static, site static, and site assets to dist.
func copyAssets(t *theme.Theme, distDir string) error {
	// Theme assets → dist/assets/ (CSS, JS from theme)
	if err := copyDirIfExists(t.AssetsDir(), filepath.Join(distDir, "assets")); err != nil {
		return err
	}
	// Theme static → dist/ (theme-level static files)
	if err := copyDirIfExists(t.StaticDir(), distDir); err != nil {
		return err
	}
	// Site static → dist/ (site-level static files)
	if err := copyDirIfExists("static", distDir); err != nil {
		return err
	}
	// Site assets → dist/assets/ (images, audio, video under assets/)
	// Merges with theme assets: theme has css/js, site has images/audio/video.
	// No conflict between the two.
	if err := copyDirIfExists("assets", filepath.Join(distDir, "assets")); err != nil {
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
