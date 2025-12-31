#!/bin/bash
# check-memory.sh - Check available memory before starting opt-switch
#
# Usage:
#   ./scripts/check-memory.sh [REQUIRED_MB]
#
# Example:
#   ./scripts/check-memory.sh 128
#
# This script checks if the system has enough available memory
# to run opt-switch. It's useful to run before starting the service.

set -e

# Default required memory (in MB)
REQUIRED_MB=${1:-128}
REQUIRED_KB=$((REQUIRED_MB * 1024))

echo "======================================"
echo "opt-switch Memory Check"
echo "======================================"
echo ""

# Get memory info (works on Linux)
if [ ! -f /proc/meminfo ]; then
    echo "Error: /proc/meminfo not found"
    echo "This script requires Linux with /proc filesystem"
    exit 1
fi

# Read memory values (in KB)
mem_total=$(grep MemTotal /proc/meminfo | awk '{print $2}')
mem_free=$(grep MemFree /proc/meminfo | awk '{print $2}')
mem_available=$(grep MemAvailable /proc/meminfo | awk '{print $2}')

# Fallback if MemAvailable is not available (older kernels)
if [ -z "$mem_available" ] || [ "$mem_available" = "0" ]; then
    # Use MemFree + Buffers + Cached as approximation
    mem_buffers=$(grep Buffers /proc/meminfo | awk '{print $2}')
    mem_cached=$(grep ^Cached /proc/meminfo | awk '{print $2}')
    mem_available=$((mem_free + mem_buffers + mem_cached))
fi

# Calculate swap (also consider available swap)
swap_total=$(grep SwapTotal /proc/meminfo | awk '{print $2}')
swap_free=$(grep SwapFree /proc/meminfo | awk '{print $2}')

# Convert to MB for display
total_mb=$((mem_total / 1024))
free_mb=$((mem_free / 1024))
available_mb=$((mem_available / 1024))
swap_total_mb=$((swap_total / 1024))
swap_free_mb=$((swap_free / 1024))

# Print memory status
echo "Memory Status:"
echo "  Total RAM:     ${total_mb} MB"
echo "  Available RAM: ${available_mb} MB"
echo "  Free RAM:      ${free_mb} MB"
if [ "$swap_total_mb" -gt 0 ]; then
    echo "  Total Swap:    ${swap_total_mb} MB"
    echo "  Free Swap:     ${swap_free_mb} MB"
fi
echo ""
echo "Required: ${REQUIRED_MB} MB"
echo "======================================"
echo ""

# Check if enough memory
if [ "$mem_available" -lt "$REQUIRED_KB" ]; then
    deficit_mb=$(( (REQUIRED_KB - mem_available) / 1024 ))
    echo "❌ ERROR: Not enough memory!"
    echo "   Deficit: ~${deficit_mb} MB"
    echo ""
    echo "Recommendations:"
    echo "  1. Close other applications"
    echo "  2. Use settings.minimal.yml for lower memory usage"
    echo "  3. Add swap space"
    echo "  4. Upgrade to a device with more RAM"
    echo ""
    exit 1
else
    surplus_mb=$(( (mem_available - REQUIRED_KB) / 1024 ))
    echo "✓ OK: Sufficient memory available"
    echo "  Surplus: ~${surplus_mb} MB"
    echo ""
    echo "You can safely start opt-switch"
    echo ""
    exit 0
fi
