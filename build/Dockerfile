FROM golang:stretch as builder

COPY . /src/rest-api
WORKDIR /src/rest-api
ENV GO111MODULE=on
RUN go mod download

# buildmode pie: https://groups.google.com/forum/#!topic/golang-nuts/Jd9tlNc6jUE
# Can't use buildmode pie in MacOS: https://stackoverflow.com/a/3801032/380599
# Adding it here gives a different error for which I haven't yet run down the root cause.
# tags 'osusergo': https://golang.org/doc/go1.11#os/user
RUN CGO_ENABLED=0 GOOS=linux go build \
  -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o jra \
  -tags 'osusergo' \
  ./cmd/jra

FROM scratch
ARG MONGO_HOST=host.docker.internal
ENV MONGO_HOST=${MONGO_HOST}
COPY --from=builder /src/rest-api/jra /
ENTRYPOINT [ "/jra" ]

# TODO
# Make some use of this:
# https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build
# And/or check out Cloud Run on GCP
