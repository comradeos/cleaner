// Package internal цей файл відповідає за cli команди вивід та взаємодію з користувачем
package internal

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

const usageText = `cleaner is a safe CLI utility for scanning and removing disposable files.

Usage:
  cleaner scan
  cleaner clean --id 3
  cleaner clean --id 3 --id 5
  cleaner clean --all --yes

Options:
  --id    Clean one or more target IDs.
  --all   Clean every target for the current platform.
  --yes   Skip the confirmation prompt.
`

type idList []int

// Повертає список id як рядок
func (ids *idList) String() string {
	if len(*ids) == 0 {
		return ""
	}

	parts := make([]string, 0, len(*ids))

	for _, id := range *ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}

	return strings.Join(parts, ",")
}

// Set додає id до списку
func (ids *idList) Set(value string) error {
	var id int

	_, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &id)

	if err != nil {
		return fmt.Errorf("invalid id %q", value)
	}

	if id <= 0 {
		return fmt.Errorf("invalid id %q", value)
	}

	*ids = append(*ids, id)

	return nil
}

// Run запускає cli та повертає код завершення
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		_, _ = fmt.Fprintf(stderr, "failed to resolve home directory: %v\n", err)
		return 1
	}

	targets, platformName, err := TargetsForCurrentOS()

	if err != nil {
		_, _ = fmt.Fprintf(stderr, "%v\n", err)
		return 1
	}

	cleaner := NewCleaner(homeDir, targets)

	switch args[0] {
	case "scan":
		printScan(stdout, platformName, cleaner.ScanAll())
		return 0
	case "clean":
		return runClean(args[1:], stdout, stderr, platformName, cleaner)
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command: %s\n\n", args[0])
		printUsage(stderr)
		return 1
	}
}

// Виконує команду очищення
func runClean(args []string, stdout, stderr io.Writer, platformName string, cleaner *Cleaner) int {
	flags := flag.NewFlagSet("clean", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var ids idList
	var cleanAll bool
	var yes bool

	flags.Var(&ids, "id", "Clean one or more target IDs")
	flags.BoolVar(&cleanAll, "all", false, "Clean every target")
	flags.BoolVar(&yes, "yes", false, "Skip confirmation prompt")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if cleanAll && len(ids) > 0 {
		_, _ = fmt.Fprintln(stderr, "use either --all or one or more --id flags")
		return 1
	}

	if !cleanAll && len(ids) == 0 {
		_, _ = fmt.Fprintln(stderr, "provide --all or at least one --id flag")
		return 1
	}

	if err := validateIDs(cleaner.Targets(), ids); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}

	scanned := cleaner.ScanAll()
	selected := filterSelected(scanned, cleanAll, ids)

	if len(selected) == 0 {
		_, _ = fmt.Fprintln(stderr, "no matching targets found")
		return 1
	}

	printSelection(stdout, platformName, selected)

	if !yes {
		confirmed, err := promptConfirmation(stdout)

		if err != nil {
			_, _ = fmt.Fprintf(stderr, "failed to read confirmation: %v\n", err)
			return 1
		}

		if !confirmed {
			_, _ = fmt.Fprintln(stdout, "Cleanup cancelled.")
			return 0
		}
	}

	var results []TargetResult
	var err error

	if cleanAll {
		results = cleaner.CleanAll()
	} else {
		results, err = cleaner.CleanByIDs(ids)

		if err != nil {
			_, _ = fmt.Fprintln(stderr, err)
			return 1
		}
	}

	printCleanResults(stdout, platformName, results)

	if hasCleanupFailures(results) {
		return 2
	}

	return 0
}

// Виводить довідку з використання
func printUsage(w io.Writer) {
	_, _ = fmt.Fprint(w, usageText)
}

// Виводить результати сканування
func printScan(w io.Writer, platformName string, results []TargetResult) {
	_, _ = fmt.Fprintf(w, "%s scan results\n\n", platformName)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "ID\tName\tSize\tMatched Paths")
	_, _ = fmt.Fprintln(tw, "--\t----\t----\t-------------")

	var total int64

	for _, result := range results {
		total += result.SizeBytes
		_, _ = fmt.Fprintf(tw, "%d\t%s\t%s\t%d\n", result.ID, result.Name, result.HumanSize, result.MatchedPath)
	}
	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTotal reclaimable size: %s\n", HumanSize(total))

	printWarnings(w, results)
}

