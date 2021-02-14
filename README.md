# Widevine
Golang Client API for Widevine Cloud.

https://godoc.org/github.com/alfg/widevine

![Unittest](https://github.com/Cooomma/widevine/workflows/Unittest/badge.svg) [![GoDoc](https://godoc.org/github.com/alfg/widevine?status.svg)](https://godoc.org/github.com/alfg/widevine)  


## Install
```
go get github.com/alfg/widevine
```

## Usage

#### Generating license keys
```golang
// Set Widevine options and create instance.
options := widevine.Options{
    Key:      []byte{key},     // Your Widevine Key as a byte array.
    IV:       []byte{iv},      // Your Widevine IV as a byte array.
    Provider: "widevine_test", // Your Widevine Provider/Portal.
}

// Create the Widevine instance.
wv := widevine.New(options)

// Your video content ID, usually a GUID.
contentID := "testing"

// Set policy options.
policy := widevine.Policy{
    ContentID: contentID,
    Tracks:    []string{"SD", "HD", "AUDIO"},
    DRMTypes:  []string{"WIDEVINE"},
    Policy:    "default",
}

// Make the request to generate or get a content key.
resp := wv.GetContentKey(contentID, policy)

// Response data from Widevine Cloud.
fmt.Println("status: ", resp.Status)
fmt.Println("drm: ", resp.DRM)
for _, v := range resp.Tracks {
    fmt.Println("key_id: ", v.KeyID)
    fmt.Println("type: ", v.Type)
    fmt.Println("drm_type: ", v.PSSH[0].DRMType)
    fmt.Println("data: ", v.PSSH[0].Data)
}
fmt.Println("already_used: ", resp.AlreadyUsed)
```

#### License Proxy
You can also use this package to create a license proxy.

See: [examples/proxy](/examples/proxy)


## Examples
See: [examples](/examples)

## Develop
TODO

`protoc.exe --go_out=. *.proto`

## TODO
* External Keys
* Custom PSSH API
* Tests
* More error handling
* Implement more Widevine features

## Resources
* https://www.widevine.com/news
* https://integration.widevine.com/
* https://developers.google.com/protocol-buffers/
* https://github.com/google/shaka-player
* https://support.google.com/widevine/answer/6048495?hl=en
* https://www.widevine.com/contact for Widevine documents (now restricted to licensees)
