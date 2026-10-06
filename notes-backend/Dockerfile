FROM golang:1.23.2-bullseye
RUN apt-get update && \
    apt-get install -y sqlite3 libsqlite3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
ENV APP_AUTH="/app/secrets" \
    SERVER_BIND_TO="0.0.0.0:8001"\
    DB_ENV="dev"

WORKDIR /app

COPY . .
RUN go mod download && go mod tidy 


RUN go build /app/cmd/main.go

EXPOSE 8001

CMD [ "/app/main" ]
