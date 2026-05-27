// цей файл відповідає за моделі сканування та очищення файлової системи
package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Описує правила для шляху очищення
type PathSpec struct {
	Pattern             string
	ClearContents       bool
	ExcludeBaseNames    map[string]struct{}
	ExcludeBasePrefixes []string
}

// Описує ціль очищення
type Target struct {
	ID    int
	Name  string
	Paths []PathSpec
}

// Описує результат сканування або очищення
type TargetResult struct {
	ID             int
	Name           string
	SizeBytes      int64
	HumanSize      string
	MatchedPath    int
	Cleaned        bool
	DeletedEntries int
	FailedEntries  int
	Status         string
	Warnings       []string
}

// Зберігає стан cleaner
type Cleaner struct {
	homeDir string
	targets []Target
}

// Створює новий екземпляр cleaner
func NewCleaner(homeDir string, targets []Target) *Cleaner {
	return &Cleaner{
		homeDir: homeDir,
		targets: targets,
	}
}

// Повертає список цілей очищення
func (c *Cleaner) Targets() []Target {
	out := make([]Target, len(c.targets))
	copy(out, c.targets)
	return out
}

// Сканує всі доступні цілі
func (c *Cleaner) ScanAll() []TargetResult {
	results := make([]TargetResult, 0, len(c.targets))

	for _, target := range c.targets {
		results = append(results, c.ScanTarget(target))
	}

	return results
}

// Сканує одну ціль
func (c *Cleaner) ScanTarget(target Target) TargetResult {
	result := TargetResult{
		ID:   target.ID,
		Name: target.Name,
	}

	for _, spec := range target.Paths {
		matches, warnings := c.resolveMatches(spec)

		result.Warnings = append(result.Warnings, warnings...)
		result.MatchedPath += len(matches)

		for _, match := range matches {
			size, err := pathSize(match, spec.ClearContents, spec.ExcludeBaseNames, spec.ExcludeBasePrefixes)
			result.SizeBytes += size

			if err != nil {
				result.Warnings = append(result.Warnings, warningsFromError(match, err)...)
			}
		}
	}

	result.HumanSize = HumanSize(result.SizeBytes)

	result.Status = "ready"

	return result
}

// Очищає вибрані цілі за id
func (c *Cleaner) CleanByIDs(ids []int) ([]TargetResult, error) {
	if len(ids) == 0 {
		return nil, errors.New("at least one --id value is required")
	}

	selected := make([]Target, 0, len(ids))

	seen := make(map[int]struct{}, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}

		target, ok := c.targetByID(id)

		if !ok {
			return nil, fmt.Errorf("unknown target id: %d", id)
		}

		seen[id] = struct{}{}
		selected = append(selected, target)
	}

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].ID < selected[j].ID
	})

	return c.cleanTargets(selected), nil
}

// Очищає всі цілі
func (c *Cleaner) CleanAll() []TargetResult {
	return c.cleanTargets(c.targets)
}

// Виконує очищення списку цілей
func (c *Cleaner) cleanTargets(targets []Target) []TargetResult {
	results := make([]TargetResult, 0, len(targets))

	for _, target := range targets {
		result := TargetResult{
			ID:   target.ID,
			Name: target.Name,
		}

		for _, spec := range target.Paths {
			matches, warnings := c.resolveMatches(spec)

			result.Warnings = append(result.Warnings, warnings...)
			result.MatchedPath += len(matches)

			for _, match := range matches {
				size, err := pathSize(match, spec.ClearContents, spec.ExcludeBaseNames, spec.ExcludeBasePrefixes)
				result.SizeBytes += size

				if err != nil {
					result.Warnings = append(result.Warnings, warningsFromError(match, err)...)
				}

				stats, err := cleanPath(match, spec.ClearContents, spec.ExcludeBaseNames, spec.ExcludeBasePrefixes)

				result.DeletedEntries += stats.DeletedEntries
				result.FailedEntries += stats.FailedEntries

				if err != nil {
					result.Warnings = append(result.Warnings, warningsFromError(match, err)...)
				}
			}
		}

		result.Cleaned = true
		result.HumanSize = HumanSize(result.SizeBytes)
		result.Status = cleanupStatus(result)

		results = append(results, result)
	}

	return results
}

// Визначає підсумковий статус очищення
func cleanupStatus(result TargetResult) string {
	switch {
	case result.FailedEntries > 0 && result.DeletedEntries > 0:
		return "partial"
	case result.FailedEntries > 0:
		return "failed"
	case result.DeletedEntries > 0:
		return "success"
	case result.MatchedPath == 0:
		return "nothing to clean"
	default:
		return "nothing to clean"
	}
}

// Шукає ціль за id
func (c *Cleaner) targetByID(id int) (Target, bool) {
	for _, target := range c.targets {
		if target.ID == id {
			return target, true
		}
	}

	return Target{}, false
}

