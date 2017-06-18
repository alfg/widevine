package widevine

import (
	"encoding/base64"
	"encoding/json"

	"github.com/alfg/widevine/pssh"
	"github.com/golang/protobuf/proto"
)

// Widevine Cloud URLs.
const (
	widevineCloudURL     = "https://license.widevine.com"
	widevineCloudURLTest = "https://license.uat.widevine.com"
)

// Widevine structure.
type Widevine struct {
	Key      []byte
	IV       []byte
	Provider string
	URL      string
}

// Options provided to Widevine{} instance.
type Options struct {
	Key      []byte
	IV       []byte
	Provider string
	URL      string
}

// GetContentKeyResponse JSON response from Widevine Cloud /cenc/getcontentkey/<provider>.
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

// LicenseResponse decoded JSON response from Widevine Cloud /cenc/getlicense.
type LicenseResponse struct {
	Status          string `json:"status"`
	License         string `json:"license"`
	LicenseMetadata []struct {
		ContentID   string `json:"content_id"`
		LicenseType string `json:"license_type"`
		RequestType string `json:"request_type"`
	}
	SupportedTracks []struct {
		Type  string `json:"type"`
		KeyID string `json:"key_id"`
	}
	Make           string `json:"make"`
	Model          string `json:"model"`
	SecurityLevel  int    `json:"security_level"`
	InternalStatus int    `json:"internal_status"`
	SessionState   struct {
		LicenseID struct {
			RequestID  string `json:"request_id"`
			SessionID  string `json:"session_id"`
			PurchaseID string `json:"purchase_id"`
			Type       string `json:"type"`
			Version    int    `json:"version"`
		}
		SigningKey     string `json:"signing_key"`
		KeyboxSystemID int    `json:"keybox_system_id"`
		LicenseCounter int    `json:"license_counter"`
	}
	DRMCertSerialNumber  string `json:"drm_cert_serial_number"`
	DeviceWhitelistState string `json:"device_whitelist_state"`
	MessageType          string `json:"message_type"`
	Platform             string `json:"platform"`
	DeviceState          string `json:"device_state"`
	PSSHData             struct {
		KeyID     string `json:"key_id"`
		ContentID string `json:"content_id"`
	}
	ClientMaxHDCPVersion string `json:"client_max_hdcp_version"`
	ClientInfo           []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	PlatformVerificationStatus string `json:"platform_verification_status"`
	ContentOwner               string `json:"content_owner"`
	ContentPRovider            string `json:"content_provider"`
}

// New returns a Widevine instance with options.
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
	resp := wp.getContentKeyRequest(msg)

	// TODO
	// Build custom PSSH from protobuf.
	// enc := wp.buildPSSH(contentID)
	// fmt.Println("pssh  build:", enc)
	return resp
}

// LicenseRequest creates a license request used with a proxy server.
func (wp *Widevine) LicenseRequest(contentID string, body string) LicenseResponse {
	msg := wp.buildLicenseMessage(contentID, body)
	resp := wp.getLicenseRequest(msg)
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

func (wp *Widevine) buildLicenseMessage(contentID string, body string) map[string]interface{} {
	enc := base64.StdEncoding.EncodeToString([]byte(contentID))

	message := map[string]interface{}{
		"payload":             body,
		"content_id":          enc,
		"provider":            wp.Provider,
		"allowed_track_types": "SD_HD",
	}
	jsonMessage, _ := json.Marshal(message)
	b64message := base64.StdEncoding.EncodeToString(jsonMessage)

	// Create signature and postBody.
	crypto := NewCrypto(wp.Key, wp.IV)
	postBody := map[string]interface{}{
		"request":   b64message,
		"signature": crypto.generateSignature(jsonMessage),
		"signer":    wp.Provider,
	}
	return postBody
}

func (wp *Widevine) getContentKeyRequest(body map[string]interface{}) GetContentKeyResponse {
	// Set production or test portal.
	var url string
	if wp.Provider == "widevine_test" {
		url = widevineCloudURLTest + "/cenc/getcontentkey/widevine_test"
	} else {
		url = widevineCloudURL + "/cenc/getcontentkey/" + wp.Provider
	}

	// Make client call.
	resp := make(map[string]string)
	client, _ := NewClient()
	client.post(url, &resp, body)

	// Decode and unmarshal the response.
	dec, _ := base64.StdEncoding.DecodeString(resp["response"])
	output := GetContentKeyResponse{}
	json.Unmarshal(dec, &output)
	return output
}

func (wp *Widevine) getLicenseRequest(body map[string]interface{}) LicenseResponse {
	// Set production or test portal.
	var url string
	if wp.Provider == "widevine_test" {
		url = widevineCloudURLTest + "/cenc/getlicense"
	} else {
		url = widevineCloudURL + "/cenc/getlicense"
	}
	// Make client call.
	resp := LicenseResponse{}
	client, _ := NewClient()
	client.post(url, &resp, body)
	return resp
}
