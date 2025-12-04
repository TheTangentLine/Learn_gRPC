# gRPC Messenger CLI

This project aims to learn gRPC through building a messenger CLI application. It demonstrates server streaming, message broadcasting, and real-time communication between multiple clients using gRPC in Go.

### Architecture

The application consists of a gRPC server that manages client connections and broadcasts messages, and CLI clients that connect to the server to send and receive messages in real-time.

### Server Side

The server starts by creating a gRPC server instance and listening on port 8080. It registers the Broadcast service which handles two RPC methods: CreateConnection for establishing client connections and BroadcastMessage for sending messages to all connected clients. The server maintains a list of active connections, each with a server streaming channel that allows pushing messages to clients.

### Client Side

The client connects to the server at localhost:8080 using insecure credentials. Upon startup, it prompts for a user name and creates a User object with a unique ID. The client then calls CreateConnection to establish a server streaming connection, which spawns a goroutine to continuously receive and display incoming messages. A separate goroutine handles reading user input from stdin and sending messages via BroadcastMessage.

### Message Flow

When a client sends a message through BroadcastMessage, the server receives the message and broadcasts it to all active connections. Each connection uses a goroutine to send the message through its server streaming channel. The server waits for all sends to complete before returning. Clients receive messages through their streaming connection and display them, filtering out their own messages to avoid duplicates.
