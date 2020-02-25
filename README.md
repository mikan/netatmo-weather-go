netatmo-weather-go
==================

Unofficial client library for Netatmo Weather Station written in Go.

## Usage 

### Setup

```
go get github.com/mikan/netatmo-weather-go
```

### Create a client

```go
client, err := netatmo.NewClient(context.Background(), clientID, clientSecret, username, password)
if err != nil {
    panic(err)
}
```

### Get stations data

```go
devices, user, err := client.GetStationsData()
if err != nil {
    panic(err)
}
fmt.Println(user)
fmt.Println(devices)
```

### Get measure

```go
value, err := client.GetMeasureByNewest(device, module)
if err != nil {
    panic(err)
}
fmt.Println(value)
```

### Example code

See `cmd/example` directory.

Usage:

```
go run cmd/example/*.go -c <CLIENT_ID> -s <CLIENT_SECRET> -u <USER> -p <PASSWORD>
```

## License

netatmo-weather-go licensed under the [BSD 3-clause](LICENSE).

## Author

- [mikan](https://github.com/mikan)
