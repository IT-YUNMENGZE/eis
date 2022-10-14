FROM debian:stretch-slim

WORKDIR /

COPY bin/edge-scheduler /usr/local/bin

CMD ["edge-scheduler"]