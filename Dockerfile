# FROM golang:1.18-alpine as builder
# COPY go.mod go.sum /home/samir/workspace/lastingdynamics/ldgo/
# WORKDIR /code
# RUN go mod download
# COPY . /code/
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/ldgo github.com/Arka-cell/ldgo

# FROM alpine
# #RUN apk add --no-cache ca-certificates && update-ca-certificates
# COPY --from=builder /home/samir/workspace/lastingdynamics/ldgo /usr/bin/ldgo
# EXPOSE 8080 8080
# ENTRYPOINT ["/home/samir/workspace/lastingdynamics/"]

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
