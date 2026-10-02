@echo off
setlocal
cd /d "%~dp0"

echo ======================================================================
echo  Editor de Adventures - Build dos fontes em Go
echo ======================================================================

powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0build.ps1"
if %ERRORLEVEL% neq 0 (
    echo.
    echo [ERRO] O build falhou!
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo Build finalizado com sucesso!
pause
