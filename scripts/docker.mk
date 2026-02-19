DOCKER_COMPOSE_PATH := ./docker/docker-compose.yaml

# --- Docker Infrastructure ---

start_all:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d --build

stop_all:
	docker compose -f $(DOCKER_COMPOSE_PATH) down

reset_all:
	docker compose -f $(DOCKER_COMPOSE_PATH) down -v

docker_logs:
	docker compose -f $(DOCKER_COMPOSE_PATH) logs -f

# --- Database & Redis ---

start_postgres:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d postgres

connect_postgres:
	docker exec -it auth_postgres psql -U kratos -d kratos

start_redis:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d redis

connect_redis:
	docker exec -it auth_redis redis-cli

# --- Kratos Specific ---

start_kratos:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d kratos

kratos_logs:
	docker compose -f $(DOCKER_COMPOSE_PATH) logs -f kratos

restart_kratos:
	docker compose -f $(DOCKER_COMPOSE_PATH) restart kratos

# --- Utils ---

open_mail:
	@echo "${CYAN}Opening Mailpit...${NC}"
	@xdg-open http://localhost:8025 || open http://localhost:8025 || echo "Manual link: http://localhost:8025"
