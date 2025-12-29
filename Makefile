PROJECT:=go-admin

# Default build with SQLite support (pure Go, no CGO required)
.PHONY: build
build:
	go build -ldflags="-w -s" -o go-admin .

# Build for Linux (Docker)
build-linux:
	@docker build -t go-admin:latest .
	@echo "build successful"

# Build for specific platforms (pure Go SQLite, no CGO needed)
build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o go-admin-linux-amd64 .

build-windows:
	env GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o go-admin.exe .

build-mac:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o go-admin-mac .

build-mac-arm:
	env GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o go-admin-mac-arm .

# Build without any special flags
build-std:
	go build -ldflags="-w -s" -o go-admin .

# =============================================================================
# Switch Device / Embedded System Build Targets
# =============================================================================
# Common switch architectures (ARMv7, ARM64, MIPSLE)
.PHONY: build-switch
build-switch: build-armv7 build-arm64 build-mipsle
	@echo ""
	@echo "Switch binaries built:"
	@ls -lh go-admin-armv7 go-admin-arm64 go-admin-mipsle 2>/dev/null || echo "Some builds may have failed"

# ARM variants (32-bit)
.PHONY: build-armv5
build-armv5:
	@echo "Building for ARMv5..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 \
		go build -ldflags="-w -s" -o go-admin-armv5 .
	@file go-admin-armv5
	@ls -lh go-admin-armv5

.PHONY: build-armv6
build-armv6:
	@echo "Building for ARMv6 (Raspberry Pi 1/Zero)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 \
		go build -ldflags="-w -s" -o go-admin-armv6 .
	@file go-admin-armv6
	@ls -lh go-admin-armv6

.PHONY: build-armv7
build-armv7:
	@echo "Building for ARMv7 (Raspberry Pi 2/3, most switches)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 \
		go build -ldflags="-w -s" -o go-admin-armv7 .
	@file go-admin-armv7
	@ls -lh go-admin-armv7

# ARM 64-bit
.PHONY: build-arm64
build-arm64:
	@echo "Building for ARM64 (Raspberry Pi 4/5, newer switches)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -ldflags="-w -s" -o go-admin-arm64 .
	@file go-admin-arm64
	@ls -lh go-admin-arm64

# MIPS variants (common in routers/switches)
.PHONY: build-mips
build-mips:
	@echo "Building for MIPS (big-endian)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=mips \
		go build -ldflags="-w -s" -o go-admin-mips .
	@file go-admin-mips
	@ls -lh go-admin-mips

.PHONY: build-mipsle
build-mipsle:
	@echo "Building for MIPSLE (little-endian, most common)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle \
		go build -ldflags="-w -s" -o go-admin-mipsle .
	@file go-admin-mipsle
	@ls -lh go-admin-mipsle

.PHONY: build-mips64
build-mips64:
	@echo "Building for MIPS64 (big-endian)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=mips64 \
		go build -ldflags="-w -s" -o go-admin-mips64 .
	@file go-admin-mips64
	@ls -lh go-admin-mips64

.PHONY: build-mips64le
build-mips64le:
	@echo "Building for MIPS64LE (little-endian)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=mips64le \
		go build -ldflags="-w -s" -o go-admin-mips64le .
	@file go-admin-mips64le
	@ls -lh go-admin-mips64le

# PowerPC variants
.PHONY: build-ppc64
build-ppc64:
	@echo "Building for PPC64 (PowerPC 64-bit)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 \
		go build -ldflags="-w -s" -o go-admin-ppc64 .
	@file go-admin-ppc64
	@ls -lh go-admin-ppc64

.PHONY: build-ppc64le
build-ppc64le:
	@echo "Building for PPC64LE (PowerPC 64-bit little-endian)..."
	@env CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le \
		go build -ldflags="-w -s" -o go-admin-ppc64le .
	@file go-admin-ppc64le
	@ls -lh go-admin-ppc64le

# Build all supported architectures
.PHONY: build-all
build-all: build-armv5 build-armv6 build-armv7 build-arm64 build-mips build-mipsle build-mips64 build-mips64le build-ppc64 build-ppc64le
	@echo ""
	@echo "All architectures built successfully"
	@ls -lh go-admin-* | grep -v "\.exe"

