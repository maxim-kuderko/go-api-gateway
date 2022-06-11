FROM golang:1
WORKDIR /application
ADD go.* ./
RUN go mod download
ADD . .
RUN  CGO_ENABLED=0 go build -o api main.go
CMD ["./api"]

FROM scratch
WORKDIR /application/
COPY --from=0 /application/api /application/api
CMD ["/application/api"]
