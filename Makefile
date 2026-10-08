# Common tasks for the whole project.
# Usage: make test | make coverage | make up | make down
.PHONY: test test-backend test-frontend coverage up down

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test -- --run

coverage:
	cd backend && go test -cover ./...
	cd frontend && npm run coverage

up:
	docker compose up --build

down:
	docker compose down