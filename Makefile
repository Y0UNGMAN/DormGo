.PHONY: dev start stop restart status logs logs-backend logs-agent logs-frontend clean help

# DormGo 一键开发环境管理
# make dev   — 启动全部服务
# make stop  — 停止全部服务
# make logs  — 查看全部日志
# make status — 查看服务状态

dev: start

start:
	@bash start.sh

stop:
	@bash stop.sh

restart: stop start

status:
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  DormGo 服务状态"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@for port in 8080 8090 3000; do \
		case $$port in \
			8080) name="Go Backend " ;; \
			8090) name="Python Agent" ;; \
			3000) name="Vue Frontend" ;; \
		esac; \
		pid=$$(lsof -ti :$$port 2>/dev/null || true); \
		if [ -n "$$pid" ]; then \
			echo "  ✅ $$name  :$$port (pid $$pid)"; \
		else \
			echo "  ❌ $$name  :$$port (not running)"; \
		fi; \
	done
	@echo ""

logs:
	@tail -n 20 logs/backend.log 2>/dev/null || echo "(no backend log)"
	@echo "---"
	@tail -n 20 logs/agent.log 2>/dev/null || echo "(no agent log)"
	@echo "---"
	@tail -n 20 logs/frontend.log 2>/dev/null || echo "(no frontend log)"

logs-backend:
	@tail -f logs/backend.log

logs-agent:
	@tail -f logs/agent.log

logs-frontend:
	@tail -f logs/frontend.log

clean:
	@bash stop.sh 2>/dev/null || true
	@rm -rf logs/ .pids/
	@echo "Cleaned logs/ and .pids/"

help:
	@echo "DormGo 开发命令："
	@echo "  make dev      启动全部服务"
	@echo "  make stop     停止全部服务"
	@echo "  make restart  重启全部服务"
	@echo "  make status   查看各服务状态"
	@echo "  make logs     查看最近日志"
	@echo "  make logs-*   实时跟踪日志 (logs-backend / logs-agent / logs-frontend)"
	@echo "  make clean    停止并清理日志和PID"
