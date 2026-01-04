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
echo Starting node3 service (Port: 13001)...
::start "node3" 
"%ACCOUNTS_EXE%" run ^
    --service-name node3 ^
    --system-name Accounts ^
    --port 13001 ^
    --gossip-port 13002 ^
    --peers-port 13003 ^
    --remoting-port 13004 ^
    --etcd-endpoints localhost:2379 ^
    --etcd-dial-timeout 5 ^
    --etcd-ttl 30 ^
    --etcd-timeout 5