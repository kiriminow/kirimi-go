package main

import (
	"fmt"
	"log"

	kirimi "github.com/yolk/kirimi-go"
)

func main() {
	// Membuat client baru
	client := kirimi.NewClient()

	// Atau dengan custom base URL
	// client := kirimi.NewClientWithBaseURL("https://custom-api.kirimi.id")

	// Set timeout jika diperlukan
	// client.SetTimeout(60 * time.Second)

	// Contoh penggunaan Health Check
	fmt.Println("=== Health Check ===")
	healthResp, err := client.HealthCheck()
	if err != nil {
		log.Printf("Health check error: %v", err)
	} else {
		fmt.Printf("Health Status: %+v\n", healthResp)
	}

	// Contoh penggunaan Generate OTP
	fmt.Println("\n=== Generate OTP ===")
	otpReq := kirimi.GenerateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		Secret:   "your-secret-key",
	}

	otpResp, err := client.GenerateOTP(otpReq)
	if err != nil {
		log.Printf("Generate OTP error: %v", err)
	} else {
		fmt.Printf("OTP Response: %+v\n", otpResp)
	}

	// Contoh penggunaan Validate OTP
	fmt.Println("\n=== Validate OTP ===")
	validateReq := kirimi.ValidateOTPRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Phone:    "628123456789",
		OTP:      "123456",
		Secret:   "your-secret-key",
	}

	validateResp, err := client.ValidateOTP(validateReq)
	if err != nil {
		log.Printf("Validate OTP error: %v", err)
	} else {
		fmt.Printf("Validate Response: %+v\n", validateResp)
	}

	// Contoh penggunaan Send Message (text only)
	fmt.Println("\n=== Send Text Message ===")
	msgReq := kirimi.SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "Hello from Kirimi Go SDK!",
		Secret:   "your-secret-key",
	}

	msgResp, err := client.SendMessage(msgReq)
	if err != nil {
		log.Printf("Send message error: %v", err)
	} else {
		fmt.Printf("Message Response: %+v\n", msgResp)
	}

	// Contoh penggunaan Send Message dengan media
	fmt.Println("\n=== Send Message with Media ===")
	msgWithMediaReq := kirimi.SendMessageRequest{
		UserCode: "USER123",
		DeviceID: "DEVICE456",
		Receiver: "628987654321",
		Message:  "Check out this image!",
		Secret:   "your-secret-key",
		MediaURL: "https://example.com/image.jpg",
	}

	msgWithMediaResp, err := client.SendMessage(msgWithMediaReq)
	if err != nil {
		log.Printf("Send message with media error: %v", err)
	} else {
		fmt.Printf("Message with Media Response: %+v\n", msgWithMediaResp)
	}

	// Contoh penggunaan helper functions
	fmt.Println("\n=== Helper Functions ===")
	fmt.Printf("Package Basic1 supports OTP: %v\n", kirimi.IsBasicOrProPackage(kirimi.PackageBasic1))
	fmt.Printf("Package Free supports media: %v\n", kirimi.IsMediaSupportedPackage(kirimi.PackageFree))
	fmt.Printf("Package Free is free package: %v\n", kirimi.IsFreePackage(kirimi.PackageFree))
}