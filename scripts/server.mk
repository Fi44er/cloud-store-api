run_app:
	go run ./cmd/main.go -migrate=${migrate_flag} -redis=${redis_mode}

gen_dock:
	swag init -d ./cmd --pdl 3

gen_mock:
	@echo "Generating mock for $(src)"
	@echo "$(dir $(dest))/mock"
	@mkdir -p $(dir $(dest))/mock
	@go run github.com/golang/mock/mockgen@latest \
		-source=$(src) \
		-destination=$(dir $(dest))/mock/$(notdir $(basename $(src)))_mock.go \
		-package=mock
