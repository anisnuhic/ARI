package main

import (
	"bufio"
	//"flag"
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
			fmt.Printf("Primljen događaj: %v\n", event)
		}
	}()

	fmt.Println("Unesite komande (dial <from> <to>, list, join <bridge_id>):")
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
				fmt.Println("Neispravan unos za dial. Koristite: dial <from> <to>")
				continue
			}
			from := parts[1]
			to := parts[2]
			err := Dial(client, from, to)
			if err != nil {
				fmt.Println("Greška pri pokretanju poziva:", err)
			}
		case "list":
			err := List(client)
			if err != nil {
				fmt.Println("Greška pri listanju poziva:", err)
			}
		case "join":
			if len(parts) != 2 {
				fmt.Println("Neispravan unos za join. Koristite: join <bridge_id>")
				continue
			}
			bridgeID := parts[1]
			err := Join(client, bridgeID)
			if err != nil {
				fmt.Println("Greška pri pridruživanju pozivu:", err)
			}
		default:
			fmt.Println("Nepoznata komanda. Koristite: dial <from> <to>, list, join <bridge_id>")
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
// List prikazuje trenutne pozive
func List(client *ari.Client) error {
	// Dohvatimo sve kanale
	channels, err := client.Channels.List()
	if err != nil {
		return fmt.Errorf("neuspješno listanje kanala: %v", err)
	}

	fmt.Println("Lista trenutnih poziva:")
	for _, channel := range channels {
		fmt.Printf("Kanal ID: %s, Endpoint: %s, Status: %s\n", channel.ID, channel.Connected, channel.State)
	}

	return nil
}


// Join pridružuje korisnika mostu poziva
func Join(client *ari.Client, bridgeID string) error {
	// Ovdje bi koristili Bridges service za pridruživanje pozivu
	fmt.Println("Pridruživanje mostu:", bridgeID)
	return nil
}
