APP = shortener
BUILD_DIR = build
REPO = $(shell go list -m)
BUILD_DATE = $(shell date +%FT%T%Z)
BUILD_COMMIT = $(shell git rev-parse HEAD)
VERSION = $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD)))

GO_BUILD_ARGS = \
  -ldflags " \
    -X '$(REPO)/internal/version.Version=$(VERSION)' \
    -X '$(REPO)/internal/version.BuildCommit=$(BUILD_COMMIT)' \
    -X '$(REPO)/internal/version.BuildDate=$(BUILD_DATE)' \
  " \

.PHONY: build
build:
	@echo "+ $@"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GO_BUILD_ARGS) -o "./$(BUILD_DIR)/$(APP)" ./cmd/server

.PHONY: test
test:
	@echo "+ $@"
	go test -cover ./...

.PHONY: test-cover
test-cover:
	@echo "+ $@"
	go test -coverprofile=profile.out ./...
	go tool cover -html=profile.out
	rm profile.out

.PHONY: check
check:
	golangci-lint run

.PHONY: run
run: clean build
	@echo "+ $@"
	./${BUILD_DIR}/${APP}

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)

.PHONY: watch
watch: go-prep-watch
	#reflex -s -r '\.go$$' make run
	reflex -r '\.go$$' -s -- sh -c "go run ./cmd/server/main.go"

go-prep-watch:
	@echo "\nPreparing environment...."
	go get github.com/cespare/reflex

.PHONY: gen-swagger
gen-swagger:
	@echo "+ $@"
	swagger generate spec -m -o ./api/swagger.json
	cp ./api/swagger.json ./web/static/swaggerui/
