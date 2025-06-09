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

RUN go mod download

CMD ["air", "-c", ".air.toml"]
#ENTRYPOINT ["./main"]

