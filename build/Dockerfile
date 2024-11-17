FROM alpine:latest
RUN apk add --no-cache libc6-compat
RUN mkdir -p /app/logs

COPY ./app /app/
WORKDIR /app
RUN chmod +x app

CMD ["./app"]