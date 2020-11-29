# sse

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]

Server-Sent Events (SSE) library for Go.

See https://www.w3.org/TR/eventsource for the technical specification.

## Features

* Simple API.
* Performant.
* Dependency-free.
* Low-level API to build a server.

## Install

Go version 1.13

```
go get github.com/cristalhq/sse
```

## Example

As a simple HTTP handler:
```go
http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
    stream, err := sse.UpgradeHTTP(r, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data := struct {
        Text string `json:"text"`
    }{
        Text: "hey there",
    }
    stream.WriteJSON("123", "msg", data)
})
```

See [example_test.go](https://github.com/cristalhq/sse/blob/master/example_test.go) for more.

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/sse/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/sse/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/sse
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/sse
[reportcard-img]: https://goreportcard.com/badge/cristalhq/sse
[reportcard-url]: https://goreportcard.com/report/cristalhq/sse
[coverage-img]: https://codecov.io/gh/cristalhq/sse/branch/master/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/sse
