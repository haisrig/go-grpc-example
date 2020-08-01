gen: 
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:proto
run:
	go run main.go