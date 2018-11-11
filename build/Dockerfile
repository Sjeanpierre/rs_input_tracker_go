from alpine:3.7

WORKDIR /app

RUN apk update

RUN apk --no-cache add \
  bash \
  wget \
  openssl \
  ca-certificates \
  curl \
  tar \
  less && \
  rm -rf /var/cache/apk/*
RUN mkdir /lib64 && ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN wget -O- https://bin.equinox.io/c/VdrWdbjqyF/cloudflared-stable-linux-amd64.tgz | tar xz


COPY ./argo_tunnel.pem.enc /app/
COPY ./startup.sh /app/
COPY ./bin/input_tracker_app /app/
EXPOSE 9080
ENTRYPOINT ["./startup.sh"]
