// Package internal цей файл відповідає за вибір цілей для поточної операційної системи
package internal

import (
	"fmt"
	"runtime"
)

// TargetsForCurrentOS повертає цілі для поточної системи
func TargetsForCurrentOS() ([]Target, string, error) {
	switch runtime.GOOS {
	case "darwin":
		return macOSTargets(), "macOS", nil
	case "linux":
		return nil, "Linux", fmt.Errorf("linux support is not implemented yet")
	case "windows":
		return nil, "Windows", fmt.Errorf("windows support is not implemented yet")
	default:
		return nil, runtime.GOOS, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
