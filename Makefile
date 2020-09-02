
BUILD_NAME      = facade
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_RELEASE   = $(shell if echo ${BUILD_VERSION} | egrep -q '^[0-9]+\.[0-9]+\.[0-9]+$$'; then echo true; else echo false; fi )
BUILD_DATE     ?= $(shell if ${BUILD_RELEASE}; then date -u +"%Y-%m-%d"; else date -u +"%Y-%m-%dT%H:%M:%S%z"; fi)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_PLATFORM}





PROTOS  = facade/facade.pb.go
SOURCES = $(filter-out ${ASSETS} , $(wildcard *.go */*.go) )
ASSETS  = facade/shaderAssets.go facade/fontAssets.go facade/assets.go

FONTS ?= Monaco.ttf RobotoMono.ttf SpaceMono.ttf VT323.ttf Adore64.ttf OCRAExt.ttf
ASSET_FONT= $(foreach x,$(FONTS),font/$(x) )

SHADERS ?= def.vert def.frag 
SHADERS += color.vert color.frag 
SHADERS += grid/def.vert grid/def.frag grid/debug.frag grid/debug2.frag 
SHADERS += grid/wave.vert grid/roll.vert grid/rows.vert grid/crawl.vert grid/disk.vert grid/drop.vert 
SHADERS += set/def.vert set/def.frag set/scroll.vert set/field.vert set/debug.frag
SHADERS += mask/def.vert mask/def.frag mask/mask.frag mask/debug.frag 
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
	@echo " make deps     # fetch golang dependencies"
	@echo " make info     # show build info"
	@echo " make assets   # build fonts and shaders"
	@echo " make proto    # rebuild protobuf code"
	@echo " make clean    # clean up"
	@echo " make demo     # try 'make demo | facade pipe lines'"
	


info: 
	@echo "#Info"
	@echo " name       ${BUILD_NAME}"
	@echo " version    ${BUILD_VERSION}"
	@echo " platform   ${BUILD_PLATFORM}"
	@echo " date       ${BUILD_DATE}"
	@echo " product    ${BUILD_PRODUCT}"
	@echo "\n#Sources"
	@echo "${SOURCES}"
	@echo "\n#Assets"
	@echo "${ASSETS}"
	@echo "\n#Shaders"
	@echo "${ASSET_SHADER}"
	@echo "\n#Fonts"   
	@echo "${ASSET_FONT}"   
	
build: ${BUILD_PRODUCT}

demo:
	@for f in ${SOURCES}; do cat $$f | while read -r line; do echo "$$line"; sleep 0.5; done; sleep 2; done

	
	


deps:
	go get -v 
	
clean:
	-rm -f ${BUILD_PRODUCT} ${ASSETS} ${BUILD_NAME}-*-*-*


${BUILD_PRODUCT}: ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}
	cp -f ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM} ${BUILD_PRODUCT}

${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}: ${SOURCES} ${ASSETS} ${PROTOS}
	go build -o ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM} ${BUILD_FLAGS} -gcflags all="${GCFLAGS}" -ldflags "${LDFLAGS}" 



proto: ${PROTOS}

facade/facade.pb.go: facade/facade.proto
	protoc -I facade $^ --go_out=plugins=grpc:facade
	



assets: ${ASSETS}

font/Monaco.ttf:
	mkdir -p font
	curl -L -o $@.zip https://www.cufonfonts.com/download/font/monaco
	unzip -j -o -b $@.zip $$(basename $@) -d font/  && unlink $@.zip

font/RobotoMono.ttf:
	mkdir -p font
	curl -L -o $@ https://raw.githubusercontent.com/TypeNetwork/RobotoMono/master/fonts/ttf/RobotoMono-Regular.ttf

font/VT323.ttf:
	mkdir -p font
	curl -L -o $@ https://raw.githubusercontent.com/phoikoi/VT323/master/fonts/ttf/VT323-Regular.ttf

font/SpaceMono.ttf:
	mkdir -p font
	curl -L -o $@ https://raw.githubusercontent.com/googlefonts/spacemono/master/fonts/SpaceMono-Regular.ttf

font/OCRAExt.ttf:
	mkdir -p font
	curl -L -o $@ https://www.wfonts.com/download/data/2014/12/31/ocr-a-extended/OCRAEXT.TTF
	
font/Adore64.ttf:
	mkdir -p font
	curl -L -o $@.zip https://dl.dafont.com/dl/?f=adore64
	unzip -j -o -b $@.zip $$(basename $@) -d font/  && unlink $@.zip

font/amiga4ever.ttf:
	mkdir -p font
	curl -L -o $@.zip https://dl.dafont.com/dl/?f=amiga_forever
	unzip -j -o -b $@.zip $$(basename $@) -d font/ && unlink $@.zip

facade/assets.go: README.md
	echo ""                  >|$@
	echo "package facade"    >>$@
	echo "var readme = \`"   >>$@
	cat $^ | base64          >>$@
	echo "\`\n\n"            >>$@
	go fmt $@


facade/shaderAssets.go: ${ASSET_SHADER}
	echo "// +build linux,arm"                          >|$@
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
	go fmt $@


facade/fontAssets.go: ${ASSET_FONT}
	echo "// +build linux,arm"                      >|$@
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
	go fmt $@



.PHONY: help build deps info assets proto demo clean

