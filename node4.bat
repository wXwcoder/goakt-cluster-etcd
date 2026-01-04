@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   GoAkt Cluster with ETCD - Start node3
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

REM Start Accounts node services - simple approach
echo Starting node4 service (Port: 14001)...
::start "node4" 
"%ACCOUNTS_EXE%" run ^
    --service-name node4 ^
    --system-name Accounts ^
    --port 14001 ^
    --gossip-port 14002 ^
    --peers-port 14003 ^
    --remoting-port 14004 ^
    --etcd-endpoints localhost:2379 ^
    --etcd-dial-timeout 5 ^
    --etcd-ttl 30 ^
    --etcd-timeout 5