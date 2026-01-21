@echo off
echo Building VniConverter...
wails build -clean -platform windows/amd64
if %errorlevel% neq 0 (
    echo Build failed!
    pause
    exit /b %errorlevel%
)

echo.
echo Build successful!
echo Executable located at: build\bin\VniConverter.exe
echo.
pause
