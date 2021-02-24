FROM golang:1.13.4-alpine3.10 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags '-extldflags "-static" -s -w' \
    -o redirector \
    main.go
RUN apk add --no-cache tzdata ca-certificates
RUN /usr/sbin/adduser -D user

FROM scratch AS url-redirector
LABEL maintainer="Caleb Woodbine <calebwoodbine.public@gmail.com>" \
  org.opencontainers.image.title="url-redirector" \
  org.opencontainers.image.description="A simple cloud-native golang + yaml URL redirector app" \
  org.opencontainers.image.url="https://gitlab.com/BobyMCbobs/url-redirector" \
  org.opencontainers.image.source="git@gitlab.com/BobyMCbobs/url-redirector" \
  org.opencontainers.image.vendor="BobyMCbobs" \
  org.opencontainers.image.licenses="AGPL-3.0" \
  org.opencontainers.image.documentation="https://bobymcbobs.gitlab.io/url-redirector" \
  org.opencontainers.image.authors="BobyMCbobs"
WORKDIR /app
ENV PATH=/app/bin
COPY robots.txt /app/robots.txt
COPY --from=build /app/redirector /app/bin/redirector
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl /etc/ssl
USER user
EXPOSE 8080 4433
CMD ["/app/bin/redirector"]
