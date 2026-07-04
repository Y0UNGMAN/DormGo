#!/bin/bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
PIDS_DIR="$ROOT_DIR/.pids"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${CYAN}[stop]${NC} $1"; }
ok()   { echo -e "${GREEN}[ok]${NC}   $1"; }
warn() { echo -e "${YELLOW}[warn]${NC}  $1"; }

kill_gracefully() {
    local pid="$1"
    local name="$2"
    if ! kill -0 "$pid" 2>/dev/null; then
        return 1  # already dead
    fi
    log "stopping $name (pid $pid)..."
    kill "$pid" 2>/dev/null || true
    # wait up to 5 seconds for graceful shutdown
    for i in $(seq 1 5); do
        if ! kill -0 "$pid" 2>/dev/null; then
            ok "$name stopped"
            return 0
        fi
        sleep 1
    done
    # force kill
    warn "$name did not exit, force-killing..."
    kill -9 "$pid" 2>/dev/null || true
    sleep 0.5
    if kill -0 "$pid" 2>/dev/null; then
        warn "$name could not be killed"
        return 1
    fi
    ok "$name force-stopped"
    return 0
}

kill_port() {
    local port="$1"
    local label="$2"
    local pids
    pids=$(lsof -ti :"$port" 2>/dev/null || true)
    if [ -n "$pids" ]; then
        log "cleaning up port $port ($label)..."
        echo "$pids" | xargs kill 2>/dev/null || true
        sleep 1
        # force-kill stragglers
        pids=$(lsof -ti :"$port" 2>/dev/null || true)
        if [ -n "$pids" ]; then
            echo "$pids" | xargs kill -9 2>/dev/null || true
        fi
        ok "port $port freed"
    fi
}

# ---- stop by PID files ----
STOPPED=0
for name in frontend agent backend; do
    pid_file="$PIDS_DIR/$name.pid"
    if [ -f "$pid_file" ]; then
        pid=$(cat "$pid_file" 2>/dev/null || true)
        if [ -n "$pid" ]; then
            kill_gracefully "$pid" "$name" && STOPPED=$((STOPPED + 1)) || true
        fi
        rm -f "$pid_file"
    fi
done

# ---- fallback: clean by known ports ----
for port in 3000 8090 8080; do
    if lsof -ti :"$port" >/dev/null 2>&1; then
        case $port in
            3000) label="frontend" ;;
            8090) label="agent" ;;
            8080) label="backend" ;;
        esac
        kill_port "$port" "$label"
    fi
done

if [ $STOPPED -gt 0 ]; then
    echo ""
    ok "All services stopped ($STOPPED via PID files)."
else
    log "No PID files found — cleaned by port fallback."
fi
