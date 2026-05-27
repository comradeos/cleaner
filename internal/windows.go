// Package internal цей файл відповідає за набір цілей очищення для windows
package internal

import (
	"os"
	"path/filepath"
)

// Повертає список цілей для windows
func windowsTargets() []Target {
	localAppData := firstNonEmpty(
		os.Getenv("LOCALAPPDATA"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local"),
	)

	roamingAppData := firstNonEmpty(
		os.Getenv("APPDATA"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming"),
	)

	tempPath := firstNonEmpty(
		os.Getenv("TEMP"),
		os.Getenv("TMP"),
		filepath.Join(localAppData, "Temp"),
	)

	return []Target{
		{
			ID:    1,
			Name:  "Temp Files",
			Paths: clearPaths(tempPath),
		},
		{
			ID:   2,
			Name: "Crash Reports and Error Reports",
			Paths: clearPaths(
				filepathJoin(localAppData, "CrashDumps"),
				filepathJoin(localAppData, "Microsoft", "Windows", "WER", "ReportArchive"),
				filepathJoin(localAppData, "Microsoft", "Windows", "WER", "ReportQueue"),
				filepathJoin(localAppData, "Microsoft", "Windows", "WER", "Temp"),
			),
		},
		{
			ID:   3,
			Name: "Browser Caches",
			Paths: clearPaths(
				filepathJoin(localAppData, "Google", "Chrome", "User Data", "*", "Cache"),
				filepathJoin(localAppData, "Google", "Chrome", "User Data", "*", "Code Cache"),
				filepathJoin(localAppData, "Google", "Chrome", "User Data", "*", "GPUCache"),
				filepathJoin(localAppData, "BraveSoftware", "Brave-Browser", "User Data", "*", "Cache"),
				filepathJoin(localAppData, "BraveSoftware", "Brave-Browser", "User Data", "*", "Code Cache"),
				filepathJoin(localAppData, "BraveSoftware", "Brave-Browser", "User Data", "*", "GPUCache"),
				filepathJoin(localAppData, "Microsoft", "Edge", "User Data", "*", "Cache"),
				filepathJoin(localAppData, "Microsoft", "Edge", "User Data", "*", "Code Cache"),
				filepathJoin(localAppData, "Microsoft", "Edge", "User Data", "*", "GPUCache"),
				filepathJoin(localAppData, "Chromium", "User Data", "*", "Cache"),
				filepathJoin(localAppData, "Chromium", "User Data", "*", "Code Cache"),
				filepathJoin(localAppData, "Chromium", "User Data", "*", "GPUCache"),
				filepathJoin(localAppData, "Mozilla", "Firefox", "Profiles", "*", "cache2"),
			),
		},
		{
			ID:   4,
			Name: "Thumbnail Cache",
			Paths: clearPaths(
				filepathJoin(localAppData, "Microsoft", "Windows", "Explorer", "thumbcache_*.db"),
				filepathJoin(localAppData, "Microsoft", "Windows", "Explorer", "iconcache_*.db"),
			),
		},
		{
			ID:   5,
			Name: "Package Manager Caches",
			Paths: clearPaths(
				filepathJoin(roamingAppData, "npm-cache"),
				filepathJoin(localAppData, "Yarn", "Cache"),
				filepathJoin(localAppData, "pnpm", "store"),
				filepathJoin(localAppData, "pip", "Cache"),
				filepathJoin(os.Getenv("USERPROFILE"), ".cargo", "registry", "cache"),
				filepathJoin(os.Getenv("USERPROFILE"), ".cargo", "git", "db"),
				filepathJoin(localAppData, "go-build"),
				filepathJoin(os.Getenv("USERPROFILE"), "go", "pkg", "mod", "cache"),
			),
		},
	}
}

// Повертає перше непорожнє значення
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
