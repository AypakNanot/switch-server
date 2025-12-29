#!/bin/sh
################################################################################
# init-db.sh - Initialize go-admin database in Docker container
################################################################################

# 运行数据库迁移
echo "Initializing go-admin database..."

/app/go-admin migrate -c /app/config/settings.yml

if [ $? -eq 0 ]; then
    echo "Database initialized successfully!"
    echo "Database location: /tmp/go-admin-db.db"
    echo ""
    echo "To start the service:"
    echo "  docker exec -it go-admin-arm64 sh"
    echo "  /app/go-admin server -c /app/config/settings.yml"
else
    echo "Database initialization failed!"
    exit 1
fi
