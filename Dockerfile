FROM golang:1.13.4-alpine3.10
WORKDIR /usr/app/

RUN apk add curl git make bash gcc musl-dev

# install kustomzie and sealed secret transformer
COPY . .
ENV XDG_CONFIG_HOME=/usr/app/
RUN git clone https://github.com/plaidev/kustomize-plugins.git \
    && cd kustomize-plugins \
    && make setup \
    && XDG_CONFIG_HOME=$XDG_CONFIG_HOME make build \
    && mv ./bin/kustomize /usr/local/bin/kustomize

# install sealed secret
RUN curl -L https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.9.5/kubeseal-linux-amd64 -o kubeseal \
    && install -m 755 kubeseal /usr/local/bin/kubeseal 

ENTRYPOINT []