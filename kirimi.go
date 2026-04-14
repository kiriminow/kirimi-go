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
	UserCode         string `json:"user_code"`
	DeviceID         string `json:"device_id"`
	Phone            string `json:"phone"`
	Secret           string `json:"secret"`
	OtpLength        int    `json:"otp_length,omitempty"`
	OtpType          string `json:"otp_type,omitempty"`
	CustomOtpMessage string `json:"customOtpMessage,omitempty"`
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
	UserCode string `json:"user_code"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	Message  string `json:"message"`
	Secret   string `json:"secret"`
	MediaURL string `json:"media_url,omitempty"`
}

// SendMessageResponse represents the response data for sending message
type SendMessageResponse struct {
	MessageLength int    `json:"message_length"`
	MediaURL      string `json:"media_url,omitempty"`
	HasMedia      bool   `json:"has_media"`
}

// SendMessageFileRequest represents the request for sending a file via multipart/form-data
type SendMessageFileRequest struct {
	UserCode string
	DeviceID string
	Phone    string
	Secret   string
	File     io.Reader
	Filename string
	Message  string
	FileName string
}

// SendWabaMessageRequest represents the request for sending a WABA message
type SendWabaMessageRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	DeviceID string `json:"device_id"`
	Phone    string `json:"phone"`
	Message  string `json:"message"`
}

// DeviceStatusRequest is used for device-status and device-status-enhanced
type DeviceStatusRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	DeviceID string `json:"device_id"`
}

// SaveContactRequest represents the request for saving a contact
type SaveContactRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Phone    string `json:"phone"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
}

// SendOtpV2Request represents the request for /v2/otp/send
type SendOtpV2Request struct {
	UserCode      string `json:"user_code"`
	Secret        string `json:"secret"`
	Phone         string `json:"phone"`
	DeviceID      string `json:"device_id"`
	Method        string `json:"method,omitempty"`
	AppName       string `json:"app_name,omitempty"`
	TemplateCode  string `json:"template_code,omitempty"`
	CustomMessage string `json:"custom_message,omitempty"`
}

// VerifyOtpV2Request represents the request for /v2/otp/verify
type VerifyOtpV2Request struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Phone    string `json:"phone"`
	OtpCode  string `json:"otp_code"`
}

// BroadcastMessageRequest represents the request for broadcast-message
type BroadcastMessageRequest struct {
	UserCode string  `json:"user_code"`
	Secret   string  `json:"secret"`
	DeviceID string  `json:"device_id"`
	Phones   string  `json:"phones"`
	Message  string  `json:"message"`
	Delay    float64 `json:"delay,omitempty"`
}

// ListDepositsRequest represents the request for list-deposits
type ListDepositsRequest struct {
	UserCode string `json:"user_code"`
	Secret   string `json:"secret"`
	Status   string `json:"status,omitempty"`
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
	_ = w.WriteField("phone", req.Phone)
	if req.Message != "" {
		_ = w.WriteField("message", req.Message)
	}
	if req.FileName != "" {
		_ = w.WriteField("fileName", req.FileName)
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

// --- WABA ---

// SendWabaMessage sends a message via WhatsApp Business API (Meta Cloud API)
func (c *Client) SendWabaMessage(req SendWabaMessageRequest) (*Response, error) {
	return c.makeRequest("POST", "/v1/waba/send-message", req)
}

// --- Devices ---

// ListDevices returns all registered devices for the account
func (c *Client) ListDevices(userCode, secret string) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-devices", c.authBody(userCode, secret))
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

// SendOtpV2 sends an OTP via WABA template or device (v2 endpoint)
func (c *Client) SendOtpV2(req SendOtpV2Request) (*Response, error) {
	return c.makeRequest("POST", "/v2/otp/send", req)
}

// VerifyOtpV2 verifies an OTP code (v2 endpoint)
func (c *Client) VerifyOtpV2(req VerifyOtpV2Request) (*Response, error) {
	return c.makeRequest("POST", "/v2/otp/verify", req)
}

// --- Broadcast ---

// BroadcastMessage sends a message to multiple recipients (phones as comma-separated string)
func (c *Client) BroadcastMessage(req BroadcastMessageRequest) (*Response, error) {
	return c.makeRequest("POST", "/v1/broadcast-message", req)
}

// --- Deposits & Packages ---

// ListDeposits returns deposit list, optionally filtered by status ("", "paid", "unpaid", "expired")
func (c *Client) ListDeposits(userCode, secret, status string) (*Response, error) {
	return c.makeRequest("POST", "/v1/list-deposits", ListDepositsRequest{
		UserCode: userCode,
		Secret:   secret,
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
