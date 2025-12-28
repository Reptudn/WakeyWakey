FROM golang:latest AS build
 
WORKDIR /bot

COPY . .

RUN CGO_ENABLED=0 go build -o wakeywakey .

FROM alpine:3.20

WORKDIR /bot

COPY --from=build /bot/wakeywakey .

ENTRYPOINT ["/bot/wakeywakey"]