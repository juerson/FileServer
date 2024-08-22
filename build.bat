@echo off

REM 编译适用于Windows的64位程序
set GOOS=windows
set GOARCH=amd64
go build -o FileServer-windows-amd64.exe

REM 编译适用于Windows的32位程序
set GOOS=windows
set GOARCH=386
go build -o FileServer-windows-386.exe

REM 编译适用于Windows ARM的64位程序
set GOOS=windows
set GOARCH=arm64
go build -o FileServer-windows-arm64.exe

REM 编译适用于Linux的64位程序
set GOOS=linux
set GOARCH=amd64
go build -o FileServer-linux-amd64

REM 编译适用于Linux的32位程序
set GOOS=linux
set GOARCH=386
go build -o FileServer-linux-386

REM 编译适用于macOS的64位程序
set GOOS=darwin
set GOARCH=amd64
go build -o FileServer-macos-amd64

REM 编译适用于macOS的ARM64位程序（适用于Apple M1/M2芯片）
set GOOS=darwin
set GOARCH=arm64
go build -o FileServer-macos-arm64

echo 编译完成!
