@echo off
echo ================================================
echo  Go Tic-Tac-Toe Deployment Package Creator
echo ================================================

REM Clean up any existing deployment artifacts
if exist deployment rmdir /S /Q deployment >nul 2>&1
if exist deployment.zip del deployment.zip >nul 2>&1

echo Step 1: Building WebAssembly client...
cd client
echo Building WebAssembly...
set GOOS=js
set GOARCH=wasm
go build -o main.wasm .
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to build WebAssembly client
    cd ..
    pause
    exit /b 1
)
echo WebAssembly build complete: main.wasm created
cd ..

echo.
echo Step 2: Creating deployment structure...
mkdir deployment
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to create deployment directory
    pause
    exit /b 1
)

echo Copying configuration files...
copy docker-compose.yml deployment\ >nul
copy nginx.conf deployment\ >nul
copy tic-tac-toe.service deployment\ >nul
copy setup-lightsail-prebuilt.sh deployment\ >nul
copy Dockerfile deployment\ >nul
if exist .dockerignore copy .dockerignore deployment\ >nul

echo Copying source code directories...
xcopy /E /I /Q server deployment\server\ >nul 2>&1
xcopy /E /I /Q client deployment\client\ >nul 2>&1
xcopy /E /I /Q shared_types deployment\shared_types\ >nul 2>&1

echo.
echo Step 3: Creating deployment package...
powershell Compress-Archive -Path deployment -DestinationPath deployment.zip -Force
if exist deployment.zip (
    echo Created deployment.zip successfully
    set PACKAGE_FILE=deployment.zip
    set EXTRACT_CMD=unzip deployment.zip
) else (
    echo Failed to create deployment package
    pause
    exit /b 1
)

echo.
echo Step 4: Verifying package...
for %%I in (%PACKAGE_FILE%) do set size=%%~zI
echo Package created: %PACKAGE_FILE%
echo Package size: %size% bytes

echo.
echo Step 5: Cleaning up...
rmdir /S /Q deployment >nul 2>&1

echo.
echo ================================================
echo SUCCESS: Deployment package ready!
echo ================================================
echo File: %PACKAGE_FILE%
echo Size: %size% bytes
echo.
echo Your Docker image is on Docker Hub: pakasfand/go-tic-tac-toe:latest
echo.
echo Next steps:
echo 1. Upload %PACKAGE_FILE% to your Lightsail instance
echo 2. SSH to your instance and run:
echo    %EXTRACT_CMD%
echo    cd deployment
echo    chmod +x setup-lightsail-prebuilt.sh
echo    sudo ./setup-lightsail-prebuilt.sh
echo.
echo Deployment will pull from Docker Hub (much faster)!
echo Your game will be available at: http://YOUR_INSTANCE_IP
echo ================================================
pause