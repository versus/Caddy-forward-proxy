FROM amazonlinux:2

RUN yum -y upgrade && yum -y install python3 python3-devel python3-pip which tar gzip unzip git gcc make wget

ENV GO_VERSION 1.13.4
RUN curl -sSL https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz | tar -C /usr/local -xz
RUN export GO111MODULE=off
ENV GOPATH /go
ENV GOCACHE /tmp/gocache
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

ENV METALINTER_VERSION=2.0.11
RUN curl -sSL https://github.com/alecthomas/gometalinter/releases/download/v$METALINTER_VERSION/gometalinter-$METALINTER_VERSION-linux-amd64.tar.gz | tar -xz && mv ./gometalinter-$METALINTER_VERSION-linux-amd64/* /usr/local/bin/

ENV MIGRATE_VERSION 3.5.4
RUN curl -sSL https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz && mv ./migrate.linux-amd64 /usr/local/bin/migrate

ENV GODEP_VERSION 0.5.1
RUN curl -sSL https://github.com/golang/dep/releases/download/v$GODEP_VERSION/dep-linux-amd64 -o /usr/local/bin/dep && chmod +x /usr/local/bin/dep

ENV CI_LINT_VERSION 1.21.0
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b /usr/local/bin v$CI_LINT_VERSION

ENV GOOSE_VERSION 2.6.5
RUN curl -sSL --output $GOPATH/bin/goose https://github.com/lonja/goose/releases/download/v$GOOSE_VERSION/goose-linux64 && chmod +x $GOPATH/bin/goose

ENV DATA_DIR=/var/data
RUN mkdir -p $DATA_DIR
RUN curl -sSL --output $DATA_DIR/user-agents.yml https://raw.githubusercontent.com/ua-parser/uap-core/master/regexes.yaml


RUN export GO111MODULE=on

# copy code
COPY . ./proxy/
RUN cd ./proxy

# build binary
ARG VERSION
WORKDIR proxy
RUN go build ./cmd/proxy/

# download and extract GeoIP DB
RUN wget -O geolite.tar https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz
RUN tar -xvf geolite.tar --strip 1

# prepare artifact
RUN mkdir proxy_app
RUN mv ./proxy ./proxy_app/proxy
RUN mv ./Caddyfile ./proxy_app/Caddyfile
RUN mv ./GeoLite2-City.mmdb ./proxy_app/GeoLite2-City.mmdb
RUN mv ./appspec.yml ./proxy_app/appspec.yml
RUN mv ./install_dependencies.sh ./proxy_app/install_dependencies.sh
RUN mv ./proxy.service ./proxy_app/proxy.service
RUN mv ./start.sh ./proxy_app/start.sh
RUN mv ./stop.sh ./proxy_app/stop.sh
RUN pwd
