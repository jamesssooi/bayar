REM Set CWD to dist/
PUSHD dist
PUSHD linux_amd64

REM Set environment variables
set GOARCH=amd64
set GOOS=linux

REM Build Go project
go build ../../pkg/bayar
go build ../../cmd/bayar

REM Return to original CWD
POPD
POPD
