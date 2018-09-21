FROM alpine:3.8

RUN apk add --no-cache ca-certificates

ADD /bin/hcloud-cloud-controller-manager /bin/

CMD ["/bin/hcloud-cloud-controller-manager"]
