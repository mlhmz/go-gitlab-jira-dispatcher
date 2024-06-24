FROM --platform=linux/amd64 golang:1.22.2-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -tags=gogitlabjiradispatcher -o app cmd/go-gitlab-jira-dispatcher/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
EXPOSE 8080
