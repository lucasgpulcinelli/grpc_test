FROM golang:alpine
WORKDIR /grpc/

COPY functions /grpc/functions/
COPY server /grpc/server/

COPY go.mod go.sum /grpc/
RUN go build -v -o /grpc/main /grpc/server/


FROM alpine:latest

COPY --from=0 /grpc/main .
ENV PORT 50501
ENTRYPOINT ["./main"]