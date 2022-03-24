from golang:1.17

WORKDIR /twitchbot

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .
RUN go build .

CMD ["./TwitchBot"]
