FROM golang:latest AS build

COPY . /app

WORKDIR /app

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod tidy

WORKDIR /app/cmd/api

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go build -o /app/main main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /app/main /app/main

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]