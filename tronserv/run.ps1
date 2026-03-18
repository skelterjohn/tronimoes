$ErrorActionPreference = "Stop"

Write-Output "Building tronserv components..."
go install ./...
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Output "Running tronserv..."
tronserv --dev $args