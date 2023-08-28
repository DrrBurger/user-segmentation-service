FROM golang:1.21

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o user-segmentation-service cmd/segmentation-service/main.go

EXPOSE 8080

CMD ["./user-segmentation-service"]
