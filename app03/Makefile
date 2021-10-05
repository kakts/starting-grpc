gen_proto:
	protoc \
		-Iproto \
		--go_out=plugins=grpc:. \
		proto/*

# リバーシサーバの起動
run_server:
	go run server/grpc/server.go

# リバーシクライアントの起動
run_client:
	go run cmd/main.go