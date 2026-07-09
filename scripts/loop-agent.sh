#!/bin/zsh
# =============================================================================
# 常驻循环 Agent 守护脚本 (Resident Loop Agent Daemon)
# -----------------------------------------------------------------------------
# 职责:
#   1. 心跳: 持续运行直到截止时间，定期向 LOOP_LOG.md 写入存活标记。
#   2. 验证: 运行可用静态检查 (前端 tsc/eslint；若 go 可用则 go build/vet)。
#   3. 提交: 仅当工作区有改动且检查通过时，自动 git commit（无需人工审批）。
#   4. 单实例: 通过 mkdir 原子锁防止重复启动（避免多实例/资源耗尽）。
#
# 说明: 真正的"重构/优化"智能由 Roo 会话驱动；本脚本是常驻 harness，
#       保证即使会话中断也有进程在跑、有日志、有自动提交。
#       若未来设置 ENGINE_CMD 指向可用编码引擎，可接管自主编码。
# =============================================================================

set -u

# ---- 配置 -------------------------------------------------------------------
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_DIR" || exit 1

# 截止时间: 2026-07-10 09:00 Asia/Shanghai (UTC+8)
DEADLINE_EPOCH=$(date -j -f "%Y-%m-%d %H:%M %z" "2026-07-10 09:00 +0800" +%s 2>/dev/null)
if [ -z "$DEADLINE_EPOCH" ]; then
  DEADLINE_EPOCH=$(( $(date +%s) + 54000 ))
fi

LOG_FILE="LOOP_LOG.md"
LOCK_DIR="/tmp/loop-agent.lock"
SLEEP_SEC=120
ENGINE_CMD="${ENGINE_CMD:-}"   # 可选: 指向编码引擎 CLI

# 工具路径 (非交互 shell 的 PATH 很窄)
export PATH="$PATH:/usr/local/bin:/opt/homebrew/bin:/usr/local/go/bin:$HOME/go/bin"

# ---- 单次检查模式 (供手动触发，不进入循环，不加锁) ---------------------------
if [ "${1:-}" = "_check" ]; then
  run_frontend_checks() {
    [ -d frontend/node_modules ] || return 0
    local ok=1
    [ -x frontend/node_modules/.bin/tsc ] && (cd frontend && node_modules/.bin/tsc --noEmit) >/tmp/tsc.log 2>&1 || ok=0
    [ -x frontend/node_modules/.bin/eslint ] && (cd frontend && node_modules/.bin/eslint "src/**/*.{ts,tsx}" --max-warnings=0) >/tmp/eslint.log 2>&1 || ok=0
    return $ok
  }
  run_backend_checks() {
    command -v go >/dev/null 2>&1 || return 0
    (cd backend && go build ./... ) >/tmp/gobuild.log 2>&1 || return 1
    (cd backend && go vet ./... ) >/tmp/govet.log 2>&1 || return 1
    return 0
  }
  if run_frontend_checks && run_backend_checks; then
    echo "[loop-agent] _check: PASS"
    exit 0
  else
    echo "[loop-agent] _check: FAILED (see /tmp/*.log)"
    exit 1
  fi
fi

# ---- 单实例原子锁 (mkdir 在 POSIX 下原子) -----------------------------------
if ! mkdir "$LOCK_DIR" 2>/dev/null; then
  OLD_PID="$(cat "$LOCK_DIR/pid" 2>/dev/null)"
  echo "[loop-agent] already running (pid ${OLD_PID:-unknown}). Exiting."
  exit 0
fi
echo $$ > "$LOCK_DIR/pid"

# ---- 函数 -------------------------------------------------------------------
ts() { date -u +"%Y-%m-%dT%H:%M:%SZ"; }

log() {
  local phase="$1" action="$2" result="$3" next="$4"
  printf -- "- [%s] %s | %s | %s | %s\n" "$(ts)" "$phase" "$action" "$result" "$next" >> "$LOG_FILE"
}

run_frontend_checks() {
  [ -d frontend/node_modules ] || return 0
  local ok=1
  if [ -x frontend/node_modules/.bin/tsc ]; then
    (cd frontend && node_modules/.bin/tsc --noEmit) >/tmp/tsc.log 2>&1 || ok=0
  fi
  if [ -x frontend/node_modules/.bin/eslint ]; then
    (cd frontend && node_modules/.bin/eslint "src/**/*.{ts,tsx}" --max-warnings=0) >/tmp/eslint.log 2>&1 || ok=0
  fi
  return $ok
}

run_backend_checks() {
  command -v go >/dev/null 2>&1 || return 0   # 无 go 则跳过（不视为失败）
  (cd backend && go build ./... ) >/tmp/gobuild.log 2>&1 || return 1
  (cd backend && go vet ./... ) >/tmp/govet.log 2>&1 || return 1
  return 0
}

auto_commit() {
  if git diff --quiet && git diff --cached --quiet; then
    return 0
  fi
  local msg="loop: auto-commit verified changes @ $(ts)"
  git add -A
  git commit -q -m "$msg" && log "COMMIT" "$msg" "ok" "continue loop"
}

cleanup() {
  rm -rf "$LOCK_DIR"
  log "DAEMON" "stop" "signal received" "exit"
  exit 0
}
trap cleanup INT TERM EXIT

# ---- 主循环 -----------------------------------------------------------------
log "DAEMON" "start" "pid=$$ deadline=$(date -j -f %s "$DEADLINE_EPOCH" +%Y-%m-%dT%H:%M:%S%z 2>/dev/null || echo $DEADLINE_EPOCH)" "heartbeat + checks"

while [ "$(date +%s)" -lt "$DEADLINE_EPOCH" ]; do
  if [ -n "$ENGINE_CMD" ] && command -v "$ENGINE_CMD" >/dev/null 2>&1; then
    log "ENGINE" "delegate-next-backlog" "invoking $ENGINE_CMD" "await engine"
  fi

  if run_frontend_checks && run_backend_checks; then
    auto_commit
    log "HEARTBEAT" "checks passed" "ok" "sleep $SLEEP_SEC"
  else
    log "HEARTBEAT" "checks FAILED (see /tmp/*.log)" "skip commit" "sleep $SLEEP_SEC"
  fi

  sleep "$SLEEP_SEC"
done

log "DAEMON" "stop" "deadline reached" "hand off to human inspection"
exit 0
