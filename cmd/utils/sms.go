package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/0xAckerMan/Savanah/internal/data"
	"github.com/edwinwalela/africastalking-go/pkg/sms"
)

func Sendmessage (to string, order *data.Order) {
	// Define Africa's Talking SMS client
	client := &sms.Client{
		ApiKey:    os.Getenv("AT_API_KEY"),
		Username:  os.Getenv("AT_USERNAME"),
		IsSandbox: true,
	}

	// Define a request for the Bulk SMS request
	bulkRequest := &sms.BulkRequest{
		To:            []string{"+254723285857"},
		Message:       fmt.Sprintf("hello, %s, your order for %s has been received and is being processed", order.CustomerID, order.Product.ProductName),
		From:          "",
		BulkSMSMode:   true,
		RetryDuration: time.Hour,
	}

	// Send SMS to the defined recipients
	response, err := client.SendBulk(bulkRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.Message)
}