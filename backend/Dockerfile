# ---- Stage 1: build the binary ----
FROM golang:1.27 AS build

WORKDIR /src

# Copy the source (filtered by .dockerignore) and compile it.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Stage 2: minimal runtime image ----
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/server /server

# The server reads PORT at startup; 8080 is the default.
ENV PORT=8080
EXPOSE 8080

# Never run as root inside the container.
USER nonroot:nonroot

ENTRYPOINT ["/server"]