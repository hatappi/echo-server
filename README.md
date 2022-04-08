# echo-server

echo-server is HTTP server.  
This server responsed Request Header and message.

## Request

### HTTP

```sh
$ curl -H "foo: bar" "localhost:3000?name=bar"
```

### gRPC

```sh
$ grpcurl -plaintext -H "foo: bar" -d '{"name": "foo"}' localhost:5000 echo.Echo.SayHello
```
