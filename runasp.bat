@echo off
REM Windows launcher for Assemplex
REM Usage: runasp program_name

set file=%1

IF "%file%"=="" (
    echo Usage: runasp program_name
    exit /b 1
)

REM Automatically append .asp if missing
echo %file% | findstr /i "\.asp$" >nul
IF %ERRORLEVEL% NEQ 0 (
    set file=%file%.asp
)

REM Call Python interpreter
python "%~dp0runasp.py" "%~dp0%file%"
