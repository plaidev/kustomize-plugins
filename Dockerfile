FROM golang:1.14.2-alpine3.10
WORKDIR /usr/app/

RUN apk add curl git make bash gcc musl-dev

# install kustomzie and sealed secret transformer
COPY . ./kustomize-plugins
ENV XDG_CONFIG_HOME=/usr/app/
RUN cd kustomize-plugins \
    && make setup \
    && XDG_CONFIG_HOME=$XDG_CONFIG_HOME make build \
    && mv ./bin/kustomize /usr/local/bin/kustomize

# install sealed secret
RUN curl -L https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.9.5/kubeseal-linux-amd64 -o kubeseal \
    && install -m 755 kubeseal /usr/local/bin/kubeseal 

ENTRYPOINT []
