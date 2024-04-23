GIT_HASH ?= $(shell git log --format="%h" -n 1)
DOCKER_REGISTRY ?= deployto.io
APPLICATION_NAME ?= $(shell basename `pwd`)

_BUILD_ARGS_TAG ?= ${GIT_HASH}
_BUILD_ARGS_RELEASE_TAG ?= latest
_BUILD_ARGS_DOCKERFILE ?= Dockerfile



_builder:
		docker build --tag ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} -f ${_BUILD_ARGS_DOCKERFILE} .
		docker tag ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} ${DOCKER_REGISTRY}/${APPLICATION_NAME}:latest
 
_pusher:
#		docker push ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG}
		minikube image rm ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG}  && minikube image load  ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} 
		minikube image rm ${DOCKER_REGISTRY}/${APPLICATION_NAME}:latest && minikube image load ${DOCKER_REGISTRY}/${APPLICATION_NAME}:latest
 
_releaser:
		docker pull ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG}
		docker tag ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_TAG} ${DOCKER_REGISTRY}/${APPLICATION_NAME}:latest
		docker push ${DOCKER_REGISTRY}/${APPLICATION_NAME}:${_BUILD_ARGS_RELEASE_TAG}

build:
		$(MAKE) _builder
 
push:
		$(MAKE) _pusher
 
release:
		$(MAKE) _releaser


