#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOGS_DIR="$ROOT_DIR/logs"
PIDS_DIR="$ROOT_DIR/.pids"

mkdir -p "$LOGS_DIR" "$PIDS_DIR"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${CYAN}[start]${NC} $1"; }
ok()   { echo -e "${GREEN}[ok]${NC}   $1"; }
warn() { echo -e "${YELLOW}[warn]${NC}  $1"; }
fail() { echo -e "${RED}[fail]${NC}  $1"; }

# ---- helpers ----
port_in_use() {
    lsof -ti :"$1" >/dev/null 2>&1
}

wait_for_port() {
    local port="$1"
    local timeout="${2:-30}"
    local elapsed=0
    while ! port_in_use "$port"; do
        sleep 1
        elapsed=$((elapsed + 1))
        if [ "$elapsed" -ge "$timeout" ]; then
            fail "timed out waiting for port $port"
            return 1
        fi
    done
    return 0
}

kill_port() {
    local port="$1"
    local pids
    pids=$(lsof -ti :"$port" 2>/dev/null || true)
    if [ -n "$pids" ]; then
        echo "$pids" | xargs kill 2>/dev/null || true
    fi
}

# ---- check if already running ----
ALREADY_RUNNING=""
for port in 8080 8090 3000; do
    if port_in_use "$port"; then
        ALREADY_RUNNING="$ALREADY_RUNNING $port"
    fi
done
if [ -n "$ALREADY_RUNNING" ]; then
    warn "ports already in use:$ALREADY_RUNNING"
    warn "run ./stop.sh first, or continue anyway? [y/N]"
    read -r ans
    if [ "$ans" != "y" ] && [ "$ans" != "Y" ]; then
        exit 0
    fi
fi

# ---- 1. Go Backend (:8080) ----
log "starting Go backend..."
cd "$ROOT_DIR/backend"
nohup go run main.go >> "$LOGS_DIR/backend.log" 2>&1 &
echo $! > "$PIDS_DIR/backend.pid"
if wait_for_port 8080; then
    ok "Go backend  running on :8080  (pid $(cat "$PIDS_DIR/backend.pid"))"
else
    fail "Go backend failed to start — check logs/backend.log"
    exit 1
fi

# ---- 2. Python Agent (:8090) ----
log "starting Python agent..."
cd "$ROOT_DIR/agent-service"
nohup python3 main.py >> "$LOGS_DIR/agent.log" 2>&1 &
echo $! > "$PIDS_DIR/agent.pid"
if wait_for_port 8090 60; then
    ok "Python agent running on :8090 (pid $(cat "$PIDS_DIR/agent.pid"))"
else
    fail "Python agent failed to start — check logs/agent.log"
    exit 1
fi

# ---- 3. Vue Frontend (:3000) ----
log "starting Vue frontend..."
cd "$ROOT_DIR/frontend"
nohup npm run dev >> "$LOGS_DIR/frontend.log" 2>&1 &
echo $! > "$PIDS_DIR/frontend.pid"
if wait_for_port 3000 20; then
    ok "Vue frontend running on :3000  (pid $(cat "$PIDS_DIR/frontend.pid"))"
else
    fail "Vue frontend failed to start — check logs/frontend.log"
    exit 1
fi

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}  All services started!${NC}"
echo -e "${GREEN}  Backend:  http://127.0.0.1:8080${NC}"
echo -e "${GREEN}  Agent:    http://127.0.0.1:8090${NC}"
echo -e "${GREEN}  Frontend: http://127.0.0.1:3000${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "  Logs: logs/backend.log | logs/agent.log | logs/frontend.log"
echo "  Stop: ./stop.sh  or  make stop"
echo ""
