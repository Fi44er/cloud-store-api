include ./scripts/server.mk
include ./scripts/docker.mk

# Цвета для текста
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
CYAN=\033[0;36m
NC=\033[0m # No Color
RED=\033[0;31m

# Флаги из твоего оригинала
migrate_flag ?= false
redis_mode ?= 0
WIDTH=25

.PHONY: help run_app gen_mock gen_dock start_postgres connect_postgres start_redis connect_redis \
        start_all stop_all reset_all docker_logs kratos_logs restart_kratos open_mail

help:
	@echo "${BLUE}=== Auth System Help ===${NC}"
	@echo ""
	@echo "${YELLOW}App Commands:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "run_app" "Start the app (migrate_flag=${migrate_flag}, redis_mode=${redis_mode})"
	@printf "  ${CYAN}%-${WIDTH}s${NC}%s\n" "" "Example: make run_app redis_mode=2 migrate_flag=true"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "gen_mock" "Generate mock for src"
	@printf "  ${CYAN}%-${WIDTH}s${NC}%s\n" "" "Example: make gen_mock src=./cmd/handler/user.go"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "gen_dock" "Generate Swagger documentation"
	@echo ""
	@echo "${YELLOW}Docker Infrastructure:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "start_all" "Start everything (Postgres, Kratos, Redis, Backend)"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "stop_all" "Down all services"
	@printf "  ${RED}%-${WIDTH}s${NC}%s\n" "reset_all" "Stop and DELETE ALL DATA (volumes)"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "start_postgres" "Start PostgreSQL"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "start_redis" "Start Redis"
	@echo ""
	@echo "${YELLOW}Kratos & Debug:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "kratos_logs" "View Kratos specific logs"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "restart_kratos" "Restart Kratos (useful after config change)"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "docker_logs" "View all logs"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "open_mail" "Open Mailpit UI (Web Check Email)"
	@echo ""
	@echo "${YELLOW}Connections:${NC}"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "connect_postgres" "Connect to Postgres (psql)"
	@printf "  ${GREEN}%-${WIDTH}s${NC}%s\n" "connect_redis" "Connect to Redis (cli)"
