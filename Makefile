GO_BUILD := go build
GO_BUILD_VARS := CGO_ENABLED=0 GOOS=linux
GO_BUILD_FLAGS := -a -tags netgo -ldflags="-w"

DOCKER_REPOSITORY_NAME := sass-infrastructure
DOCKER_IMAGE_TAG := $(shell ./test-services/scripts/image-tag)

DOCKER_PACKAGE_CMD := docker build -t $(DOCKER_REPOSITORY_NAME)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -t $(DOCKER_REPOSITORY_NAME)/$(DOCKER_IMAGE_NAME):latest docker/

image:
	$(GO_BUILD_VARS) $(GO_BUILD) $(GO_BUILD_FLAGS) -o docker/operator github.com/mdevilliers/k8s-sass-operator/cmd/operator
	docker build -t $(DOCKER_REPOSITORY_NAME)/operator:$(DOCKER_IMAGE_TAG) -t $(DOCKER_REPOSITORY_NAME)/operator:latest docker/

.PHONY:image

deploy:
	kubectl delete -f k8s/ 2>/dev/null; true
	kubectl create -f k8s/

.PHONY: deploy

deploy-services:
	$(MAKE) -C test-services all

