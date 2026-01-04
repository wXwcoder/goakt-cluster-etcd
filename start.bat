@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   GoAkt Cluster with ETCD - Start Script
echo ========================================

REM Set paths
set BIN_DIR=bin
set ETCD_EXE=%BIN_DIR%\etcd.exe
set ACCOUNTS_EXE=%BIN_DIR%\goakt-cluster-etcd.exe

REM Check if binary files exist
if not exist "%ETCD_EXE%" (
    echo Error: Cannot find etcd.exe at %ETCD_EXE%
    pause
    exit /b 1
)

if not exist "%ACCOUNTS_EXE%" (
    echo Error: Cannot find goakt-cluster-etcd.exe at %ACCOUNTS_EXE%
    pause
    exit /b 1
)

echo Starting services...
echo.

REM Start ETCD service
echo [1/3] Starting ETCD service...
start "ETCD Server" "%ETCD_EXE%" ^
    --name etcd-1 ^
    --data-dir etcd-data ^
    --listen-client-urls http://0.0.0.0:2379 ^
    --advertise-client-urls http://localhost:2379 ^
    --listen-peer-urls http://0.0.0.0:2380 ^
    --initial-advertise-peer-urls http://localhost:2380 ^
    --initial-cluster etcd-1=http://localhost:2380 ^
    --initial-cluster-state new ^
    --initial-cluster-token etcd-cluster-1

echo   Waiting for ETCD service to start...
timeout /t 3 /nobreak >nul

REM Start Accounts node services - simple approach
echo [2/3] Starting node0 service (Port: 8000)...
start "node0" "%ACCOUNTS_EXE%" run ^
    --service-name node0 ^
    --system-name Accounts ^
    --port 50051 ^
    --gossip-port 3322 ^
    --peers-port 3320 ^
    --remoting-port 50052 ^
    --etcd-endpoints localhost:2379 ^
    --etcd-dial-timeout 5 ^
    --etcd-ttl 30 ^
    --etcd-timeout 5

echo   Waiting 3 seconds to start next node...
timeout /t 3 /nobreak >nul

echo [3/3] Starting node1 service (Port: 8001)...
start "node1" "%ACCOUNTS_EXE%" run ^
    --service-name node1 ^
    --system-name Accounts ^
    --port 50052 ^
    --gossip-port 3323 ^
    --peers-port 3321 ^
    --remoting-port 50053 ^
    --etcd-endpoints localhost:2379 ^
    --etcd-dial-timeout 5 ^
    --etcd-ttl 30 ^
    --etcd-timeout 5

echo.
echo ========================================
echo   All services started successfully!
echo ========================================
echo.
echo Service Status:
echo   ETCD:      http://localhost:2379
echo   Node0:     http://localhost:50051
echo   Node1:     http://localhost:50052
echo.
echo Each service is running in a separate command window.
echo Close the window to stop the corresponding service.
echo.
pause