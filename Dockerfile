FROM golang:stretch as builder

# Set up for modules
ENV GO111MODULE=on
WORKDIR /src/rest-api

# These layers should stay cached unless dependencies change, speeding up image rebuilds
COPY go.mod /src/rest-api
COPY go.sum /src/rest-api
RUN go mod download

# Copy the rest of the code
COPY . /src/rest-api

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
ARG MONGO_CS=mongodb://host.docker.internal:27017
ENV MONGO_CS=${MONGO_CS}
COPY --from=builder /src/rest-api/jra /
ENTRYPOINT [ "/jra" ]

# TODO
# Make some use of this:
# https://cloud.google.com/appengine/docs/flexible/custom-runtimes/build
# And/or check out Cloud Run on GCP
