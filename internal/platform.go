// цей файл відповідає за вибір цілей для поточної операційної системи
package internal

import (
	"fmt"
	"runtime"
)

// повертає цілі для поточної системи
func TargetsForCurrentOS() ([]Target, string, error) {
	switch runtime.GOOS {
	case "darwin":
		return macOSTargets(), "macOS", nil
	case "linux":
		return linuxTargets(), "Linux", nil
	case "windows":
		return windowsTargets(), "Windows", nil
	default:
		return nil, runtime.GOOS, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
