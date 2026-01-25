# Device Interaction Layer - Testing Guide

## Quick Start

### 1. Start the Application

```bash
# Option 1: Run directly with go
go run cmd/api/main.go

# Option 2: Run the compiled binary
./go-admin.exe

# Option 3: Build and run
go build -o go-admin.exe ./cmd/api && ./go-admin.exe
```

The application will start on `http://localhost:8000`

### 2. Run the Test Script

```powershell
# Run the automated test script
./test-api.ps1
```

## Manual Testing with curl

### Prerequisites

The device must have SSH service enabled and accessible. Update `config/settings.yml`:

```yaml
device:
  connection:
    host: 192.168.1.1    # Your device IP
    port: 22
    protocol: ssh
    username: admin      # Your username
    password: your-password  # Your password
```

### Test Cases

#### 1. Health Check (No Auth)

```bash
curl http://localhost:8000/api/v1/device
```

Expected response:
```json
{
  "code": 200,
  "data": {
    "status": "online",
    "type": "switch"
  }
}
```

#### 2. Login (Get Token)

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Save the `token` from response for subsequent requests.

#### 3. Get Device Status

```bash
curl http://localhost:8000/api/v1/device/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 4. Execute Single Command

```bash
curl -X POST http://localhost:8000/api/v1/device/command/execute \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "command": "show version",
    "timeout": 10
  }'
```

#### 5. Execute Batch Commands

```bash
curl -X POST http://localhost:8000/api/v1/device/command/batch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commands": [
      "show system",
      "show running-config",
      "show interfaces"
    ],
    "timeout": 30
  }'
```

#### 6. Get Command History

```bash
curl "http://localhost:8000/api/v1/device/command/history?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Troubleshooting

### Device Not Connected

If status shows `"connected": false`:

1. **Check Configuration**: Verify `config/settings.yml` device section
2. **Network**: Ping the device IP address
3. **SSH Service**: Ensure SSH is enabled on device
4. **Credentials**: Verify username and password
5. **Logs**: Check `logs/command.log` for detailed error messages

### Connection Timeout

```
Error: Device connection failed: dial tcp 127.0.0.1:22: connectex: No connection could be made
```

**Solution**: The device is not running on localhost. Update the `host` in settings.yml to your actual device IP.

### Authentication Failed

```
Error: Device connection failed: ssh: handshake failed: ssh: unable to authenticate
```

**Solution**: Verify username and password in settings.yml.

### Command Timeout

```
Error: Command execution timeout
```

**Solution**: Increase `command_timeout` in settings.yml or reduce command complexity.

## Testing Without Real Device

For testing without a real network device, you can:

### Option 1: Use Local SSH Server

Install and start OpenSSH Server on Windows:

```powershell
# Install OpenSSH Server
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0

# Start service
Start-Service sshd
Set-Service -Name sshd -StartupType 'Automatic'

# Confirm firewall rule
Get-NetFirewallRule -Name *ssh*
```

Then update `settings.yml`:
```yaml
device:
  connection:
    host: 127.0.0.1
    port: 22
    username: YOUR_WINDOWS_USERNAME
    password: YOUR_WINDOWS_PASSWORD
```

### Option 2: Use Docker SSH Container

```bash
docker run -d -p 2222:22 \
  -e PASSWORD=admin \
  --name ssh-server \
  panubo/sshd
```

Then update `settings.yml`:
```yaml
device:
  connection:
    host: 127.0.0.1
    port: 2222
    username: root
    password: admin
```

## Expected Results

When properly configured and connected:

1. **Device Status**: Shows `connected: true`
2. **Single Command**: Returns output with `success: true`
3. **Batch Commands**: All commands execute successfully
4. **History**: Shows recent command executions with timestamps

## Log Files

- **Application Logs**: `temp/logs/` directory
- **Command Execution Log**: `logs/command.log` (if enabled in settings)

## Performance Testing

Test concurrent connections:

```powershell
# Run 10 parallel command executions
1..10 | ForEach-Object {
    Start-Job -ScriptBlock {
        param($baseUrl, $token)
        Invoke-RestMethod -Uri "$baseUrl/device/command/execute" -Method POST -Headers @{
            Authorization = "Bearer $token"
        } -Body @{
            command = "show version"
            timeout = 10
        }
    } -ArgumentList "http://localhost:8000/api/v1", $token
} | Wait-Job | Receive-Job
```

Expected behavior:
- First 3 requests succeed (max_connections: 3)
- Remaining requests queue and execute as connections free up
- All requests eventually succeed without errors
