# PubSub package

## Features
* Support to RabbitMQ

## How to setup?
Set the RabbitMQ connection string as environment variable

```sh
export AMQP_URL="amqps://..."
```

Then,

```sh
go run main.go
```