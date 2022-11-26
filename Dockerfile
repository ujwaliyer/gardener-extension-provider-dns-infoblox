############# builder
FROM golang:1.19 AS builder

WORKDIR /gardener-extension-provider-dns-infoblox
COPY . .
RUN ls -a
RUN chown -R 777 .
RUN make install

