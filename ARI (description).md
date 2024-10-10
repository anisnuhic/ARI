```
import (
	"bufio"
	"fmt"
	"github.com/abourget/ari"
	"os"
	"strings"
	"sync"
	"time"
)
```

- `bufio` : Ovaj paket se koristi za citanje unosa iz konzole (stdin) efikasno.
- `fmt` : Ovaj paket se koristi za formatirano ispisivanje u konzolu.
- `github.../ari`: Ovaj paket se koristi za interakciju s Asterisk serverom putem REST APi-ja.
- `os` : Ovaj paket se koristi za interakciju sa operativnim sistemom
- `strings` : Ovaj paket koristi se za mainpulaciju stringovima, poput razdvajanja unosa na dijelove.
- `sync` : Ovaj paket se koristi za rad sa "wait group"-om (grupama cekanja) kako bi goroutine mogle sinhronizirano cekati.
- `time` : Ovaj paket se koristi za kontrolu cekanja (npr. cekanje izmedju provjera statusa kanala).

```
	func main(){
		username := "main"
		password := "pass"
		hostname := "localhost"
		port := 8088
		appName := "my_app"
		client := ari.NewClient(username, password, hostname, port, appName)
		eventsChannel := client.LaunchListener()
	}
```

- Ovdje kreiramo ARI klijenta koji se povezuje na Asterisk server. To radimo pomocu funkcije ari.NewClient() sa navedenim korisnickim imenom, lozinkom, hostname-om, portom i imenom aplikacije.
- Nakon kreiranja klijenta, pozivamo funkciju `LaunchListener()`, koja pokrece slusanje dogadjaja s Asterisk servera. To nam vraca kanal (`eventsChannel`), gdje ce se svaki novi dogadjaj proslijediti.

```
	go func(){
		for event := range eventsChannel {
			fmt.Printf("Event received: %v\n", event)
		}
	}()
```

- Gorutina je funkcija koja se izvrsava u pozadini, paralelno s ostatom programa, izvrsava se na drugom threadu procesora.
- U ovoj gorutini, dogadjaji primljeni putem `eventsChannel` se kontinuirano ispisuju u konzolu. Svaki novi dogadjaj iz Asterisk servera bice prikazan u obliku tekstualnog zapisa.

```
	go func(){
		bridges, err := client.Bridges.List()
		if err != nil {
			return 
		}
	for _, bridge := range bridges {
		if len(bridge.Channels) == 0{
			bridge.Destroy()
		}
	}
	}()
```

- Ova gorutina dohvata listu trenutnih bridgova sa servera pomocu `client.Bridges.List()`.
- Ako funkcija `List()` vrati gresku (`err`), izlazi iz funkcije.
- Petlja prolazi kroz sve mostove i provjerava da li bridge ima  kanale (`bridge.Channels`). Ako je lista prazna (`len(bridge.Channels) == 0` ) most se unistava pomocu `bridge.Destroy().` 

```
	fmt.Println("Enter commands: (dial <extension1> <extension2> ...))
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		input := scanner.Text()

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0]{
			...
		}
	}
```

- Ovdje se kreira `scanner` za unos komandi iz konzole.
- Koristi se beskonacna `for` petlja koja ceka unos korisnika.
- `scanner.Scan()` cita unos korisnika, dok `input := scanner.Text()` uzima tekstualni unos
- `strings.Fields(input)` razdvaja unijeti tekst u dijelove (prema praznim mjestima).
- Prvi dio (`parts[0]`) se koristi za prepoznavanje komande, a onda se kroz `switch` odredjuje koja komanda se izvrsava (npr. `dial`, `list`, `join`).

```
	func Dial(client *ari.Client, extensions [] string) error {
		if len(extensions) == 2 {
			return DirectCall (client , extensions[0], extensions[1])
		}

		createParams := ari.CreateBridgeParams{
			Type: "mixing",
			Name: "Conference_ " + strings.Join(extensions, "_"),
		}
		bridge, err := client.Bridges.Create(CreateParams)
		...
		}
```

- Ako su zadane dvije ekstenzije, poziva se `DirectCall()` za uspostavljanje direktnog poziva.
- Ako je vise od dvije ekstenzije, kreira se novi bridge za konferenciju. Parametri za kreiranje bridga se prosljedjuju preko `ari.CreateBridgeParams`. Ime mosta je spoj ekstenzija.

```
	 var wg sync.WaitGroup
	 for _, ext := range extensions {
		 wg.Add(1)
		 go func(ext string) {
			 defer wg.Done()
		 }
		params := ari.OriginateParams {
			Endpoint: "PJSIP/" + ext,
			Extension: ext,
			Context: "sets",
			...
		}
		channel, err := client.Channels.Create(params)
		...
	 }(ext)
	}
	wg.Wait()
```

- Koristi se `sync.WaitGroup` kako bi se sinhroniziralo izvrsavanje gorutina koje uspostavljaju kanale za svaku ekstenziju.
- Za svaku ekstenziju, kreira se gorutina koja pokrece kanal koristeci `client.Channels.Create()`.
- Gorutine cekaju dok svi kanali ne budu kreiranji, koristeci `wg.Wait().`

```
	err = bridge.AddChannel(channel.ID, ari.Participant)
	if err != nil {
		fmt.Printf("failed to add channel %s to bridge %s: %v\n", channel.ID)
	}
```

- Nakon sto kanal postane aktivan , kanal se dodaje u kreirani most pomocu `bridge.AddChannel()`.

```
	func DirectCall(client * ari.Client, ext1, ext2 string) error {
		params1 := ari.OriginateParams {
			Endpoint: "PJSIP/" + ext1, 
			Extension: ext1,
			...
		}
		channel1, err := client.Channels.Create(params1)
		...
	}
```

 - Ako se radi o dvije ekstenzije, kreira se kanal za svaku ekstenziju i uspostavlja se direktan poziv izmedju njih. Koriste se parametri za kreiranje kanala (`ari.OriginateParams`). 

```
	func List(client *ari.Client) error {
		bridges, err := client.Bridges.List()
		if err != nil {
			return fmt.Errorf("failed to list bridges: %v", err)
		}

		fmt.Println("list of current bridges: ")
		for _, bridge := range bridges {
			if len(bridge.Channels) == 0 {
				bridge.Destroy()
			} else {
				fmt.Printf("Bridge ID: %s", bridge.ID)
			}
			}
			return nil
	}
```

- Ova funkcija lista sve trenutne mostove. Ako most nema kanala, on se unistava. Inace, ispisuje se njegov ID.

```
	func Join(client *ari.Client, channelID string, extensions []string) error {
		bridge, err := client.Bridges.Get(channelID)
		...
		for _, ext : range extensions {
			params := ari.OriginateParams{
				Endpoint: "PJSIP/" + ext,
				...
			}
			channel, err := client.Channels.Create(params)
			...
			bridge.AddChannel(channel.ID, ari.Participant)
		}
		...
	}
```

 - Funkcija `Join` omogucava korisniku da dodaje nove ekstenzije u vec postojeci most. Kreira nove kanale za svaku ekstenziju i dodaje ih u most.


UKRATKO:
	Ovaj kod omogucava interakciju s Asterisk serverom putem konzolnih komandi za kreiranje konferencijskih poziva, listanje i unistavanje mostova, te dodavanje novih ekstenzija u bridge-ove.
