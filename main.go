package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/abourget/ari"
)

func main() {
	// Kreiramo klijenta
	username := "main"
	password := "pass"
	hostname := "localhost"
	port := 8088
	appName := "my_app"
	client := ari.NewClient(username, password, hostname, port, appName)

	// Pokrećemo listener za događaje
	eventsChannel := client.LaunchListener()

	// Pratimo događaje iz kanala u zasebnoj gorutini
	go func() {
		for event := range eventsChannel {
			fmt.Printf("Event received: %v\n", event)
		}
	}()

	fmt.Println("Enter commands: (dial <from> <to>, list, join <channel_id> <extension1> <extension2> ...):")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        // Čekamo korisnički unos
        scanner.Scan()
        input := scanner.Text()

        // Parsiramo unos
        parts := strings.Fields(input)
        if len(parts) == 0 {
            continue
        }

        switch parts[0] {
        case "dial":
            if len(parts) != 3 {
                fmt.Println("Incorrect input for dial. Use: dial <from> <to>")
                continue
            }
            from := parts[1]
            to := parts[2]
            err := Dial(client, from, to)
            if err != nil {
                fmt.Println("Error starting call:", err)
            }
        case "list":
            err := List(client)
            if err != nil {
                fmt.Println("Error listing calls: ", err)
            }
        case "join":
            if len(parts) < 3 {
                fmt.Println("Incorrect input for join. Use: join <channel_id> <extension1> <extension2> ...")
                continue
            }
            channelID := parts[1]
            extensions := parts[2:]
            err := Join(client, channelID, extensions)
            if err != nil {
                fmt.Println("Error joining call:", err)
            }
        default:
            fmt.Println("Unknown command. Use: dial <from> <to>, list, join <channel_id> <extension1> <extension2> ...")
        }
    }
}

// Dial pokreće novi poziv između dve ekstenzije
func Dial(client *ari.Client, from, to string) error {
	// Kreiramo parametre za Originate (pozivanje)
	paramsFrom := ari.OriginateParams{
		Endpoint:  "PJSIP/" + from, // Pretpostavljam da koristiš SIP, promijeni prema potrebi
		Extension: to,
		Context:   "sets", // Postavi svoj kontekst
		Priority:  1,
		CallerID:  to, // Postavi svoj CallerID
		Timeout:   30,
	}

	paramsTo := ari.OriginateParams{
		Endpoint:  "PJSIP/" + to,
		Extension: from,
		Context:   "sets",
		Priority:  1,
		CallerID:  from,
		Timeout:   30,
	}

	// Iniciramo novi kanal (poziv) između dve ekstenzije
	_, err := client.Channels.Create(paramsFrom)
	if err != nil {
		return fmt.Errorf("neuspješno pokretanje poziva sa ekstenzije %s ka %s: %v", from, to, err)
	}

	_, err = client.Channels.Create(paramsTo)
	if err != nil {
		return fmt.Errorf("neuspješno pokretanje poziva sa ekstenzije %s ka %s: %v", to, from, err)
	}

	fmt.Printf("Poziv pokrenut između ekstenzija: %s i %s\n", from, to)
	return nil
}

// List prikazuje trenutne pozive
func List(client *ari.Client) error {
	// Dohvatimo sve kanale
	channels, err := client.Channels.List()
	if err != nil {
		return fmt.Errorf("neuspješno listanje kanala: %v", err)
	}

	fmt.Println("Lista trenutnih poziva:")
	for _, channel := range channels {
		fmt.Printf("Kanal ID: %s, Endpoint: %s, Status: %s\n", channel.ID, channel.Caller, channel.State)
	}

	return nil
}

// Join kreira most (konferenciju) i pridružuje kanale
func Join(client *ari.Client, channelID string, extensions []string) error {
    // Kreiraj novi most (konferenciju)
    createParams := ari.CreateBridgeParams{
        Type: "mixed", // Tip mosta za konferenciju
        Name: "Conference_" + channelID, // Naziv mosta
    }
    bridge, err := client.Bridges.Create(createParams)
    if err != nil {
        return fmt.Errorf("neuspješno kreiranje mosta: %v", err)
    }

    // Dodaj početni kanal u most
    err = bridge.AddChannel(channelID, ari.Participant)
    if err != nil {
        return fmt.Errorf("neuspješno dodavanje početnog kanala %s u most %s: %v", channelID, bridge.ID, err)
    }

    // Dodaj sve ekstenzije kao kanale u most
    for _, ext := range extensions {
        params := ari.OriginateParams{
            Endpoint:  "PJSIP/" + ext,
            Context:   "sets",
            Priority:  1,
            CallerID:  ext,
            Timeout:   30,
        }
        // Kreiraj kanal za svaku ekstenziju
        channel, err := client.Channels.Create(params)
        if err != nil {
            return fmt.Errorf("neuspješno kreiranje kanala za ekstenziju %s: %v", ext, err)
        }

        // Dodaj kanal u most
        err = bridge.AddChannel(channel.ID, ari.Participant)
        if err != nil {
            return fmt.Errorf("neuspješno dodavanje kanala %s u most %s: %v", channel.ID, bridge.ID, err)
        }
    }

    fmt.Printf("Konferencija kreirana sa mostom %s i pridruženi kanali: %s\n", bridge.ID, append([]string{channelID}, extensions...))
    return nil
}
