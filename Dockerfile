FROM golang:1.13.4-alpine3.10 AS api
WORKDIR /opt/redirector
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o redirector main.go

FROM alpine:3.10
WORKDIR /opt/redirector
ENV PATH=/opt/redirector
COPY --from=api /opt/redirector/redirector .
COPY robots.txt /opt/redirector/robots.txt
RUN /usr/sbin/adduser -D redirector
RUN /bin/chown -R redirector /opt/redirector
USER redirector
EXPOSE 8080
CMD ["/opt/redirector/redirector"]
