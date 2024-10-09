***ASTERISK REST INTERFACE(ARI)***
___________________________________

***WHAT TO DO ?***


Create an ARI application<br> The application should be written in Golang.<br> Create an ARI CLI (terminal) application that will wait for the user’s input and based on that input, it will execute some actions. Once an action is registered, the terminal should be free for entering a new command. The commands that need to be implemented are “dial”, “list”, and “join”.<br> Dial initiates a call between X endpoints:<br> Examples:<br> “dial 100 101” will initiate a call between these two endpoints. Once one of these endpoints hangs up, the call ends.<br> ”dial 100 101 102 103” initiates a call between all 4 endpoints, creating a conference. A conference call ends when all users leave the conference. It is possible to have only one user inside a conference.<br> List prints all ongoing calls<br> The information shown should be the call ID (can be anything) and the participants (Extension numbers in the call).<br> Join allows joining an ongoing call<br> The syntax for this command is: “join %CALLID% 100 101”<br> The %CALLID% is a call ID shown by executing the list command. Extensions 100 and 101 then join that ongoing call.<br> Examples:<br> ”dial 100 101” - initiates a call between two Extensions.<br> //find call id with list<br> ”join CALLID 102 103” - Extensions 102 and 103 join the call, making it a conference.<br>

____________________________________

***ABOUT THE IMPLEMENTATION***

For this implementation we used several softphones like baresip, zoiper5 and twinkle.<br> All the functionality we implemented in Golang programming language.<br> All code is wriiten in the file main.go and the explanation of that code is written in the ASTERISK REST INTERFACE(ARI).md file<br>

____________________________________
