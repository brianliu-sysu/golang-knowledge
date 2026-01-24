PROTO_DIR := proto
GEN_DIR := internal/gen

.PHONY: proto test run tidy

proto:
	mkdir -p $(GEN_DIR)
	protoc \
		-I $(PROTO_DIR) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		$(PROTO_DIR)/im/v1/im.proto

test:
	go test ./...

run:
	go run main.go

tidy:
	go mod tidy
