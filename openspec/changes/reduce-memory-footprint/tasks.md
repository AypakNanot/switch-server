# Tasks: Reduce Memory Footprint to 50MB

## Phase 1: Configuration Files

### 1.1 Create minimal configuration
- [ ] Create `config/settings.minimal.yml`
  - maxOpenConns: 2, maxIdleConns: 1
  - poolSize: 5
  - enableFrontend: false
  - Enable middleware flags (all false)
  - Document when to use this config

### 1.2 Update existing configurations
- [ ] Add comments to `config/settings.switch.yml` about memory options
- [ ] Document memory requirements in each config

## Phase 2: Runtime Memory Tuning

### 2.1 Add runtime initialization
- [ ] Create `cmd/api/runtime.go` with `initRuntime()` function
  - Set GOMAXPROCS=1 (configurable)
  - Set GOGC=200 (configurable)
  - SetMemoryLimit=60MB (Go 1.19+)
  - SetMaxThreads=100

### 2.2 Integrate runtime tuning
- [ ] Call `initRuntime()` in `main()` before server starts
- [ ] Add config options for runtime tuning
  - runtime.gomaxprocs
  - runtime.gogc
  - runtime.memoryLimit

## Phase 3: Conditional Features

### 3.1 Conditional frontend routes
- [ ] Add `enableFrontend` config option
- [ ] Wrap frontend route registration in if-block
- [ ] Test with and without frontend

### 3.2 Conditional middleware
- [ ] Add `enableMiddleware` config section
  - sentinel
  - requestID
  - metrics
  - logger (keep always)
- [ ] Update `cmd/api/server.go` to conditionally enable
- [ ] Test each middleware independently

## Phase 4: Tools and Scripts

### 4.1 Memory measurement script
- [ ] Create `scripts/measure-memory.sh`
  - Monitor RSS/VSZ of opt-switch process
  - Output CSV format
  - Support custom interval

### 4.2 Memory check script
- [ ] Create `scripts/check-memory.sh`
  - Check available memory before start
  - Compare against required memory
  - Exit with error if insufficient

### 4.3 Memory monitoring helper
- [ ] Create `scripts/memory-test.sh`
  - Start opt-switch with config
  - Monitor memory for 60 seconds
  - Report average, min, max RSS

## Phase 5: Code Changes

### 5.1 Database connection optimization
- [ ] Review `common/database/initialize.go`
- [ ] Ensure connection pool settings are respected
- [ ] Add connection pool metrics (optional)

### 5.2 Queue optimization
- [ ] Review `common/storage/initialize.go`
- [ ] Ensure queue poolSize is respected
- [ ] Test with poolSize=5

### 5.3 Gin engine optimization
- [ ] Review `cmd/api/server.go` initRouter()
- [ ] Remove unnecessary middleware for minimal config
- [ ] Consider reducing buffer sizes

## Phase 6: Testing

### 6.1 Baseline measurement
- [ ] Measure memory with current `settings.switch.yml`
  - Record idle RSS
  - Record under load (1, 5, 10 concurrent users)

### 6.2 Minimal config testing
- [ ] Start with `settings.minimal.yml`
  - Measure idle RSS (target: < 50MB)
  - Test core functionality (login, menu, CRUD)
  - Measure under load

### 6.3 Docker testing
- [ ] Build Docker image with minimal config
  - Measure container memory usage
  - Verify functionality

### 6.4 Comparison table
- [ ] Document memory usage for each config:
  - settings.yml (full)
  - settings.switch.yml (normal)
  - settings.minimal.yml (minimal)

## Phase 7: Documentation

### 7.1 Update README
- [ ] Add memory optimization section
- [ ] Document when to use each config
- [ ] Add memory measurement instructions

### 7.2 Update deployment guides
- [ ] Update `deploy/docs/DEPLOYMENT_GUIDE.md`
- [ ] Update `deploy/docs/MANUAL_DEPLOYMENT.md`
- [ ] Add memory requirements per config

### 7.3 Create tuning guide
- [ ] Document runtime tuning options
- [ ] Provide examples for different scenarios
- [ ] Add troubleshooting for OOM issues

## Phase 8: Validation

### 8.1 Success criteria verification
- [ ] ✅ Idle memory < 50MB with minimal config
- [ ] ✅ Core functions work (login, menu, CRUD)
- [ ] ✅ Docker image tested
- [ ] ✅ Documentation complete

### 8.2 Performance testing
- [ ] Response time comparison (normal vs minimal)
- [ ] Throughput comparison
- [ ] GC pause time measurement

### 8.3 Stability testing
- [ ] Run for 24 hours with minimal config
- [ ] Monitor for memory leaks
- [ ] Check GC logs

## Dependencies

- **Phase 1** must complete before Phase 3
- **Phase 2** must complete before Phase 6
- **Phase 4** can be done in parallel with Phase 3
- **Phase 7** depends on Phase 6 results

## Estimated Effort

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| 1: Configuration | 2 | 30 min |
| 2: Runtime Tuning | 2 | 1 hour |
| 3: Conditional Features | 3 | 2 hours |
| 4: Tools | 3 | 1 hour |
| 5: Code Changes | 3 | 1 hour |
| 6: Testing | 4 | 2 hours |
| 7: Documentation | 3 | 1 hour |
| 8: Validation | 3 | 2 hours |
| **Total** | **23** | **~10 hours** |
