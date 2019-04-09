# TODO

- Configure auth on the MongoDB server
- Revisit a `decimal` package, instead of using floats for money amounts
  - Wrap `decimal` with similar funcs to `mongoUUID`, to un/marshal as Decimal128?
- ~~Database & collection defaults on Server constructor~~ [DONE](https://github.com/jlucktay/rest-api/commit/1f69608b3b9c6ea1c31aa6b62a3ff0944152d05c)
- Logging
- Contexts
- ~~Better UUID handling ([ref](https://groups.google.com/forum/#!topic/mongodb-go-driver/vNHkY2EZq70))~~ [DONE](https://github.com/jlucktay/rest-api/commit/1f69608b3b9c6ea1c31aa6b62a3ff0944152d05c#diff-64e14639fdc8f8bdee63201031217aef)
- Revisit/leverage functionality in new chi router package
  - [Pagination example](https://github.com/go-chi/chi/blob/a86787d732a6ebbe0b7a70f61cd74c1ef9d88bd9/_examples/rest/main.go#L83)
  - [Middleware](https://godoc.org/github.com/go-chi/chi/middleware)

Also tracking TODOs on [this Trello board](https://trello.com/b/e4ZeAJp4/restful-http-api) which can be scraped on the
CLI with [this script](../scripts/trello.sh).
