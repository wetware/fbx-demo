.PHONY: all clean gen build

FLAGS := CGO_ENABLED=0

LLM_HEADER_FILE := cap/header_llm.txt
LLM_SRC_FILE := ../caps/llm.capnp
LLM_CAP_FILE := cap/llm.capnp

TIKTOK_HEADER_FILE := cap/header_tiktok.txt
TIKTOK_SRC_FILE := ../caps/tiktok.capnp
TIKTOK_CAP_FILE := cap/tiktok.capnp

all: clean gen build

gen:
	cat ${LLM_HEADER_FILE} ${LLM_SRC_FILE} > ${LLM_CAP_FILE}
	cat ${TIKTOK_HEADER_FILE} ${TIKTOK_SRC_FILE} > ${TIKTOK_CAP_FILE}
	go generate ./...
	rm -f ${CAP_FILE}

build:
	go build -o app .

clean:
	rm -f ${LLM_CAP_FILE} ${TIKTOK_CAP_FILE}
