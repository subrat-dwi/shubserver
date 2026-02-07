# Stage 1: Build

FROM golang:1.25.4 AS builder

# sets working directory
WORKDIR /app

# copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# copy rest of the code
COPY . .

# build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: Run

# starts new stage using lightweight Alpine linux
FROM alpine:3

# sets working directory inside runtime container
WORKDIR /app

# install certificates
RUN apk add --no-cache ca-certificates

# copy the final binary
COPY --from=builder /app/server .

EXPOSE 8080

# defines the default command to run when container starts
CMD [ "./server" ]