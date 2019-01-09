REM Set CWD to dist/
PUSHD dist

REM Build Go project
go build ../pkg/bayar
go build ../cmd/bayar

REM Return to original CWD
POPD
