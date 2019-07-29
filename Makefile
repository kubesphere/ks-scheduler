# Image URL to use all building/pushing image targets
IMG ?= zhuxiaoyang/ks-pipeline-scheduler:v1

# Build manager binary
manager: fmt manifests vet
	GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o ks-pipeline-scheduler /root/cmd

# Build the docker image
docker-build:
	docker build -t $(IMG) .
	docker push $(IMG)

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./cmd/manager/main.go

e2e-test: manager
	./hack/e2etest.sh

.PHONY : clean test


