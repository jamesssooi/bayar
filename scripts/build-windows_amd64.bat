REM Set CWD to dist/
PUSHD dist
PUSHD win_amd64

REM Set environment variables
set GOARCH=amd64
set GOOS=windows

REM Build Go project
go build ../../pkg/bayar
go build ../../cmd/bayar

REM Return to original CWD
POPD
POPD