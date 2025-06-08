.PHONY: all build build-arm clean gen install

FLAGS := CGO_ENABLED=0

HEADER_FILE := cap/header.txt
SRC_FILE := ../tiktok/api/cap/tiktok.capnp
CAP_FILE := cap/tiktok.capnp

all: clean build

build:
	cat ${HEADER_FILE} ${SRC_FILE} > ${CAP_FILE}
	go generate ./...
	env ${FLAGS} go build .
	rm -f ${CAP_FILE}

clean:
	rm -f ${CAP_FILE}
