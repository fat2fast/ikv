FROM golang:1.24.4-alpine AS build

RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata

WORKDIR /app

COPY ./app /app


RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

FROM alpine
WORKDIR /var/
COPY --from=build /app .

COPY build/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

COPY build/wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

EXPOSE 3000
CMD ["docker-entrypoint.sh", "--", "./app"]

