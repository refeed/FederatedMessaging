FROM golang:1.16-alpine

LABEL maintainer="Rafid Aslam"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /federatedMessagingServer

EXPOSE 80

ENV PORT 80

CMD [ "/federatedMessagingServer" ]
