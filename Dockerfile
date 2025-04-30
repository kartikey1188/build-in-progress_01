FROM golang:1.24.2-bookworm

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "cmd/build-in-progress_01/main.go"]