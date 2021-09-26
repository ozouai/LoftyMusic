FROM gitpod/workspace-base
LABEL dazzle/layer=tool-docker
LABEL dazzle/test=tests/tool-docker.yaml
USER root
ENV TRIGGER_REBUILD=2
# https://docs.docker.com/engine/install/ubuntu/
RUN curl -o /var/lib/apt/dazzle-marks/docker.gpg -fsSL https://download.docker.com/linux/ubuntu/gpg \
    && apt-key add /var/lib/apt/dazzle-marks/docker.gpg \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
    && install-packages docker-ce=5:19.03.15~3-0~ubuntu-focal docker-ce-cli=5:19.03.15~3-0~ubuntu-focal containerd.io

RUN curl -o /usr/bin/slirp4netns -fsSL https://github.com/rootless-containers/slirp4netns/releases/download/v1.1.11/slirp4netns-$(uname -m) \
    && chmod +x /usr/bin/slirp4netns

RUN curl -o /usr/local/bin/docker-compose -fsSL https://github.com/docker/compose/releases/download/1.29.2/docker-compose-Linux-x86_64 \
    && chmod +x /usr/local/bin/docker-compose

# https://github.com/wagoodman/dive
RUN curl -o /tmp/dive.deb -fsSL https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_linux_amd64.deb \
    && apt install /tmp/dive.deb \
    && rm /tmp/dive.deb

WORKDIR /
ENV HOME="/root"
RUN apt-get update
RUN apt-get install -y nginx
RUN rm /etc/nginx/sites-enabled/default
RUN wget -q -c https://golang.org/dl/go1.17.1.linux-amd64.tar.gz -O- | tar -C /usr/local -xzvf -
ENV PATH="/usr/local/go/bin:${PATH}"
RUN go version
RUN mkdir -p /usr/local/nodejs
RUN wget -q -c https://nodejs.org/dist/v14.17.6/node-v14.17.6-linux-x64.tar.xz -O- | tar -C /usr/local/nodejs --strip-components 1 -xJvf -
ENV PATH="/usr/local/nodejs/bin:${PATH}"
RUN node --version
RUN npm install -g yarn
WORKDIR /usr/local
RUN wget -q -c https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip -O /tmp/protoc.zip
RUN unzip /tmp/protoc.zip
RUN chmod +x /usr/local/bin/protoc
RUN protoc --version
WORKDIR /
USER gitpod
ENV HOME="/home/gitpod"
RUN sudo chown -R gitpod:gitpod /home/gitpod
USER gitpod
ENV PATH="${PATH}:/home/gitpod/go/bin"
RUN go get -v golang.org/x/tools/gopls
RUN go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs
RUN go get -v github.com/ramya-rao-a/go-outline
RUN go get -v github.com/go-delve/delve/cmd/dlv
RUN go get -v github.com/go-delve/delve/cmd/dlv@master
RUN go get -v honnef.co/go/tools/cmd/staticcheck
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1