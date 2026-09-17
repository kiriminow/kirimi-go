package kirimi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

// Constants
const (
	DefaultBaseURL   = "https://api.kirimi.id"
	ContentType      = "application/json"
	MaxMessageLength = 1200
)

// Package Types
const (
	PackageFree   = 1
	PackageLite1  = 2
	PackageBasic1 = 3
	PackagePro1   = 4
	PackageLite2  = 6
	PackageBasic2 = 7
	PackagePro2   = 8
	PackageLite3  = 9
	PackageBasic3 = 10
	PackagePro3   = 11
)

// OTP v2 delivery methods
const (
	OtpMethodWhatsApp = "whatsapp"
	OtpMethodWaba     = "waba"
	OtpMethodDevice   = "device"
	OtpMethodWabaUser = "waba_user"
)

// Client represents the Kirimi API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Kirimi API client
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientWithBaseURL creates a new client with custom base URL
func NewClientWithBaseURL(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetTimeout sets the HTTP client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.HTTPClient.Timeout = timeout
}

// Response represents the standard API response format
type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

// GenerateOTPRequest represents the request for generating OTP
type GenerateOTPRequest struct {
	UserCode           string `json:"user_code"`
	DeviceID           string `json:"device_id"`
	Phone              string `json:"phone"`
	Secret             string `json:"secret"`
	OtpLength          int    `json:"otp_length,omitempty"`
	OtpType            string `json:"otp_type,omitempty"`
	CustomOtpMessage   string `json:"customOtpMessage,omitempty"`
	CustomOtpText      string `json:"customOtpText,omitempty"`
	EnableTypingEffect *bool  `json:"enableTypingEffect,omitempty"`
	TypingSpeedMs      int    `json:"typingSpeedMs,omitempty"`
}

// GenerateOTPResponse represents the response data for OTP generation
type GenerateOTPResponse struct {
	Phone     string `json:"phone"`
	Message   string `json:"message"`
	ExpiresIn string `json:"expires_in"`
}

