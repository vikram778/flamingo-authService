# flamingo-authService

1. This Repo has both the services auth-service in ./auth directory and otp service in ./otpService directory
2. Auth service is grpc service for profile creating and verification via otp.
3. Otp service is consumer service which consumes otp message and sends otp to required number using twilio api.

#Service communication Flow

1. Upon user signup auth service receives signup request with user data.
2. Auth service creates profile in database and generates otp.
3. Auth service publishes generated otp to verification queue of rabbitmq along with the mobile no
4. Otp service consumes message from verification queue and sends otp to user's mobile using twilio's api.
5. User enters otp and sends verification request to auth service and auth service validates otp.

#Database Schema

Database schema is put under auth/migrations.

#Application startup and set up

1. To run and build auth service :
   go run auth/cmd/auth_service/main.go
   go build ./auth/cmd/auth_service/main.go

2. To run otp service :
   go run otpService/cmd/otp_service/main.go
   go build ./otpService/cmd/otp_service/main.go

#Prerequisites

1. PostgresSQl
2. RabbitMq

Above dependencies are added in docker-compose file .

#Starting dependencies in local

docker-compose -f docker-compose.yml up --build

#Configs

config .yml files are in ./config/ directory

1. Would need to update twilio authid and token in .yml file