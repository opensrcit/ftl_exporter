FROM golang:alpine as builder

ARG OS=linux
ARG ARCH=amd64
ARG VERSION=development

WORKDIR /go/src/github.com/opensrcit/ftl-exporter
RUN apk --no-cache add git
COPY . .
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -trimpath -ldflags "-s -w -X version/version.Version=${VERSION}" -o ftl-exporter .


FROM scratch

WORKDIR /
COPY --from=builder /go/src/github.com/opensrcit/ftl-exporter/ftl-exporter ftl-exporter

ENTRYPOINT ["/ftl-exporter"]