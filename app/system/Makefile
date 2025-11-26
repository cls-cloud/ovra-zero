APP_NAME=app-system
TAR_NAME=$(APP_NAME).tar.gz
BIN_DIR := ../bin

init:
	goctl env -w GOCTL_EXPERIMENTAL=off
	go install github.com/zeromicro/go-zero/tools/goctl@latest

api:
	goctl api go --api desc/desc.api --style "go_zero" --dir .

grpc:
	goctl rpc protoc desc/desc.proto --style "go_zero" \
	  --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m

db:
	gentool -c "./gen.yaml"
	./replace.sh

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags no_k8s \
	  -o $(APP_NAME) system.go
	tar -czvf $(TAR_NAME) $(APP_NAME) Makefile etc/
	mkdir -p $(BIN_DIR)
	mv $(TAR_NAME) $(BIN_DIR)/
	mv $(APP_NAME) $(BIN_DIR)/

back_run:
	nohup ./$(APP_NAME) > $(APP_NAME).log 2>&1 &

run:
	go run system.go -f etc/system.yaml

clean:
	rm -f $(BIN_DIR)/$(APP_NAME) $(BIN_DIR)/$(TAR_NAME)
