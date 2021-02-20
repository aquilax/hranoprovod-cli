FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR .

COPY . .

ENV GOPATH /tmp

RUN go get -d -v

ENV CGO_ENABLED 0

RUN go build -o /go/bin/hranoprovod-cli


FROM scratch

USER 1001

COPY --from=builder /go/bin/hranoprovod-cli /app/hranoprovod-cli

ENTRYPOINT ["/app/hranoprovod-cli"]