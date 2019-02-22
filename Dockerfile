############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder
MAINTAINER ariefdfaltah

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates

# Create appuser
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/ariefdfaltah-golang-graphql
WORKDIR $GOPATH/src/ariefdfaltah-golang-graphql

# Fetch dependencies.
# Using go mod with go 1.11
RUN go get ariefdfaltah-golang-graphql
RUN go install
# RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/ariefdfaltah-golang-graphql

############################
# STEP 2 running image
############################
FROM scratch
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/ariefdfaltah-golang-graphql /go/bin/ariefdfaltah-golang-graphql
# Use an unprivileged user.
USER appuser

# Run the hello binary.
ENTRYPOINT ["/go/bin/ariefdfaltah-golang-graphql"]

