FROM golang:latest AS build-env

RUN mkdir /app
WORKDIR /app
COPY . .
ENV CGO_ENABLED 0
RUN go build

FROM alpine:latest
RUN apk add ca-certificates
COPY --from=build-env /app/tezosign /
COPY --from=build-env /app/.secrets /
COPY --from=build-env /app/repos/migrations /repos/migrations
COPY --from=build-env /app/resources /resources

CMD ["/tezosign"]