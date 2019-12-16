FROM golang:1.13.4-alpine3.10 AS api
WORKDIR /opt/redirector
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o redirector ./src
RUN /usr/sbin/adduser -D redirector

FROM scratch
WORKDIR /opt/redirector
ENV PATH=/opt/redirector
COPY robots.txt /opt/redirector/robots.txt
COPY --from=api /opt/redirector/redirector .
COPY --from=api /etc/passwd /etc/passwd
USER redirector
EXPOSE 8080
CMD ["/opt/redirector/redirector"]
