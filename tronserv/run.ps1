echo "Building tronserv..."
go build -o tronserv.exe .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
go install .
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
.\tronserv.exe $args