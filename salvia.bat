@echo off
cd C:\salvia
echo [%DATE% %TIME%] Iniciando Salvia...
.\salvia.exe
if %errorlevel% neq 0 (
    echo.
    echo [ERROR] La aplicacion se detuvo inesperadamente con codigo: %errorlevel%
)
echo.
echo Presione una tecla para cerrar esta ventana...
pause > nul
