#!/bin/bash
# memory-test.sh - Automated memory testing for opt-switch
#
# Usage:
#   ./scripts/memory-test.sh <config_file> [test_duration]
#
# Example:
#   ./scripts/memory-test.sh config/settings.minimal.yml 60
#
# This script:
# 1. Checks available memory
# 2. Starts opt-switch with specified config
# 3. Monitors memory usage
# 4. Stops the service
# 5. Reports statistics

set -e

# Check arguments
if [ -z "$1" ]; then
    echo "Usage: $0 <config_file> [test_duration]"
    echo ""
    echo "Example:"
    echo "  $0 config/settings.minimal.yml 60"
    echo ""
    echo "Arguments:"
    echo "  config_file    - Path to config file (required)"
    echo "  test_duration  - Test duration in seconds (default: 60)"
    exit 1
fi

CONFIG_FILE=$1
TEST_DURATION=${2:-60}
BINARY="./opt-switch"

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    # Try with .exe extension (Windows)
    if [ ! -f "${BINARY}.exe" ]; then
        echo "Error: Binary not found: $BINARY"
        echo "Please build opt-switch first"
        exit 1
    fi
    BINARY="${BINARY}.exe"
fi

# Check if config exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Config file not found: $CONFIG_FILE"
    exit 1
fi

echo "======================================"
echo "opt-switch Memory Test"
echo "======================================"
echo "Config:    $CONFIG_FILE"
echo "Duration:  ${TEST_DURATION}s"
echo "Binary:    $BINARY"
echo "======================================"
echo ""

# 1. Check available memory
echo "[1/5] Checking available memory..."
./scripts/check-memory.sh 128
if [ $? -ne 0 ]; then
    echo "Warning: Memory check failed, but continuing..."
fi
echo ""

# 2. Start opt-switch in background
echo "[2/5] Starting opt-switch..."
$BINARY server -c "$CONFIG_FILE" &
PID=$!
echo "Started with PID: $PID"
sleep 3

# Check if process is still running
if ! kill -0 $PID 2>/dev/null; then
    echo "Error: Process failed to start"
    exit 1
fi
echo "✓ Process running"
echo ""

# 3. Monitor memory for test duration
echo "[3/5] Monitoring memory for ${TEST_DURATION}s..."
echo "Time,RSS_MB,VSZ_MB,%CPU,%MEM"

# Create temp file for results
TEMP_FILE=$(mktemp)
trap "rm -f $TEMP_FILE" EXIT

# Run measurement in background
./scripts/measure-memory.sh $PID 5 > $TEMP_FILE &
MEASURE_PID=$!

# Wait for test duration
sleep $TEST_DURATION

# Stop measurement
kill $MEASURE_PID 2>/dev/null || true
wait $MEASURE_PID 2>/dev/null || true

# 4. Stop opt-switch
echo ""
echo "[4/5] Stopping opt-switch..."
kill $PID 2>/dev/null || true
wait $PID 2>/dev/null || true
echo "✓ Stopped"
echo ""

# 5. Analyze results
echo "[5/5] Analyzing results..."
echo "======================================"
echo "Memory Usage Statistics"
echo "======================================"

# Calculate stats from CSV
# Skip header and process data
tail -n +2 $TEMP_FILE | awk -F',' '
{
    rss += $2
    vsz += $3
    cpu += $4
    mem += $5
    count++
    if (NR == 2 || $2 > max_rss) max_rss = $2
    if (NR == 2 || $2 < min_rss) min_rss = $2
}
END {
    if (count > 0) {
        printf "  RSS (Resident Set Size):\n"
        printf "    Average: %d MB\n", rss / count
        printf "    Min:     %d MB\n", min_rss
        printf "    Max:     %d MB\n", max_rss
        printf "\n"
        printf "  VSZ (Virtual Memory):\n"
        printf "    Average: %d MB\n", vsz / count
        printf "\n"
        printf "  CPU Usage:\n"
        printf "    Average: %.1f%%\n", cpu / count
        printf "\n"
        printf "  Memory %%:\n"
        printf "    Average: %.1f%%\n", mem / count
        printf "\n"
        printf "  Samples: %d\n", count
    }
}'

echo "======================================"
echo ""
echo "Test complete!"
echo ""

# Check if target memory was met
AVG_RSS=$(tail -n +2 $TEMP_FILE | awk -F',' '{sum+=$2; count++} END {printf "%.0f", sum/count}')
echo "Average RSS: ${AVG_RSS} MB"

if [ "$AVG_RSS" -lt 50 ]; then
    echo "✓ Target met (< 50MB)"
    exit 0
elif [ "$AVG_RSS" -lt 60 ]; then
    echo "⚠ Close to target (50-60MB)"
    exit 0
else
    echo "✗ Target not met (> 60MB)"
    echo ""
    echo "Suggestions:"
    echo "  1. Review settings.minimal.yml configuration"
    echo "  2. Ensure enableFrontend: false"
    echo "  3. Check if middleware are disabled"
    echo "  4. Verify runtime settings are applied"
    exit 1
fi
