
BUILD_NAME      = facade
BUILD_VERSION  ?= $(shell git describe --tags)
BUILD_RELEASE   = $(shell if echo ${BUILD_VERSION} | egrep -q '^[0-9]+\.[0-9]+\.[0-9]+$$'; then echo true; else echo false; fi )
BUILD_DATE     ?= $(shell if ${BUILD_RELEASE}; then date -u +"%Y-%m-%d"; else date -u +"%Y-%m-%dT%H:%M:%S%z"; fi)
BUILD_PLATFORM ?= $(shell go env GOOS )-$(shell go env GOARCH)
BUILD_PRODUCT   = ${BUILD_NAME}-${BUILD_PLATFORM}
BUILD_PACKAGE   = ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}.tgz


PROTOS  = facade/facade.pb.go
SOURCES = $(filter-out ${ASSETS} , $(wildcard *.go */*.go) )
ASSETS  = facade/shaderAssets.go facade/fontAssets.go facade/assets.go
EXTRAS  = README.md

FONTS ?= Monaco.ttf RobotoMono.ttf SpaceMono.ttf VT323.ttf Adore64.ttf OCRAExt.ttf
ASSET_FONT= $(foreach x,$(FONTS),font/$(x))

SHADERS ?= def.vert def.frag 
SHADERS += color.vert color.frag 
SHADERS += grid/def.vert grid/def.frag grid/debug.frag grid/debug2.frag 
SHADERS += grid/wave.vert grid/roll.vert grid/rows.vert grid/crawl.vert grid/disk.vert grid/vortex.vert 
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
	@echo "#FACADE Help"
	@echo " make info     # show build info"
	@echo " make build    # build static executable"
	@echo " make get      # fetch golang packages"
	@echo " make assets   # build fonts and shaders"
	@echo " make proto    # rebuild protobuf code"
	@echo " make clean    # remove binaries and golang objects"
	@echo " make package  # build platform package"
	@echo " make demo     # for 'make demo | facade pipe lines'"
	


info: 
	@echo "#FACADE Info"
	@echo " name       ${BUILD_NAME}"
	@echo " version    ${BUILD_VERSION}"
	@echo " platform   ${BUILD_PLATFORM}"
	@echo " date       ${BUILD_DATE}"
	@echo " product    ${BUILD_PRODUCT}"
	@echo " package    ${BUILD_PACKAGE}"
	@echo "\n#FACADE Sources"
	@echo "${SOURCES}"
	@echo "\n#FACADE Assets"
	@echo "${ASSETS}"
	@echo "\n#FACADE Shaders"
	@echo "${ASSET_SHADER}"
	@echo "\n#FACADE Fonts"
	@echo "${ASSET_FONT}"   
	@echo "\n#FACADE Extras"
	@echo "${EXTRAS}"
	
build: ${BUILD_PRODUCT}
	@echo "#FACADE built ${BUILD_PRODUCT}"

package: clean ${BUILD_PACKAGE}
	@echo "#FACADE packaged ${BUILD_PACKAGE}"

demo:
	@for f in ${SOURCES}; do cat $$f | while read -r line; do echo "$$line"; sleep 0.5; done; sleep 2; done

get:
	go get -v -u
	@echo "#FACADE got dependencies"

clean:
	-rm -f ${BUILD_PRODUCT} ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}
	-rm -rf package/${BUILD_PLATFORM}/
	go clean -r 
	@echo "#FACADE cleaned up"

touch:
	touch ${ASSET_SHADER} ${ASSET_FONT} ${EXTRAS}
	@echo "#FACADE touched assets"

assets: ${ASSETS}
	@echo "#FACADE built assets ${ASSETS}"

proto: ${PROTOS}
	@echo "#FACADE built proto ${PROTOS}"

	

${BUILD_PRODUCT}: ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}
	cp -f ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM} ${BUILD_PRODUCT}

${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM}: ${SOURCES} ${PROTOS} ${ASSETS}
	go build -o ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM} ${BUILD_FLAGS} -gcflags all="${GCFLAGS}" -ldflags "${LDFLAGS}" 

${BUILD_PACKAGE}: ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM} ${EXTRAS}
	@if ${BUILD_RELEASE}; then true; else { echo "REFUSE TO RELEASE UNTAGGED VERSION ${BUILD_VERSION}"; false; }; fi;
	mkdir -p package/${BUILD_PLATFORM}/
	cp -f ${BUILD_NAME}-${BUILD_VERSION}-${BUILD_PLATFORM} package/${BUILD_PLATFORM}/${BUILD_NAME}
	cp -f ${EXTRAS} package/${BUILD_PLATFORM}/
	cd package/${BUILD_PLATFORM}/ \
    && tar cfz ../${BUILD_PACKAGE} ${BUILD_NAME} ${EXTRAS} \
    && cd ..

facade/facade.pb.go: facade/facade.proto
	protoc -I facade $^ --go_out=plugins=grpc:facade

font/Monaco.ttf:
	mkdir -p font/
	curl -L -o $@.zip https://www.cufonfonts.com/download/font/monaco
	unzip -j -o -b $@.zip $$(basename $@) -d font/  && unlink $@.zip

font/RobotoMono.ttf:
	mkdir -p font/
	curl -L -o $@ https://raw.githubusercontent.com/TypeNetwork/RobotoMono/master/fonts/ttf/RobotoMono-Regular.ttf

font/VT323.ttf:
	mkdir -p font/
	curl -L -o $@ https://raw.githubusercontent.com/phoikoi/VT323/master/fonts/ttf/VT323-Regular.ttf

font/SpaceMono.ttf:
	mkdir -p font/
	curl -L -o $@ https://raw.githubusercontent.com/googlefonts/spacemono/master/fonts/SpaceMono-Regular.ttf

font/OCRAExt.ttf:
	mkdir -p font/
	curl -L -o $@ https://www.wfonts.com/download/data/2014/12/31/ocr-a-extended/OCRAEXT.TTF
	
font/Adore64.ttf:
	mkdir -p font/
	curl -L -o $@.zip https://dl.dafont.com/dl/?f=adore64
	unzip -j -o -b $@.zip $$(basename $@) -d font/  && unlink $@.zip

font/amiga4ever.ttf:
	mkdir -p font/
	curl -L -o $@.zip https://dl.dafont.com/dl/?f=amiga_forever
	unzip -j -o -b $@.zip $$(basename $@) -d font/ && unlink $@.zip

facade/assets.go: README.md
	echo ""                                 >|$@
	echo "package facade"                   >>$@
	echo "var Asset = map[string]string{"   >>$@
	echo "\n\"README\":\`"                  >>$@
	cat README.md | base64                  >>$@
	echo "\`,\n\n"                          >>$@
	echo "}"                                >>$@
	go fmt $@


facade/shaderAssets.go: ${ASSET_SHADER}
	echo ""                                             >|$@
	echo "package facade"                               >>$@
	echo "var ShaderAsset = map[string]string{"         >>$@
	for src in ${ASSET_SHADER}; do \
      name=$$(echo $$src) \
      name=$$(echo $$name | tr "[:upper:]" "[:lower:]") \
      name=$$(echo $$name | sed -e 's:shader/::'); \
      echo "\n\n\"$${name}\":\`";\
      if [ "${BUILD_PLATFORM}" = "linux-arm" ]; then cat $$src | base64; fi; \
      echo "\`,\n\n"; \
    done                                                >>$@
	echo "}"                                            >>$@
	go fmt $@


facade/fontAssets.go: ${ASSET_FONT}
	echo ""                                         >|$@
	echo "package facade"                           >>$@
	echo "var FontAsset = map[string]string{"       >>$@
	for src in ${ASSET_FONT}; do \
      name=$$(echo $$src ) \
      name=$$(echo $$name | tr "[:upper:]" "[:lower:]") \
      name=$$(echo $$name | sed -e 's:font/::;s:\.[tT][tT][fFcC]::' ); \
      echo "\n\n\"$${name}\":\`";\
      if [ "${BUILD_PLATFORM}" = "linux-arm" ]; then cat $$src | base64; fi; \
      echo "\`,\n\n"; \
    done                                            >>$@
	echo "}"                                        >>$@
	go fmt $@



.PHONY: help build package get info assets proto demo touch clean

