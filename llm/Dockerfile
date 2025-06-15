FROM golang:1.24.3-alpine
RUN apk add --no-cache git make capnproto-dev libc6-compat gcompat
RUN go install capnproto.org/go/capnp/v3/capnpc-go@latest
RUN mkdir -p /go/src/capnproto.org/go/capnp && \
    git clone https://github.com/capnproto/go-capnp /go/src/capnproto.org/go/capnp
COPY . /llm
RUN cd /llm && \
    make clean build
