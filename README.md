# jlucktay's RESTful API (`jra`)

## Description

This is my implementation of a RESTful HTTP API capable of CRUD operations and persisting resource state to a database.

### Implementation guidelines

- Follow best practices, for example TDD/BDD, with a focus on full-stack testing.
- Prioritize correctness, robustness, and extensibility over extra features and optimizations.
- Write code with the quality bar one would expect to see in production.
- Try to simplify by using open source frameworks and libraries where possible.

## Badges

[![Build Status](https://travis-ci.org/jlucktay/rest-api.svg?branch=master)][badge-travis]
[![codecov](https://codecov.io/gh/jlucktay/rest-api/branch/master/graph/badge.svg)][badge-codecov]
[![GoDoc](https://godoc.org/github.com/jlucktay/rest-api?status.svg)][badge-godoc]
[![Go Report Card](https://goreportcard.com/badge/github.com/jlucktay/rest-api)][badge-goreportcard]
[![License](https://img.shields.io/github/license/jlucktay/rest-api.svg)][badge-license]

## Installation

### Prerequisites

You should have a [working Go environment](https://golang.org/doc/install) and have `$GOPATH/bin` in your `$PATH`.

### Compiling

To download the source, compile, and install the relevant binaries, run:

``` shell
go get github.com/jlucktay/rest-api/...
```

The source code will be located in `$GOPATH/src/github.com/jlucktay/rest-api`.

Newly compiled `jra` and `jrams` binaries will be in `$GOPATH/bin/`.

## Usage

Launching the API server:

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

Here is [the full design doc][design-doc] for this API, which describes the various endpoints, how to call them, and
what to expect in return.

## Roadmap

Features and functionality yet to be implemented are captured in [the TODO markdown file in this repo](./docs/TODO.md)
as well as on [the Trello board][trello].

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)

[badge-codecov]: https://codecov.io/gh/jlucktay/rest-api
[badge-godoc]: https://godoc.org/github.com/jlucktay/rest-api
[badge-goreportcard]: https://goreportcard.com/report/github.com/jlucktay/rest-api
[badge-license]: https://github.com/jlucktay/rest-api/blob/master/LICENSE
[badge-travis]: https://travis-ci.org/jlucktay/rest-api
[design-doc]: https://docs.google.com/document/d/1xtqwQDhdwTe3BUEyf3lGWycPIvl66uxDdJgHLqa9hz4
[trello]: https://trello.com/b/e4ZeAJp4
