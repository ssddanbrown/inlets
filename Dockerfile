FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.16-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG GIT_COMMIT
ARG VERSION

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPATH=/go/src/
WORKDIR /go/src/github.com/inlets/inlets

COPY .git               .git
COPY go.mod             .
COPY go.sum             .
COPY pkg                pkg
COPY cmd                cmd
COPY main.go            .

RUN test -z "$(gofmt -l $(find . -type f -name '*.go' ))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; } \
    && CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go test ./... -cover

# add user in this stage because it cannot be done in next stage which is built from scratch
# in next stage we'll copy user and group information from this stage
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -a -installsuffix cgo -o /usr/bin/inlets \
    && addgroup -S app \
    && adduser -S -g app app

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

ARG REPO_URL

LABEL org.opencontainers.image.source $REPO_URL

COPY --from=builder /usr/bin/inlets /usr/bin/

EXPOSE 8000
EXPOSE 8123

VOLUME /tmp/

ENTRYPOINT ["/usr/bin/inlets"]
CMD ["--help"]
