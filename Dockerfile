FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD app /
ENTRYPOINT ["/app"]
