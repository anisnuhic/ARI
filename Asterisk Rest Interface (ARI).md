Let's break down the code line by line

***Importing packages*** 
```
	import (
		"bufio"
		"fmt"
		"os"
		"strings"
		"github.com/abourget/ari"
	)
```

- `bufio` : provides buffered I/O for more efficient reading and writing.
- `fmt` : Provides formatted I/O functions (like `Printf`).
- `os` : Provides functions for interacting with the operating system (like reading from standard input).
- `strings` : Provides functions for strings manipulation.
- `github.com/abourget/ari` : A Go package for interacting with Asterisk Rest Interface (ARI).

***Main function***
```
	func main (){
		username := "main"
		password := "pass"
		hostname := "localhost"
		port := 8088
		appName := "my_app"
		client := ari.NewClient(username, password, hostname, port, appName)
	}
```

- `username`, `password`, `hostname`, `port`, `appName` : Variables for ARI client connection parameters.
- `client := ari.NewClient(...)` : Creates a new ARI client instance with the provided credentials  and connection details.

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
	func Dial(client *ari.Client, extensions[] string) error {...}
```

- `func Dial(...)` : Defines the `dial` function to start a call between more than two extensions.

```
	if len(extensions) == 2 {
		return DirectCall(client, extensions[0], extensions[1])
	}
```

- This checks if there are only two extensions in the list. If so, it calls `DirectCall` function to create a direct call between the two extensions.

```
	createParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: "Conference_" + strings.Join(extensions, "_"),
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
```

- This ensures that the function waits for all extensions (channels) to be connected before proceeding.

***DirectCall Function***
```
	func DirectCall(client * ari.Client, ext1, ext2 string) error {...}
```

- This function handles a simple direct call between two extensions.

```
	params1 := ari.OriginateParams{ 
		Endpoint: "PJSIP/" + ext1, 
		Extension: ext1, 
		Context: "sets", 
		Priority: 1, 
		CallerID: ext1, 
		Timeout: 30, 
		App: "my_app", 
	}
```

- A channel is created for each extension using ari.OriginateParams, which specifies the endpoint, context, priority, etc.

```
	channel1, err := client.Channels.Create(params1)
	if err != nil{
		return fmt.Errorf("failed to create channel %s: %v", channel1, err)
	}
```

 - If there is an error in creating the channel, it prints an error message.

```
	fmt.Printf("Call between extensions %s and %s was successfully established.      \n", ext1, ext2")
```

- Once both channels are created, the call is established between the two extensions.
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
	func Join (client *ari.Client, channelID string, extensions []string) error{
```

- This function allows existing channels to join an existing bridge.

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
	err = bridge.AddChannel(channel.ID, ari.Participant)
```

- The created channel is add to the existing bridge, turning the current call into a conference









