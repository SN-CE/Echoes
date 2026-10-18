# Echoes

A self-hosted, terminal-based group chat system built in Go.


## What is Echoes?

Echoes is a minimal, self-hosted chat system consisting of two binaries — a server and a client. There are no accounts, no cloud, no database, no third-party services. Just two binaries, a port, an IP address and your friends. Completely Private.


## How It Works

1. One person runs the server binary. It asks for a port number and starts listening.
2. Everyone — including the server host — runs the client binary. It asks for the server's IP address or domain name.
3. Each client picks a username on connect. The list of connected clients will be shown, so pick a unique one.

Messages are broadcast in real time to all connected clients. When a client disconnects, everyone is notified. When the server shuts down, all connections are closed cleanly.


## Requirements

No runtime required. Just download and run the binary for your platform.

- **Linux** : `Echoes` (client), `Echoes_Server` (server)
- **Windows** : `Echoes.exe` (client), `Echoes_Server.exe` (server)


## Usage

### Server

```
./Echoes_Server
```

The server will prompt for a port number, then start listening.

### Client

```
./Echoes
```

The client will prompt for the server's IP address or domain name. Both raw IPs and domain names work — Go resolves DNS automatically.

**Example:**

```
Enter the IP or URL of the Echoes Server you wish to connect to - 192.168.1.10:8080
```

or

```
Enter the IP or URL of the Echoes Server you wish to connect to - chat.example.com:8080
```


## Running Over The Internet

If the person hosting the server has a **public IP**, no additional setup is needed. Share your IP and port and you're done.

If the host is behind **CGNAT** (common with most home ISPs), they won't have a public IP. You have a few options, though:

- Use a **VPS** to run the server on a machine with a public IP.
- Use a **tunnel service** like [Pagekite](https://pagekite.net/) to expose the server through the CGNAT.


## Technical Details

- **Protocol** : WebSocket over HTTP (`gorilla/websocket`)
- **Concurrency** : Each client connection runs in its own goroutine. A `sync.RWMutex` protects the shared client map.
- **TUI** : Built with `tview` and `tcell`
- **Persistence** : None for now. Chats exist only for the duration of the session. Nothing is written to disk. Later on there might be an option added to store chat history to a text file and recall later. But it will be entirely optional.

## Roadmap

- [ ] End-to-end encryption
- [ ] Adding persistence options on both ends, server and client. For now, server only displays a simple log that is not saved.
- [ ] Tor routing support
- [ ] Adding a separate channel for commands



---



## Why

Most chat platforms sit in the middle of your conversations. They own the infrastructure, they hold the keys, and they can be compelled to hand over your data.

Echoes removes the middleman entirely. Either you or your friends/partners/colleagues own it. No company is involved. No data is stored. If you don't like a server or it's host, you can easily start hosting yourself. 

Once encryption is implemented, the contents of every conversation will be completely opaque to anyone who isn't in it.




## Built With And Thank You To

- [gorilla/websocket](https://github.com/gorilla/websocket)
- [tview](https://github.com/rivo/tview)
- [tcell](https://github.com/gdamore/tcell)
- [encoding](https://github.com/gdamore/encoding)
- [go-colorful](https://github.com/lucasb-eyer/go-colorful)
- [uniseg](https://github.com/rivo/uniseg)
- [sys](https://pkg.go.dev/golang.org/x/sys)
- [term](https://pkg.go.dev/golang.org/x/term)
- [text](https://pkg.go.dev/golang.org/x/text)
