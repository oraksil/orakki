# Orakki

`Orakki` means an arcade gaming console in Korean. It is a core instance that provides clients with game streaming via WebRTC. It fetches encoded video/audio frames from Gipan and sends those packets to client. And it receives player controller input and forwards it to Gipan.

`Orakki` instance should be provisioned on demand in a cluster mode. It might need TURN server in some cases where client cannot communicate with server directly. It's built with Golang.

# Prerequites

`RabbitMQ` is required for `Orakki` to communicate with `Azumma`. Please refer to [this](https://github.com/oraksil/azumma#message-queue-rabbitmq) for more details.

# Run

```bash
$ go run cmd/app.go
```