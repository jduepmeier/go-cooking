FROM docker.io/library/golang:1.20.6-alpine AS builder

RUN apk add sqlite make gcc musl-dev
COPY . /app
WORKDIR /app
RUN make

FROM docker.io/library/alpine:latest

RUN apk add sqlite

RUN addgroup -g 10000 gocooking && adduser -G gocooking -S gocooking

COPY --from=builder /app/bin/gocooking /gocooking
COPY templates /app/templates
COPY static /app/static

USER gocooking

CMD [ "/gocooking" "serve" ]
