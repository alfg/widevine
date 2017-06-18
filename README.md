# Widevine
Golang Widevine Cloud Client API and License Proxy for Widevine.

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

// Make the request to generate or get a content key.
resp := wv.GetContentKey(contentID)

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
* Custom PSSH API
* Tests
* More error handling
* Implement more Widevine features

## Resources
* https://www.widevine.com/product_news.html
* https://storage.googleapis.com/wvdocs/Widevine_DRM_Getting_Started.pdf
* https://storage.googleapis.com/wvdocs/Widevine_DRM_Architecture_Overview.pdf
* https://storage.googleapis.com/wvdocs/Widevine_DRM_Encryption_API.pdf
* https://storage.googleapis.com/wvdocs/Widevine_DRM_Proxy_Integration.pdf
* https://developers.google.com/protocol-buffers/
* https://github.com/google/shaka-player
* https://support.google.com/widevine/answer/6048495?hl=en