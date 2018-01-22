package widevine

import (
	"encoding/base64"
	"encoding/json"

	"github.com/alfg/widevine/proto"
	protobuf "github.com/golang/protobuf/proto"
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

// Policy struct to set policy options for a ContentKey request.
type Policy struct {
	ContentID string
	Tracks    []string
	DRMTypes  []string
	Policy    string
}

// LicenseOptions struct to set license options for GetLicense.
type LicenseOptions struct {
	ContentID string
	Body      string
}

// GetContentKeyResponse JSON response from Widevine Cloud.
// /cenc/getcontentkey/<provider>
type GetContentKeyResponse struct {
	Status      string   `json:"status"`
	DRM         []drm    `json:"drm"`
	Tracks      []tracks `json:"tracks"`
	AlreadyUsed bool     `json:"already_used"`
}

type drm struct {
	Type     string `json:"type"`
	SystemID string `json:"system_id"`
}

type tracks struct {
	Type  string `json:"type"`
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
	PSSH  []pssh `json:"pssh"`
}

type pssh struct {
	DRMType string `json:"drm_type"`
	Data    string `json:"data"`
}

// GetLicenseResponse decoded JSON response from Widevine Cloud.
// /cenc/getlicense
type GetLicenseResponse struct {
	Status                     string            `json:"status"`
	License                    string            `json:"license"`
	Make                       string            `json:"make"`
	Model                      string            `json:"model"`
	SecurityLevel              int               `json:"security_level"`
	InternalStatus             int               `json:"internal_status"`
	DRMCertSerialNumber        string            `json:"drm_cert_serial_number"`
	DeviceWhitelistState       string            `json:"device_whitelist_state"`
	MessageType                string            `json:"message_type"`
	Platform                   string            `json:"platform"`
	DeviceState                string            `json:"device_state"`
	ClientMaxHDCPVersion       string            `json:"client_max_hdcp_version"`
	PlatformVerificationStatus string            `json:"platform_verification_status"`
	ContentOwner               string            `json:"content_owner"`
	ContentPRovider            string            `json:"content_provider"`
	SessionState               sessionState      `json:"session_state"`
	LicenseMetadata            licenseMetadata   `json:"license_metadata"`
	SupportedTracks            []supportedTracks `json:"supported_tracks"`
	PSSHData                   psshData          `json:"pssh_data"`
	ClientInfo                 []clientInfo      `json:"client_info"`
}
type licenseMetadata struct {
	ContentID   string `json:"content_id"`
	LicenseType string `json:"license_type"`
	RequestType string `json:"request_type"`
}

type supportedTracks struct {
	Type  string `json:"type"`
	KeyID string `json:"key_id"`
}
type sessionState struct {
	LicenseID      licenseID `json:"license_id"`
	SigningKey     string    `json:"signing_key"`
	KeyboxSystemID int       `json:"keybox_system_id"`
	LicenseCounter int       `json:"license_counter"`
}
type licenseID struct {
	RequestID  string `json:"request_id"`
	SessionID  string `json:"session_id"`
	PurchaseID string `json:"purchase_id"`
	Type       string `json:"type"`
	Version    int    `json:"version"`
}

type psshData struct {
	KeyID     string `json:"key_id"`
	ContentID string `json:"content_id"`
}
type clientInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
func (wp *Widevine) GetContentKey(contentID string, policy Policy) GetContentKeyResponse {
	p := wp.setPolicy(contentID, policy)
	msg := wp.buildCKMessage(p)
	resp := wp.getContentKeyRequest(msg)

	// TODO
	// Build custom PSSH from protobuf.
	// enc := wp.buildPSSH(contentID)
	// fmt.Println("pssh  build:", enc)
	return resp
}

// GetExternalContentKey creates a content key giving a contentID.
func (wp *Widevine) GetExternalContentKey(contentID string, policy Policy) GetContentKeyResponse {
	p := wp.setPolicy(contentID, policy)
	msg := wp.buildCKMessage(p)
	resp := wp.getContentKeyRequest(msg)

	// TODO
	// Build custom PSSH from protobuf.
	// enc := wp.buildPSSH(contentID)
	// fmt.Println("pssh  build:", enc)
	return resp
}

// GetLicense creates a license request used with a proxy server.
func (wp *Widevine) GetLicense(options LicenseOptions) GetLicenseResponse {
	msg := wp.buildLicenseMessage(options)
	resp := wp.getLicenseRequest(msg)
	return resp
}

func (wp *Widevine) buildPSSH(contentID string) string {
	wvpssh := &proto.WidevineCencHeader{
		Provider:  protobuf.String(wp.Provider),
		ContentId: []byte(contentID),
	}
	p, _ := protobuf.Marshal(wvpssh)
	return base64.StdEncoding.EncodeToString(p)
}

func (wp *Widevine) buildCKMessage(policy map[string]interface{}) map[string]interface{} {
	// Marshal and encode payload.
	jsonPayload, _ := json.Marshal(policy)
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

func (wp *Widevine) setPolicy(contentID string, policy Policy) map[string]interface{} {
	enc := base64.StdEncoding.EncodeToString([]byte(contentID))

	// Build tracks []interface.
	var tracks []interface{}
	for _, track := range policy.Tracks {
		tracks = append(tracks, map[string]string{"type": track})
	}

	// Build policy interface.
	// TODO: Set defaults.
	p := map[string]interface{}{
		"content_id": enc,
		"tracks":     tracks,
		"drm_types":  policy.DRMTypes,
		"policy":     policy.Policy,
	}
	return p
}

func (wp *Widevine) buildLicenseMessage(options LicenseOptions) map[string]interface{} {

	message := map[string]interface{}{
		"payload":             options.Body,
		"provider":            wp.Provider,
		"allowed_track_types": "SD_UHD1",
	}

	// Add the content ID to message if provided.
	if options.ContentID != "" {
		enc := base64.StdEncoding.EncodeToString([]byte(options.ContentID))
		message["content_id"] = enc
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
	// fmt.Println(string(dec))
	output := GetContentKeyResponse{}
	json.Unmarshal(dec, &output)
	return output
}

func (wp *Widevine) getLicenseRequest(body map[string]interface{}) GetLicenseResponse {
	// Set production or test portal.
	var url string
	if wp.Provider == "widevine_test" {
		url = widevineCloudURLTest + "/cenc/getlicense"
	} else {
		url = widevineCloudURL + "/cenc/getlicense"
	}
	// Make client call.
	resp := GetLicenseResponse{}
	client, _ := NewClient()
	client.post(url, &resp, body)
	// fmt.Println(body)
	return resp
}
