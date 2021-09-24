
# Image URL to use all building/pushing image targets
IMG ?= cpservice-annotator:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: webhook

# Build manager binary
webhook:
	go build -o bin/webhook main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: webhook
	go run ./main.go

uninstall:
	kubectl delete -f deployments/
	kubectl delete -f certs/

deploy-cert: webhook
	kubectl apply -f certs/

deploy: webhook
	kubectl apply -f deployments/

# Build the docker image
docker-build: webhook
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

install-cert-manager:
	kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.2/cert-manager.yaml

