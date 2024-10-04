package main

import (
	"bufio"
	"fmt"
	"github.com/abourget/ari"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

func main() {
	// Kreiramo klijenta
	username := "main"
	password := "pass"
	hostname := "localhost"
	port := 8088
	appName := "my_app"
	client := ari.NewClient(username, password, hostname, port, appName)
	eventsChannel := client.LaunchListener()

	mapa := make(map[string]int)

	go func() {
		for event := range eventsChannel {
			fmt.Printf("Event received: %v\n", event.GetType())
			if event.GetType() == "ChannelDestroyed" {
				bridges, _ := client.Bridges.List()
				for _, bridge := range bridges {
					if len(bridge.Channels) == 1 && mapa[bridge.ID] == 0 {
						c1, _ := client.Channels.Get(bridge.Channels[0])
						c1.Hangup()
						bridge.Destroy()
					}
					if len(bridge.Channels) == 0 {
						bridge.Destroy()
					}
				}

			}
		}
	}()

	fmt.Println("Enter commands: (dial <extension1> <extension2> ... for conference, list, join <channel_id> <extension1> <extension2> ...):")
	scanner := bufio.NewScanner(os.Stdin)

	for {

		scanner.Scan()
		input := scanner.Text()

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "dial":
			if len(parts) < 3 {
				fmt.Println("Incorrect input for dial. Use: dial <extension1> <extension2> ...")
				continue
			}
			extensions := parts[1:]
			err := Dial(client, extensions, mapa)
			if err != nil {
				fmt.Println("Error starting conference call:", err)
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
			err := Join(client, channelID, extensions, mapa)
			if err != nil {
				fmt.Println("Error joining call:", err)
			}
		case "exit":
			return
		default:
			fmt.Println("Unknown command. Use: dial <extension1> <extension2> ..., list, join <channel_id> <extension1> <extension2> ...")
		}
	}
}

// dial function for more than 2 extensions
func Dial(client *ari.Client, extensions []string, mapa map[string]int) error {
	name := "Conference_"
	flag := 1
	if len(extensions) == 2 {
		name = "Call_"
		flag = 0
	}
	// Kreiramo novi bridge
	createParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: name,
	}
	bridge, err := client.Bridges.Create(createParams)
	if err != nil {
		return fmt.Errorf("failed to create bridge: %v ", err)
	}

	// Map za čuvanje kanala i njihovih stanja
	channels := make(map[string]*ari.Channel)

	// Gorutine za pravljenje kanala i praćenje njihovih statusa
	var wg sync.WaitGroup

	for _, ext := range extensions {
		wg.Add(1)
		go func(ext string) {
			defer wg.Done()

			params := ari.OriginateParams{
				Endpoint:  "PJSIP/" + ext,
				Extension: ext,
				Context:   "sets",
				Priority:  1,
				CallerID:  ext,
				Timeout:   30,
				App:       "my_app",
				AppArgs:   "dial",
			}

			// Kreiramo kanal
			channel, err := client.Channels.Create(params)
			if err != nil {
				fmt.Printf("failed to create chennel for extension %s: %v\n", ext, err)
				return
			}

			// Čuvamo kanal u mapu
			channels[channel.ID] = channel

			// Pratimo status kanala
			for {
				channel, err = client.Channels.Get(channel.ID)
				if err != nil {
					fmt.Printf("error getting channel %s: %s\n", channel.ID, err)
					return
				}

				if channel.State == "Up" {
					break
				}

				time.Sleep(time.Millisecond * 100)
			}

			// Dodajemo kanal u most
			err = bridge.AddChannel(channel.ID, ari.Participant)
			if err != nil {
				fmt.Printf("failed to add channel %s to bridge %s: %v\n", channel.ID, bridge.ID, err)
				return
			}
		}(ext)
	}

	// Čekamo da sve gorutine završe
	wg.Wait()

	fmt.Printf("Conference call created for extensions: %v\n", extensions)
	mapa[bridge.ID] = flag
	return nil
}

// List function
func List(client *ari.Client) error {
	bridges, err := client.Bridges.List()
	if err != nil {
		return fmt.Errorf("failed to list bridges: %v", err)
	}
	fmt.Println("list of current bridges:")
	for _, bridge := range bridges {
		if len(bridge.Channels) == 0 {
			bridge.Destroy()
		} else {
			fmt.Printf("bridges: %s\n", bridge.ID)
		}
	}
	return nil
}

// Join function
func Join(client *ari.Client, channelID string, extensions []string, mapa map[string]int) error {
	// Pokušaj da dobiješ most po ID-u
	if strings.IndexFunc(channelID, unicode.IsLetter) >= 0 {
		bridge, err := client.Bridges.Get(channelID)
		if err != nil {
			return fmt.Errorf("failed to get a bridge: %v", err)
		}

		for _, ext := range extensions {
			params := ari.OriginateParams{
				Endpoint:  "PJSIP/" + ext,
				Extension: ext,
				Context:   "sets",
				Priority:  1,
				CallerID:  ext,
				Timeout:   30,
				App:       "my_app",
			}

			channel, err := client.Channels.Create(params)
			if err != nil {
				return fmt.Errorf("failed to create channel for extension %s: %v", ext, err)
			}
			if mapa[bridge.ID] == 0 {
				mapa[bridge.ID] = 1
			}
			err = bridge.AddChannel(channel.ID, ari.Participant)
			if err != nil {
				return fmt.Errorf("failed to add channel %s to bridge %s: %v", channel.ID, bridge.ID, err)
			}

		}

		fmt.Printf("Extesnion %s was successfully added to bridge %s\n", extensions, bridge.ID)
	}
	return nil
}
