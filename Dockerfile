FROM golang:1 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o api-gateway .

FROM scratch
COPY --from=builder /app/api-gateway /app/api-gateway
CMD ["/app/api-gateway"]