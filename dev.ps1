param(
  [string]$Root = "C:\Users\treen\recruitflow",
  [string]$DbUser = "recruitflow",
  [string]$DbPass = "recruitflow",
  [string]$DbName = "recruitflow",
  [int]$DbPort = 5433,
  [string]$DbContainer = "api-db-1",
  [string]$JwtSecret = "super_dev_secret_change_me"
)

$ErrorActionPreference = "Stop"

function Info($m){ Write-Host "[INFO] $m" -ForegroundColor Cyan }
function Ok($m){ Write-Host "[OK]   $m" -ForegroundColor Green }
function Warn($m){ Write-Host "[WARN] $m" -ForegroundColor Yellow }

$apiDir = Join-Path $Root "apps\api"
$webDir = Join-Path $Root "apps\web"
if (!(Test-Path $apiDir)) { throw "Missing: $apiDir" }
if (!(Test-Path $webDir)) { throw "Missing: $webDir" }

# --- Ensure the right DB container is running (the one bound to 5433)
Info "Ensuring DB container '$DbContainer' is running..."
$names = docker ps -a --format "{{.Names}}"
if (-not ($names -contains $DbContainer)) {
  throw "Docker container '$DbContainer' not found. You currently have api-db-1 and recruitflow-db. Use -DbContainer api-db-1"
}
docker start $DbContainer *> $null
Ok "DB container running."

# --- Build DB URL
$dbUrl = "postgres://$DbUser`:$DbPass@localhost:$DbPort/$DbName?sslmode=disable"
Info "DATABASE_URL = $dbUrl"

# --- Start API
$apiOneLiner = "cd `"$apiDir`"; `$env:DATABASE_URL=`"$dbUrl`"; `$env:JWT_SECRET=`"$JwtSecret`"; go run .\cmd\api\main.go"
Start-Process powershell -ArgumentList "-NoExit","-Command",$apiOneLiner

# --- Start Web
$webOneLiner = "cd `"$webDir`"; npm install; npm run dev"
Start-Process powershell -ArgumentList "-NoExit","-Command",$webOneLiner

Ok "Launched!"
Info "API: http://localhost:8080/health"
Info "Web: http://localhost:5173 (or next free port)"
