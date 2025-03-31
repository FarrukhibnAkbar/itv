.PHONY: api_gateway movie_service

api_gateway:
	@echo "Running api gateway..."
	go run monorepo/api_gateway/main.go

movie_service:
	@echo "Running movie service..."
	go run monorepo/movie_service/main.go

proto-gen:
	protoc --proto_path=monorepo/proto/movie_service --go_out=monorepo/proto/movie_service --go-grpc_out=monorepo/proto/movie_service monorepo/proto/movie_service/movie.proto
swag-init:
	cd monorepo/api_gateway && swag init --parseDependency main.go --output docs
	@echo "Swagger docs generated at monorepo/api_gateway/docs"
