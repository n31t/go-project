FROM golang:1.21.0 as builder
WORKDIR /usr/src/app
COPY . .
RUN go mod download

EXPOSE 8081

CMD ["go", "run", "/usr/src/app/cmd/app", "-fill=true"]


