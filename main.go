package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: legalmail recipient@example.com")
		os.Exit(2)
	}
	client, err := NewEmailClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	matter := Matter{ID: "MAT-1042", ClientEmail: os.Args[1], ClientName: "Case team", Deadline: time.Now().Add(48 * time.Hour), SignedDocumentURL: "https://files.example.test/MAT-1042/signed.pdf"}
	id, err := deliverMatter(context.Background(), client, matter, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("receipt sent; message_id:", id)
}
