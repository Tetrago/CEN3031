FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLE=0 GOOS=linux go build -o /api

CMD ["/api"]