#!/bin/bash
################################################################################
# deploy-to-switch.sh - 部署 go-admin 到交换机/嵌入式设备
################################################################################
#
# 用法:
#   ./scripts/deploy-to-switch.sh --host=192.168.1.1 --user=root --arch=armv7
#   ./scripts/deploy-to-switch.sh --host=192.168.1.1 --action=restart
#   ./scripts/deploy-to-switch.sh --host=192.168.1.1 --action=stop
#
################################################################################

set -e

# 默认值
HOST=""
USER="root"
PORT="22"
ARCH=""
ACTION="deploy"
BINARY="./go-admin"
REMOTE_PATH="/usr/bin"
REMOTE_CONFIG="/etc/go-admin"
REMOTE_DB="/tmp/go-admin-db.db"
SERVICE_NAME="go-admin"
SSH_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
Usage: $0 [OPTIONS]

部署 go-admin 到网络交换机或嵌入式设备

OPTIONS:
  --host HOST       目标设备 IP 地址 (必需)
  --user USER       SSH 用户名 (默认: root)
  --port PORT       SSH 端口 (默认: 22)
  --arch ARCH       目标架构 (armv5, armv6, armv7, arm64, mips, mipsle, 等)
  --binary BINARY    本地二进制文件路径 (默认: ./go-admin-<arch>)
  --action ACTION   操作类型: deploy(默认), restart, stop, status, rollback

ACTIONS:
  deploy   部署新版本（备份、上传、启动）
  restart  重启服务
  stop     停止服务
  status   查看服务状态
  rollback 回滚到上一个版本

EXAMPLES:
  # 部署到 ARMv7 设备
  $0 --host=192.168.1.1 --arch=armv7

  # 部署自定义二进制
  $0 --host=192.168.1.1 --arch=armv7 --binary=./custom-go-admin

  # 仅重启服务
  $0 --host=192.168.1.1 --action=restart

  # 查看服务状态
  $0 --host=192.168.1.1 --action=status

  # 回滚到上一版本
  $0 --host=192.168.1.1 --action=rollback

SUPPORTED ARCHITECTURES:
  ARM: armv5, armv6, armv7, arm64
  MIPS: mips, mipsle, mips64, mips64le
  PowerPC: ppc64, ppc64le

EOF
}

# 解析命令行参数
parse_args() {
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi

    while [[ $# -gt 0 ]]; do
        case $1 in
            --host)
                HOST="$2"
                shift 2
                ;;
            --user)
                USER="$2"
                shift 2
                ;;
            --port)
                PORT="$2"
                shift 2
                ;;
            --arch)
                ARCH="$2"
                shift 2
                ;;
            --binary)
                BINARY="$2"
                shift 2
                ;;
            --action)
                ACTION="$2"
                shift 2
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                log_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # 验证必需参数
    if [ -z "$HOST" ]; then
        log_error "必须指定 --host 参数"
        show_help
        exit 1
    fi

    # 如果没有指定 binary，根据 arch 推断
    if [ "$ACTION" = "deploy" ] && [ -z "$BINARY" ]; then
        if [ -z "$ARCH" ]; then
            log_error "deploy 操作需要指定 --arch 或 --binary 参数"
            show_help
            exit 1
        fi
        BINARY="./go-admin-$ARCH"
        if [ ! -f "$BINARY" ]; then
            # 尝试不带架构后缀的文件名
            if [ ! -f "./go-admin" ]; then
                log_error "找不到二进制文件: $BINARY"
                log_error "请先编译: make build-$ARCH"
                exit 1
            fi
            BINARY="./go-admin"
        fi
    fi
}

# SSH 执行命令
ssh_exec() {
    ssh -p "$PORT" "$USER@$HOST" $SSH_OPTS "$@"
}

# 检查服务状态
check_status() {
    log_info "检查服务状态..."

    # 检查进程是否运行
    if ssh_exec "pgrep -x $SERVICE_NAME > /dev/null" 2>/dev/null; then
        log_info "服务正在运行"

        # 显示进程信息
        ssh_exec "ps | grep -E '[PID|$SERVICE_NAME]' | head -5" 2>/dev/null || true

        # 显示内存使用
        ssh_exec "free" 2>/dev/null || true

        # 检查端口监听
        ssh_exec "netstat -tlnp 2>/dev/null | grep :8000" || ssh_exec "netstat -an 2>/dev/null | grep :8000" || true

        return 0
    else
        log_warn "服务未运行"
        return 1
    fi
}

# 停止服务
stop_service() {
    log_info "停止服务..."

    # 尝试优雅停止
    ssh_exec "PID=\$(pgrep -x $SERVICE_NAME); [ -n \"\$PID\" ] && kill -TERM \$PID && sleep 2 && kill -0 \$PID 2>/dev/null && kill -KILL \$PID; true" 2>/dev/null

    # 等待进程结束
    for i in {1..10}; do
        if ssh_exec "! pgrep -x $SERVICE_NAME > /dev/null" 2>/dev/null; then
            log_info "服务已停止"
            return 0
        fi
        sleep 1
    done

    log_warn "服务可能仍在运行"
}

