FROM golang:alpine
WORKDIR /grpc/

COPY go.mod go.sum /grpc/
COPY functions /grpc/functions/

COPY client /grpc/client/
RUN go build -v -o /grpc/main /grpc/client/


FROM alpine:latest
ENV ADDR grpc-server:50501

COPY --from=0 /grpc/main .
ENTRYPOINT ["./main"]
CMD ["--echotype=all", "world", "5 world", "000", "111", "222"]