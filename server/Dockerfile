FROM golang:alpine
RUN apk add git
WORKDIR /server
COPY go.mod go.sum server.go ./
RUN go mod download
RUN go build .

FROM alpine:latest
COPY --from=0 /server/server .
ENV PORT 50501
CMD ["./server"]
