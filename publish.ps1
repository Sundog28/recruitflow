# ==========================
# RecruitFlow publish script
# ==========================
param(
  [string]$Root = "$env:USERPROFILE\recruitflow",
  [string]$VercelDomain = "YOUR-PROJECT.vercel.app"  # <- change later
)

$ErrorActionPreference = "Stop"

function Info($m){ Write-Host $m -ForegroundColor Cyan }
function Ok($m){ Write-Host $m -ForegroundColor Green }
function Warn($m){ Write-Host $m -ForegroundColor Yellow }

$apiDir = Join-Path $Root "apps\api"
$webDir = Join-Path $Root "apps\web"

if (!(Test-Path $apiDir)) { throw "API dir not found: $apiDir" }
if (!(Test-Path $webDir)) { throw "Web dir not found: $webDir" }

# --------------------------
# 1) WEB: Use VITE_API_BASE
# --------------------------
Info "Patching web API base..."
$apiTs = Join-Path $webDir "src\lib\api.ts"
if (!(Test-Path $apiTs)) { throw "Missing: $apiTs" }

$webText = Get-Content $apiTs -Raw
# Replace only if it still has localhost hardcoded
if ($webText -match 'const\s+API_BASE\s*=\s*".*localhost:8080\/v1";') {
  $webText = $webText -replace 'const\s+API_BASE\s*=\s*".*localhost:8080\/v1";',
  'const API_BASE = (import.meta.env.VITE_API_BASE ?? "http://localhost:8080") + "/v1";'
  Ok "Updated API_BASE to use VITE_API_BASE."
} else {
  Warn "API_BASE line not matched. Ensure api.ts has the correct API_BASE."
}
Set-Content -Path $apiTs -Value $webText -Encoding utf8

# -----------------------------------
# 2) API: Add Vercel domain to CORS
# -----------------------------------
Info "Patching API CORS allowed origins..."
$mainGo = Join-Path $apiDir "cmd\api\main.go"
if (!(Test-Path $mainGo)) { throw "Missing: $mainGo" }

$goText = Get-Content $mainGo -Raw
# Try to find AllowedOrigins array and inject Vercel domain if not present.
if ($goText -match 'AllowedOrigins:\s*\[\]string\s*\{') {
  if ($goText -notmatch [regex]::Escape("https://$VercelDomain")) {
    $goText = $goText -replace '(AllowedOrigins:\s*\[\]string\s*\{\s*)',
      ('$1' + "`n      ""https://$VercelDomain""," )
    Ok "Added https://$VercelDomain to AllowedOrigins."
  } else {
    Warn "Vercel domain already in AllowedOrigins."
  }
} else {
  Warn "Could not find AllowedOrigins block. You may need to add it manually."
}
Set-Content -Path $mainGo -Value $goText -Encoding utf8

# --------------------------
# 3) Sanity: Build checks
# --------------------------
Info "Running quick checks..."
Push-Location $webDir
npm install | Out-Null
npm run build | Out-Null
Pop-Location
Ok "Web builds."

Push-Location $apiDir
go mod tidy | Out-Null
go build ./cmd/api | Out-Null
Pop-Location
Ok "API builds."

# --------------------------
# 4) Git commit + push
# --------------------------
Info "Git commit + push..."
Push-Location $Root

if (!(Test-Path (Join-Path $Root ".git"))) {
  Warn "No git repo found. Initializing..."
  git init | Out-Null
  git add . | Out-Null
  git commit -m "Initial RecruitFlow" | Out-Null
  Warn "Now add a GitHub remote (see instructions below) then re-run this script."
  Pop-Location
  exit 0
}

git add . | Out-Null
git commit -m "Prepare RecruitFlow for deployment" 2>$null
git push
Pop-Location
Ok "Pushed to GitHub."

# --------------------------
# 5) Print Render/Vercel setup
# --------------------------
Write-Host ""
Write-Host "============================" -ForegroundColor Magenta
Write-Host "PASTE THESE INTO DASHBOARDS" -ForegroundColor Magenta
Write-Host "============================" -ForegroundColor Magenta
Write-Host ""
Write-Host "Render (API Web Service)" -ForegroundColor Cyan
Write-Host "  Repo: your GitHub repo"
Write-Host "  Root Directory: apps/api"
Write-Host "  Build Command:"
Write-Host '    go install github.com/pressly/goose/v3/cmd/goose@latest && goose -dir db/migrations postgres "$DATABASE_URL" up && go build -o server ./cmd/api'
Write-Host "  Start Command:"
Write-Host "    ./server"
Write-Host ""
Write-Host "Render Env Vars:" -ForegroundColor Cyan
Write-Host "  DATABASE_URL = (Render Postgres URL)"
Write-Host "  JWT_SECRET   = (generate a secret)"
Write-Host "  PORT         = (Render sets this automatically)"
Write-Host ""
Write-Host "Vercel (Web)" -ForegroundColor Cyan
Write-Host "  Root Directory: apps/web"
Write-Host "  Env Var:"
Write-Host "    VITE_API_BASE = https://YOUR-RENDER-API.onrender.com"
Write-Host ""
Write-Host "After Vercel gives you a URL, re-run this script with:" -ForegroundColor Yellow
Write-Host "  .\publish.ps1 -VercelDomain YOUR-PROJECT.vercel.app"
Write-Host ""
Ok "Done."
