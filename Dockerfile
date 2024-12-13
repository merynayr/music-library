FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOPATH=/


RUN go build -o music-library ./cmd/app/main.go

CMD ["./music-library"]