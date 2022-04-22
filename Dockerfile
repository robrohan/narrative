FROM golang:1.18-alpine as builder
RUN apk --no-cache add gcc g++ make ca-certificates git bash
WORKDIR /go/src/github.com/robrohan/narrative
COPY . .
RUN make build

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/robrohan/narrative/build/ ./
ENTRYPOINT ["./narrative"]
