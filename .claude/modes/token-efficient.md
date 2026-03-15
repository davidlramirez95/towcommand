# Token-Efficient Mode

## Description

Compressed output mode for high-volume work sessions. Reduces token usage by 30-70% while maintaining code quality. Use during batch sprint execution, large refactors, or when approaching context limits.

## When to Use

- Batch agent spawning (sprint execution)
- Large refactors touching many files
- High-volume test writing
- When context window is getting full
- Cost-sensitive sessions

---

## Behavior

### Communication Rules
- **No preamble** — skip "I'll now...", "Let me...", "Here's what I'm doing..."
- **No summaries** — don't recap what was just done; the diff speaks
- **No explanations for obvious code** — only explain non-obvious decisions
- **No trailing status reports** — the tool output is the status
- **Code over prose** — show the implementation, not a description of it
- **Abbreviate file references** — `booking.go:42` not "in the file internal/usecase/booking/booking.go at line 42"

### Output Compression
| Standard | Token-Efficient |
|----------|----------------|
| "I'll create a new handler function that processes the webhook payload and validates the signature before..." | `cmd/payment-webhook/main.go` — webhook handler with sig validation |
| Full file content with comments | Only the changed/new sections |
| Step-by-step explanation | Numbered action list |
| Detailed test descriptions | Test names + assertion summary |

### What to KEEP (never compress)
- Error messages and failure diagnostics
- Architecture decisions that affect other components
- 2nd-order warnings (things that could break elsewhere)
- Security concerns
- Questions that need user input

### Savings Targets
| Content Type | Target Savings |
|---|---|
| Code explanations | 60-70% |
| Status updates | 80-90% |
| File operations | 40-50% |
| Error handling | 0% (never compress) |
| Architecture decisions | 20-30% |

---

## Example

### Standard Mode
```
I'll now create the booking cancel handler. This handler needs to:
1. Parse the booking ID from the path parameters
2. Validate the JWT token to get the user ID
3. Look up the booking to verify ownership
4. Check that the booking is in a cancellable state
5. Update the booking status to CANCELLED
6. Publish a booking.cancelled event

Let me implement this step by step...

[200 lines of explanation + code]

The handler is now complete. I've added proper error handling for each step,
including cases where the booking doesn't exist, the user doesn't own the
booking, or the booking is in a non-cancellable state. The event is published
after the status update to ensure consistency.
```

### Token-Efficient Mode
```
`cmd/booking-cancel/main.go`:
[code]

`internal/usecase/booking/cancel.go`:
[code]

`internal/usecase/booking/cancel_test.go`:
[code]

Tests: 5/5 pass. Lint clean.
```

---

## Activation

```
/mode token-efficient
```

Or per-command:
```
/ship --format=concise
/spawn --format=ultra "implement handler"
```

## Combines Well With
- `/spawn` — batch agent dispatch with minimal overhead
- Implementation mode — code-focused + compressed = maximum efficiency
- Sprint execution workflow — ship more PRs per session
