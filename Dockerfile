FROM golang:1
WORKDIR /application
ADD go.* ./
RUN go mod download
ADD . .
RUN  go build -o api .

FROM debian
WORKDIR /application/
COPY --from=0 /application/api /application/api
CMD ["/application/api"]
