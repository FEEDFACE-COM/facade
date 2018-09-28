
BUILD_NAME      = fcd
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_RELEASE   = $(shell if echo ${BUILD_VERSION} | egrep -q '^[0-9]+\.[0-9]+\.[0-9]+$$'; then echo true; else echo false; fi )
BUILD_DATE     ?= $(shell if ${BUILD_RELEASE}; then date -u +"%Y-%m-%d"; else date -u +"%Y-%m-%dT%H:%M:%S%z"; fi)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_PLATFORM}



SOURCE_FILES=$(wildcard *.go */*.go) gfx/shaderFragment.go gfx/shaderVertex.go
SHADER_FILES=$(wildcard shader/*.vert shader/*.frag)



LDFLAGS = "-X main.BUILD_NAME=${BUILD_NAME} -X main.BUILD_VERSION=${BUILD_VERSION} -X main.BUILD_PLATFORM=${BUILD_PLATFORM} -X main.BUILD_DATE=${BUILD_DATE}"


default: build

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
	@echo " shader     ${SHADER_FILES}"
	
build: ${BUILD_PRODUCT}

demo: ${BUILD_PRODUCT}
	for f in ${SOURCE_FILES}; do cat $$f | while read l; do sleep 0.7; echo $$l | ./${BUILD_PRODUCT} pipe grid -h 10 -s; done; done
# for f in gfx/*.go; do cat $f | while read l; do sleep 1; echo $l | fcd; done; done


get:
	go get -v 
	
clean:
	-rm -f ${BUILD_PRODUCT} ${BUILD_NAME} gfx/shaderFragment.go gfx/shaderVertex.go

${BUILD_NAME}: ${BUILD_PRODUCT}
	cp -f ${BUILD_PRODUCT} ${BUILD_NAME}

${BUILD_PRODUCT}: ${SOURCE_FILES} ${SHADER_FILES}
	go build -v -o ${BUILD_PRODUCT} -v -ldflags ${LDFLAGS} $(shell go list -f '{{.GoFiles}}' | tr -d '[]' )


gfx/shaderFragment.go: shader/*.frag
	echo "" >|$@
	echo "// +build linux,arm" >>$@
	echo "package gfx" >>$@
	echo "var FragmentShader = map[string]string{" >>$@
	for src in shader/*.frag; do \
      name=$$(basename $$src | cut -d. -f -1); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src; \
      echo "\`,\n\n"; \
    done >>$@
	echo "}" >>$@
	
	
gfx/shaderVertex.go: shader/*.vert
	echo "" >|$@
	echo "// +build linux,arm" >>$@
	echo "package gfx" >>$@
	echo "var VertexShader = map[string]string{" >>$@
	for src in shader/*.vert; do \
      name=$$(basename $$src | cut -d. -f -1); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src; \
      echo "\`,\n\n"; \
    done >>$@
	echo "}" >>$@




.PHONY: help info build clean get 

