@echo off
setlocal

set "BINARY=build\cleaner.exe"
set "GOCACHE_DIR=%CD%\.gocache"

if "%~1"=="" goto :help
if /I "%~1"=="help" goto :help
if /I "%~1"=="build" goto :build
if /I "%~1"=="run" goto :run
if /I "%~1"=="scan" goto :scan
if /I "%~1"=="test" goto :test
if /I "%~1"=="fmt" goto :fmt
if /I "%~1"=="clean" goto :clean

echo Unknown command: %~1
echo.
goto :help

:help
echo Available commands:
echo   build.bat build          Build the cleaner binary
echo   build.bat run [args...]  Build and run the CLI
echo   build.bat scan           Build and run "cleaner scan"
echo   build.bat test           Run Go tests
echo   build.bat fmt            Format Go sources
echo   build.bat clean          Remove build artifacts
exit /b 0

:build
if not exist build mkdir build
set "GOCACHE=%GOCACHE_DIR%"
go build -o "%BINARY%" .
exit /b %ERRORLEVEL%

:run
call "%~f0" build
if errorlevel 1 exit /b %ERRORLEVEL%
shift
"%BINARY%" %*
exit /b %ERRORLEVEL%

:scan
call "%~f0" build
if errorlevel 1 exit /b %ERRORLEVEL%
"%BINARY%" scan
exit /b %ERRORLEVEL%

:test
set "GOCACHE=%GOCACHE_DIR%"
go test ./...
exit /b %ERRORLEVEL%

:fmt
gofmt -w main.go
if errorlevel 1 exit /b %ERRORLEVEL%
for /R internal %%f in (*.go) do (
    gofmt -w "%%f"
    if errorlevel 1 exit /b %ERRORLEVEL%
)
exit /b 0

:clean
if exist build rmdir /S /Q build
if exist bin rmdir /S /Q bin
if exist .gocache rmdir /S /Q .gocache
if exist cleaner del /Q cleaner
if exist cleaner.exe del /Q cleaner.exe
exit /b 0