# 启动服务
start_service() {
    log_info "启动服务..."

    # 创建配置目录（如果不存在）
    ssh_exec "mkdir -p $REMOTE_CONFIG 2>/dev/null || true"

    # 启动服务
    ssh_exec "nohup $REMOTE_PATH/$SERVICE_NAME server -c $REMOTE_CONFIG/settings.yml > /tmp/$SERVICE_NAME.log 2>&1 &"

    # 等待服务启动
    sleep 2

    # 检查是否启动成功
    if ssh_exec "pgrep -x $SERVICE_NAME > /dev/null" 2>/dev/null; then
        log_info "服务启动成功"

        # 显示日志
        log_info "最近的日志:"
        ssh_exec "tail -20 /tmp/$SERVICE_NAME.log" 2>/dev/null || true

        return 0
    else
        log_error "服务启动失败"
        ssh_exec "cat /tmp/$SERVICE_NAME.log" 2>/dev/null || true
        return 1
    fi
}

# 部署
deploy() {
    log_info "开始部署到 $HOST:$PORT (架构: $ARCH)..."

    # 检查本地二进制文件
    if [ ! -f "$BINARY" ]; then
        log_error "找不到二进制文件: $BINARY"
        exit 1
    fi

    # 显示二进制信息
    log_info "二进制文件信息:"
    file "$BINARY"
    ls -lh "$BINARY"

    # 测试 SSH 连接
    log_info "测试 SSH 连接..."
    if ! ssh_exec "echo 'Connection successful'" 2>/dev/null; then
        log_error "SSH 连接失败"
        exit 1
    fi

    # 检查目标设备架构
    log_info "检查目标设备架构..."
    REMOTE_ARCH=$(ssh_exec "uname -m" 2>/dev/null || echo "unknown")
    log_info "目标设备架构: $REMOTE_ARCH"

    # 备份现有二进制
    log_info "备份现有二进制..."
    ssh_exec "[ -f $REMOTE_PATH/$SERVICE_NAME ] && cp $REMOTE_PATH/$SERVICE_NAME $REMOTE_PATH/$SERVICE_NAME.backup || echo 'No existing binary to backup'" 2>/dev/null

    # 上传新二进制
    log_info "上传新二进制..."
    scp -P "$PORT" "$BINARY" "$USER@$HOST:$REMOTE_PATH/$SERVICE_NAME" 2>/dev/null || {
        log_error "上传失败"
        exit 1
    }

    # 设置权限
    log_info "设置权限..."
    ssh_exec "chmod +x $REMOTE_PATH/$SERVICE_NAME"

    # 停止现有服务
    if ssh_exec "pgrep -x $SERVICE_NAME > /dev/null" 2>/dev/null; then
        log_info "停止现有服务..."
        stop_service
    fi

    # 启动新服务
    if start_service; then
        log_info "部署成功！"

        # 显示服务状态
        echo ""
        check_status

        # 访问提示
        echo ""
        log_info "应用已部署，可通过以下方式访问:"
        echo "   Web UI: http://$HOST:8000/"
        echo "   API: http://$HOST:8000/api/v1/"
        echo ""
        log_info "查看日志: ssh $USER@$HOST 'tail -f /tmp/$SERVICE_NAME.log'"

        return 0
    else
        # 部署失败，回滚
        log_error "部署失败，回滚..."
        rollback
        return 1
    fi
}

# 重启服务
restart() {
    log_info "重启服务..."

    stop_service
    start_service

    echo ""
    check_status
}

# 查看状态
status() {
    check_status
}

# 回滚
rollback() {
    log_info "回滚到上一版本..."

    # 检查备份文件是否存在
    if ! ssh_exec "[ -f $REMOTE_PATH/$SERVICE_NAME.backup ]" 2>/dev/null; then
        log_error "没有找到备份文件"
        exit 1
    fi

    # 停止当前服务
    stop_service

    # 恢复备份
    log_info "恢复备份文件..."
    ssh_exec "cp $REMOTE_PATH/$SERVICE_NAME.backup $REMOTE_PATH/$SERVICE_NAME"
    ssh_exec "chmod +x $REMOTE_PATH/$SERVICE_NAME"

    # 启动服务
    if start_service; then
        log_info "回滚成功"
        check_status
    else
        log_error "回滚失败"
        exit 1
    fi
}

# 主函数
main() {
    parse_args "$@"

    case "$ACTION" in
        deploy)
            deploy
            ;;
        restart)
            restart
            ;;
        stop)
            stop_service
            ;;
        status)
            status
            ;;
        rollback)
            rollback
            ;;
        *)
            log_error "未知操作: $ACTION"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