// ValidateOTPRequest represents the request for validating OTP
type ValidateOTPRequest struct {
	UserCode string `json:"user_code"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	OTP      string `json:"otp"`
	Secret   string `json:"secret"`
}

// ValidateOTPResponse represents the response data for OTP validation
type ValidateOTPResponse struct {
	Phone      string    `json:"phone"`
	Verified   bool      `json:"verified"`
	VerifiedAt time.Time `json:"verified_at"`
}

// SendMessageRequest represents the request for sending a message
type SendMessageRequest struct {
	UserCode           string `json:"user_code"`
	DeviceID           string `json:"device_id"`
	Receiver           string `json:"receiver"`
	Message            string `json:"message"`
	Secret             string `json:"secret"`
	MediaURL           string `json:"media_url,omitempty"`
	FileName           string `json:"fileName,omitempty"`
	EnableTypingEffect *bool  `json:"enableTypingEffect,omitempty"`
	TypingSpeedMs      int    `json:"typingSpeedMs,omitempty"`
	QuotedMessageID    string `json:"quotedMessageId,omitempty"`
}

// SendMessageResponse represents the response data for sending message
type SendMessageResponse struct {
	MessageLength int    `json:"message_length"`
	MediaURL      string `json:"media_url,omitempty"`
	HasMedia      bool   `json:"has_media"`
}

// SendMessageFileRequest represents the request for sending a file via multipart/form-data
type SendMessageFileRequest struct {
	UserCode        string
	DeviceID        string
	Receiver        string
	Secret          string
	File            io.Reader
	Filename        string
	Message         string
	Caption         string
	FileName        string
	QuotedMessageID string
}

// SendWabaMessageRequest represents the request for sending a WABA message
type SendWabaMessageRequest struct {
	UserCode     string              `json:"user_code"`
	Secret       string              `json:"secret"`
	WabaID       string              `json:"waba_id"`
	To           string              `json:"to"`
	TemplateName string              `json:"template_name"`
	Variables    []string            `json:"variables,omitempty"`
	Header       *WabaTemplateHeader `json:"header,omitempty"`
	Buttons      []interface{}       `json:"buttons,omitempty"`
}

// WabaTemplateHeader represents a WABA template header component
type WabaTemplateHeader struct {
	Type     string `json:"type"`
	Link     string `json:"link,omitempty"`
	ID       string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	Text     string `json:"text,omitempty"`
}

// WabaReplyMessage represents a free-form WABA reply. Only the fields relevant
// to Type are serialized; the rest stay empty thanks to omitempty.
type WabaReplyMessage struct {
	Type        string                 `json:"type"`
	Text        string                 `json:"text,omitempty"`
	MediaURL    string                 `json:"media_url,omitempty"`
	Caption     string                 `json:"caption,omitempty"`
	Filename    string                 `json:"filename,omitempty"`
	Interactive map[string]interface{} `json:"interactive,omitempty"`
}

// WabaReplyRequest represents the request for /v1/waba/messages/reply
type WabaReplyRequest struct {
	UserCode string           `json:"user_code"`
	Secret   string           `json:"secret"`
	WabaID   string           `json:"waba_id"`
	To       string           `json:"to"`
	Message  WabaReplyMessage `json:"message"`
}

// WabaConversationsRequest represents the request for /v1/waba/conversations
type WabaConversationsRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Limit    int    `json:"limit,omitempty"`
	Page     int    `json:"page,omitempty"`
}

// WabaTemplateSyncRequest represents the request for /v1/waba/templates/sync
type WabaTemplateSyncRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	WabaID   string `json:"waba_id"`
}

// WabaSendOtpRequest represents the request for /v1/waba/send-otp
type WabaSendOtpRequest struct {
	UserCode     string `json:"user_code"`
	Secret       string `json:"secret"`
	WabaID       string `json:"waba_id"`
	To           string `json:"to"`
	TemplateName string `json:"template_name"`
}

// WabaVerifyOtpRequest represents the request for /v1/waba/verify-otp
type WabaVerifyOtpRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	WabaID   string `json:"waba_id"`
	To       string `json:"to"`
	OtpCode  string `json:"otp_code"`
}

// DeviceStatusRequest is used for device-status and device-status-enhanced
type DeviceStatusRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	DeviceID string `json:"device_id"`
}

// CreateDeviceRequest represents the request for /v1/create-device
type CreateDeviceRequest struct {
	UserCode    string      `json:"user_code"`
	Secret      string      `json:"secret"`
	PackageID   interface{} `json:"package_id,omitempty"`
	VoucherCode string      `json:"voucher_code,omitempty"`
}

// ConnectDeviceRequest represents the request for /v1/connect-device
type ConnectDeviceRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	DeviceID string `json:"device_id"`
}

// RenewDeviceRequest represents the request for /v1/renew-device
type RenewDeviceRequest struct {
	UserCode    string      `json:"user_code"`
	Secret      string      `json:"secret"`
	DeviceID    string      `json:"device_id"`
	PackageID   interface{} `json:"package_id,omitempty"`
	VoucherCode string      `json:"voucher_code,omitempty"`
}

// SaveContactRequest represents the request for saving a contact
type SaveContactRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Nama     string `json:"nama"`
	Nomor    string `json:"nomor"`
	DeviceID string `json:"device_id,omitempty"`
}

// BulkContact represents a single contact for a bulk save
type BulkContact struct {
	Nama  string `json:"nama"`
	Nomor string `json:"nomor"`
}

// SaveContactsBulkRequest represents the request for /v1/save-contacts-bulk
type SaveContactsBulkRequest struct {
	UserCode string        `json:"user_code"`
	Secret   string        `json:"secret"`
	Contacts []BulkContact `json:"contacts"`
	DeviceID string        `json:"device_id,omitempty"`
}

// SendOtpV2Request represents the request for /v2/otp/send
type SendOtpV2Request struct {
	UserCode      string `json:"user_code"`
	Secret        string `json:"secret"`
	Phone         string `json:"phone"`
	Method        string `json:"method,omitempty"`
	AppName       string `json:"app_name,omitempty"`
	DeviceID      string `json:"device_id,omitempty"`
	WabaID        string `json:"waba_id,omitempty"`
	TemplateName  string `json:"template_name,omitempty"`
	CustomMessage string `json:"custom_message,omitempty"`
}

// VerifyOtpV2Request represents the request for /v2/otp/verify
type VerifyOtpV2Request struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Phone    string `json:"phone"`
	OtpCode  string `json:"otp_code"`
}

// OtpReverseCreateRequest represents the request for /v2/otp-reverse/create
type OtpReverseCreateRequest struct {
	UserCode       string `json:"user_code"`
	Secret         string `json:"secret"`
	Phone          string `json:"phone"`
	DeviceID       string `json:"device_id"`
	AppName        string `json:"app_name,omitempty"`
	CallbackURL    string `json:"callback_url,omitempty"`
	CustomMessage  string `json:"custom_message,omitempty"`
	SuccessMessage string `json:"success_message,omitempty"`
	FailureMessage string `json:"failure_message,omitempty"`
}

// OtpReverseStatusRequest represents the request for /v2/otp-reverse/status
type OtpReverseStatusRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Token    string `json:"token"`
}

// BroadcastMessageRequest represents the request for broadcast-message
type BroadcastMessageRequest struct {
	UserCode           string   `json:"user_code"`
	Secret             string   `json:"secret"`
	DeviceID           string   `json:"device_id"`
	Label              string   `json:"label"`
	Numbers            []string `json:"numbers"`
	Message            string   `json:"message"`
	Delay              float64  `json:"delay,omitempty"`
	DelayMin           float64  `json:"delayMin,omitempty"`
	DelayMax           float64  `json:"delayMax,omitempty"`
	MediaURL           string   `json:"media_url,omitempty"`
	FileName           string   `json:"fileName,omitempty"`
	StartedAt          string   `json:"started_at,omitempty"`
	EnableTypingEffect *bool    `json:"enableTypingEffect,omitempty"`
	TypingSpeedMs      int      `json:"typingSpeedMs,omitempty"`
}

// CreateDepositRequest represents the request for /v1/create-deposit
type CreateDepositRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Nominal  int    `json:"nominal"`
}

// DepositRefRequest represents the request for deposit-status and cancel-deposit
type DepositRefRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Ref      string `json:"ref"`
}

// ListDepositsRequest represents the request for list-deposits
type ListDepositsRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Status   string `json:"status,omitempty"`
}

// ListDevicesRequest represents the request for /v1/list-devices
type ListDevicesRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
	Response   *Response
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// --- Internal helpers ---

// makeRequest makes a JSON POST request using a full path (e.g. "/v1/send-message")
func (c *Client) makeRequest(method, path string, body interface{}) (*Response, error) {
	url := c.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", ContentType)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &apiResp, &APIError{
			StatusCode: resp.StatusCode,
			Message:    apiResp.Message,
			Response:   &apiResp,
		}
	}

	return &apiResp, nil
}

// makeGetRequest makes a GET request (for health check)
func (c *Client) makeGetRequest(path string) (*Response, error) {
	url := c.BaseURL
	if path != "" {
		url = c.BaseURL + "/" + path
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &apiResp, nil
}

// unmarshalData decodes Response.Data into v
func unmarshalData(resp *Response, v interface{}) error {
	if resp.Data == nil {
		return nil
	}
	if err := json.Unmarshal(resp.Data, v); err != nil {
		return fmt.Errorf("failed to unmarshal response data: %w", err)
	}
	return nil
}

// authBody builds a minimal JSON body with user_code + secret
func (c *Client) authBody(userCode, secret string) map[string]string {
	return map[string]string{"user_code": userCode, "secret": secret}
}

// --- WhatsApp Unofficial ---

// SendMessage sends a WhatsApp message with optional media
func (c *Client) SendMessage(req SendMessageRequest) (*SendMessageResponse, error) {
	if len(req.Message) > MaxMessageLength {
		return nil, fmt.Errorf("message length exceeds maximum of %d characters", MaxMessageLength)
	}

	resp, err := c.makeRequest("POST", "/v1/send-message", req)
	if err != nil {
		return nil, err
	}

	var out SendMessageResponse
	if err := unmarshalData(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendMessageFast sends a WhatsApp message without typing effect
func (c *Client) SendMessageFast(req SendMessageRequest) (*SendMessageResponse, error) {
	if len(req.Message) > MaxMessageLength {
		return nil, fmt.Errorf("message length exceeds maximum of %d characters", MaxMessageLength)
	}

	resp, err := c.makeRequest("POST", "/v1/send-message-fast", req)
	if err != nil {
		return nil, err
	}

	var out SendMessageResponse
	if err := unmarshalData(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendMessageFile sends a WhatsApp message with a file via multipart/form-data (max 50MB)
func (c *Client) SendMessageFile(req SendMessageFileRequest) (*Response, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("user_code", req.UserCode)
	_ = w.WriteField("secret", req.Secret)
	_ = w.WriteField("device_id", req.DeviceID)
	_ = w.WriteField("receiver", req.Receiver)
	if req.Message != "" {
		_ = w.WriteField("message", req.Message)
	}
	if req.Caption != "" {
		_ = w.WriteField("caption", req.Caption)
	}
	if req.FileName != "" {
		_ = w.WriteField("fileName", req.FileName)
	}
	if req.QuotedMessageID != "" {
		_ = w.WriteField("quotedMessageId", req.QuotedMessageID)
	}

	filename := req.Filename
	if filename == "" {
		filename = "file"
	}

	part, err := w.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}
	w.Close()

	url := c.BaseURL + "/v1/send-message-file"
	httpReq, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	httpResp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if httpResp.StatusCode >= 400 {
		return &apiResp, &APIError{
			StatusCode: httpResp.StatusCode,
			Message:    apiResp.Message,
			Response:   &apiResp,
		}
	}

	return &apiResp, nil
}

// BroadcastMessage sends a message to multiple recipients (numbers as a JSON array)
func (c *Client) BroadcastMessage(req BroadcastMessageRequest) (*Response, error) {
	return c.makeRequest("POST", "/v1/broadcast-message", req)
}

// --- WABA ---

// SendWabaMessage sends a message via WhatsApp Business API (Meta Cloud API)
func (c *Client) SendWabaMessage(req SendWabaMessageRequest) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/send-message", req)
}

// WabaReply sends a free-form reply to a WABA conversation
func (c *Client) WabaReply(userCode, secret, wabaID, to string, message WabaReplyMessage) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/messages/reply", WabaReplyRequest{
		UserCode: userCode,
		Secret:   secret,
		WabaID:   wabaID,
		To:       to,
		Message:  message,
	})
}

// WabaConversations lists WABA conversations still inside the 24h window
func (c *Client) WabaConversations(userCode, secret string, limit, page int) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/conversations", WabaConversationsRequest{
		UserCode: userCode,
		Secret:   secret,
		Limit:    limit,
		Page:     page,
	})
}

// WabaTemplatesSync refreshes template status from Meta for one WABA
func (c *Client) WabaTemplatesSync(userCode, secret, wabaID string) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/templates/sync", WabaTemplateSyncRequest{
		UserCode: userCode,
		Secret:   secret,
		WabaID:   wabaID,
	})
}

// WabaSendOtp sends an OTP through a WABA + AUTHENTICATION template
func (c *Client) WabaSendOtp(userCode, secret, wabaID, to, templateName string) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/send-otp", WabaSendOtpRequest{
		UserCode:     userCode,
		Secret:       secret,
		WabaID:       wabaID,
		To:           to,
		TemplateName: templateName,
	})
}

// WabaVerifyOtp verifies an OTP previously sent through WabaSendOtp
func (c *Client) WabaVerifyOtp(userCode, secret, wabaID, to, otpCode string) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/verify-otp", WabaVerifyOtpRequest{
		UserCode: userCode,
		Secret:   secret,
		WabaID:   wabaID,
		To:       to,
		OtpCode:  otpCode,
	})
}

// --- Devices ---

// ListDevices returns all registered devices for the account
func (c *Client) ListDevices(userCode, secret string) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-devices", c.authBody(userCode, secret))
}

// ListDevicesPaged returns a page of registered devices for the account
func (c *Client) ListDevicesPaged(userCode, secret string, page, limit int) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-devices", ListDevicesRequest{
		UserCode: userCode,
		Secret:   secret,
		Page:     page,
		Limit:    limit,
	})
}

// CreateDevice creates a new device
func (c *Client) CreateDevice(userCode, secret string, packageID interface{}, voucherCode string) (*Response, error) {
	return c.makeRequest("POST", "/v1/create-device", CreateDeviceRequest{
		UserCode:    userCode,
		Secret:      secret,
		PackageID:   packageID,
		VoucherCode: voucherCode,
	})
}

// ConnectDevice connects a device and returns its QR/session state
func (c *Client) ConnectDevice(userCode, secret, deviceID string) (*Response, error) {
	return c.makeRequest("POST", "/v1/connect-device", ConnectDeviceRequest{
		UserCode: userCode,
		Secret:   secret,
		DeviceID: deviceID,
	})
}

// RenewDevice renews a device subscription
func (c *Client) RenewDevice(userCode, secret, deviceID string, packageID interface{}, voucherCode string) (*Response, error) {
	return c.makeRequest("POST", "/v1/renew-device", RenewDeviceRequest{
		UserCode:    userCode,
		Secret:      secret,
		DeviceID:    deviceID,
		PackageID:   packageID,
		VoucherCode: voucherCode,
	})
}

// DeviceStatus checks the connection status of a device
func (c *Client) DeviceStatus(userCode, secret, deviceID string) (*Response, error) {
	return c.makeRequest("POST", "/v1/device-status", DeviceStatusRequest{
		UserCode: userCode,
		Secret:   secret,
		DeviceID: deviceID,
	})
}

// DeviceStatusEnhanced returns detailed status of a device
func (c *Client) DeviceStatusEnhanced(userCode, secret, deviceID string) (*Response, error) {
	return c.makeRequest("POST", "/v1/device-status-enhanced", DeviceStatusRequest{
		UserCode: userCode,
		Secret:   secret,
		DeviceID: deviceID,
	})
}

// --- User ---

// UserInfo returns account information
func (c *Client) UserInfo(userCode, secret string) (*Response, error) {
	return c.makeRequest("POST", "/v1/user-info", c.authBody(userCode, secret))
}

// --- Contacts ---

// SaveContact saves a contact to the account
func (c *Client) SaveContact(req SaveContactRequest) (*Response, error) {
	return c.makeRequest("POST", "/v1/save-contact", req)
}

// SaveContactsBulk saves up to 1000 contacts in one request
func (c *Client) SaveContactsBulk(userCode, secret string, contacts []BulkContact, deviceID string) (*Response, error) {
	return c.makeRequest("POST", "/v1/save-contacts-bulk", SaveContactsBulkRequest{
		UserCode: userCode,
		Secret:   secret,
		Contacts: contacts,
		DeviceID: deviceID,
	})
}

// --- OTP v1 ---

// GenerateOTP generates an OTP and sends it to the specified phone number
func (c *Client) GenerateOTP(req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	resp, err := c.makeRequest("POST", "/v1/generate-otp", req)
	if err != nil {
		return nil, err
	}

	var out GenerateOTPResponse
	if err := unmarshalData(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ValidateOTP validates the provided OTP
func (c *Client) ValidateOTP(req ValidateOTPRequest) (*ValidateOTPResponse, error) {
	resp, err := c.makeRequest("POST", "/v1/validate-otp", req)
	if err != nil {
		return nil, err
	}

	var out ValidateOTPResponse
	if err := unmarshalData(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// --- OTP v2 ---

// SendOtpV2 sends an OTP via the Kirimi provider, your own device, or your own WABA
func (c *Client) SendOtpV2(req SendOtpV2Request) (*Response, error) {
	return c.makeRequest("POST", "/v2/otp/send", req)
}

// VerifyOtpV2 verifies an OTP code (v2 endpoint)
func (c *Client) VerifyOtpV2(req VerifyOtpV2Request) (*Response, error) {
	return c.makeRequest("POST", "/v2/otp/verify", req)
}

// --- OTP Reverse ---

// OtpReverseCreate creates a reverse OTP token and the message the customer must send back
func (c *Client) OtpReverseCreate(userCode, secret string, req OtpReverseCreateRequest) (*Response, error) {
	req.UserCode = userCode
	req.Secret = secret
	return c.makeRequest("POST", "/v2/otp-reverse/create", req)
}

// OtpReverseStatus checks the status of a reverse OTP token
func (c *Client) OtpReverseStatus(userCode, secret, token string) (*Response, error) {
	return c.makeRequest("POST", "/v2/otp-reverse/status", OtpReverseStatusRequest{
		UserCode: userCode,
		Secret:   secret,
		Token:    token,
	})
}

// --- Deposits & Packages ---

// CreateDeposit creates a deposit payment link. Nominal minimum is 100.
func (c *Client) CreateDeposit(userCode, secret string, nominal int) (*Response, error) {
	return c.makeRequest("POST", "/v1/create-deposit", CreateDepositRequest{
		UserCode: userCode,
		Secret:   secret,
		Nominal:  nominal,
	})
}

// DepositStatus checks a deposit's status by reference
func (c *Client) DepositStatus(userCode, secret, ref string) (*Response, error) {
	return c.makeRequest("POST", "/v1/deposit-status", DepositRefRequest{
		UserCode: userCode,
		Secret:   secret,
		Ref:      ref,
	})
}

// CancelDeposit cancels an unpaid deposit
func (c *Client) CancelDeposit(userCode, secret, ref string) (*Response, error) {
	return c.makeRequest("POST", "/v1/cancel-deposit", DepositRefRequest{
		UserCode: userCode,
		Secret:   secret,
		Ref:      ref,
	})
}

// ListDeposits returns deposit list, optionally filtered by status ("", "paid", "unpaid", "expired")
func (c *Client) ListDeposits(userCode, secret, status string) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-deposits", ListDepositsRequest{
		UserCode: userCode,
		Secret:   secret,
		Status:   status,
	})
}

// ListDepositsPaged returns a page of deposits, optionally filtered by status
func (c *Client) ListDepositsPaged(userCode, secret string, page, limit int, status string) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-deposits", ListDepositsRequest{
		UserCode: userCode,
		Secret:   secret,
		Page:     page,
		Limit:    limit,
		Status:   status,
	})
}

// ListPackages returns available packages
func (c *Client) ListPackages(userCode, secret string) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-packages", c.authBody(userCode, secret))
}

// --- Health Check ---

// HealthCheck checks the API health status
func (c *Client) HealthCheck() (*HealthCheckResponse, error) {
	resp, err := c.makeGetRequest("")
	if err != nil {
		return nil, err
	}

	var out HealthCheckResponse
	if resp.Data != nil {
		_ = json.Unmarshal(resp.Data, &out)
	}
	if out.Status == "" && out.Message == "" {
		out.Message = resp.Message
	}

	return &out, nil
}

// --- Helper functions ---

// IsBasicOrProPackage checks if the package supports OTP features
func IsBasicOrProPackage(packageID int) bool {
	return packageID == PackageBasic1 || packageID == PackagePro1 ||
		packageID == PackageBasic2 || packageID == PackagePro2 ||
		packageID == PackageBasic3 || packageID == PackagePro3
}

// IsMediaSupportedPackage checks if the package supports media features
func IsMediaSupportedPackage(packageID int) bool {
	return packageID != PackageFree
}

// IsFreePackage checks if the package is free (text only with watermark)
func IsFreePackage(packageID int) bool {
	return packageID == PackageFree
}
