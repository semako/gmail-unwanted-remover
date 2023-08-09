FROM golang:1.21.0 as build

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gmail-unwanted-remover cmd/gmail-unwanted-remover/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/gmail-unwanted-remover ./

RUN chmod +x ./gmail-unwanted-remover

CMD "./gmail-unwanted-remover"
