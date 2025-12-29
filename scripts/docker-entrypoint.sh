#!/bin/sh
################################################################################
# docker-entrypoint.sh - go-admin Docker entrypoint
# 自动检测并初始化数据库
################################################################################

set -e

CONFIG_FILE="${CONFIG_FILE:-/app/config/settings.yml}"
DB_FILE="/tmp/go-admin-db.db"

echo "=================================================="
echo "  go-admin Docker Container"
echo "  Version: 2.2.0"
echo "  Platform: ARM64"
echo "=================================================="
echo ""

# 检查数据库是否需要初始化
if [ ! -f "$DB_FILE" ] || [ ! -s "$DB_FILE" ]; then
    echo "📦 Database not found. Initializing..."
    /app/go-admin migrate -c "$CONFIG_FILE"

    if [ $? -eq 0 ]; then
        echo "✅ Database initialized successfully!"
    else
        echo "❌ Database initialization failed!"
        exit 1
    fi
else
    echo "📊 Database exists at $DB_FILE"
    echo "   Checking for updates..."
    /app/go-admin migrate -c "$CONFIG_FILE" || echo "⚠️  Migration completed with warnings"
fi

echo ""
echo "🚀 Starting go-admin server..."
echo "   Config: $CONFIG_FILE"
echo "   Port: 8000"
echo ""

# 启动服务
exec /app/go-admin server -c "$CONFIG_FILE"
