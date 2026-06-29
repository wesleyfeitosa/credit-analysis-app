# syntax=docker/dockerfile:1

FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /bin/api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /bin/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
