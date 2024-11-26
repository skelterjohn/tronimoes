echo "Building tronserv..."
go install .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
go build -o tronserv.exe .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
.\tronserv.exe $args