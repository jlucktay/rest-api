FROM golang:stretch as builder

RUN go get -d -u -v github.com/go-chi/chi
RUN go get -d -u -v github.com/gofrs/uuid
RUN go get -d -u -v go.mongodb.org/mongo-driver/...

COPY . /go/src/github.com/jlucktay/rest-api/
WORKDIR /go/src/github.com/jlucktay/rest-api

ENV CGO_ENABLED 0
ENV GOOS linux

# buildmode pie: https://groups.google.com/forum/#!topic/golang-nuts/Jd9tlNc6jUE
# Can't use buildmode pie: https://stackoverflow.com/a/3801032/380599
#
# tags 'osusergo': https://golang.org/doc/go1.11#os/user
RUN go build \
    -a \
    -installsuffix cgo \
    -ldflags '-extldflags "-static"' \
    -o main \
    -tags 'osusergo' \
    ./cmd/api

FROM scratch
COPY --from=builder /go/src/github.com/jlucktay/rest-api/main /
ENTRYPOINT [ "/main", "--mongo-host", "host.docker.internal" ]

# TODO
# Make some use of this:
# https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build
