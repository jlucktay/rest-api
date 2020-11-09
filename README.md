# jlucktay's RESTful API (`jra`)

[![Travis][Travis-badge]][Travis]
[![Codecov][Codecov-badge]][Codecov]
[![Go Report Card][Go Report Card-badge]][Go Report Card]
[![GoDoc][GoDoc-badge]][GoDoc]
[![License][License-badge]][License]

## Description

This is my implementation of a RESTful HTTP API capable of CRUD operations and persisting resource state to a database.

### Implementation guidelines

- Follow best practices, for example TDD/BDD, with a focus on full-stack testing.
- Prioritize correctness, robustness, and extensibility over extra features and optimizations.
- Write code with the quality bar one would expect to see in production.
- Try to simplify by using open source frameworks and libraries where possible.

## Installation

### Prerequisites

You should have a [working Go environment](https://golang.org/doc/install) and have `$GOPATH/bin` in your `$PATH`.
Mage is being used for various build/run/test tasks, and should also [be installed](https://magefile.org).

### Compiling

To download the source, compile, and install the relevant binaries, run:

``` shell
go get go.jlucktay.dev/rest-api/...
```

The source code will be located in `$GOPATH/src/go.jlucktay.dev/rest-api`.

Newly compiled `jra` and `jrams` binaries will be in `$GOPATH/bin/`.

#### Mage

The Magefile contains targets for various tasks, which can be listed out with `mage -l`.

## Usage

Launching the API server:

### Running with Docker

```shell
$ docker-compose up
Starting rest-api_web_1   ... done
Starting rest-api_mongo_1 ... done
Attaching to rest-api_mongo_1, rest-api_web_1
web_1    | Connected to MongoDB!
web_1    | Collection 'payments' contains 0 records.
...
```

The API server will be listening at <http://localhost:8080/v1/>.

### Running the server directly

``` shell
$ jra
Connected to MongoDB!
Collection 'payments' contains 0 records.
```

Seeding the server with some sample records:

``` shell
$ jrams
Continuing will delete ALL payment records in MongoDB (database: rest-api, collection: payments)
Press 'Enter' to continue, or CTRL+C to cancel...
Connected to MongoDB!
Collection 'payments' contains 14 records.
Collection 'payments' dropped.
Disconnected from MongoDB.
Connected to MongoDB!
Collection 'payments' contains 0 records.
Added payment with ID '4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43'.
...
```

Accessing the server:

``` shell
$ curl --silent --request GET http://localhost:8080/v1/payments
{
  "data": [
    {
      "attributes": {
...
```

### Documentation

Here is [the full design doc] for this API, which describes the various endpoints, how to call them, and what to expect
in return.

## Testing

There is a Mage target to run tests across all packages in the repo:

``` shell
$ mage test
?       go.jlucktay.dev/rest-api/cmd/jra    [no test files]
?       go.jlucktay.dev/rest-api/internal/cmd/jrams [no test files]
?       go.jlucktay.dev/rest-api/pkg/org    [no test files]
ok      go.jlucktay.dev/rest-api/pkg/server 0.187s
ok      go.jlucktay.dev/rest-api/pkg/storage        0.262s
?       go.jlucktay.dev/rest-api/pkg/storage/inmemory       [no test files]
?       go.jlucktay.dev/rest-api/pkg/storage/mongo  [no test files]
ok      go.jlucktay.dev/rest-api/test       0.128s
```

For more details on how Go itself discovers and executes tests, and the various flags with which to alter behaviour
when doing so, run `go help test` and `go help testflag`.

## Roadmap

Features and functionality yet to be implemented are captured in [the TODO markdown file in this repo](./docs/TODO.md)
as well as on [the Trello board].

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)

[Codecov-badge]: https://codecov.io/gh/jlucktay/rest-api/branch/master/graph/badge.svg
[Codecov]: https://codecov.io/gh/jlucktay/rest-api

[Go Report Card-badge]: https://goreportcard.com/badge/go.jlucktay.dev/rest-api
[Go Report Card]: https://goreportcard.com/report/go.jlucktay.dev/rest-api

[GoDoc-badge]: https://godoc.org/github.com/jlucktay/rest-api?status.svg
[GoDoc]: https://godoc.org/github.com/jlucktay/rest-api

[License-badge]: https://img.shields.io/github/license/jlucktay/rest-api.svg
[License]: https://github.com/jlucktay/rest-api/blob/master/LICENSE

[Travis-badge]: https://travis-ci.org/jlucktay/rest-api.svg?branch=master
[Travis]: https://travis-ci.org/jlucktay/rest-api

[the full design doc]: https://docs.google.com/document/d/1xtqwQDhdwTe3BUEyf3lGWycPIvl66uxDdJgHLqa9hz4
[the Trello board]: https://trello.com/b/e4ZeAJp4
