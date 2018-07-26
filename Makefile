
BUILD_NAME     = fcd
BUILD_VERSION ?= $(shell git describe --tags)
BUILD_DATE     ?= $(shell date -u +"%Y-%m-%d")
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_PLATFORM}




SOURCE_FILES=$(wildcard *.go)

LDFLAGS = "-X main.BUILD_NAME=${BUILD_NAME} -X main.BUILD_VERSION=${BUILD_VERSION} -X main.BUILD_PLATFORM=${BUILD_PLATFORM} -X main.BUILD_DATE=${BUILD_DATE}"


help:
	@echo "### Usage ###"
	@echo " make build    # build static binary"
	@echo " make info     # show build info"
	@echo " make clean    # clean up"
	@echo " make clean    # this help"


info: 
	@echo "### Version Info ###"
	@echo " name       ${BUILD_NAME}"
	@echo " version    ${BUILD_VERSION}"
	@echo " platform   ${BUILD_PLATFORM}"
	@echo " date       ${BUILD_DATE}"
	@echo " product    ${BUILD_PRODUCT}"


build: ${BUILD_PRODUCT}

${BUILD_NAME}: ${BUILD_PRODUCT}
	cp -f ${BUILD_PRODUCT} ${BUILD_NAME}

${BUILD_PRODUCT}: ${SOURCE_FILES}
	go build -o ${BUILD_PRODUCT} -v -ldflags ${LDFLAGS} ${SOURCE_FILES}
	