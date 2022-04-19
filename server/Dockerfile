FROM golang:alpine
ARG GRPC_TEST_VERSION=v0.0.1
RUN go install github.com/lucasgpulcinelli/grpc_test/server@$GRPC_TEST_VERSION

FROM alpine:latest
COPY --from=0 /go/bin/server .
ENV PORT 50501
CMD ["./server"]