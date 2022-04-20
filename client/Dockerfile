FROM golang:alpine
WORKDIR /grpc/

COPY go.mod go.sum /grpc/
COPY functions /grpc/functions/

COPY client /grpc/client/
RUN go build -v -o /grpc/main /grpc/client/


FROM alpine:latest
ENV ADDR grpc_server:50501

COPY --from=0 /grpc/main .
CMD ["sh", "-c", "./main --addr $ADDR"]
