FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/vngitSub

FROM alpine:latest
WORKDIR /run
COPY --from=builder /app/vngitSub /apprun/vngitSub
EXPOSE 8000
CMD ["/apprun/vngitSub"]