FROM golang:1
WORKDIR /plugins_builder/
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -buildmode=plugin -o /plugins/plugin_a plugin.go