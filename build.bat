@echo off

setlocal EnableDelayedExpansion

go build -mod=mod -ldflags "-X main.Version=1.0.0" -o playground/app.exe
go build -mod=mod -ldflags "-X main.Version=1.1.0" -o playground/updates/app.exe

REM read hash of the exe used as an update
FOR /F "tokens=1 delims= " %%i IN ('powershell "(get-filehash playground\\updates\\app.exe).Hash"') DO (
    set hash=%%i
)

echo { > playground/updates/update.json
echo   "lastVer": "1.1.0", >> playground/updates/update.json
echo   "releasedAt": 1623790082, >> playground/updates/update.json
echo   "sha256": "%hash%" >> playground/updates/update.json
echo } >> playground/updates/update.json
