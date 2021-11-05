FROM golang:1.16.3-alpine3.12 AS builder

ENV GOPATH ""

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /bin/echo-server .

FROM alpine:3.12

RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/echo-server /bin/echo-server

RUN addgroup -g 1001 echo-server && adduser -D -G echo-server -u 1001 echo-server
USER 1001

CMD ["/bin/echo-server"]

