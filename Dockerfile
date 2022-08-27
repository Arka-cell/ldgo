FROM golang:1.18-alpine

LABEL maintainer="Samir Ahmane"

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base
RUN mkdir /app
WORKDIR /app
COPY . .
COPY .env .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build
EXPOSE 8080
# Run the executable
CMD [ "/build" ]
