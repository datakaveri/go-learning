# Teaching image. A real service pins approved image digests in its repository.
FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/dx-example-go ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/dx-example-go /dx-example-go
USER nonroot:nonroot
ENTRYPOINT ["/dx-example-go"]
