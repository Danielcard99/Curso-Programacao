FROM golang:1.25-alpine
COPY . .
RUN go build -o server .
CMD [ "./server" ]