// Виводить вибрані цілі перед очищенням
func printSelection(w io.Writer, platformName string, results []TargetResult) {
	_, _ = fmt.Fprintf(w, "%s cleanup targets\n\n", platformName)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "ID\tName\tSize")
	_, _ = fmt.Fprintln(tw, "--\t----\t----")

	var total int64

	for _, result := range results {
		total += result.SizeBytes
		_, _ = fmt.Fprintf(tw, "%d\t%s\t%s\n", result.ID, result.Name, result.HumanSize)
	}

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nEstimated reclaimable size: %s\n", HumanSize(total))
}

// Виводить результати очищення
func printCleanResults(w io.Writer, platformName string, results []TargetResult) {
	_, _ = fmt.Fprintf(w, "%s cleanup completed\n\n", platformName)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "ID\tName\tEstimated Size\tStatus")
	_, _ = fmt.Fprintln(tw, "--\t----\t--------------\t------")

	var total int64

	for _, result := range results {
		total += result.SizeBytes
		_, _ = fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", result.ID, result.Name, result.HumanSize, result.Status)
	}
	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nEstimated requested cleanup size: %s\n", HumanSize(total))

	printWarnings(w, results)
}

// Виводить попередження
func printWarnings(w io.Writer, results []TargetResult) {
	var generalWarnings []string

	permissionRestricted := false

	for _, result := range results {
		permissionCount := 0

		for _, warning := range result.Warnings {
			if isPermissionWarning(warning) {
				permissionCount++
				permissionRestricted = true
				continue
			}

			generalWarnings = append(generalWarnings, warning)
		}

		if permissionCount > 0 {
			generalWarnings = append(generalWarnings,
				fmt.Sprintf("%s: skipped %d restricted path(s) because access was denied", result.Name, permissionCount),
			)
		}
	}

	if len(generalWarnings) == 0 {
		return
	}

	sort.Strings(generalWarnings)

	_, _ = fmt.Fprintln(w, "\nWarnings:")

	for _, warning := range generalWarnings {
		_, _ = fmt.Fprintf(w, "- %s\n", warning)
	}

	if permissionRestricted {
		_, _ = fmt.Fprintln(w, "\nNote: Some macOS paths may still require additional permissions or elevated access.")
	}
}

// Відбирає результати за id
func filterSelected(results []TargetResult, cleanAll bool, ids []int) []TargetResult {
	if cleanAll {
		return results
	}

	selected := make([]TargetResult, 0, len(ids))
	index := make(map[int]TargetResult, len(results))

	for _, result := range results {
		index[result.ID] = result
	}

	seen := make(map[int]struct{}, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}

		if result, ok := index[id]; ok {
			selected = append(selected, result)
			seen[id] = struct{}{}
		}
	}

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].ID < selected[j].ID
	})

	return selected
}

// Перевіряє що всі id існують
func validateIDs(targets []Target, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	known := make(map[int]struct{}, len(targets))

	for _, target := range targets {
		known[target.ID] = struct{}{}
	}

	var invalid []string

	seen := make(map[int]struct{}, len(ids))

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}

		if _, ok := known[id]; !ok {
			invalid = append(invalid, fmt.Sprintf("%d", id))
		}
	}

	if len(invalid) == 0 {
		return nil
	}

	sort.Strings(invalid)
	return fmt.Errorf("unknown target id(s): %s", strings.Join(invalid, ", "))
}

// Перевіряє чи були помилки під час очищення
func hasCleanupFailures(results []TargetResult) bool {
	for _, result := range results {
		if result.Status == "partial" || result.Status == "failed" {
			return true
		}
	}

	return false
}

// Запитує підтвердження на видалення
func promptConfirmation(w io.Writer) (bool, error) {
	_, _ = fmt.Fprint(w, "Type YES to permanently delete the selected files: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	return strings.TrimSpace(input) == "YES", nil
}

// Визначає чи попередження пов'язане з доступом
func isPermissionWarning(warning string) bool {
	lower := strings.ToLower(warning)
	return strings.Contains(lower, "operation not permitted") || strings.Contains(lower, "permission denied")
}
