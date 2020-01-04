
BUILD_NAME      = facade
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_RELEASE   = $(shell if echo ${BUILD_VERSION} | egrep -q '^[0-9]+\.[0-9]+\.[0-9]+$$'; then echo true; else echo false; fi )
BUILD_DATE     ?= $(shell if ${BUILD_RELEASE}; then date -u +"%Y-%m-%d"; else date -u +"%Y-%m-%dT%H:%M:%S%z"; fi)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}





ASSETS  = facade/shaderAssets.go facade/fontAssets.go
PROTOS  = facade/facade.pb.go
SOURCES = $(filter-out ${ASSETS} , $(wildcard *.go */*.go) )

FONTS ?= RobotoMono.ttf VT323.ttf Menlo.ttc OCRAEXT.TTF MONACO.TTF Adore64.ttf
ASSET_FONT= $(foreach x,$(FONTS),font/$(x) )

SHADERS ?= def.vert def.frag 
SHADERS += color.vert color.frag 
SHADERS += grid/def.vert grid/def.frag grid/debug.frag grid/debug2.frag 
SHADERS += grid/bent.vert grid/crawl.vert grid/disk.vert grid/pipe.vert grid/roll.vert grid/wave.vert grid/rowz.vert 
SHADERS += mask/def.vert mask/def.frag mask/mask.frag mask/debug.frag 
#ASSET_SHADER=$(wildcard shader/*.vert shader/*/*.vert shader/*.frag shader/*/*.frag)
ASSET_SHADER = $(foreach x,$(SHADERS),shader/$(x))


GCFLAGS ?= 
# REM, for debug: GCFLAGS="-N -l"; make


LDFLAGS ?=
LDFLAGS += -X main.BUILD_NAME=${BUILD_NAME} -X main.BUILD_VERSION=${BUILD_VERSION} -X main.BUILD_PLATFORM=${BUILD_PLATFORM} -X main.BUILD_DATE=${BUILD_DATE}


BUILD_FLAGS ?= 
BUILD_FLAGS +=-v

default: build 

help:
	@echo "#Usage"
	@echo " make build    # build static executable"
	@echo " make deps     # fetch go dependencies"
	@echo " make info     # show build info"
	@echo " make asset    # build fonts and shaders"
	@echo " make proto    # rebuild protobuf code"
	@echo " make clean    # clean up"
	


info: 
	@echo "#Info"
	@echo " name       ${BUILD_NAME}"
	@echo " version    ${BUILD_VERSION}"
	@echo " platform   ${BUILD_PLATFORM}"
	@echo " date       ${BUILD_DATE}"
	@echo " product    ${BUILD_PRODUCT}"
	@echo ""
	@echo "#Sources"
	@echo "${SOURCES}"
	@echo "#Assets"
	@echo "${ASSETS}"
	@echo "#Shaders"
	@echo "${ASSET_SHADER}"
	@echo "#Fonts"   
	@echo "${ASSET_FONT}"   
	
build: ${BUILD_PRODUCT}


demo:
	@for f in ${SOURCES}; do cat $$f | while read l; do sleep 0.5; echo $$l ; done; done


deps:
	go get -v 
	
clean:
	-rm -f ${BUILD_PRODUCT} ${ASSETS} ${BUILD_NAME}-*-*-*

${BUILD_NAME}: ${BUILD_PRODUCT}
	cp -f ${BUILD_PRODUCT} ${BUILD_NAME}

${BUILD_PRODUCT}: ${SOURCES} ${ASSETS} ${PROTOS}
	go build -o ${BUILD_PRODUCT} ${BUILD_FLAGS} -gcflags all="${GCFLAGS}" -ldflags "${LDFLAGS}" 



proto: ${PROTOS}

facade/facade.pb.go: facade/facade.proto
	protoc -I facade $^ --go_out=plugins=grpc:facade
	



asset: ${ASSETS}

font/RobotoMono.ttf:
	mkdir -p font
	curl -o $@ https://raw.githubusercontent.com/google/fonts/master/apache/robotomono/RobotoMono-Regular.ttf

font/VT323.ttf:
	mkdir -p font
	curl -o $@ https://raw.githubusercontent.com/google/fonts/master/ofl/vt323/VT323-Regular.ttf

facade/shaderAssets.go: ${ASSET_SHADER}
	echo ""                                             >|$@
#	echo "// +build linux,arm"                          >>$@
	echo "package facade"                               >>$@
	echo "var ShaderAsset = map[string]string{"         >>$@
	for src in ${ASSET_SHADER}; do \
      name=$$(echo $$src) \
      name=$$(echo $$name | tr "[:upper:]" "[:lower:]") \
      name=$$(echo $$name | sed -e 's:shader/::'); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src | base64; \
      echo "\`,\n\n"; \
    done                                                >>$@
	echo "}"                                            >>$@


facade/fontAssets.go: ${ASSET_FONT}
	echo ""                                         >|$@
#	echo "// +build linux,arm"                      >>$@
	echo "package facade"                           >>$@
	echo "var FontAsset = map[string]string{"       >>$@
	for src in ${ASSET_FONT}; do \
      name=$$(echo $$src ) \
      name=$$(echo $$name | tr "[:upper:]" "[:lower:]") \
      name=$$(echo $$name | sed -e 's:font/::;s:\.[tT][tT][fFcC]::' ); \
      echo "\n\n\"$${name}\":\`";\
      cat $$src | base64; \
      echo "\`,\n\n"; \
    done                                            >>$@
	echo "}"                                        >>$@



.PHONY: help info build clean fetch asset demo proto

