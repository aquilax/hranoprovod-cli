FROM golang:alpine3.15 AS builder

RUN apk update && apk add --no-cache git

RUN mkdir /build

WORKDIR /build

COPY . .

ENV GOPATH /tmp

RUN cd cmd/hranoprovod-cli && go get -d -v

ENV CGO_ENABLED 0

RUN cd cmd/hranoprovod-cli && go build -o /go/bin/hranoprovod-cli


FROM scratch

USER 1001

COPY --from=builder /go/bin/hranoprovod-cli /app/hranoprovod-cli

ENTRYPOINT ["/app/hranoprovod-cli"]