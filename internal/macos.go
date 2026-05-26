// Package internal цей файл відповідає за набір цілей очищення для macOS
package internal

// Повертає список цілей для macOS
func macOSTargets() []Target {
	return []Target{
		{
			ID:   1,
			Name: "System and User Logs",
			Paths: []PathSpec{
				{
					Pattern:       "~/Library/Logs",
					ClearContents: true,
					ExcludeBaseNames: map[string]struct{}{
						"DiagnosticReports": {},
					},
				},
				{
					Pattern:       "/Library/Logs",
					ClearContents: true,
					ExcludeBaseNames: map[string]struct{}{
						"DiagnosticReports": {},
					},
				},
			},
		},
		{
			ID:   2,
			Name: "Crash Reports and Diagnostic Reports",
			Paths: []PathSpec{
				{Pattern: "~/Library/Logs/DiagnosticReports", ClearContents: true},
				{Pattern: "/Library/Logs/DiagnosticReports", ClearContents: true},
			},
		},
		{
			ID:   3,
			Name: "Browser Caches",
			Paths: []PathSpec{
				{Pattern: "~/Library/Caches/Google/Chrome", ClearContents: true},
				{Pattern: "~/Library/Caches/Google/Chrome Canary", ClearContents: true},
				{Pattern: "~/Library/Caches/Chromium", ClearContents: true},
				{Pattern: "~/Library/Caches/BraveSoftware/Brave-Browser", ClearContents: true},
				{Pattern: "~/Library/Caches/Microsoft Edge", ClearContents: true},
				{Pattern: "~/Library/Caches/Firefox", ClearContents: true},
			},
		},
		{
			ID:   4,
			Name: "Xcode Derived Data",
			Paths: []PathSpec{
				{Pattern: "~/Library/Developer/Xcode/DerivedData", ClearContents: true},
			},
		},
		{
			ID:   5,
			Name: "Homebrew Cache",
			Paths: []PathSpec{
				{Pattern: "~/Library/Caches/Homebrew", ClearContents: true},
			},
		},
		{
			ID:   6,
			Name: "Package Manager Caches",
			Paths: []PathSpec{
				{Pattern: "~/.npm", ClearContents: true},
				{Pattern: "~/Library/Caches/Yarn", ClearContents: true},
				{Pattern: "~/Library/Caches/pnpm", ClearContents: true},
				{Pattern: "~/Library/pnpm/store", ClearContents: true},
				{Pattern: "~/.cache/pip", ClearContents: true},
				{Pattern: "~/Library/Caches/pip", ClearContents: true},
				{Pattern: "~/.cargo/registry/cache", ClearContents: true},
				{Pattern: "~/.cargo/git/db", ClearContents: true},
				{Pattern: "~/Library/Caches/go-build", ClearContents: true},
				{Pattern: "~/go/pkg/mod/cache", ClearContents: true},
			},
		},
	}
}
