# Implementation Tasks

## 1. Multi-Architecture Build Targets
- [x] 1.1 Add `build-armv5` target for ARMv5 devices (GOARM=5)
- [x] 1.2 Add `build-armv6` target for ARMv6 devices (GOARM=6, Raspberry Pi 1/Zero)
- [x] 1.3 Add `build-armv7` target for ARMv7 devices (GOARM=7, Raspberry Pi 2/3)
- [x] 1.4 Add `build-arm64` target for ARM64 devices (aarch64, Raspberry Pi 4/5)
- [x] 1.5 Add `build-mips` target for MIPS (big-endian)
- [x] 1.6 Add `build-mipsle` target for MIPS (little-endian, common in routers)
- [x] 1.7 Add `build-mips64` target for MIPS64 (big-endian)
- [x] 1.8 Add `build-mips64le` target for MIPS64 (little-endian)
- [x] 1.9 Add `build-ppc64` target for PowerPC 64-bit
- [x] 1.10 Add `build-all` target to build all architectures
- [x] 1.11 Add `build-switch` target to build most common switch architectures

## 2. Switch-Optimized Configuration
- [x] 2.1 Create `config/settings.switch.yml` based on `settings.yml`
- [x] 2.2 Reduce database connection pool (maxOpenConns: 5, maxIdleConns: 2)
- [x] 2.3 Increase read/write timeouts to reduce connection churn
- [x] 2.4 Set mode to `prod` to reduce debug overhead
- [x] 2.5 Reduce queue pool size (queue.memory.poolSize: 20)
- [x] 2.6 Disable or reduce logging level (level: warn, enableddb: false)
- [x] 2.7 Set stdout logging to avoid file I/O (stdout: '1')

## 3. Deployment Scripts
- [x] 3.1 Create `scripts/deploy-to-switch.sh` for SSH-based deployment
- [x] 3.2 Add command-line arguments: host, user, port, architecture
- [x] 3.3 Implement backup functionality (backup existing binary)
- [x] 3.4 Implement stop/start/restart commands
- [x] 3.5 Add health check after deployment
- [x] 3.6 Create rollback functionality if deployment fails

## 4. Build Optimization
- [x] 4.1 Ensure CGO_ENABLED=0 in all cross-compile targets
- [ ] 4.2 Use UPX compression to reduce binary size (optional - left for user)
- [x] 4.3 Add verification step to check binary with `file` command
- [x] 4.4 Add size output for each compiled binary

## 5. Documentation
- [x] 5.1 Create `docs/switch-deployment.md` with deployment guide
- [x] 5.2 Document supported architectures and typical devices
- [x] 5.3 Document manual deployment steps (SCP, SSH)
- [x] 5.4 Document configuration customization
- [x] 5.5 Document common troubleshooting (out of memory, permission issues)
- [x] 5.6 Document performance expectations and limitations

## 6. Validation & Testing
- [x] 6.1 Cross-compile for ARMv7 (most common for testing)
- [x] 6.2 Verify binary with `file` and `ldd` commands
- [ ] 6.3 Test startup with switch configuration on low-memory VM (256MB)
- [ ] 6.4 Verify database operations work correctly
- [ ] 6.5 Verify API endpoints respond correctly
- [ ] 6.6 Verify static file serving works

## 7. Makefile Integration
- [x] 7.1 Add help target showing all available build targets
- [x] 7.2 Add `list-arch` target to show supported architectures
- [x] 7.3 Ensure all targets are idempotent
- [x] 7.4 Add clean-sw-artifacts target for switch-specific cleaning

## 8. Optional Enhancements
- [ ] 8.1 Add OpenWrt IPK package definition (for OpenWrt-based switches)
- [x] 8.2 Create systemd service file template
- [x] 8.3 Create init.d script template (for older systems)
- [ ] 8.4 Add Docker multi-arch build support
