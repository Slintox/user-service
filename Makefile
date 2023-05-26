DEFAULT_GOAL := help
.PHONY:

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=./migrations
LOCAL_MIGRATION_DSN="host=localhost port=54322 dbname=user user=user-user password=user-password sslmode=disable"

# HELP =================================================================================================================
# This will output the help for each task with comment
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

install-go-deps: ## Install deps
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2


gen: ## Generates on linux
	mkdir -p pkg/swagger/
	make generate_user_grpc
	statik -f -src=pkg/swagger/ -include='*.css,*html,*js,*json,*png'

gen_win: ## Generates on windows
	mkdir -p pkg/swagger/
	make generate_user_grpc_win
	statik -f -src=pkg/swagger/ -include='*.css,*html,*js,*json,*png'

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

generate_user_grpc: ## Generate grpc api
	mkdir -p "pkg/user_v1"
	protoc --proto_path=api/user_v1 --proto_path vendor.protogen \
  	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
  	--plugin=protoc-gen-go=./bin/protoc-gen-go \
  	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
  	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
  	--validate_out lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
  	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
  	--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
  	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-go-gateway \
  	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
  	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
  	api/user_v1/service.proto

generate_user_grpc_win: ## Generate grpc api using windows protoc binaries
	mkdir -p "pkg/user_v1"
	protoc --proto_path=api/user_v1 --proto_path vendor.protogen \
 		--go_out=pkg/user_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go.exe \
		--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
    	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc.exe \
    	--validate_out lang=go:pkg/user_v1 --validate_opt=paths=source_relative \
    	--plugin=protoc-gen-validate=bin/protoc-gen-validate.exe \
    	--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
    	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway.exe \
    	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
    	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2.exe \
    	api/user_v1/service.proto

local-migration-status: ## Migration status
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up: ## Migration up
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down: ## Migration down
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

vendor-proto: ## Clones vendor
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p vendor.protogen/google &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
  		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
  		git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
  		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
  		rm -rf vendor.protogen/openapiv2 ;\
  	fi
