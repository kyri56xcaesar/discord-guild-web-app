# Variables
$appName = "dadservice.exe"
$src = @("cmd/api/main.go")

# Function to build the application
function Build {
    Write-Host "Building the application..."
    $buildCmd = "go build -o $appName $($src -join ' ')"
    Invoke-Expression $buildCmd
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Build failed!" -ForegroundColor Red
        exit $LASTEXITCODE
    } else {
        Write-Host "Build succeeded!" -ForegroundColor Green
    }
}

# Function to run the application
function Run {
    if (-Not (Test-Path $appName)) {
        Write-Host "Executable not found! Run 'build' first." -ForegroundColor Yellow
        exit 1
    }
    Write-Host "Running the application..."
    Start-Process -NoNewWindow -Wait "./$appName"
}

# Function to clean the build
function Clean {
    Write-Host "Cleaning up..."
    if (Test-Path $appName) {
        Remove-Item $appName
        Write-Host "Executable removed."
    } else {
        Write-Host "Nothing to clean."
    }
}

[string]$Target = "all"


switch ($Target) {
    "build" { Build }
    "run" { Run }
    "clean" { Clean }
    "all"  { Build; Run }
    default { Write-Host "Unknown target: $Target"; exit 1 }
}
