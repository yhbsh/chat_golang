# Chat app in go

- A simple but elegant multi threaded chat application in golang


# Run

```console
$ go run main.go
```

- This command will run the server.
- To test that the server is working you can do

```console
$ telnet localhost 8080
```

- This will connect to the server and allow you to start sending messages
- You can connect with multiple terminal sessions to emulator multips clients at the same time.
- The server is using goroutines, which are lightweight concurrent threads managed by the go runtime, they are not fully-fledged OS-threads, but they allow to simulate parallel programming in a very simple way.
