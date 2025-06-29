FROM golang:1.24.0

RUN echo "Asia/Ho_Chi_Minh" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata

WORKDIR /app

COPY ./app /app

# Copy the air configuration
COPY ./app/.air.toml /app/.air.toml
# Cài đặt air để auto reload
RUN go install github.com/air-verse/air@latest
# Cài đặt golang-migrate CLI
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#ENV GOPATH /go
#ENV PATH $PATH:/go/bin:$GOPATH/bin
#ENV CGO_ENABLED 0

COPY build/wait-for-it.sh /usr/local/bin/wait-for-it.sh
COPY build/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

RUN go mod download

# wait-for-it.sh $MEB_HOST:$MEB_PORT --strict --timeout=5 -- echo "Service Bus is up"
CMD ["docker-entrypoint.sh", "--", "air", "-c", ".air.toml"]

#CMD ["air", "-c", ".air.toml"]
#ENTRYPOINT ["./main"]

