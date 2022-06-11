FROM golang:1
WORKDIR /plugins_builder
COPY . .
RUN for dir in /plugins_builder/*; do (cd $dir && go build -buildmode=plugin -o /plugins/ .); done
RUN ls /plugins/
FROM api-gw
ENV PLUGINS_DIR=/plugins
COPY --from=0 /plugins/ /plugins/