// Знаходить файлові шляхи для цілі
func (c *Cleaner) resolveMatches(spec PathSpec) ([]string, []string) {
	pattern := expandHome(spec.Pattern, c.homeDir)

	if !hasGlob(pattern) {
		if _, err := os.Lstat(pattern); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, nil
			}

			return nil, warningsFromError(pattern, err)
		}

		return []string{pattern}, nil
	}

	matches, err := filepath.Glob(pattern)

	if err != nil {
		return nil, warningsFromError(pattern, err)
	}

	existing := make([]string, 0, len(matches))

	for _, match := range matches {
		if _, err := os.Lstat(match); err == nil {
			existing = append(existing, match)
		}
	}

	return existing, nil
}

// Перетворює байти у зручний формат
func HumanSize(size int64) string {
	const unit = 1024

	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	value := float64(size)

	suffixes := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	index := 0

	for value >= unit && index < len(suffixes)-1 {
		value /= unit
		index++
	}

	if value >= 10 {
		return fmt.Sprintf("%.1f %s", value, suffixes[index])
	}
	return fmt.Sprintf("%.2f %s", value, suffixes[index])
}

// Рахує розмір шляху
func pathSize(path string, clearContents bool, excluded map[string]struct{}, excludedPrefixes []string) (int64, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	if clearContents && info.IsDir() {
		return dirContentsSize(path, excluded, excludedPrefixes)
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	return treeSize(path)
}

// Рахує розмір вмісту директорії
func dirContentsSize(dir string, excluded map[string]struct{}, excludedPrefixes []string) (int64, error) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return 0, err
	}

	var total int64
	var errs []error

	for _, entry := range entries {
		if isExcluded(entry.Name(), excluded, excludedPrefixes) {
			continue
		}

		child := filepath.Join(dir, entry.Name())
		size, err := pathSize(child, false, nil, nil)
		total += size

		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", child, err))
		}
	}

	return total, errors.Join(errs...)
}

// Рахує розмір дерева файлів
func treeSize(root string) (int64, error) {
	var total int64
	var errs []error

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, walkErr))

			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()

		if err != nil {
			return err
		}

		total += info.Size()

		return nil
	})

	if err != nil {
		errs = append(errs, err)
	}

	return total, errors.Join(errs...)
}

// Описує статистику очищення
type CleanupStats struct {
	DeletedEntries int
	FailedEntries  int
}

// Видаляє шлях або його вміст
func cleanPath(path string, clearContents bool, excluded map[string]struct{}, excludedPrefixes []string) (CleanupStats, error) {
	info, err := os.Lstat(path)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return CleanupStats{}, nil
		}

		return CleanupStats{FailedEntries: 1}, err
	}

	if clearContents && info.IsDir() {
		entries, err := os.ReadDir(path)

		if err != nil {
			return CleanupStats{FailedEntries: 1}, err
		}

		stats := CleanupStats{}

		var errs []error

		for _, entry := range entries {
			if isExcluded(entry.Name(), excluded, excludedPrefixes) {
				continue
			}

			child := filepath.Join(path, entry.Name())

			if err := os.RemoveAll(child); err != nil {
				stats.FailedEntries++
				errs = append(errs, fmt.Errorf("%s: %w", child, err))
				continue
			}

			stats.DeletedEntries++
		}

		return stats, errors.Join(errs...)
	}

	if err := os.RemoveAll(path); err != nil {
		return CleanupStats{FailedEntries: 1}, err
	}

	return CleanupStats{DeletedEntries: 1}, nil
}

// Розгортає домашню директорію в повний шлях
func expandHome(path, homeDir string) string {
	if path == "~" {
		return homeDir
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}

	return path
}

// Збирає шлях з частин
func filepathJoin(parts ...string) string {
	filtered := make([]string, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		filtered = append(filtered, part)
	}

	if len(filtered) == 0 {
		return ""
	}

	return filepath.Join(filtered...)
}

// Перевіряє чи містить шлях glob символи
func hasGlob(path string) bool {
	return strings.ContainsAny(path, "*?[")
}

// Перевіряє чи шлях треба виключити
func isExcluded(base string, excluded map[string]struct{}, excludedPrefixes []string) bool {
	if len(excluded) > 0 {
		if _, ok := excluded[base]; ok {
			return true
		}
	}

	for _, prefix := range excludedPrefixes {
		if strings.HasPrefix(base, prefix) {
			return true
		}
	}

	return false
}

// Формує текст попередження
func formatWarning(path string, err error) string {
	return fmt.Sprintf("%s: %v", path, err)
}

// Розгортає складені помилки в список попереджень
func warningsFromError(path string, err error) []string {
	type multiUnwrapper interface {
		Unwrap() []error
	}

	var multi multiUnwrapper

	if errors.As(err, &multi) {
		unwrapped := multi.Unwrap()
		warnings := make([]string, 0, len(unwrapped))

		for _, item := range unwrapped {
			warnings = append(warnings, warningsFromError(path, item)...)
		}

		return warnings
	}

	return []string{formatWarning(path, err)}
}
