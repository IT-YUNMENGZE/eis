all: local

local: fmt vet
	GOOS=linux GOARCH=amd64 go build  -o=bin/edge-scheduler ./cmd/scheduler

build:  local
	docker build --no-cache . -t docker push registry.cn-guangzhou.aliyuncs.com/yunmengze/edge-scheduler:0.01

push:   build
	docker push registry.cn-guangzhou.aliyuncs.com/yunmengze/edge-scheduler:0.01

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

clean: fmt vet
	sudo rm -f edge-scheduler