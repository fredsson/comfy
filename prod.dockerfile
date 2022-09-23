FROM golang:1.19-alpine

RUN apk add make

WORKDIR /app
COPY ./ ./

RUN go mod download

RUN make build

CMD ["./bin/comfy"]
