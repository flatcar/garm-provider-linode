FROM docker.io/golang:alpine

WORKDIR /root
USER root

RUN apk add --quiet \
    autoconf \
    gcc \
    git \
    g++ \
    libblkid \
    libtool \
    linux-headers \
    make \
    musl-dev \
    m4 \
    util-linux-dev

RUN wget --quiet http://musl.cc/aarch64-linux-musl-cross.tgz -O /tmp/aarch64-linux-musl-cross.tgz && \
    tar --strip-components=1 -C /usr/local -xzf /tmp/aarch64-linux-musl-cross.tgz && \
    rm /tmp/aarch64-linux-musl-cross.tgz

ADD ./scripts/build-static.sh /build-static.sh
RUN chmod +x /build-static.sh

CMD ["/bin/sh"]
