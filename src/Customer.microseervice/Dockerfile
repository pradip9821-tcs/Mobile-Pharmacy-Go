FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

ENV PORT 8082

RUN go build

CMD [ "./customer.microservice" ]