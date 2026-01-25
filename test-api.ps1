# Device Interaction Layer API Test Script
# Run this after starting the go-admin.exe application

$baseUrl = "http://localhost:8000/api/v1"
$token = ""

Write-Host "=== Device Interaction Layer API Test ===" -ForegroundColor Cyan
Write-Host ""

# Step 1: Login to get token
Write-Host "[1/6] Testing Login..." -ForegroundColor Yellow
$loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body @{
    username = "admin"
    password = "admin123"
} -ContentType "application/json"

$token = $loginResponse.data.token
Write-Host "✓ Login successful, token acquired" -ForegroundColor Green
Write-Host "Token: $($token.Substring(0, 20))..." -ForegroundColor DarkGray
Write-Host ""

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# Step 2: Test Device Info (no auth required)
Write-Host "[2/6] Testing Device Info (health check)..." -ForegroundColor Yellow
try {
    $deviceInfo = Invoke-RestMethod -Uri "$baseUrl/device" -Method GET
    Write-Host "✓ Device Info: $($deviceInfo.data.status), $($deviceInfo.data.type)" -ForegroundColor Green
} catch {
    Write-Host "✗ Device Info failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Step 3: Test Device Status
Write-Host "[3/6] Testing Device Status..." -ForegroundColor Yellow
try {
    $statusResponse = Invoke-RestMethod -Uri "$baseUrl/device/status" -Method GET -Headers $headers
    $status = $statusResponse.data
    Write-Host "✓ Device Status:" -ForegroundColor Green
    Write-Host "  Connected: $($status.connected)" -ForegroundColor DarkGray
    Write-Host "  Total Connections: $($status.total_connections)" -ForegroundColor DarkGray
    Write-Host "  Active Connections: $($status.active_connections)" -ForegroundColor DarkGray
    Write-Host "  Queue Size: $($status.queue_size)" -ForegroundColor DarkGray
    Write-Host "  Max Connections: $($status.max_connections)" -ForegroundColor DarkGray

    if (-not $status.connected) {
        Write-Host "  ⚠ Device not connected. Check configuration!" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ Status check failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Step 4: Test Single Command Execution
Write-Host "[4/6] Testing Single Command Execution..." -ForegroundColor Yellow
try {
    $commandResponse = Invoke-RestMethod -Uri "$baseUrl/device/command/execute" -Method POST -Headers $headers -Body @{
        command = "show version"
        timeout = 10
    } | ConvertTo-Json -Depth 10

    $result = $commandResponse | ConvertFrom-Json
    if ($result.data.success) {
        Write-Host "✓ Command executed successfully" -ForegroundColor Green
        Write-Host "  Command: $($result.data.command)" -ForegroundColor DarkGray
        Write-Host "  Duration: $($result.data.duration)ms" -ForegroundColor DarkGray
        Write-Host "  Output (first 100 chars): $($result.data.output.Substring(0, [Math]::Min(100, $result.data.output.Length)))..." -ForegroundColor DarkGray
    } else {
        Write-Host "✗ Command execution failed" -ForegroundColor Red
        Write-Host "  Error: $($result.data.error)" -ForegroundColor DarkGray
    }
} catch {
    Write-Host "✗ Command execution failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Step 5: Test Batch Command Execution
Write-Host "[5/6] Testing Batch Command Execution..." -ForegroundColor Yellow
try {
    $batchResponse = Invoke-RestMethod -Uri "$baseUrl/device/command/batch" -Method POST -Headers $headers -Body @{
        commands = @("show system", "show running-config", "show interfaces")
        timeout = 30
    } | ConvertTo-Json -Depth 10

    $result = $batchResponse | ConvertFrom-Json
    Write-Host "✓ Batch executed: $($result.data.total) commands" -ForegroundColor Green
    Write-Host "  Success: $($result.data.success)" -ForegroundColor DarkGray
    Write-Host "  Failed: $($result.data.failed)" -ForegroundColor DarkGray

    foreach ($item in $result.data.results) {
        $statusIcon = if ($item.success) { "✓" } else { "✗" }
        Write-Host "  $statusIcon $($item.command): $($item.duration)ms" -ForegroundColor DarkGray
    }
} catch {
    Write-Host "✗ Batch execution failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Step 6: Test Command History
Write-Host "[6/6] Testing Command History..." -ForegroundColor Yellow
try {
    $historyResponse = Invoke-RestMethod -Uri "$baseUrl/device/command/history?limit=10&offset=0" -Method GET -Headers $headers
    $history = $historyResponse.data
    Write-Host "✓ History retrieved: $($history.total) records" -ForegroundColor Green

    foreach ($item in $history.history) {
        $statusIcon = if ($item.success) { "✓" } else { "✗" }
        Write-Host "  $statusIcon [$($item.timestamp)] $($item.username): $($item.command)" -ForegroundColor DarkGray
    }
} catch {
    Write-Host "✗ History retrieval failed: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Test Complete ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Note: If connection failed, check:" -ForegroundColor Yellow
Write-Host "  1. Device configuration in config/settings.yml" -ForegroundColor DarkGray
Write-Host "  2. SSH service is running on the target device" -ForegroundColor DarkGray
Write-Host "  3. Network connectivity to the device" -ForegroundColor DarkGray
Write-Host "  4. Correct username and password" -ForegroundColor DarkGray
