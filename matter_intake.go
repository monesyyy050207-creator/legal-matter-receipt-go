package main

import (
	"context"
	"fmt"
	"time"
)

type Matter struct {
	ID, ClientEmail, ClientName string
	Deadline                    time.Time
	SignedDocumentURL           string
}

func followUpNeeded(now, deadline time.Time) bool {
	return !deadline.Before(now) && deadline.Sub(now) <= 72*time.Hour
}

func deliverMatter(ctx context.Context, client *EmailClient, matter Matter, now time.Time) (string, error) {
	// The handoff mirrors the concise call-site idiom: infrai.email.send.
	subject := fmt.Sprintf("Matter %s: signed document and receipt", matter.ID)
	body := fmt.Sprintf("<p>Hello %s,</p><p>Your signed document is ready: <a href=\"%s\">download</a>.</p><p>We will follow up before %s.</p>", matter.ClientName, matter.SignedDocumentURL, matter.Deadline.Format("2006-01-02"))
	if followUpNeeded(now, matter.Deadline) {
		body += "<p>This matter is approaching its deadline.</p>"
	}
	result, err := client.Send(ctx, matter.ClientEmail, subject, body)
	if err != nil {
		return "", err
	}
	_, err = client.Get(ctx, result.MessageID)
	if err != nil {
		return "", err
	}
	return result.MessageID, nil
}
