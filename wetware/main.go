package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/Dstack-TEE/dstack/sdk/go/dstack"
)

// accept only POST requests. Copy the request body to the response body.
func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	io.Copy(w, r.Body)
	defer r.Body.Close()
}

func dstackCheck() error {
	client := dstack.NewDstackClient(
		// dstack.WithEndpoint("http://localhost"),
		dstack.WithLogger(slog.Default()),
	)

	// Get information about the dstack client instance
	info, err := client.Info(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(info.AppID)   // Application ID
	fmt.Println(info.TcbInfo) // Access TCB info directly
	// These lines from the example don't seem valid...
	// fmt.Println(info.TcbInfo.Mrtd)              // Access TCB info directly
	// fmt.Println(info.TcbInfo.EventLog[0].Event) // Access event log entries

	path := "/"
	purpose := "test" // or leave empty

	// Derive a key with optional path and purpose
	deriveKeyResp, err := client.GetKey(context.Background(), path, purpose)
	if err != nil {
		return err
	}
	fmt.Println(deriveKeyResp.Key)

	// Generate TDX quote
	tdxQuoteResp, err := client.GetQuote(context.Background(), []byte("test"))
	if err != nil {
		return err
	}
	fmt.Println(tdxQuoteResp.Quote) // 0x0000000000000000000 ...

	rtmrs, err := tdxQuoteResp.ReplayRTMRs()
	if err != nil {
		return err
	}
	fmt.Println(rtmrs) // map[0:00000000000000000 ...
	return nil
}

func main() {
	err := dstackCheck()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	http.HandleFunc("/", echoHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