# =============================================================================
# Utility Targets
# =============================================================================
.PHONY: list-arch
list-arch:
	@echo "Supported architectures for switch/embedded device deployment:"
	@echo ""
	@echo "ARM (32-bit):"
	@echo "  make build-armv5    - ARMv5 (old devices)"
	@echo "  make build-armv6    - ARMv6 (Raspberry Pi 1/Zero)"
	@echo "  make build-armv7    - ARMv7 (Raspberry Pi 2/3, most switches)"
	@echo ""
	@echo "ARM (64-bit):"
	@echo "  make build-arm64     - ARM64 (Raspberry Pi 4/5, newer switches)"
	@echo ""
	@echo "MIPS:"
	@echo "  make build-mips      - MIPS big-endian"
	@echo "  make build-mipsle    - MIPS little-endian (common in routers)"
	@echo "  make build-mips64    - MIPS64 big-endian"
	@echo "  make build-mips64le  - MIPS64 little-endian"
	@echo ""
	@echo "PowerPC:"
	@echo "  make build-ppc64     - PPC64 big-endian"
	@echo "  make build-ppc64le   - PPC64 little-endian"
	@echo ""
	@echo "Combined targets:"
	@echo "  make build-switch    - Build common switch architectures (armv7, arm64, mipsle)"
	@echo "  make build-all       - Build all supported architectures"

.PHONY: help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Standard builds:"
	@echo "  make build           - Build for current platform"
	@echo "  make build-linux-amd64 - Build for Linux x86_64"
	@echo "  make build-windows   - Build for Windows"
	@echo "  make build-mac       - Build for macOS x86_64"
	@echo "  make build-mac-arm   - Build for macOS ARM64"
	@echo ""
	@echo "Switch/embedded builds:"
	@echo "  make build-switch    - Build common switch architectures"
	@echo "  make build-all       - Build all architectures"
	@echo "  make list-arch       - List all supported architectures"
	@echo ""
	@echo "Individual architectures:"
	@echo "  make build-armv5/6/7/64 - ARM variants"
	@echo "  make build-mips/mipsle/mips64/mips64le - MIPS variants"
	@echo "  make build-ppc64/ppc64le - PowerPC variants"

# Clean switch build artifacts
.PHONY: clean-sw
clean-sw:
	@echo "Cleaning switch build artifacts..."
	@rm -f go-admin-armv5 go-admin-armv6 go-admin-armv7 go-admin-arm64
	@rm -f go-admin-mips go-admin-mipsle go-admin-mips64 go-admin-mips64le
	@rm -f go-admin-ppc64 go-admin-ppc64le
	@echo "Clean completed"

# make run
run:
    # delete go-admin-api container
	@if [ $(shell docker ps -aq --filter name=go-admin --filter publish=8000) ]; then docker rm -f go-admin; fi

    # 启动方法一 run go-admin-api container  docker-compose 启动方式
    # 进入到项目根目录 执行 make run 命令
	@docker-compose up -d

	# 启动方式二 docker run  这里注意-v挂载的宿主机的地址改为部署时的实际决对路径
    #@docker run --name=go-admin -p 8000:8000 -v /home/code/go/src/go-admin/go-admin/config:/go-admin-api/config  -v /home/code/go/src/go-admin/go-admin-api/static:/go-admin/static -v /home/code/go/src/go-admin/go-admin/temp:/go-admin-api/temp -d --restart=always go-admin:latest

	@echo "go-admin service is running..."

	# delete Tag=<none> 的镜像
	@docker image prune -f
	@docker ps -a | grep "go-admin"

stop:
    # delete go-admin-api container
	@if [ $(shell docker ps -aq --filter name=go-admin --filter publish=8000) ]; then docker-compose down; fi
	#@if [ $(shell docker ps -aq --filter name=go-admin --filter publish=8000) ]; then docker rm -f go-admin; fi
	#@echo "go-admin stop success"


#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t go-admin:latest

# make deploy
deploy:

	#@git checkout master
	#@git pull origin master
	make build-linux
	make run
