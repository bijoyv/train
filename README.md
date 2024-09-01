# Train Reservation gRPC Service

## Overview

This project implements a gRPC-based train reservation system. It consists of a server that provides various services like purchasing tickets, retrieving tickets, getting seats by section, removing users, and modifying seat information. There is also a client application that interacts with the server using command-line options.

## Project Structure

```
├── cmd
│   ├── client
│   │   └── main.go       # Client application
│   └── server
│       └── main.go       # Server application
├── pkg
│   ├── proto
│   │   ├── train.pb.go   # Generated gRPC code
│   │   └── train_grpc.pb.go # Generated gRPC server and client interfaces
│   └── train
│       ├── reserv.go     # Logic implementation
│       └── reserv_test.go # Unit tests
├── proto
│   └── train.proto       # Protobuf definition file
└── README.md             # This file
```

## Prerequisites

- Go 1.22.3
- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go` plugin for generating Go code from `.proto` files

## Getting Started

### 1. Generate gRPC Code

If you make changes to the `.proto` files, regenerate the Go code by running:

```bash
protoc --go_out=pkg --go_opt=paths=source_relative --go-grpc_out=pkg --go-grpc_opt=paths=source_relative proto/train.proto
```

### 2. Running the Server

To start the server, run:

```bash
go run cmd/server/main.go
```

The server will start on `localhost:50051`.

### 3. Running the Client

The client application allows you to interact with the server using different commands. Below are the available commands and their options:

```bash
go run cmd/client/main.go --cmd=<command> [options]
```

### Command-Line Options

- **purchase**: Purchase a ticket.
  ```bash
  go run cmd/client/main.go --cmd=purchase --from=<origin> --to=<destination> --email=<user_email>
  ```
  Example:
  ```bash
  go run cmd/client/main.go --cmd=purchase --from=London --to=Paris --email=john.doe@example.com
  ```

- **getticket**: Retrieve a ticket using the user's email.
  ```bash
  go run cmd/client/main.go --cmd=getticket --email=<user_email>
  ```
  Example:
  ```bash
  go run cmd/client/main.go --cmd=getticket --email=john.doe@example.com
  ```

- **getseats**: Get available seats by section.
  ```bash
  go run cmd/client/main.go --cmd=getseats --section=<seat_section>
  ```
  Example:
  ```bash
  go run cmd/client/main.go --cmd=getseats --section=FirstClass
  ```

- **removeuser**: Remove a user by email.
  ```bash
  go run cmd/client/main.go --cmd=removeuser --email=<user_email>
  ```
  Example:
  ```bash
  go run cmd/client/main.go --cmd=removeuser --email=john.doe@example.com
  ```

- **modifyseat**: Modify a user's seat.
  ```bash
  go run cmd/client/main.go --cmd=modifyseat --email=<user_email> --newseat=<new_seat>
  ```
  Example:
  ```bash
  go run cmd/client/main.go --cmd=modifyseat --email=john.doe@example.com --newseat=2B
  ```

### 4. Error Handling

The client and server applications include basic error handling.

### 5. Running Tests

To run the unit tests for the service implementation, use:

```bash
go test ./pkg/train
```