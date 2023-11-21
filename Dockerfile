FROM golang:1.21.1-alpine as base

FROM base as builder

RUN apk add --no-cache make git

WORKDIR /app
COPY . .

RUN make build

FROM base as app
COPY --from=builder /app/bin/ /app/bin/

WORKDIR /app/bin

EXPOSE 8080
CMD ["./FlexiProxyHub"]