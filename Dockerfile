############# builder
FROM golang:1.19.1 AS builder

WORKDIR /go/src/gardener-extension-provider-dns-infoblox
COPY . .
RUN go mod tidy && go mod vendor
RUN chmod -R 777 vendor/
RUN make install

############# base image
FROM alpine:3.13.7 AS base

############# gardener-extension-provider-dns-infoblox
FROM base AS gardener-extension-provider-dns-infoblox

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-provider-dns-infoblox /gardener-extension-provider-dns-infoblox
ENTRYPOINT ["/gardener-extension-provider-dns-infoblox"]
