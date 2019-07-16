FROM golang:latest
MAINTAINER gravityslave92
ENV BOT_ID=784720809:AAGQBCIdvrtzbCLW2pxwHt1j0N93bUiMlfU
RUN mkdir /src
ADD . /src
WORKDIR /src
RUN go build -o main .
CMD ["/src/main"]

