
GO_BUILD := go build
GO_BUILD_VARS := CGO_ENABLED=0 GOOS=linux
GO_BUILD_FLAGS := -a -tags netgo -ldflags="-w"

DOCKER_REPOSITORY_NAME := sass-infrastructure
DOCKER_IMAGE_TAG := $(shell ./scripts/image-tag)

DOCKER_PACKAGE_CMD := docker build -t $(DOCKER_REPOSITORY_NAME)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -t $(DOCKER_REPOSITORY_NAME)/$(DOCKER_IMAGE_NAME):latest docker/

all: fe-image store-service-image user-service-image

.PHONY: all

fe-image:
	$(GO_BUILD_VARS) $(GO_BUILD) $(GO_BUILD_FLAGS) -o docker/fe/fe github.com/mdevilliers/k8s-sass-operator/test-services/cmd/fe
	docker build -t $(DOCKER_REPOSITORY_NAME)/fe:$(DOCKER_IMAGE_TAG) -t $(DOCKER_REPOSITORY_NAME)/fe:latest docker/fe/

.PHONY: fe-image

store-service-image:
	$(GO_BUILD_VARS) $(GO_BUILD) $(GO_BUILD_FLAGS) -o docker/store-service/store-service github.com/mdevilliers/k8s-sass-operator/test-services/cmd/store-service
	docker build -t $(DOCKER_REPOSITORY_NAME)/store-service:$(DOCKER_IMAGE_TAG) -t $(DOCKER_REPOSITORY_NAME)/store-service:latest docker/store-service/

.PHONY: store-service-image

user-service-image:
	$(GO_BUILD_VARS) $(GO_BUILD) $(GO_BUILD_FLAGS) -o docker/user-service/user-service github.com/mdevilliers/k8s-sass-operator/test-services/cmd/user-service
	docker build -t $(DOCKER_REPOSITORY_NAME)/user-service:$(DOCKER_IMAGE_TAG) -t $(DOCKER_REPOSITORY_NAME)/user-service:latest docker/user-service/

.PHONY: user-service-image

