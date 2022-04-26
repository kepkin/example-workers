FROM golang:1.18-alpine as build-stage
WORKDIR /opt
RUN apk add git

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/server


FROM alpine:3.15

WORKDIR /opt/app
RUN apk add libc6-compat

COPY --from=build-stage /opt/server .
COPY ./repository/migrations ./repository/migrations
CMD ["/opt/app/server"]