package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"

	"github.com/alfg/widevine"
)

// AES key and IV for the provider "widevine_test".
// Use these test keys for testing or integration tests.
var (
	key = []byte{
		0x1a, 0xe8, 0xcc, 0xd0, 0xe7, 0x98, 0x5c, 0xc0,
		0xb6, 0x20, 0x3a, 0x55, 0x85, 0x5a, 0x10, 0x34,
		0xaf, 0xc2, 0x52, 0x98, 0x0e, 0x97, 0x0c, 0xa9,
		0x0e, 0x52, 0x02, 0x68, 0x9f, 0x94, 0x7a, 0xb9}

	iv = []byte{
		0xd5, 0x8c, 0xe9, 0x54, 0x20, 0x3b, 0x7c, 0x9a,
		0x9a, 0x9d, 0x46, 0x7f, 0x59, 0x83, 0x92, 0x49}
)

// Test contentID.
const contentID = "fkj3ljaSdfalkr3j"

func main() {

	// Create handler and http server.
	http.HandleFunc("/proxy", proxy)
	http.ListenAndServe(":8000", nil)
}

func proxy(w http.ResponseWriter, r *http.Request) {

	// Read bytes from license request.
	buf, _ := ioutil.ReadAll(r.Body)
	body := base64.StdEncoding.EncodeToString(buf)

	// Set Widevine options and create instance.
	options := widevine.Options{
		Key:      key,
		IV:       iv,
		Provider: "widevine_test",
	}
	wv := widevine.New(options)

	// Create license request.
	data := wv.LicenseRequest(contentID, body)
	b, _ := base64.StdEncoding.DecodeString(data.License)

	// CORS required for Javascript players.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write decoded license bytes back to player.
	w.Write(b)
}
