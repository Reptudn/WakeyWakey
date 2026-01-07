FROM golang:1.25 AS build
 
WORKDIR /bot
COPY ./src .
RUN CGO_ENABLED=0 go build -o wakeywakey .

FROM alpine:3.20

WORKDIR /bot
RUN apk add --no-cache awake

COPY --from=build /bot/wakeywakey .

RUN mkdir -p /bot/data && chmod 777 /bot/data

ENV PRODUCTION="true"

ENTRYPOINT ["/bot/wakeywakey"]