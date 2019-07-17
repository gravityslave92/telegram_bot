FROM golang:latest
MAINTAINER gravityslave92
ENV BOT_ID=784720809:AAGQBCIdvrtzbCLW2pxwHt1j0N93bUiMlfU
ENV PROXY_URL="//105.234.154.69:8080"
RUN mkdir /src
ADD . /src
WORKDIR /src
RUN go build -o ./main cmd/telegram_bot.go
CMD ["/src/main"]

