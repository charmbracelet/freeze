//go:build !darwin && !linux && !windows

package font

func init() {
	// No default font directories for this platform.
	// System fonts will not be loaded automatically.
	// Users can still use --font.file to specify a custom font.
	DefaultFontsDirs = nil
}
