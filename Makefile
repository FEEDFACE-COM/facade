
BUILD_NAME      = fcd
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_RELEASE   = $(shell if echo ${BUILD_VERSION} | egrep -q '^[0-9]+\.[0-9]+\.[0-9]+$$'; then echo true; else echo false; fi )
BUILD_DATE     ?= $(shell if ${BUILD_RELEASE}; then date -u +"%Y-%m-%d"; else date -u +"%Y-%m-%dT%H:%M:%S%z"; fi)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_PLATFORM}



SOURCES=$(wildcard *.go */*.go) 
ASSETS=gfx/shaderAssets.go gfx/fontAssets.go




ASSET_FONT=font/RobotoMono-Regular.ttf font/VT323-Regular.ttf
ASSET_VERT=$(wildcard shader/*.vert shader/*/*.vert)
ASSET_FRAG=$(wildcard shader/*.frag shader/*/*.frag)

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
	@echo " source     ${SOURCES}"
	@echo "  asset     ${ASSETS}"
	@echo ""
	@echo "   vert     ${ASSET_VERT}"
	@echo "   frag     ${ASSET_FRAG}"
	@echo "   font     ${ASSET_FONT}"   
	
build: ${BUILD_PRODUCT}

demo: ${BUILD_PRODUCT}
	for f in ${SOURCES}; do cat $$f | while read l; do sleep 0.7; echo $$l | ./${BUILD_PRODUCT} pipe grid -h 10 -s; done; done
# for f in gfx/*.go; do cat $f | while read l; do sleep 1; echo $l | fcd; done; done


get:
	go get -v 
	
clean:
	-rm -f ${BUILD_PRODUCT} ${BUILD_NAME} ${ASSETS}

${BUILD_NAME}: ${BUILD_PRODUCT}
	cp -f ${BUILD_PRODUCT} ${BUILD_NAME}

${BUILD_PRODUCT}: ${SOURCES} ${ASSETS}
	go build -v -o ${BUILD_PRODUCT} -v -ldflags ${LDFLAGS} $(shell go list -f '{{.GoFiles}}' | tr -d '[]' )

assets: ${ASSETS}


	
gfx/shaderAssets.go: ${ASSET_VERT} ${ASSET_FRAG}
	echo ""                                         >|$@
#	echo "// +build linux,arm"                      >>$@
	echo "package gfx"                              >>$@
	echo "var VertexShader = map[string]string{"    >>$@
	for src in ${ASSET_VERT}; do \
      name=$$(echo $$src | sed -e 's:shader/::;s/.vert//'); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src; \
      echo "\`,\n\n"; \
    done                                            >>$@
	echo "}\n\n"                                    >>$@
	echo "var FragmentShader = map[string]string{"  >>$@
	for src in ${ASSET_FRAG}; do \
      name=$$(echo $$src | sed -e 's:shader/::;s/.frag//'); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src; \
      echo "\`,\n\n"; \
    done                                            >>$@
	echo "}"                                        >>$@


gfx/fontAssets.go: ${ASSET_FONT}
	echo ""                                         >|$@
#	echo "// +build linux,arm"                      >>$@
	echo "package gfx"                              >>$@
	echo "var VectorFont = map[string]string{"      >>$@
	for src in ${ASSET_FONT}; do \
      name=$$( echo $$src | sed -e 's:font/::;s:\.[tT][tT][fFcC]::' ); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src | base64 ; \
      echo "\`,\n\n"; \
    done                                            >>$@
	echo "}"                                        >>$@



.PHONY: help info build clean get assets demo

