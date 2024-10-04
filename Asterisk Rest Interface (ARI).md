Let's break down the code line by line

***Importing packages*** 
```
	import (
		"bufio"
		"fmt"
		"os"
		"strings"
		"github.com/abourget/ari"
		"sync"
		"time"
		"unicode"
	)
```

- `bufio` : provides buffered I/O for more efficient reading and writing.
- `fmt` : Provides formatted I/O functions (like `Printf`).
- `os` : Provides functions for interacting with the operating system (like reading from standard input).
- `strings` : Provides functions for strings manipulation.
- `github.com/abourget/ari` : A Go package for interacting with Asterisk Rest Interface (ARI).
- `sync` : Package for synchronizing goroutines
- `time` : Package for time-related functions
- `unicode` : Package for working with unicode characters

***Main function***
```
	func main (){
		username := "main"
		password := "pass"
		hostname := "localhost"
		port := 8088
		appName := "my_app"
		client := ari.NewClient(username, password, hostname, port, appName)
		mapa := make(map[string]int)
	}
```

- `username`, `password`, `hostname`, `port`, `appName` : Variables for ARI client connection parameters.
- `client := ari.NewClient(...)` : Creates a new ARI client instance with the provided credentials  and connection details.
- `mapa := make(map[string]int)` : This is a dictionary used to store bridge states (whether it's a conference call or regular call).
```
	eventsChannel := client.LaunchListener()
```
 - `eventsChannel := client.LaunchListener()` : Starts listening for events from the ARI server an returns a channel through which events will be sent.

```
	go func() {
		for event := range eventsChannel {
			fmt.Printf("Event received: %v\n", event)
		}
	}()
```

- `go func() {...}` : Launches a new goroutine to handle incoming events asynchronously.
- `for event := range eventsChannel` : Loops over events received on the `eventsChannel.` 
- `fmt.Printf("Event received: %v\n", event)` : Prints each received event.
- `if eveny.GetType() == "ChannelDestroyed" {...}` : It checks for the `ChannelDestroyed` event, signaling the end of a call. If the bridge is nearly empty (with 0 or 1 channel in the regular call or 0 channel in conference call), it destroys the bridge.

```
	fmt.Println("Enter commands: (dial <from> <to>, list, join <bridge_id> ):")
	scanner := bufio.NewScanner(os.Stdin)
	
```

- `fmt.Println(...)` : Prints instructions for user commands.
- `scanner := bufio.NewScanner(os.Stdin)` : Creates a new Scanner to read input from the standard input (keyboard).

```
	for {
		scanner.Scan()
		input := scanner.Text()
	}
```

- `for {...}` : Starts an infinite loop to continuously read user input.
- `scanner.Scan()`: Reads the next line of input from the user.
- `input := scanner.Text()` : Retrieves the text of the input line.

```
	parts := strings.Fields(input)
	if len(parts) == 0 {
		continue
	}
```

- `parts := strings.Fields(input)` : Splits the input text into fields (words) based on whitespace.
- `if len(parts) == 0 {continue}` : If no input is provided, continues to the next iteration of the loop.

``` 
	switch parts[0]{
		case "dial":
			if len (parts) != 3{
				fmt.Println("Incorrect input for dial. Use: dial <from> <to>")
				continue
			}
			from := parts[1]
			to := parts[2]
			err := Dial(client, from, to)
			if (err != nil){
				fmt.Println("Error starting call:", err)
			}
	}
```
- `switch parts[0]` : Checks the first part of the input to determine the command.
- `case "dial":` : Handles the `dial` command.
	 - `if len(parts) != 3` : Validates the input length.
	 - `from := parts[1]`, `to := parts[2]` : Retrieves the source and destination extensions.
	 - `err := Dial(client, from, to)` : Calls the `Dial` function to start a call.
	 - `if err != nil {...}` : Prints an error message if the call fails.

```
	case "list":
		err := List(client)
		if err != nil{
			fmt .Println("Error listing calls: ", err)
		}
```

 - `case "List": `: Handles the `list` command.
	- `err := List(client)` : Calls the `List` function to list current calls.
	- `if err != nil {...}` : Prints an error message if listing fails.

```
	case "join" :
		if len(parts) < 3{
			fmt.Println("Incorrect input for join. Use: join <bridge_id>)
			continue
		}
		channelID := parts[1]
		extensions := parts[2:]
		err := Join(client, channelID, extenisions)
		if err != nil {
			fmt.Println("Error joining call: ", err)
		}
```

- `case "join": ` : Handles the `join` command.
	- `if len(parts) < 3` : Ensures the `join` command has at least two arguments(`channel_id` and one or more extensions).
	- `channelID := parts[1]` : Extracts the `channel_id` from the input.
	- `extensions := parts[2:]` : Extracts the list of extensions to join into conference.
	- `err := Join(client, channelID, extensions)` : Calls the `Join` function to join a call bridge, creating a conference and add the specified extensions.
	- `if err != nil {...}` : Prints an error message if joining fails.

```
	default:
		fmt.Println("Unknown command. Use dial <from> <to>,list,join <bridge_id")
		}
	}
	}
```

- `default:` : Handles unknown commands.
- `fmt.Println(...)` : Prints a message for unknown commands.

***Dial Function***
```
	func Dial(client *ari.Client, extensions[] string, mapa map[string]int)           error {...}
```

- `func Dial(...)` : Defines the `dial` function to start a call between more than two extensions.

```
	name := "Conference_"
	flag := 1
```

- Creates two variables which is used to set bridge parameters 

```
	if len(extensions) == 2 {
		name = "Call_"
		flag = 0 
	}
```

- This checks if there are only two extensions in the list. If so, change the values of variables

```
	createParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: name + strings.Join(extensions, "_"),
	}
	bridge, err := client.Bridges.Create(createParams)
```

- If there are more than two extensions, it creates a **new conference call** (bridge). A `bridge` in ARI is a mechanism that allows multiple channels to communicate (mixing).
-  `createParams` is a struct that defines the bridge creation parameters (type: `mixing` and a name fir the bridge).

```
	channels := make(map[string]*ari.Channel)
```

- This creates a map to store the channels created for each extension. The key is the channel ID, and the value is the channel itself.

```
	var wg sync.WaitGroup
```

- The `sync.WaitGroup` is used to wait for all goroutines to finish. Goroutines are lightweight threads that allow us to run tasks concurrently.

```
	go func(ext string){
		deger wg.Done()
	}(ext)
```

- For each extension, a goroutine is created that runs the logic for originating a call (creating a channel for the extension).

```
	params := ari.OriginateParams{ 
		Endpoint: "PJSIP/" + ext,
		Extension: "s",
		Context: "sets",
		Priority: 1,
		CallerID: ext,
		Timeout: 30, 
		App: "my_app", 
		AppArgs: "dial", 
	}
	channel, err := client.Channels.Create(params)
```

- A channel is created using the `OriginateParams`, specifying the Endpoint (e.g., "PJSIP/1001"), context, priority, and CallerID for the call.
- After originating the call, the program continuously checks the state of the channel until it goes Up (i.e., the call is connected).

```
	err = bridge.AddChannel(channel.ID, ari.Participant)
```

- Once the channel is ready (state is "Up"), the channel is added to the bridge (conference call).

```
	wg.Wait()
	mapa[bridge.ID] = flag
```

- This ensures that the function waits for all extensions (channels) to be connected before proceeding.
- `mapa[bridge.ID] = flag` : Sets key and value of map

***List function***
```
	func List(client *ari.Client) error {
		bridges, err := client.Bridges.List()
		if err != nil {
			return fmt.Errorf("failed to list bridges: %v", err)
		}

		fmt.Println("list of current bridges:")
		for _, bridge := range bridges {
			fmt.Printf("Bridge ID %s", bridge.ID )
		}
		return nil
	}
```
- `func List (...)` : Defines the `List` function to display current calls.
- `channels, err := client.Channels.List()` : Retrieves  a list of all current bridges.
- `fmt.Errorf(...)` : Returns an error if listing fails.
- `for _, bridge := range bridges {...}` : It prints out the ID of each bridge found.

***Join function***

```
	func Join (client *ari.Client, channelID string, extensions []string, mapa        map[string]int) error{
```

- This function allows existing channels to join an existing bridge.

```
	if strings.IndexFunc(channelID, unicode.IsLetter) >= 0 {...}
```

- This checks if one letter exists in the `channelId` string

```
	bridge, err := client.Bridge.Get(channelID)
```

- It retrieves the bridge using the provided `channelID`.

```
	params := ari.OriginateParams{ 
		Endpoint: "PJSIP/" + ext, 
		Extension: ext, 
		Context: "sets", 
		Priority: 1, 
		CallerID: ext, 
		Timeout: 30, 
		App: "my_app", 
	} 
	channel, err := client.Channels.Create(params)
```

- For each extension, it creates a new channel using `OriginateParams`, similar to the `Dial` function.

```
	if mapa[bridge.ID] == 0 {
		mapa.[bridge.ID] = 1
	}
```

- If someone joins the regular call set the flag to 1 and that indicates that becomes conference call

```
	err = bridge.AddChannel(channel.ID, ari.Participant)
```

- The created channel is add to the existing bridge, turning the current call into a conference









