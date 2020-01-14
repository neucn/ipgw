BUILD_ENV := CGO_ENABLED=0
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X ipgw/base.Version=${VERSION} -X ipgw/base.Build=${BUILD} -X ipgw/base.SavePath=${SAVEPATH}"
TARGET_DIR = build/${VERSION}

TARGET_EXEC := ipgw

.PHONY: all clean setup build-linux build-osx build-windows setup-linux setup-osx setup-windows pack-linux pack-osx pack-windows

all: clean setup build-linux build-osx build-windows

release: all pack-linux pack-osx pack-windows

clean:
	rm -rf ${TARGET_DIR}

setup: setup-linux setup-osx setup-windows

setup-linux:
	mkdir -p ${TARGET_DIR}/linux

setup-osx:
	mkdir -p ${TARGET_DIR}/osx

setup-windows:
	mkdir -p ${TARGET_DIR}/win


build-linux: setup-linux
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o ${TARGET_DIR}/linux/${TARGET_EXEC}

build-osx: setup-osx
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o ${TARGET_DIR}/osx/${TARGET_EXEC}

build-windows: setup-windows
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o ${TARGET_DIR}/win/${TARGET_EXEC}.exe

pack-linux:
	upx ${TARGET_DIR}/linux/${TARGET_EXEC}

pack-osx:
	upx ${TARGET_DIR}/osx/${TARGET_EXEC}

pack-windows:
	upx ${TARGET_DIR}/win/${TARGET_EXEC}.exe && cp install.bat ${TARGET_DIR}/win