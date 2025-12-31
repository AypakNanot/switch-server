#!/bin/bash
# measure-memory.sh - Monitor opt-switch memory usage over time
#
# Usage:
#   ./scripts/measure-memory.sh <PID> [INTERVAL]
#
# Example:
#   ./scripts/measure-memory.sh $(pgrep opt-switch) 5
#
# Arguments:
#   PID      - Process ID to monitor (required)
#   INTERVAL - Update interval in seconds (default: 5)
#
# Output format: CSV (Time,RSS_MB,VSZ_MB,%CPU,%MEM)
# Press Ctrl+C to stop monitoring

set -e

# Check arguments
if [ -z "$1" ]; then
    echo "Usage: $0 <PID> [INTERVAL]"
    echo ""
    echo "Example:"
    echo "  $0 \$(pgrep opt-switch) 5"
    echo ""
    echo "Arguments:"
    echo "  PID      - Process ID to monitor (required)"
    echo "  INTERVAL - Update interval in seconds (default: 5)"
    exit 1
fi

PID=$1
INTERVAL=${2:-5}

# Verify process exists
if [ ! -d "/proc/$PID" ]; then
    echo "Error: Process $PID not found"
    echo "Use 'ps aux | grep opt-switch' to find the PID"
    exit 1
fi

# Check if it's actually opt-switch
COMM=$(cat /proc/$PID/comm 2>/dev/null || echo "")
if [[ ! "$COMM" =~ opt-switch|go-admin ]]; then
    echo "Warning: Process $PID appears to be '$COMM', not opt-switch"
    echo "Continuing anyway..."
fi

# Print header
echo "Time,RSS_MB,VSZ_MB,%CPU,%MEM,Goroutines"
echo "======================================"

# Monitor loop
while true; do
    # Check if process still exists
    if [ ! -d "/proc/$PID" ]; then
        echo "$(date '+%H:%M:%S'),Process exited"
        break
    fi

    # Get memory stats from /proc
    # VmRSS: Resident Set Size (physical memory used)
    # VmSize: Virtual Memory Size (total virtual memory)
    stats=$(cat /proc/$PID/status 2>/dev/null | grep -E "VmRSS|VmSize")

    if [ -z "$stats" ]; then
        echo "$(date '+%H:%M:%S'),Error reading stats"
        break
    fi

    # Extract values (in KB)
    rss=$(echo "$stats" | grep VmRSS | awk '{print $2}')
    vsz=$(echo "$stats" | grep VmSize | awk '{print $2}')

    # Convert to MB
    rss_mb=$((rss / 1024))
    vsz_mb=$((vsz / 1024))

    # Get CPU and memory percentage from ps
    ps_stats=$(ps -p $PID -o %cpu,%mem --no-headers 2>/dev/null || echo "0.0 0.0")
    cpu=$(echo "$ps_stats" | awk '{print $1}')
    mem_percent=$(echo "$ps_stats" | awk '{print $2}')

    # Try to get goroutine count (if Go process and accessible)
    goroutines="N/A"
    # This would need to query the process, skip for now

    # Print CSV row
    echo "$(date '+%H:%M:%S'),$rss_mb,$vsz_mb,$cpu,$mem_percent,$goroutines"

    sleep $INTERVAL
done
