package widevine

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/alfg/widevine/pssh"
	"github.com/golang/protobuf/proto"
)

const (
	getLicenseURL    = "https://license.uat.widevine.com/cenc/getlicense"
	getContentKeyURL = "http://license.uat.widevine.com/cenc/getcontentkey/widevine_test"
)

// Widevine structure.
type Widevine struct {
	Key      []byte
	IV       []byte
	Provider string
	URL      string
}

// Options structure.
type Options struct {
	Key      []byte
	IV       []byte
	Provider string
	URL      string
}

// GetContentKeyResponse JSON response from Widevine.
type GetContentKeyResponse struct {
	Status string `json:"status"`
	DRM    []struct {
		Type     string `json:"type"`
		SystemID string `json:"system_id"`
	}
	Tracks []struct {
		Type  string `json:"type"`
		KeyID string `json:"key_id"`
		PSSH  []struct {
			DRMType string `json:"drm_type"`
			Data    string `json:"data"`
		}
	}
	AlreadyUsed bool `json:"already_used"`
}

// New returns a Widevine instance.
func New(opts Options) *Widevine {

	wv := &Widevine{
		Key:      opts.Key,
		IV:       opts.IV,
		Provider: opts.Provider,
	}
	return wv
}

// GetContentKey creates a content key giving a contentID.
func (wp *Widevine) GetContentKey(contentID string) GetContentKeyResponse {
	msg := wp.buildMessage(contentID)
	resp := wp.sendRequest(msg)

	enc := wp.buildPSSH(contentID)
	fmt.Println("pssh  build:", enc)
	return resp
}

func (wp *Widevine) buildPSSH(contentID string) string {
	wvpssh := &pssh.WidevineCencHeader{
		Provider:  proto.String(wp.Provider),
		ContentId: []byte(contentID),
	}
	p, _ := proto.Marshal(wvpssh)

	return base64.StdEncoding.EncodeToString(p)
}

func (wp *Widevine) buildMessage(contentID string) map[string]interface{} {
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

	// Marshal and encode payload.
	jsonPayload, _ := json.Marshal(payload)
	b64payload := base64.StdEncoding.EncodeToString([]byte(jsonPayload))

	// Create signature and postBody.
	crypto := NewCrypto(wp.Key, wp.IV)
	postBody := map[string]interface{}{
		"request":   b64payload,
		"signature": crypto.generateSignature(jsonPayload),
		"signer":    wp.Provider,
	}
	return postBody
}

func (wp *Widevine) sendRequest(body map[string]interface{}) GetContentKeyResponse {
	// Make client call.
	resp := make(map[string]string)
	client, _ := NewClient()
	client.post(getContentKeyURL, &resp, body)

	// Decode and unmarshal the response.
	dec, _ := base64.StdEncoding.DecodeString(resp["response"])
	output := GetContentKeyResponse{}
	json.Unmarshal(dec, &output)
	return output
}
