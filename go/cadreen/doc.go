// Package cadreen is the official Go SDK for the Cadreen Intelligence API.
//
// Cadreen provides intelligence-as-a-service: send messages describing what
// you want done, and Cadreen reasons, connects, recalls, and sometimes escalates
// to deliver an outcome. The SDK handles authentication, retries, idempotency,
// streaming, sandbox mode, and error classification.
//
// Quick start:
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"os"
//
//		cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"
//	)
//
//	func main() {
//		c := cadreen.NewClient(cadreen.CadreenConfig{
//			APIKey: os.Getenv("CADREEN_API_KEY"),
//		})
//
//		result, err := c.IntentInvoke(context.Background(), cadreen.IntentRequest{
//			Messages: []cadreen.IntentMessage{
//				{Role: "user", Content: "Refund invoice inv_123"},
//			},
//		})
//		if err != nil {
//			panic(err)
//		}
//
//		switch result.Type {
//		case cadreen.IntentResultDirect:
//			fmt.Println(result.Message.Content)
//		case cadreen.IntentResultBlocked:
//			fmt.Printf("Blocked: %s (policy %s)\n", result.ReasonCode, result.PolicyID)
//		case cadreen.IntentResultConnectRequired:
//			fmt.Printf("Connect required: %s\n", result.Connection.Endpoint)
//		}
//	}
//
// Sandbox mode:
//
//	c := cadreen.NewClient(cadreen.CadreenConfig{
//		Sandbox: true,
//		Fixtures: map[string]any{
//			"POST /api/v1/cadreen/intent": cadreen.IntentResult{
//				Type:    cadreen.IntentResultDirect,
//				Status:  "answered",
//				TraceID: "sandbox-trace",
//				Message: &cadreen.ResponseMessage{Role: "assistant", Content: "It's done."},
//			},
//		},
//	})
package cadreen
