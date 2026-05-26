# cleaner

A free and open-source terminal utility for scanning, detecting, and removing unnecessary files.

## Current scope

The first version focuses on `macOS` and safely scans only known disposable targets:

- `System and User Logs`
- `Crash Reports and Diagnostic Reports`
- `Browser Caches` (cache only)
- `Xcode Derived Data`
- `Homebrew Cache`
- `Package Manager Caches`

The cleaner does not use `sudo` by default and only removes files from an explicit whitelist of locations.
Protected macOS locations that typically require `Full Disk Access`, such as `Trash` and Safari cache folders, are intentionally excluded in this version.

## Make targets

```bash
make build
make run ARGS="scan"
make scan
make test
make clean
```

## Windows script

Use `build.bat` on Windows instead of `make`:

```bat
build.bat build
build.bat run scan
build.bat scan
build.bat test
build.bat clean
```

## Usage

Scan what can be removed:

```bash
make scan
```

On Windows:

```bat
build.bat scan
```

The binary is created in `build/cleaner` on macOS/Linux and `build/cleaner.exe` on Windows.

Clean a single target:

```bash
make run ARGS="clean --id 3"
```

```bat
build.bat run clean --id 3
```

Clean multiple targets:

```bash
make run ARGS="clean --id 1 --id 3 --id 6"
```

```bat
build.bat run clean --id 1 --id 3 --id 6
```

Clean everything without an interactive confirmation prompt:

```bash
make run ARGS="clean --all --yes"
```

```bat
build.bat run clean --all --yes
```

## Notes

- IDs are stable within the current platform target list.
- Sizes are displayed in floating human-readable units.
- Some system locations may produce warnings if the current user does not have permission to remove them.
- `clean` returns exit code `2` when cleanup only partially succeeds.
