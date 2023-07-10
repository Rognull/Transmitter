all: server

gen:
	protoc -I=api/proto \
    	--go_out=api/gen \
    	--go-grpc_out=api/gen \
    	transmitter.proto

gorm:
	go get -u gorm.io/gorm
	go get -u gorm.io/driver/postgres


docker:
	docker compose up -d

client:
	go run cmd/client/main.go -f 2.5

server:
	go run cmd/server/main.go

clean:



