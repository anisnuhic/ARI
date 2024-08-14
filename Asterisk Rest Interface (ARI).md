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
			fmt.Printf("Primljen dogadjaj: %v\n", event)
		}
	}()
```

- `go func() {...}` : Launches a new goroutine to handle incoming events asynchronously.
- `for event := range eventsChannel` : Loops over events received on the `eventsChannel.` 
- `fmt.Printf("Primljen dogadjaj: %v\n", event)` : Prints each received event.

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
		if len(parts) != 2{
			fmt.Println("Incorrect input for join. Use: join <bridge_id>)
			continue
		}
		bridgeID := parts[1]
		err := Join(client, bridgeID)
		if err != nil {
			fmt.Println("Greska pri pridruzivanju pozivu: ", err)
		}
```

- `case "join": ` : Handles the `join` command.
	- `if len(parts) != 2` : Validates the input length.
	- `bridgeID := parts[1]` : Retrieves the bridge ID.
	- `err := Join(client, bridgeID)` : Calls the `Join` function to join a call bridge.
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
	func Dial(client *ari.Client, from, to string) error {
		paramsFrom := ari.OriginateParams {
			Endpoint: "PJSIP/" + from,
			Extension: to,
			Context: "sets"
			Priority: 1,
			CallerID: to,
			Timeout: 30,
		}
```

- `func Dial(...)` : Defines the `dial` function to start a call between two extensions.
- `paramsFrom` : Sets parameters for originating a call from `from` to `to`.

```
	paramsTo := ari.OriginateParams {
		Endpoint: "PJSIP/" + to,
		Extension: from,
		Context: "sets",
		Priority: 1, 
		CallerID: from,
		Timeout: 30, 
	}
```

- `paramsTo` : Sets parameters for originating a call from `to` to `from`.

```
	_, err := client.Channels.Create(paramsFrom)
	if err != nil {
		return fmt.Errorf("neuspjesno pokretanje poziva sa ekstenzije %s ka %s:          %v", from, to, err)
		}
	_, err = client.Channels.Create(paramsTo)
	if err != nil{
		return fmt.Errorf("neuspjesno pokretanje poziva sa ekstenzije %s ka %s:          %v", to, from, err)
		}
	fmt.Printf("poziv pokrenut izmedju ekstenzija: %s i %s\n", from, to)
	return nil
	}
```

- `client.Channels.Create(paramsFrom)`, `client.Channels.Create(paramsTo)` : Creates channels (calls) with the specified parameters.
- `return fmt.Errorf(...)` : Returns an error if the call fails.
- `fmt.Printf(...)` : Prints a message indicating the the call has been started.

***List function***
```
	func List(client *ari.Client) error {
		channels, err := client.Channels.List()
		if err != nil {
			return fmt.Errorf("neuspjesno listanje kanala: %v", err)
		}

		fmt.Println("Lista trenutnih poziva:")
		for_, channel := range channels {
			fmt.Printf("Kanal ID %s, Endpoint: %s, Status: %s\n", channel.ID,                channel.Connected, channel.State)
		}
		return nil
	}
