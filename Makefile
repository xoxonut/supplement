DOCKER_IMAGE_OWNER = 'free5gc'
DOCKER_IMAGE_NAME = 'base'
DOCKER_IMAGE_TAG = 'latest'

.PHONY: base
all: base amf ausf udm scp

base:
	docker build -t ${DOCKER_IMAGE_OWNER}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG} ./base
	docker image ls ${DOCKER_IMAGE_OWNER}/${DOCKER_IMAGE_NAME}:${DOCKER_IMAGE_TAG}


amf: base
	docker build --build-arg F5GC_MODULE=amf -t ${DOCKER_IMAGE_OWNER}/amf-base:${DOCKER_IMAGE_TAG} -f ./base/Dockerfile.nf ./base
udm: base
	docker build --build-arg F5GC_MODULE=udm -t ${DOCKER_IMAGE_OWNER}/udm-base:${DOCKER_IMAGE_TAG} -f ./base/Dockerfile.nf ./base
ausf: base
	docker build --build-arg F5GC_MODULE=ausf -t ${DOCKER_IMAGE_OWNER}/ausf-base:${DOCKER_IMAGE_TAG} -f ./base/Dockerfile.nf ./base
scp: base
	docker build --build-arg F5GC_MODULE=scp -t ${DOCKER_IMAGE_OWNER}/scp-base:${DOCKER_IMAGE_TAG} -f ./base/Dockerfile.nf ./base
