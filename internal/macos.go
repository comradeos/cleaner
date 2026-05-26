// Package internal цей файл відповідає за набір цілей очищення для macOS
package internal

const (
	userLogsPath              = "~/Library/Logs"
	systemLogsPath            = "/Library/Logs"
	diagnosticReportsDir      = "DiagnosticReports"
	userDiagnosticReportsPath = userLogsPath + "/" + diagnosticReportsDir
	systemDiagnosticReports   = systemLogsPath + "/" + diagnosticReportsDir

	googleChromeCachePath  = "~/Library/Caches/Google/Chrome"
	chromeCanaryCachePath  = "~/Library/Caches/Google/Chrome Canary"
	chromiumCachePath      = "~/Library/Caches/Chromium"
	braveBrowserCachePath  = "~/Library/Caches/BraveSoftware/Brave-Browser"
	microsoftEdgeCachePath = "~/Library/Caches/Microsoft Edge"
	firefoxCachePath       = "~/Library/Caches/Firefox"

	xcodeDerivedDataPath = "~/Library/Developer/Xcode/DerivedData"
	homebrewCachePath    = "~/Library/Caches/Homebrew"

	npmCachePath        = "~/.npm"
	yarnCachePath       = "~/Library/Caches/Yarn"
	pnpmCachePath       = "~/Library/Caches/pnpm"
	pnpmStorePath       = "~/Library/pnpm/store"
	pipCachePath        = "~/.cache/pip"
	pipLibraryCachePath = "~/Library/Caches/pip"
	cargoRegistryPath   = "~/.cargo/registry/cache"
	cargoGitDBPath      = "~/.cargo/git/db"
	goBuildCachePath    = "~/Library/Caches/go-build"
	goModuleCachePath   = "~/go/pkg/mod/cache"
)

// Повертає список цілей для macOS
func macOSTargets() []Target {
	return []Target{
		{
			ID:    1,
			Name:  "System and User Logs",
			Paths: clearPathsWithExcludedNames([]string{userLogsPath, systemLogsPath}, diagnosticReportsDir),
		},
		{
			ID:    2,
			Name:  "Crash Reports and Diagnostic Reports",
			Paths: clearPaths(userDiagnosticReportsPath, systemDiagnosticReports),
		},
		{
			ID:   3,
			Name: "Browser Caches",
			Paths: clearPaths(
				googleChromeCachePath,
				chromeCanaryCachePath,
				chromiumCachePath,
				braveBrowserCachePath,
				microsoftEdgeCachePath,
				firefoxCachePath,
			),
		},
		{
			ID:    4,
			Name:  "Xcode Derived Data",
			Paths: clearPaths(xcodeDerivedDataPath),
		},
		{
			ID:    5,
			Name:  "Homebrew Cache",
			Paths: clearPaths(homebrewCachePath),
		},
		{
			ID:   6,
			Name: "Package Manager Caches",
			Paths: clearPaths(
				npmCachePath,
				yarnCachePath,
				pnpmCachePath,
				pnpmStorePath,
				pipCachePath,
				pipLibraryCachePath,
				cargoRegistryPath,
				cargoGitDBPath,
				goBuildCachePath,
				goModuleCachePath,
			),
		},
	}
}

// Створює список шляхів для очищення
func clearPaths(patterns ...string) []PathSpec {
	paths := make([]PathSpec, 0, len(patterns))

	for _, pattern := range patterns {
		paths = append(paths, PathSpec{
			Pattern:       pattern,
			ClearContents: true,
		})
	}

	return paths
}

// Створює список шляхів з виключенням за назвами
func clearPathsWithExcludedNames(patterns []string, excludedNames ...string) []PathSpec {
	paths := make([]PathSpec, 0, len(patterns))
	excluded := newNameSet(excludedNames...)

	for _, pattern := range patterns {
		paths = append(paths, PathSpec{
			Pattern:          pattern,
			ClearContents:    true,
			ExcludeBaseNames: excluded,
		})
	}

	return paths
}

// Створює набір назв для виключення
func newNameSet(names ...string) map[string]struct{} {
	values := make(map[string]struct{}, len(names))

	for _, name := range names {
		values[name] = struct{}{}
	}

	return values
}
