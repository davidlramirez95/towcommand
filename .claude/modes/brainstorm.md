# Brainstorm Mode

## Description

Creative exploration mode for ideation, architecture decisions, and design sessions. Optimized for divergent thinking, trade-off analysis, and 2nd-order consequence mapping.

## When to Use

- Initial feature exploration and scoping
- Architecture and system design decisions
- Problem definition and approach selection
- When stuck on implementation approach
- Evaluating technology choices
- Pre-sprint planning sessions

---

## Behavior

### Communication
- Ask probing questions before converging — one at a time, validate before next
- Present 3+ alternatives for every decision point
- Explore edge cases verbally before coding
- Use "what if" and "what breaks when" scenarios
- Map 2nd-order effects of each option explicitly

### Problem Solving
- Divergent thinking first — generate options broadly
- Delay convergence — resist jumping to the first viable solution
- Consider unconventional approaches (different data model, different service boundary, different UX flow)
- Map trade-offs with explicit dimensions: complexity, performance, maintainability, cost, time-to-ship

### 2nd-Order Thinking (ALWAYS in Brainstorm)
For every proposed approach, answer:
1. **What does this make easier later?** (positive cascade)
2. **What does this make harder later?** (hidden cost)
3. **What assumption does this bake in that we can't easily change?** (lock-in risk)
4. **Who else is affected by this choice?** (blast radius — other services, mobile app, ops)
5. **What failure mode does this introduce?** (new ways things can break)

### Output Format
- Decision matrices with weighted criteria
- Pro/con lists with severity ratings (minor/moderate/critical)
- ASCII/Mermaid diagrams for architecture options
- Explicit "I recommend X because..." with reasoning chain
- "2nd-order effects" section for each major option

---

## Example Behaviors

### Architecture Decision
```
Before we pick an approach, let me map the design space:

**Option A: Event-driven saga (EventBridge)**
- Pros: Decoupled, auditable, matches existing patterns
- Cons: Eventually consistent, harder to debug
- 2nd-order: Makes adding new steps trivial; makes transaction rollback complex

**Option B: Synchronous orchestrator (single Lambda)**
- Pros: Simple, strongly consistent, easy to debug
- Cons: Single point of failure, harder to extend
- 2nd-order: Fast to ship now; becomes a bottleneck when we add payment retry logic

**Option C: Step Functions**
- Pros: Visual workflow, built-in retry, audit trail
- Cons: AWS lock-in, cold start latency, cost at scale
- 2nd-order: Great for complex flows; overkill for simple ones, adds infra complexity

**Recommendation**: Option A — aligns with our existing EventBridge pattern and the booking state machine. The eventual consistency trade-off is acceptable for this use case because [reason].

What's your instinct? Should we explore any of these deeper?
```

### Question-First Design
```
Before I can recommend an approach, I need to understand constraints:

1. **Volume**: How many concurrent bookings do we expect at peak?
   (This determines whether we need Redis caching or DDB is enough)

Let me know, and I'll ask the next question based on your answer.
```

---

## TowCommand-Specific Brainstorm Patterns

When brainstorming for this project, always consider:
- **PH market constraints**: Intermittent connectivity, low-end devices, GCash/Maya payment rails
- **Safety implications**: Every design choice in the booking/SOS flow has safety consequences
- **DynamoDB access patterns**: New features must work within single-table design + 5 GSIs
- **Mobile-first**: Backend decisions must consider mobile UX impact (latency, offline, battery)
- **Regulatory**: LTFRB compliance, data privacy (PH Data Privacy Act)

## Combines Well With
- `/plan` — brainstorm first, then plan in detail
- `/index` — load project structure before brainstorming
- `/spawn` — after converging, spawn parallel implementation
