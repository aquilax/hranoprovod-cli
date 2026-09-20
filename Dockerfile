FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache git

RUN mkdir /build

WORKDIR /build

COPY . .

ENV GOPATH=/tmp \
	CGO_ENABLED=0 \
	GOOS=$TARGETOS \
	GOARCH=$TARGETARCH

RUN cd cmd/hranoprovod-cli && go build -o /go/bin/hranoprovod-cli


FROM scratch

USER 1001

COPY --from=builder /go/bin/hranoprovod-cli /app/hranoprovod-cli

ENTRYPOINT ["/app/hranoprovod-cli"]