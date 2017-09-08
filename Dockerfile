FROM golang:1.9
WORKDIR /go/src/github.com/fank/jira-tempo-importer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN adduser -D -u 999 jira
USER jira

# Add api
COPY --from=0 /go/src/github.com/fank/jira-tempo-importer/app /apps

# This container will be executable
ENTRYPOINT ["/app"]
