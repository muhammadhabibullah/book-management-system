FROM golang:1.15-alpine AS build

RUN apk add --no-cache git gcc libc-dev curl

WORKDIR /build

ADD . .

RUN go get -v

RUN go test -v -race ./...

RUN CGO_ENABLED=0 go build -o book_management_system

FROM alpine

WORKDIR /usr/local/bin

COPY --from=build /build/book_management_system .
COPY --from=build /build/configs ./configs

RUN chmod +x book_management_system

CMD ["book_management_system", "serve"]
