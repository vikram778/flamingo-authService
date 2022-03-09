FROM golang:1.17 as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download


# Build the Go app
RUN go build ./otpService/cmd/otp_service/main.go

##### new stage to copy the artifact #####

FROM alpine:3.11

RUN mkdir -p /flamingo

# Set the Current Working Directory inside the container
WORKDIR /flamingo

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/config .

CMD ["./main"]