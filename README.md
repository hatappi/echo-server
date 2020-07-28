# echo-server

echo-server is HTTP server.
This server responsed Request Header and message.

## Usage

```sh
$ docker run --rm -it -p 3000:3000 hatappi/echo-server:v0.1

$ curl localhost:3000
# or
$ curl localhost:3000?message=bar
```