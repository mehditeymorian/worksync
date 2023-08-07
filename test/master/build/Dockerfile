FROM golang:1.20-alpine

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GO111MODULE=on go build -o main

RUN chmod +x ./main

CMD ./main