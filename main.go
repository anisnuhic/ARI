package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

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

	eventsChannel := client.LaunchListener()

	go func() {
		for event := range eventsChannel {
			fmt.Printf("Event received: %v\n", event)
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
			err := Dial(client, extensions)
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
			err := Join(client, channelID, extensions)
			if err != nil {
				fmt.Println("Error joining call:", err)
			}
		default:
			fmt.Println("Unknown command. Use: dial <extension1> <extension2> ..., list, join <channel_id> <extension1> <extension2> ...")
		}
	}
}

func Dial(client *ari.Client, extensions []string) error {

	createParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: "Conference_" + strings.Join(extensions, "_"),
	}
	bridge, err := client.Bridges.Create(createParams)
	if err != nil {
		return fmt.Errorf("neuspješno kreiranje mosta: %v", err)
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
			AppArgs:   "dial",
		}

		channel, err := client.Channels.Create(params)
		if err != nil {
			return fmt.Errorf("neuspješno kreiranje kanala za ekstenziju %s: %v", ext, err)
		}
		time.Sleep(2 * time.Second)
		err = bridge.AddChannel(channel.ID, ari.Participant)
		if err != nil {
			return fmt.Errorf("neuspješno dodavanje kanala %s u most %s: %v", channel.ID, bridge.ID, err)
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Printf("Konferencijski poziv kreiran za ekstenzije: %v\n", extensions)
	return nil
}

func List(client *ari.Client) error {

	channels, err := client.Channels.List()
	if err != nil {
		return fmt.Errorf("neuspješno prikazivanje kanala: %v", err)
	}

	fmt.Println("Lista trenutnih poziva:")
	for _, channel := range channels {
		fmt.Printf("Kanal ID: %s, Endpoint: %s, Status: %s\n", channel.ID, channel.Caller, channel.State)
	}

	return nil
}

func Join(client *ari.Client, channelID string, extensions []string) error {

	createParams := ari.CreateBridgeParams{
		Type: "mixed",
		Name: "Conference_" + channelID,
	}
	bridge, err := client.Bridges.Create(createParams)
	if err != nil {
		return fmt.Errorf("neuspješno kreiranje mosta: %v", err)
	}

	err = bridge.AddChannel(channelID, ari.Participant)
	if err != nil {
		return fmt.Errorf("neuspješno dodavanje početnog kanala %s u most %s: %v", channelID, bridge.ID, err)
	}

	for _, ext := range extensions {
		params := ari.OriginateParams{
			Endpoint: "PJSIP/" + ext,
			Context:  "sets",
			Priority: 1,
			CallerID: ext,
			Timeout:  30,
		}

		channel, err := client.Channels.Create(params)
		if err != nil {
			return fmt.Errorf("neuspješno kreiranje kanala za ekstenziju %s: %v", ext, err)
		}

		err = bridge.AddChannel(channel.ID, ari.Participant)
		if err != nil {
			return fmt.Errorf("neuspješno dodavanje kanala %s u most %s: %v", channel.ID, bridge.ID, err)
		}
	}

	fmt.Printf("Konferencija kreirana sa mostom %s i pridruženi kanali: %s\n", bridge.ID, append([]string{channelID}, extensions...))
	return nil
}

// funkcija za dvije extenzije
func DirectCall(client *ari.Client, ext1, ext2 string) error {
	params1 := ari.OriginateParams{
		Endpoint:  "PJSIP/" + ext1,
		Extension: ext1,
		Context:   "sets",
		Priority:  1,
		CallerID:  ext1,
		Timeout:   30,
		App:       "my_app",
	}

	params2 := ari.OriginateParams{
		Endpoint:  "PJSIP/" + ext2,
		Extension: ext2,
		Context:   "sets",
		Priority:  1,
		CallerID:  ext2,
		Timeout:   30,
		App:       "my_app",
	}

	channel1, err := client.Channels.Create(params1)
	if err != nil {
		return fmt.Errorf("neuspješno kreiranje kanala %s: %v", channel1, err)
	}

	channel2, err := client.Channels.Create(params2)
	if err != nil {
		return fmt.Errorf("neuspješno kreiranje kanala  %s: %v", channel2, err)
	}

	
	_, err = client.Channels.Create(ari.OriginateParams{
		Endpoint:  "PJSIP/" + ext2,
		Extension: ext1,
		Context:   "sets",
		Priority:  1,
		CallerID:  ext2,
		Timeout:   30,
		App:       "my_app",
	})

	if err != nil {
		return fmt.Errorf("neuspješno uspostavljanje poziva između ekstenzija %s i %s: %v", ext1, ext2, err)
	}

	fmt.Printf("Poziv između ekstenzija %s i %s je uspešno uspostavljen.\n", ext1, ext2)
	return nil
}


