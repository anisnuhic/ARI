package main

import (
    "fmt"
    "os"
    "time"
    "github.com/nvisibleinc/go-ari-library/client"
    "github.com/nvisibleinc/go-ari-library/ari"
)

var ariClient *client.Client

func main() {
    err := initARIClient()
    if err != nil {
        fmt.Println("Error initializing ARI client:", err)
        return
    }

    // Example of handling commands
    args := os.Args[1:]
    handleDial(args)
}

func initARIClient() error {
    var err error
    ariClient, err = client.New(&client.Options{
        Application: "your-app-name",
        Username:    "asterisk",
        Password:    "asterisk",
        URL:         "http://localhost:8088/ari",
    })
    return err
}

func handleDial(args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: dial <endpoint1> <endpoint2> [<endpoint3> ...]")
        return
    }

    callID := fmt.Sprintf("call-%d", time.Now().UnixNano())
    endpoints := args

    for _, endpoint := range endpoints {
        err := originateCall(endpoint, callID)
        if err != nil {
            fmt.Println("Error initiating call:", err)
            return
        }
    }
    fmt.Println("Call initiated with ID:", callID)
}

func originateCall(endpoint, callID string) error {
    // Example implementation of originateCall
    _, err := ariClient.Channels.Originate(&ari.OriginateRequest{
        Endpoint: endpoint,
        App:      "your-app-name",
        ChannelID: callID,
    })
    return err
}
