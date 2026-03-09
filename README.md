# protanki-networking

A library of MITM proxy capabilities for intercepting, analyzing, and modifying network packets of the game ProTanki.

## Installation
```bash
go get https://github.com/the-new-day/protanki-networking
```

## Usage Example
```go
// Intercepting lobby chat messages 
// and sending a fake system message to the chat

proxy.OnServerToClient(func(packet packets.Packet) packets.Packet {
    if packet.ID() == packets.ReceiveLobbyChatID {
        // ReceiveLobbyChat packet contains a single attribute - a vector of messages.
        messages := packets.Attr[[]map[string]any]("messages", packet)

        // ReceiveLobbyChat packet can contain either one message (when a new message arrives in the chat),
        // or multiple messages (when loading the entire chat).
        if len(messages) != 1 {
            return packet
        }

        for _, message := range messages {
            if strings.Contains(strings.ToLower(message["text"].(string)), "волки вотана") {
                fake := chat.NewReceiveLobbyChatPacket()
                attributes := map[string]any{
                    // Boolshortern => attribute is replaced by emptiness flag = 1
                    "authorStatus":  packets.Boolshortern(),
                    "systemMessage": true,
                    "targetStatus":  packets.Boolshortern(),
                    "text":          "Волки Вотана MENTIONED",
                    "warning":       false,
                }
                
                // Unpack and populate the messages attribute (one message vector)
                fake.UnwrapValues([]map[string]any{attributes})

                // Send the original packet first so our message appears below it
                err = proxy.SendToClient(packet)
                if err != nil {
                    log.Printf("[ERROR]: %s", err)
                }

                err = proxy.SendToClient(fake) // Send fake packet directly to the client
                if err != nil {
                    log.Printf("[ERROR]: %s", err)
                }

                return nil // Cancel packet forwarding since we've already sent it ourselves
            }
        }
    }
    return packet // Forward the original packet further
})
```

## Entities

### Proxy
A classic MITM proxy: it mimics a client to the server and mimics a server to the client.

Allows subscribing to events such as `OnServerToClient`, `OnClientToServer`. Also provides methods like `SendToClient`, `SendToServer`, and others.

The proxy uses 2 instances of `Protection` - one server-side and one client-side. The server-side decrypts packets from the server and encrypts packets when sending to the server, while the client-side does the same for the client. This allows working with packets independently of the protection state on the server and client, without breaking their synchronization.

### PacketStream
An abstraction over a connection that provides raw bytes. `PacketStream` converts the incoming byte stream into packets and can send packets over the same connection.

### PacketHandler
A wrapper around PacketStream that provides the ability to subscribe to incoming and outgoing packets. Subscribers can analyze, modify, cancel sending, or send packets themselves.

### Protection
A stream cipher used by the game to protect data from simple sniffing. It is activated with keys sent by the server in the first packet.

### Packets
The library implements most of the packets used by the game. Packet structure: `[Length (4 bytes)] [ID (4 bytes)] [Data]`.

### Codecs
Codecs are used to decrypt the `[Data]` block - they convert a sequence of bytes into a sequence of entities of a specific type and vice versa. There are 4 types of codecs:

* Primitive (`int`, `long`, `short`, `bool`, `byte`, `float`)
* Complex (`string`, `json`, etc.)
* Collections (`MultiCodec`, `VectorCodec`)
* Custom (codecs for more complex structures, such as player statistics, using the previous 3 types)

## Client Connection
To make the client connect to the proxy instead of the game server, you need to replace the `Prelauncher.swf` file that `ProTanki.exe` downloads to launch the launcher.

To do this, you can run an HTTP server, for example at `127.0.0.1:1211`, which will serve your modified `Prelauncher.swf`. In the `META-INF/AIR/application.xml` file, you need to replace the address from `http://146.59.110.103/Prelauncher.swf` to `http://127.0.0.1:1211/Prelauncher.swf`. Now `ProTanki.exe` will download and run your own prelauncher.

Prelauncher replacement is necessary because it contains the game server address that the client connects to. Download the official `Prelauncher.swf` from the address above, decompile it (for example, using FFDec), and replace the line `http://146.59.110.103/config.xml` with `http://127.0.0.1:1211/config.xml`, place `config.xml` next to `Prelauncher.swf`. The `config.xml` file specifies the server address that the client connects to during the game.

Ready-to-use `Prelauncher.swf` and `config.xml` files are located [here](/prelauncher).

In [/cmd/default](/cmd/default/) there is an example application.

## Notes
Using traffic interception tools may violate the game's rules and could result in account suspension. This library is created solely for educational purposes.

This library is based on [Teinc3's ProboTanki-Lib](https://github.com/Teinc3/ProboTanki-Lib), but it's not an exact port.