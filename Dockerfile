FROM golang:1.22.1 as development

# use a single directory for all Go caches to simpliy RUN --mount commands below
ENV GOPATH /cache/gopath
ENV GOCACHE /cache/gocache
ENV GOMODCACHE /cache/gomodcache

ENV CGO_ENABLED=0
ENV GOOS=linux


COPY go.mod go.sum /src/

WORKDIR /src

RUN --mount=type=cache,target=/cache <<EOF
set -ex

go mod download
go mod verify
EOF

# build stage
FROM golang:1.22.1 AS build-stage

# use the same directories for Go caches as above
ENV GOPATH /cache/gopath
ENV GOCACHE /cache/gocache
ENV GOMODCACHE /cache/gomodcache

# modules are already downloaded
ENV GOPROXY off

WORKDIR /src
COPY . .

# to add a dependency
COPY --from=development /src/go.mod /src/go.sum /src/

RUN --mount=type=cache,target=/cache <<EOF
set -ex

go build -v -o=bin/server -tags=streaming-server ./cmd/server/

EOF

FROM golang:1.22.1 AS deployment

# use static secrets (For testing purposes only)
ENV AWS_ACCESS_KEY_ID=AKIAQE3ROSJBD3TC7LGX
ENV AWS_SECRET_ACCESS_KEY=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
ENV port=8080

COPY --from=build-stage /src/bin/server /server

ENTRYPOINT [ "/server" ]
