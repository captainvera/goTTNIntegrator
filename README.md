# goTTNIntegrator
Integrator between TTN([The Things Network](https://www.thethingsnetwork.com)) device uplink data feed and an OpenSensors.io endpoint.

## Dependencies
This script depends on the TTN [go-sdk](https://www.thethingsnetwork.org/docs/applications/golang/quick-start.html)

## How to Use

To use the integrator you need to setup the config file, an example is included in this repository. It should have the following format:
```
{
    "TTN":{
        "AppId":"",
        "ApiKey":""
        },
    "OpenSensors":{
        "ApiKey":"",
        "ClientId":,
        "Password":"",
        "Topic":""
    }
}
```

The easiest way to run the integrator is the following:
```
go run integrator.go
```

If you do not want to share the entire payload received from your device, it's possible to filter the uplink payload. This will only forward certain data ex:[light, sound]:
```
go run integrator.go -f light sound
```

To see other options use:
```
go run integrator.go -h
```
