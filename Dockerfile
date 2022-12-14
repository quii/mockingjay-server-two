FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/mockingjay/main.go

FROM alpine:3.16.2
COPY --from=builder /app/specifications/examples ./examples
COPY --from=builder /app/svr .
RUN mkdir /tmp/testresources/
CMD [ "./svr" ]