# OpenTracing x Golang

## Prerequirement
* Go
* Docker

## Document(Blog)

[here(TBD)](./)

## Run

1. Run jaeger
    ```
    docker run --rm -p 6831:6831/udp -p 6832:6832/udp -p 16686:16686 jaegertracing/all-in-one:1.7 --log-level=debug
    ```
1. Run server
    ```
    go run  ./server/*.go
    ```
1.  Run client
    ```
    go run  ./client/*.go
    ```
1. See `localhost:16686` in any browser.


# License

This is released under the MIT License. See [License](./LICENSE)
