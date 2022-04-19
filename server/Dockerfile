FROM golang:alpine
RUN apk update && apk add git
WORKDIR /server
COPY build.sh .
RUN sh build.sh

FROM alpine:latest
COPY --from=0 /server/server .
ENV PORT 50501
CMD ["./server"]
