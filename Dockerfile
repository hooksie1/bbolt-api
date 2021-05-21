FROM golang:1.16-alpine as builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates make
RUN update-ca-certificates 
RUN mkdir /database && touch /database/bolt.db
ADD . /app/
RUN make linux


FROM scratch

COPY --from=builder /app/bbolt-api .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /database /database
ENV DATABASE_PATH=/database/bolt.db
ENV SERVER_PORT=8080

CMD ["./bbolt-api", "start"]
