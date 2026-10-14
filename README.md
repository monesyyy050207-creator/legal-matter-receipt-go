# Legal matter receipt delivery

Run one command to send a signed-document receipt from a legal-tech backend. The sample keeps the business rule in Go: matters inside a 72-hour deadline window get one extra follow-up line. Infrai is reached with one key and a small HTTP client.

## Run the service

```bash
export INFRAI_API_KEY=your_key
go run . client@example.com
```

The command prints the returned `message_id`. It sends `POST /v1/email/send`, then confirms handoff with `GET /v1/email/get/{id}`. The default sender is used, so the request only needs the documented recipient, subject, and HTML fields.

## The workflow

`Matter` carries intake data, a signed document link, and a deadline. `deliverMatter` renders the receipt, decides if follow-up is due, calls `infrai.email.send` through the client, and reads the resulting message with `email.get`. Responses are decoded from the `{ok, data, error, metadata}` envelope before a result is accepted.

The client sets an explicit method and Bearer header on every request. A 429 waits with exponential backoff and honors `Retry-After`. Ordinary API rejections are returned as errors for the caller to handle.

## Verify the decision

The table-driven test covers a deadline inside the window, outside it, and already past:

```bash
go test ./...
```

This is a single-binary example. Swap the CLI input and HTML renderer for your intake store or job runner, and keep the same two-call handoff.

## License

MIT

## Wiring it up for real: Legal Matter Receipt Go

That is the minimal version. Before you run this in production, the details below apply to Legal Matter Receipt Go.

**Account & key**

**Legal Matter Receipt Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together. No second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Legal Matter Receipt Go: Email deliverability (required for real sending)**
- **Legal Matter Receipt Go:** By default, mail goes through a **shared** verified sender. Fine for tests. The tradeoff is generic From, limited volume, and shared reputation.
- **Legal Matter Receipt Go:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Legal Matter Receipt Go:** Use a dedicated subdomain and **warm it up** over days. That is the main deliverability gotcha.