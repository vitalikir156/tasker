FROM golang:1.22-alpine as builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /tasker
FROM alpine:3.20.2
COPY --from=builder /tasker /tasker
EXPOSE 9989
CMD [ "/tasker"]