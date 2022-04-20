FROM golang:alpine
WORKDIR /grpc/
COPY functions /grpc/functions/
COPY client /grpc/client/
COPY go.mod go.sum /grpc/
RUN go build -v -o /grpc/main /grpc/client/

FROM alpine:latest
COPY --from=0 /grpc/main .
ENV PORT 50501
CMD ["./main"]
