FROM golang:1.22.2

WORKDIR /app
COPY . .
RUN go mod download
RUN apt-get update && apt-get install -y git
RUN go install github.com/cosmtrek/air@latest

ENV PATH="/go/bin:${PATH}"
CMD ["air", "-c", ".air.toml"]