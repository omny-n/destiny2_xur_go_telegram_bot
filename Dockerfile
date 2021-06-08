FROM golang:alpine

COPY . /app 
WORKDIR /app 
RUN go get github.com/go-telegram-bot-api/telegram-bot-api
RUN go build -o main . 

CMD ["/app/main"]