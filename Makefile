SERVER_DIR := stats
SERVER_FILES := $(SERVER_DIR)/wg.pb.go \
		$(SERVER_DIR)/wg_grpc.pb.go

default: $(SERVER_FILES)

$(SERVER_DIR):
	mkdir -p $@

$(SERVER_DIR)/%.pb.go: proto/%.proto | $(SERVER_DIR)
	protoc -I proto/ $< --go_out=module=github.com/jspc/wg-sock-stats/$(SERVER_DIR):$(SERVER_DIR)

$(SERVER_DIR)/%_grpc.pb.go: proto/%.proto | $(SERVER_DIR)
	protoc -I proto/ $< --go-grpc_out=module=github.com/jspc/wg-sock-stats/$(SERVER_DIR):$(SERVER_DIR)
