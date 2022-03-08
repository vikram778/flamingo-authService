FROM golang:1.17

#ENV TZ=Europe/Moscow
#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download


# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o auth/cmd/auth-service/main .

##### new stage to copy the artifact #####

FROM alpine:3.11

RUN mkdir -p /flamingo

# Set the Current Working Directory inside the container
WORKDIR /flamingo

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

CMD ["./main"]