```
- `func List (...)` : Defines the `List` function to display current calls.
- `channels, err := client.Channels.List()` : Retrieves  a list of channels (calls).
- `fmt.Errorf(...)` : Returns an error if listing fails.
- `fmt.Println(...)`, `fmt.Printf(...)` : Prints details of each channel.

***Join function***

```
	func Join (client *ari.Client, channelID string, extensions []string) error{
```

- `func Join` : This defines a new function named `Join`
- `client *ari.Client` : The function takes an `ari.Client` pointer as an argument. This is the ARI client used to interact with the Asterisk REST Interface.
- `channelID string` : This is a string argument representing the ID of the channel you want to add to the conference.
- `extensions []string` : This is a slice of strings representing the extensions you want to add to the conference.
- `error` : The function returns an error if anything goes wrong.

***Creating the Bridge (Conference)
```
	createParams := ari.CreateBridgeParams{
	Type: "mixed",
	Name: "Conference_"	+ ChannelID,
```

- `createParams := ari.CreateBridgeParams{}` : This creates an instance of `ari.CreateBridgeParams`, which is used to specify parameters for creating new bridge (conference). 
- `Type: "mixed"`: The bridge type is set to `"mixed"`, which means it's a conference where multiple participants can speak and hear each other.
- `Name: "Conference_" + channelID` : The name of the bridge is set to `"Conference_"` followed by the `channelID`. This gives the bridge a unique name based on the initial channel.

```
	bridge, err := client.Bridges.Create(createParams)
```

- `bridge, err := client.Bridges.Create(createParams)` : This line creates the bridge using the ARI client's `Bridge.Create` method with the parameters defined above. It returns the created bridge object and any error that occurs during creation.

```
	if err != nil {
		return fmt.Errorf("Neuspjesno kreiranje mosta: %v", err)
	}
```

- `if err != nil` : This checks if there was an error in creating the bridge.
- `return fmt.Errorf(...)`: If an error occurred, the function returns a formatted error message indicating the failure to create the bridge.

***Adding the Initial Channel to the Bridge***
```
	err = bridge.AddChannel(channelID, ari.Participant)
```

 - `err - bridge.AddChannel(channelID, ari.Participant)`: This line adds the initial channel (identified by `channelID`) to the bridge. The channel is added with the role `Participant` , which means it can speak and hear others.

```
	if err != nil{
		return fmt.Errorf("Neuspjesno dodavanje pocetnog kanala %s u most %s:            %v), channelID, bridge.ID, err)
	}
```

- `if err != nil` : This checks if there was an error in adding the initial channel to the bridge.
- `return fmt.Errorf(...)` : IF an error occurred, the function returns a formatted error message indicating the failure to add the channel to the bridge.

***Adding Extensions as Channels to the Bridge***
```
	for _, ext := range extensions {
```

- `for _, ext := range extensions`: This starts a loop that iterates over extension in the `extensions` slice.

```
	params := ari.OriginateParams{
		Endpoint: "PJSIP/" + ext ,
		Context: "sets",
		Priority: 1,
		CallerID: ext,
		Timeout: 30,
	}
```

- `params := ari.OriginateParams{}` : For each extension, an `OriginateParams` object is created. This defines the parameters for creating a new channel to connect extension to the conference.
- `Endpoint: "PJSIP/" + ext` : The `Endpoint` is set to `"PJSIP/` followed by the extension. This specifies the SIP endpoint to call.
- `Context: "sets"` : The `Context` in the Asterisk dial plan where the call should be routed.
- `Priority: 1` : The priority of the call in the dial plan.
- `CallerID: ext` : The caller ID for the new channel is set to the extension itself.
- `Timeout: 30` : The call will time out of it is not answered within 30 seconds.

```
	channel, err := client.Channels.Create(params)
	if err != nil {
		return fmt.Errorf("neuspjesno kreiranje kanala za extenziju %s: %v",ext,         err)
	}
```

- `channel, err := client.Channels.Create(params)` : This line creates a new channel for the extension using the `Create` method of the `Channels` service.
- `if err != nil` : This checks if there was an error in creating the channel.
- `return fmt.Errorf(...)` : If an error occurred, the function returns a formatted error message indicating the failure to create the channel for the extension.

```
	err = bridge.AddChannel(channel.ID, ari.Participant)
	if err != nil {
		return fmt.Errorf("neuspjesno dodavanje kanala %s u most %s: %v",                channel.ID, bridge.ID, err)
	}
```

- `err = bridge.AddChannel(channel.ID, ari.Participant)` : This line adds the newly created channel to the bridge as a participant.
- `if err != nil` : This checks if there was an error in adding the channel to the bridge.
- `return fmt.Errorf(...)` : If an error occurred, the function returns a formatted error message indicating the failure to add the channel to the bridge.

***Final Output***
```
	fmt.Printf("Konferencija kreirana sa mostom %s i pridruzeni kanali: %s\n",       bridge.ID, append([]string{channelID}, extensions...))
```

- `fmt.Printf(...)` : This prints a message indicating that the conference (bridge) was successfully created and lists all the channels (including the initial channel and the extensions) that were added it.

```
	return nil
```

- `return nil` : If everything was successful, the function returns `nil` indicating no error occurred.









