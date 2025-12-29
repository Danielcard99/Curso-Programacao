FROM golang:tip-alpine3.23 AS builder

RUN apk add --no-cache upx

WORKDIR /usr/src/app

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Gera um binário estático, ideal para rodar em scratch.
RUN go build -trimpath -ldflags="-s -w" -o app .

#Compila o Go, removendo metadados e deixando o binário menor.
RUN upx --best --lzma /usr/src/app/app

FROM scratch
COPY --from=builder /usr/src/app/app /app
EXPOSE 8080 
ENTRYPOINT ["/app"]
