package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/thetangentline/grpc/proto"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// -------------------------------------------------------------------->

var (
	wait sync.WaitGroup
	rl   *readline.Instance
)

// -------------------------------------------------------------------->

func GetName() string {
	fmt.Print("Enter your name: ")
	var name string
	fmt.Scanln(&name)
	return strings.TrimSpace(name)
}

func Connect(user *proto.User, broadcastClient proto.BroadcastClient) error {
	var streamerror error

	// Create connection with server
	stream, err := broadcastClient.CreateConnection(
		context.Background(),
		&proto.Connect{
			User:   user,
			Active: true,
		})

	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)

	// Create a goroutine to receive the message
	go func(str proto.Broadcast_CreateConnectionClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("error reading message: %v", err)
				break
			}

			if msg.OrgName != user.Name {
				fmt.Fprintf(rl.Stdout(), "%v: %s\n", msg.OrgName, msg.Content)
			}
		}
	}(stream)

	return streamerror
}

func SendMessage(user *proto.User, broadcastClient proto.BroadcastClient) error {
	defer wait.Done()

	for {
		text, err := rl.Readline()
		if err != nil {
			break
		}

		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		msg := &proto.Message{
			Id:      user.Id,
			OrgName: user.Name,
			Content: text,
		}

		// Send message to server and let it broadcast the message
		_, err = broadcastClient.BroadcastMessage(context.Background(), msg)
		if err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
	}

	return nil
}

// -------------------------------------------------------------------->

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("couldn't connect to service: %v", err)
	}
	broadcastClient := proto.NewBroadcastClient(conn)

	name := GetName()
	user := &proto.User{
		Id:   uuid.New().String(),
		Name: name,
	}

	rl, err = readline.New(fmt.Sprintf("%s: ", name))
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	err = Connect(user, broadcastClient)
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
		return
	}

	wait.Add(1)
	go SendMessage(user, broadcastClient)

	wait.Wait()
}
