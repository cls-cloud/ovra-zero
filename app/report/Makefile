init:
	goctl env -w GOCTL_EXPERIMENTAL=off
	go install github.com/zeromicro/go-zero/tools/goctl@latest
	go install github.com/zeromicro/goctl-go-compact@latest


api:
	goctl api go --api desc/desc.api --style "go_zero" --dir .

api-compact:
	goctl api plugin -p goctl-go-compact -api desc/desc.api --style "go_zero" -dir .

grpc:
	goctl rpc protoc desc/desc.proto --style "go_zero" --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m

db:
	gentool -c "./gen.yaml"
	./replace.sh