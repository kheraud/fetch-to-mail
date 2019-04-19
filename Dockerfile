FROM scratch

MAINTAINER Karim Heraud <karim@sowefund.com>

WORKDIR /opt

# Must be built with CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
ADD fetch-to-mail /opt/

ADD ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/opt/fetch-to-mail"]
CMD [""]
