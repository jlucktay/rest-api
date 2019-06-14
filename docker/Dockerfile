FROM golang:stretch as builder

COPY . /src/rest-api/
WORKDIR /src/rest-api
RUN go mod download

# buildmode pie: https://groups.google.com/forum/#!topic/golang-nuts/Jd9tlNc6jUE
# Can't use buildmode pie in MacOS: https://stackoverflow.com/a/3801032/380599
# Adding it here gives a different error for which I haven't yet run down the root cause.
# tags 'osusergo': https://golang.org/doc/go1.11#os/user
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o rest-api \
    -tags 'osusergo' \
    ./cmd/api

FROM scratch
COPY --from=builder /src/rest-api/rest-api /
ENTRYPOINT [ "/rest-api", "--mongo-host", "host.docker.internal" ]

# TODO
# Make some use of this:
# https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build
