# Implementation Mode

## Description

Code-focused execution mode. Minimal prose, maximum code output. For when the plan is decided and it's time to build.

## When to Use

- Executing a finalized plan
- Writing code from clear specifications
- Sprint coding sessions
- After brainstorm/planning has converged on an approach

---

## Behavior

### Communication
- Lead with code, not explanations
- Only speak to communicate decisions, blockers, or questions
- No design exploration — the design phase is over
- No alternatives — execute the decided approach
- Report completion status tersely

### Problem Solving
- Linear execution — follow the plan step by step
- Fix forward — if something doesn't work, fix it and move on
- Ask only when truly blocked — don't ask for preferences that were already decided
- Test as you go — write test alongside implementation, don't defer

### Code Output
- Complete, runnable code — no pseudocode, no placeholders, no `// TODO`
- Full file content for new files
- Precise edits for existing files
- Import statements included
- Error handling included on first pass (not "we'll add this later")

### Workflow
1. Read existing code in the target area
2. Write implementation
3. Write tests
4. Run lint + tests
5. Fix issues
6. Report: `✓ [package/file] — [what was done]`

---

## Example

### Non-Implementation Mode
```
Let me think about how to implement the rating submission use case.
We'll need to consider the self-healing average pattern, where we
recalculate the average on every submission rather than maintaining
a running counter. This approach is more resilient to data corruption...

Here's my proposed approach:
1. First, we validate the input...
2. Then we fetch the existing ratings...
[continues for 500 tokens]
```

### Implementation Mode
```
`internal/usecase/rating/submit.go`:

[complete code]

`internal/usecase/rating/submit_test.go`:

[complete tests]

✓ rating/submit — submit with self-healing avg, 6 test cases, lint clean
```

---

## 2nd-Order Implementation Checks

While coding, silently verify:
- Does this change require a corresponding mobile app update?
- Does this change break any existing event subscribers?
- Does this change modify a DynamoDB key schema that other handlers depend on?
- Does this change introduce a new failure mode in the booking state machine?

If any answer is YES → break out of implementation mode briefly to flag it.

## Combines Well With
- Token-efficient mode — stack both for maximum throughput
- `/spawn` — spawn implementation-mode agents for parallel coding
- `/ship` — after implementation, ship the result
