
BUILD_NAME      = fcd
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_RELEASE   = $(shell if echo ${BUILD_VERSION} | egrep -q '^[0-9]+\.[0-9]+\.[0-9]+$$'; then echo true; else echo false; fi )
BUILD_DATE     ?= $(shell if ${BUILD_RELEASE}; then date -u +"%Y-%m-%d"; else date -u +"%Y-%m-%dT%H:%M:%S%z"; fi)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_PLATFORM}



PROTOC_FILES=$(wildcard proto/*.proto)
SOURCE_FILES=$(shell go list -f '{{.GoFiles}}' | tr -d '[]') $(PROTOC_FILES:.proto=.pb.go)


LDFLAGS = "-X main.BUILD_NAME=${BUILD_NAME} -X main.BUILD_VERSION=${BUILD_VERSION} -X main.BUILD_PLATFORM=${BUILD_PLATFORM} -X main.BUILD_DATE=${BUILD_DATE}"


help:
	@echo "### Usage ###"
	@echo " make build    # build static executable"
	@echo " make proto    # rebuild protocol files"
	@echo " make info     # show build info"
	@echo " make clean    # clean up"


info: 
	@echo "### Version Info ###"
	@echo " name       ${BUILD_NAME}"
	@echo " version    ${BUILD_VERSION}"
	@echo " platform   ${BUILD_PLATFORM}"
	@echo " date       ${BUILD_DATE}"
	@echo " product    ${BUILD_PRODUCT}"
	@echo ""
	@echo "### Build Variables ###"
	@echo " source     ${SOURCE_FILES}"
	@echo " protoc     ${PROTOC_FILES}"
	
build: ${BUILD_PRODUCT}

get:
	go get -v 
	
clean:
	-rm -f ${BUILD_PRODUCT} ${BUILD_NAME} 

${BUILD_NAME}: ${BUILD_PRODUCT}
	cp -f ${BUILD_PRODUCT} ${BUILD_NAME}

${BUILD_PRODUCT}: ${SOURCE_FILES}
	go build -o ${BUILD_PRODUCT} -v -ldflags ${LDFLAGS} $(shell go list -f '{{.GoFiles}}' | tr -d '[]' )



proto: $(PROTOC_FILES:.proto=.pb.go)

proto/%.pb.go: proto/%.proto
	protoc -I proto $^ --go_out=plugins=grpc:proto


.PHONY: help info build clean get proto