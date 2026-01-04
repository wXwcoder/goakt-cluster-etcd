@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

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
