BUILD_ENV := CGO_ENABLED=0
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X ipgw/base/cfg.Version=${VERSION} -X ipgw/base/cfg.Build=${BUILD}"

TARGET_EXEC := ipgw

.PHONY: all clean setup build-linux build-osx build-windows setup-linux setup-osx setup-windows

all: clean setup build-linux build-osx build-windows

release: all pack-linux pack-osx pack-windows

clean:
	rm -rf build

setup: setup-linux setup-osx setup-windows

setup-linux:
	mkdir -p build/linux

setup-osx:
	mkdir -p build/osx

setup-windows:
	mkdir -p build/windows


build-linux: setup-linux
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/linux/${TARGET_EXEC}

build-osx: setup-osx
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o build/osx/${TARGET_EXEC}

build-windows: setup-windows
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o build/windows/${TARGET_EXEC}.exe

pack-linux:
	upx build/linux/${TARGET_EXEC}

pack-osx:
	upx build/osx/${TARGET_EXEC}

pack-windows:
	upx build/windows/${TARGET_EXEC}.exe