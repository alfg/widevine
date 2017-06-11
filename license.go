package widevine

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

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

const (
	// KEY      = "1ae8ccd0e7985cc0b6203a55855a1034afc252980e970ca90e5202689f947ab9"
	// IV       = "d58ce954203b7c9a9a9d467f59839249"
	getLicenseURL    = "https://license.uat.widevine.com/cenc/getlicense"
	getContentKeyURL = "http://license.uat.widevine.com/cenc/getcontentkey/widevine_test"
	provider         = "widevine_test"
)

// Widevine structure.
type Widevine struct {
	Key      []byte
	IV       []byte
	Provider string
	URL      string
	// ContentID  string
	// ContentKey string
	// KeyID      string
}

// Options structure.
type Options struct {
	Key      []byte
	IV       []byte
	Provider string
	URL      string
}

type GetContentKeyResponse struct {
	// Response string `json:"response"`
	Status string `json:"status"`
}

// New returns a Widevine instance.
func New(opts ...Options) *Widevine {

	wv := &Widevine{
		Key:      key,
		IV:       iv,
		URL:      getContentKeyURL,
		Provider: provider,
	}
	return wv
}

// GetContentKey creates a content key giving a contentID.
func (wp *Widevine) GetContentKey(contentID string) GetContentKeyResponse {
	enc := base64.StdEncoding.EncodeToString([]byte(contentID))
	payload := map[string]interface{}{
		"content_id": enc,
		"tracks": []interface{}{
			map[string]string{"type": "SD"},
			map[string]string{"type": "HD"},
			map[string]string{"type": "AUDIO"},
		},
		"drm_types": []string{"WIDEVINE"},
		"policy":    "default",
	}

	jsonPayload, _ := json.Marshal(payload)

	b64payload := base64.StdEncoding.EncodeToString([]byte(jsonPayload))

	signature := generateSignature(jsonPayload)

	postBody := map[string]interface{}{
		"request":   b64payload,
		"signature": signature,
		"signer":    provider,
	}

	// Make client call.
	resp := make(map[string]string)
	client, _ := NewClient()
	client.post(getContentKeyURL, &resp, postBody, "json")

	// Decode the response.
	dec, _ := base64.StdEncoding.DecodeString(resp["response"])
	fmt.Println(string(dec))

	// Unmarshal to response struct.
	output := GetContentKeyResponse{}
	json.Unmarshal(dec, &output)
	fmt.Println(output.Status)
	return output
}

func buildRequest() {

	// payload := buildMessage(body)

}

func buildMessage() {

}
