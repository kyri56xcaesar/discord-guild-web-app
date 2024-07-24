:: Welcome to my bot handler

@echo off

:main
cls
echo -------------------------
echo Discord Bot Handelr Shell
echo -------------------------
echo 1. Run all Bots scripts
echo 2. Run specific Python scripts
echo 3. Run one Python script
echo 4. List all bots.
echo 5. Exit
echo -------------------------
set /p choice=Enter your choice: 

if "%choice%"=="1" (
    call :run_all_scripts
    goto main
)

if "%choice%"=="2" (
    call :run_specific_scripts
    goto main
)

if "%choice%"=="3" (
    call :run_one_script
    goto main
)

if "%choice%"=="4" (
    call :list_all_scripts
    goto main
)

if "%choice%"=="5" (
    exit /b
)

echo Invalid choice. Press any key to continue.
pause >nul
goto main

:run_all_scripts
echo Running all Python scripts...
REM Add your commands to run all Python scripts here
echo All scripts executed.
pause
exit /b

:run_specific_scripts
echo Running specific Python scripts...
REM Add your commands to run specific Python scripts here
echo Specific scripts executed.
pause
exit /b

:run_one_script
echo Running one Python script...
REM Add your commands to run one Python script here
echo One script executed.
pause
exit /b

:list_all_scripts
echo Listing all python scripts.
pause
exit /b

