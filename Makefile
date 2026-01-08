# ========================
# 模块列表
# ========================
SUBMODULES := system monitor resource demo

# ========================
# 模块路径
# ========================
SYSTEM_PATH := app/system
MONITOR_PATH := app/monitor
RESOURCE_PATH := app/resource
DEMO_PATH := app/demo

# ========================
# 模块名
# ========================
SYSTEM_NAME := system
MONITOR_NAME := monitor
RESOURCE_NAME := resource
DEMO_NAME := demo

# ========================
# 初始化工具
# ========================
.PHONY: init
# init        初始化 goctl 工具和整理依赖
init:
	goctl env -w GOCTL_EXPERIMENTAL=off
	go install github.com/zeromicro/go-zero/tools/goctl@latest
	go mod tidy

# ========================
# 模块通用规则模板
# ========================
define MODULE_RULES

# -------- API 生成 --------
.PHONY: api-$(1)
# api-$(1)      生成 $(1) 模块的 API 代码
api-$(1):
	goctl api go \
		--api desc/$(1)/api/desc.api \
		--style "go_zero" \
		--dir $($(shell echo $(1) | tr a-z A-Z)_PATH)
	rm -rf $($(shell echo $(1) | tr a-z A-Z)_PATH)/etc

# -------- gRPC 生成 --------
.PHONY: grpc-$(1)
# grpc-$(1)     生成 $(1) 模块的 gRPC 代码
grpc-$(1):
	goctl rpc protoc desc/$(1)/rpc/desc.proto \
		--style "go_zero" \
		--go_out=$($(shell echo $(1) | tr a-z A-Z)_PATH)/pb \
		--go-grpc_out=$($(shell echo $(1) | tr a-z A-Z)_PATH)/pb \
		--zrpc_out=$($(shell echo $(1) | tr a-z A-Z)_PATH) \
		--client=true -m
	rm -rf $($(shell echo $(1) | tr a-z A-Z)_PATH)/desc.go $($(shell echo $(1) | tr a-z A-Z)_PATH)/etc

# -------- 数据库生成 --------
.PHONY: db-$(1)
# db-$(1)       生成 $(1) 模块的数据库代码
db-$(1):
	gentool -c "gen/$(1)/gen.yaml"

# -------- 构建模块 --------
.PHONY: build-$(1)
# build-$(1)    构建 $(1) 模块
build-$(1):
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags no_k8s \
		-o $($(shell echo $(1) | tr a-z A-Z)_NAME) \
		$($(shell echo $(1) | tr a-z A-Z)_PATH)/$($(shell echo $(1) | tr a-z A-Z)_NAME).go

# -------- 运行模块 --------
.PHONY: run-$(1)
# run-$(1)      运行 $(1) 模块
run-$(1):
	go run $($(shell echo $(1) | tr a-z A-Z)_PATH)/$($(shell echo $(1) | tr a-z A-Z)_NAME).go \
		-f etc/$($(shell echo $(1) | tr a-z A-Z)_NAME).yaml

# -------- 后台运行模块 --------
.PHONY: back-$(1)
# back-$(1)     后台运行 $(1) 模块
back-$(1):
	nohup ./$$($(shell echo $(1) | tr a-z A-Z)_NAME) \
		> $$($(shell echo $(1) | tr a-z A-Z)_NAME).log 2>&1 &

endef

# ========================
# 遍历模块生成规则
# ========================
$(foreach mod,$(SUBMODULES),$(eval $(call MODULE_RULES,$(mod))))

# ========================
# 构建所有模块
# ========================
.PHONY: build-all
# build-all     构建所有模块
build-all:
	@for mod in $(SUBMODULES); do \
		$(MAKE) build-$$mod; \
	done

# ========================
# 批量生成所有模块 API / gRPC / DB
# ========================
.PHONY: gen-all
# gen-all       批量生成所有模块的 API / gRPC / DB
gen-all:
	@for mod in $(SUBMODULES); do \
		echo ">>> Generating API for $$mod"; \
		$(MAKE) api-$$mod; \
		echo ">>> Generating gRPC for $$mod"; \
		$(MAKE) grpc-$$mod; \
		echo ">>> Generating DB for $$mod"; \
		$(MAKE) db-$$mod; \
	done

# ========================
# 后台运行所有模块
# ========================
.PHONY: back-all
# back-all      后台运行所有模块
back-all:
	@for mod in $(SUBMODULES); do \
		$(MAKE) back-$$mod; \
	done

# ========================
# Traefik 运行
# ========================
.PHONY: traefik-run
# traefik-run   启动 Traefik
traefik-run:
	./bin/traefik/traefik --configfile=./bin/traefik/traefik.yaml

# ========================
# 显示 help
# ========================
.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH-2); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

# ========================
# 默认目标
# ========================
.DEFAULT_GOAL := help
