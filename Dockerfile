FROM golang:1.13.4-alpine3.10 AS api
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o redirector ./src
RUN /usr/sbin/adduser -D redirector

FROM scratch
WORKDIR /app
ENV PATH=/app
COPY robots.txt /app/robots.txt
COPY --from=api /app/redirector .
COPY --from=api /etc/passwd /etc/passwd
USER redirector
EXPOSE 8080 4433
CMD ["/app/redirector"]
