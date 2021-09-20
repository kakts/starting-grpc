.PHONY: gen_proto
gen_proto: 
	bundle exec grpc_tools_ruby_protoc \
		-I ../proto \
		--ruby_out=app/gen/pb/image/upload \
		--grpc_out=app/gen/pb/image/upload \
		../proto/image_uploader.proto