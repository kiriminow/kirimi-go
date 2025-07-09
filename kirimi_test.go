package kirimi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test helper function to create a test server
func createTestServer(statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}))
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client.BaseURL != DefaultBaseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", DefaultBaseURL, client.BaseURL)
	}
	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", client.HTTPClient.Timeout)
	}
}

func TestNewClientWithBaseURL(t *testing.T) {
	customURL := "https://custom.api.com"
	client := NewClientWithBaseURL(customURL)
	if client.BaseURL != customURL {
		t.Errorf("Expected BaseURL to be %s, got %s", customURL, client.BaseURL)
	}
}

func TestSetTimeout(t *testing.T) {
	client := NewClient()
	newTimeout := 60 * time.Second
	client.SetTimeout(newTimeout)
	if client.HTTPClient.Timeout != newTimeout {
		t.Errorf("Expected timeout to be %v, got %v", newTimeout, client.HTTPClient.Timeout)
	}
}

func TestHealthCheck(t *testing.T) {
	// Test successful health check
	successResponse := Response{
		Success: true,
		Data:    map[string]interface{}{},
		Message: "Kirimi API v1",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	resp, err := client.HealthCheck()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Message != "Kirimi API v1" {
		t.Errorf("Expected message 'Kirimi API v1', got %s", resp.Message)
	}
}

func TestGenerateOTP(t *testing.T) {
	// Test successful OTP generation
	successResponse := Response{
		Success: true,
		Data: map[string]interface{}{
			"phone":      "628123456789",
			"message":    "OTP berhasil dikirim",
			"expires_in": "5 menit",
		},
		Message: "OTP berhasil digenerate dan dikirim",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := GenerateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		Secret:   "test-secret",
	}

	resp, err := client.GenerateOTP(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Phone != "628123456789" {
		t.Errorf("Expected phone '628123456789', got %s", resp.Phone)
	}
	if resp.ExpiresIn != "5 menit" {
		t.Errorf("Expected expires_in '5 menit', got %s", resp.ExpiresIn)
	}
}

func TestValidateOTP(t *testing.T) {
	// Test successful OTP validation
	successResponse := Response{
		Success: true,
		Data: map[string]interface{}{
			"phone":       "628123456789",
			"verified":    true,
			"verified_at": "2024-01-15T10:30:00.000Z",
		},
		Message: "OTP berhasil divalidasi",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := ValidateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		OTP:      "123456",
		Secret:   "test-secret",
	}

	resp, err := client.ValidateOTP(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.Phone != "628123456789" {
		t.Errorf("Expected phone '628123456789', got %s", resp.Phone)
	}
	if !resp.Verified {
		t.Errorf("Expected verified to be true, got %v", resp.Verified)
	}
}

func TestSendMessage(t *testing.T) {
	// Test successful message sending
	successResponse := Response{
		Success: true,
		Data: map[string]interface{}{
			"message_length": 25,
			"media_url":      "https://example.com/image.jpg",
			"has_media":      true,
		},
		Message: "Berhasil mengirim pesan dengan media",
	}

	server := createTestServer(200, successResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "Hello from test!",
		Secret:   "test-secret",
		MediaURL: "https://example.com/image.jpg",
	}

	resp, err := client.SendMessage(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp.MessageLength != 25 {
		t.Errorf("Expected message_length 25, got %d", resp.MessageLength)
	}
	if !resp.HasMedia {
		t.Errorf("Expected has_media to be true, got %v", resp.HasMedia)
	}
}

func TestSendMessageTooLong(t *testing.T) {
	client := NewClient()
	req := SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  string(make([]byte, MaxMessageLength+1)), // Message too long
		Secret:   "test-secret",
	}

	_, err := client.SendMessage(req)

	if err == nil {
		t.Error("Expected error for message too long, got nil")
	}
}

func TestAPIError(t *testing.T) {
	// Test API error response
	errorResponse := Response{
		Success: false,
		Data:    map[string]interface{}{},
		Message: "Parameter tidak boleh kosong",
	}

	server := createTestServer(400, errorResponse)
	defer server.Close()

	client := NewClientWithBaseURL(server.URL)
	req := GenerateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		Secret:   "test-secret",
	}

	_, err := client.GenerateOTP(req)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("Expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Parameter tidak boleh kosong" {
		t.Errorf("Expected message 'Parameter tidak boleh kosong', got %s", apiErr.Message)
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test IsBasicOrProPackage
	if !IsBasicOrProPackage(PackageBasic1) {
		t.Error("Expected PackageBasic1 to support OTP")
	}
	if !IsBasicOrProPackage(PackagePro1) {
		t.Error("Expected PackagePro1 to support OTP")
	}
	if IsBasicOrProPackage(PackageFree) {
		t.Error("Expected PackageFree to not support OTP")
	}
	if IsBasicOrProPackage(PackageLite1) {
		t.Error("Expected PackageLite1 to not support OTP")
	}

	// Test IsMediaSupportedPackage
	if IsMediaSupportedPackage(PackageFree) {
		t.Error("Expected PackageFree to not support media")
	}
	if !IsMediaSupportedPackage(PackageLite1) {
		t.Error("Expected PackageLite1 to support media")
	}
	if !IsMediaSupportedPackage(PackageBasic1) {
		t.Error("Expected PackageBasic1 to support media")
	}
	if !IsMediaSupportedPackage(PackagePro1) {
		t.Error("Expected PackagePro1 to support media")
	}

	// Test IsFreePackage
	if !IsFreePackage(PackageFree) {
		t.Error("Expected PackageFree to be free package")
	}
	if IsFreePackage(PackageLite1) {
		t.Error("Expected PackageLite1 to not be free package")
	}
	if IsFreePackage(PackageBasic1) {
		t.Error("Expected PackageBasic1 to not be free package")
	}
	if IsFreePackage(PackagePro1) {
		t.Error("Expected PackagePro1 to not be free package")
	}
}

func TestConstants(t *testing.T) {
	if DefaultBaseURL != "https://api.kirimi.id" {
		t.Errorf("Expected DefaultBaseURL to be 'https://api.kirimi.id', got %s", DefaultBaseURL)
	}
	if APIVersion != "v1" {
		t.Errorf("Expected APIVersion to be 'v1', got %s", APIVersion)
	}
	if ContentType != "application/json" {
		t.Errorf("Expected ContentType to be 'application/json', got %s", ContentType)
	}
	if MaxMessageLength != 1200 {
		t.Errorf("Expected MaxMessageLength to be 1200, got %d", MaxMessageLength)
	}
}