.PHONY: all clean gen build

FLAGS := CGO_ENABLED=0

HEADER_FILE := cap/header.txt
SRC_FILE := ../caps/llm.capnp
CAP_FILE := cap/llm.capnp

all: clean clean build

gen:
	cat ${HEADER_FILE} ${SRC_FILE} > ${CAP_FILE}
	go generate ./...
	rm -f ${CAP_FILE}

build:
	go build -o llm .

clean:
	rm -f ${CAP_FILE}
