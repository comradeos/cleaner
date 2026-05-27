// Package internal цей файл відповідає за набір цілей очищення для linux
package internal

const (
	linuxUserTrashPath      = "~/.local/share/Trash/files"
	linuxThumbnailCachePath = "~/.cache/thumbnails"
	linuxLegacyThumbsPath   = "~/.thumbnails"
	linuxUserLogsPath       = "~/.local/state"

	linuxChromeCachePath     = "~/.cache/google-chrome/*/Cache"
	linuxChromeCodeCachePath = "~/.cache/google-chrome/*/Code Cache"
	linuxChromeGPUCachePath  = "~/.cache/google-chrome/*/GPUCache"

	linuxChromiumCachePath     = "~/.cache/chromium/*/Cache"
	linuxChromiumCodeCachePath = "~/.cache/chromium/*/Code Cache"
	linuxChromiumGPUCachePath  = "~/.cache/chromium/*/GPUCache"

	linuxBraveCachePath     = "~/.cache/BraveSoftware/Brave-Browser/*/Cache"
	linuxBraveCodeCachePath = "~/.cache/BraveSoftware/Brave-Browser/*/Code Cache"
	linuxBraveGPUCachePath  = "~/.cache/BraveSoftware/Brave-Browser/*/GPUCache"

	linuxEdgeCachePath     = "~/.cache/microsoft-edge/*/Cache"
	linuxEdgeCodeCachePath = "~/.cache/microsoft-edge/*/Code Cache"
	linuxEdgeGPUCachePath  = "~/.cache/microsoft-edge/*/GPUCache"

	linuxFirefoxCachePath = "~/.cache/mozilla/firefox/*/cache2"

	linuxNPMCachePath      = "~/.npm"
	linuxYarnCachePath     = "~/.cache/yarn"
	linuxPNPMCachePath     = "~/.cache/pnpm"
	linuxPNPMStorePath     = "~/.local/share/pnpm/store"
	linuxAltPNPMStorePath  = "~/.pnpm-store"
	linuxPipCachePath      = "~/.cache/pip"
	linuxCargoRegistryPath = "~/.cargo/registry/cache"
	linuxCargoGitDBPath    = "~/.cargo/git/db"
	linuxGoBuildCachePath  = "~/.cache/go-build"
	linuxGoModuleCachePath = "~/go/pkg/mod/cache"
)

// Повертає список цілей для linux
func linuxTargets() []Target {
	return []Target{
		{
			ID:    1,
			Name:  "User Trash",
			Paths: clearPaths(linuxUserTrashPath),
		},
		{
			ID:    2,
			Name:  "Thumbnail Cache",
			Paths: clearPaths(linuxThumbnailCachePath, linuxLegacyThumbsPath),
		},
		{
			ID:   3,
			Name: "Browser Caches",
			Paths: clearPaths(
				linuxChromeCachePath,
				linuxChromeCodeCachePath,
				linuxChromeGPUCachePath,
				linuxChromiumCachePath,
				linuxChromiumCodeCachePath,
				linuxChromiumGPUCachePath,
				linuxBraveCachePath,
				linuxBraveCodeCachePath,
				linuxBraveGPUCachePath,
				linuxEdgeCachePath,
				linuxEdgeCodeCachePath,
				linuxEdgeGPUCachePath,
				linuxFirefoxCachePath,
			),
		},
		{
			ID:    4,
			Name:  "User Logs",
			Paths: clearPaths(filepathJoin(linuxUserLogsPath, "*", "logs"), filepathJoin(linuxUserLogsPath, "*", "log")),
		},
		{
			ID:   5,
			Name: "Package Manager Caches",
			Paths: clearPaths(
				linuxNPMCachePath,
				linuxYarnCachePath,
				linuxPNPMCachePath,
				linuxPNPMStorePath,
				linuxAltPNPMStorePath,
				linuxPipCachePath,
				linuxCargoRegistryPath,
				linuxCargoGitDBPath,
				linuxGoBuildCachePath,
				linuxGoModuleCachePath,
			),
		},
	}
}
