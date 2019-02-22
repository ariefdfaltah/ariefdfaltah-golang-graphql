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
COPY . $GOPATH/src/ariefdfaltah_golang_graphql
WORKDIR $GOPATH/src/ariefdfaltah_golang_graphql

# Fetch dependencies.
# Using go mod with go 1.11
RUN go get ariefdfaltah_golang_graphql
RUN go install
# RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/ariefdfaltah_golang_graphql

############################
# STEP 2 run image
############################
FROM scratch
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/ariefdfaltah_golang_graphql /go/bin/ariefdfaltah_golang_graphql
# Use an unprivileged user.
USER appuser

# Run the hello binary.
ENTRYPOINT ["/go/bin/ariefdfaltah_golang_graphql"]

