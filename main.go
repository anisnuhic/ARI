package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/client/native"
)
var ariClient *ari.Client
func main() {
	client, err := native.Connect(&native.Options{
		Application: "main",
		Username:    "asterisk",
		Password:    "pass",
		URL:         "http://localhost:8088/ari",
	})
	if err != nil {
		log.Fatalf("Error connecting to server ARI: %v", err)
	}

	fmt.Println("Connected to ARI server")
}
func handleDial(args []string){
    if len(args) < 2 {
        fmt.Println("Usage: dial <endpoint1> <endpoint2> <endpoint3> ...")
        return
    }

    callId := fmt.Sprintf("call-%d", time.Now().UnixNano())
    endpoints := args

    for i :=0 ; i<len(endpoints); i++ {
        source := endpoints[i]
        destination := endpoints[i+i]
        err := originateCall(source, destination, callId)
        if err != nil {
            fmt.Println("Error initiating call:", err)
            return
        }
    }
    fmt.Println("Call initiated with ID: ", callId)
}

func originateCall(source, destination, callId string) error {
    _, err := ariClient.Channel().Create(context.Background(), &ari.ChannelCreateRequest{
        Endpoint: source,
        App: "main",
        ChannelID : callId,
    })
    if err != nil {
        return fmt.Errorf("Error creating channel for call %s: %v", source, err)
        
    }
    _, err := ariClient.Channel().Originate(context.Background(), &ari.OriginateRequest{
        Endpoint: destination,
        App: "main",
        ChannelID : callId,
    }) 
    if err != nil {
        return fmt.Errorf("Error originating call to destination %s %v", destination,err)
    }
    return nil
}
func handleList(){
    //implement this function
}
func handleJoin(args []string){
    //implement this function
}

func chooseOption(input string){
    parts := strings.Split(input, " ")
    command := parts[0]
    args :=parts[1:]

    switch command {
    case "dial":
        handleDial(args);
    case "list":
        handleList();
    case "join":
        handleJoin(args);
    default:
        fmt.Println("Invalid command. Use 'dial', 'list', or 'join' followed by phone number.");
    }
}