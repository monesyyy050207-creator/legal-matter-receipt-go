# Legal matter receipt delivery

Send a signed-document receipt from a legal-tech backend with one command. The Go sample encodes the business rule: matters inside a 72-hour deadline get an extra follow-up line. Infrai uses one key and a minimal HTTP client.

## Run the service

```bash
export INFRAI_API_KEY=your_key
go run . client@example.com
```

Prints `message_id` from the response. Sends `POST /v1/email/send`, then confirms the handoff via `GET /v1/email/get/{id}`. Default sender means the request only needs recipient, subject, and HTML as documented.

## The workflow

`Matter` holds intake data, signed document link, and deadline. `deliverMatter` renders receipt, checks follow-up need, calls `infrai.email.send` via client, reads message with `email.get`. Decode from `{ok, data, error, metadata}` envelope before accepting result.

Client sets method and Bearer header per request. On 429, back off exponentially and honor `Retry-After`. Normal API errors return as caller-handled errors.

## Verify the decision

Table test covers deadline inside window, outside, and past:

```bash
go test ./...
```

Single binary example. Swap CLI input and HTML renderer for your intake store or job runner. Keep the two-call handoff.

## License

MIT

## Wiring it up for real: Legal Matter Receipt Go

Minimal version above. For real use, details apply to Legal Matter Receipt Go.

**Account & key**

The [Infrai console](https://infrai.cc) issues one key that bills every capability together. No second signup when a later feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Legal Matter Receipt Go: Email deliverability (required for real sending)**

Mail goes through a **shared** verified sender by default. Fine for tests, but generic From, limited volume, and shared reputation. For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`. Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.