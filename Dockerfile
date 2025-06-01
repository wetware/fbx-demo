# TODO: create unprivileged user

FROM golang:1.23-alpine

WORKDIR /fbx-demo

COPY * ./

RUN go mod tidy && go mod download

# Copy source code
COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
