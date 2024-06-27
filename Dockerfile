FROM --platform=linux/amd64 golang:1.22.2-alpine as builder

ENV CGO_ENABLED=1
WORKDIR /app
COPY . .

RUN apk add git build-base

RUN go mod download

RUN  go build -tags=gogitlabjiradispatcher -o app cmd/go-gitlab-jira-dispatcher/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/views ./views

CMD ["./app"]
EXPOSE 8080
