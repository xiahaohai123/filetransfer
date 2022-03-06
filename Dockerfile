FROM golang AS builder
WORKDIR /build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"
COPY . .
RUN go build -o filetransfer web/main.go

FROM golang
WORKDIR /usr/local/bin/
COPY --from=builder /build/filetransfer .
CMD ./filetransfer