# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/songs

#Step 3: Final
FROM scratch
COPY --from=builder /app/.env /.env
COPY --from=builder /app/schema /schema
COPY --from=builder /bin/app /app
CMD ["/app"]
