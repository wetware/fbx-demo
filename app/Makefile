.PHONY: all clean gen build

FLAGS := CGO_ENABLED=0

HEADER_FILE := cap/header.txt
SRC_FILE := ../tiktok/api/cap/tiktok.capnp
CAP_FILE := cap/tiktok.capnp

all: clean update-cap build

gen:
	cat ${HEADER_FILE} ${SRC_FILE} > ${CAP_FILE}
	go generate ./...
	rm -f ${CAP_FILE}

build:
	go build -o app .

clean:
	rm -f ${CAP_FILE}
