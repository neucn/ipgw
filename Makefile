BUILD_ENV := CGO_ENABLED=0
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X ipgw/base/cfg.Version=${VERSION} -X ipgw/base/cfg.Build=${BUILD} -X ipgw/base/cfg.SavePath=${SAVEPATH}"

TARGET_EXEC := ipgw

.PHONY: all clean setup build-linux build-osx build-windows setup-linux setup-osx setup-windows pack-linux pack-osx pack-windows

all: clean setup build-linux build-osx build-windows

release: all pack-linux pack-osx pack-windows

clean:
	rm -rf build/${VERSION}

setup: setup-linux setup-osx setup-windows

setup-linux:
	mkdir -p build/${VERSION}/linux

setup-osx:
	mkdir -p build/${VERSION}/osx

setup-windows:
	mkdir -p build/${VERSION}/win


build-linux: setup-linux
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/${VERSION}/linux/${TARGET_EXEC}

build-osx: setup-osx
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/${VERSION}/osx/${TARGET_EXEC}

build-windows: setup-windows
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/${VERSION}/win/${TARGET_EXEC}.exe

pack-linux:
	upx build/${VERSION}/linux/${TARGET_EXEC}

pack-osx:
	upx build/${VERSION}/osx/${TARGET_EXEC}

pack-windows:
	upx build/${VERSION}/win/${TARGET_EXEC}.exe