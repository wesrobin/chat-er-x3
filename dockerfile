FROM golang:1.22

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN ls

EXPOSE 8071

WORKDIR /usr/src/app/server

RUN CGO_ENABLED=0 GOOS=linux go build -o /chat-er-x3

CMD ["/chat-er-x3"]

