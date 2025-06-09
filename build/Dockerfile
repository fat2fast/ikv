FROM golang:1.24.4-alpine AS build

RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata

WORKDIR /app

COPY ./app /app

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

FROM alpine
WORKDIR /var/
COPY --from=build /app .
EXPOSE 3000
CMD ["./app"]

