FROM golang:1.12

WORKDIR /dionysus

ENV GO111MODULES=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT [ "go", "run", "main.go" ]