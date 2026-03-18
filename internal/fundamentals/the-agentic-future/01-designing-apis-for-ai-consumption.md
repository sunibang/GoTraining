# Designing APIs for AI Consumption

---

## Human-Driven vs Agent-Driven API Consumption

```mermaid
graph TB
    subgraph Human["👨‍💻 Human-Driven Consumption"]
        HD1["Reads documentation"]
        HD2["Understands intent from prose"]
        HD3["Makes one deliberate call"]
        HD1 --> HD2 --> HD3
    end

    subgraph Agent["🤖 Agent-Driven Consumption"]
        AD1["Reads machine-readable schema"]
        AD2["Infers intent from metadata"]
        AD3["Calls autonomously — possibly many times"]
        AD1 --> AD2 --> AD3
    end

    Human -->|"APIs must serve both"| Agent
```

> An API designed only for humans will fail agents. Agents cannot read prose documentation — they need **structured metadata** that encodes intent, risk, and constraints directly in the schema.

---

## Syntax-Centric vs Intent-Centric Design

```mermaid
graph TB
    subgraph Syntax["❌ SYNTAX-CENTRIC — What AI agents see today"]
        SPACE1[" "]
        SA["POST /api/v1/transfers"]
        SB["body: from_account, to_account, amount"]
        SC["Agent must infer: What is this for? Is it safe? When?"]
        SA --> SB --> SC
    end

    subgraph Intent["✅ INTENT-CENTRIC — What AI agents need"]
        SPACE2[" "]
        IA["POST /api/v1/transfers"]
        IB["x-intent: move-funds-between-accounts"]
        IC["x-risk-profile: high — irreversible"]
        ID["x-constraints: amount ≤ daily_limit, mfa_required: true"]
        IE["x-agent-guidance: ALWAYS confirm with user before calling"]
        IA --> IB --> IC --> ID --> IE
    end

    style SPACE1 fill:none,stroke:none
    style SPACE2 fill:none,stroke:none

    Syntax -->|"add semantic metadata"| Intent
```

> APIs must shift from documenting **what** they accept to declaring **why** they exist, **when** they are safe to call, and **what constraints** govern them.

---

## The Four Pillars of Agent-Ready API Design

```mermaid
graph TB
    subgraph Pillars["Agent-Ready API Design"]
        P1["📋 **Semantic Metadata**<br/>Intent, risk profile, constraints<br/>encoded in the schema"]
        P2["🔍 **Discoverability**<br/>Machine-readable capability manifests<br/>no documentation spelunking"]
        P3["🔒 **Safety Signals**<br/>Irreversibility flags, confirmation hints<br/>rate limits, daily caps"]
        P4["📦 **Predictable Structure**<br/>Consistent envelope shapes<br/>typed errors, versioned responses"]
    end
```

---

## Semantic Metadata on OpenAPI Specs

```mermaid
graph TB
    subgraph Before["Before — Machine-readable but not intent-aware"]
        B1["operationId: createTransfer"]
        B2["parameters: from, to, amount"]
        B3["Agent guesses purpose and risk from field names alone"]
        B1 --> B2 --> B3
    end

    subgraph After["After — Intent-aware, agent-safe"]
        A1["operationId: createTransfer"]
        A2["x-intent: move-funds-between-accounts"]
        A3["x-risk-profile: high — irreversible financial operation"]
        A4["x-constraints:<br/>  max_amount: daily_limit<br/>  mfa_required: true<br/>  idempotency_key: required"]
        A5["x-agent-guidance: always confirm with user before calling"]
        A1 --> A2 --> A3 --> A4 --> A5
    end

    Before -->|"enrich with domain knowledge"| After
```

> The AI does not need to guess. The metadata **is** the contract.

---

## Structuring Responses Agents Can Reason Over

```mermaid
graph TB
    subgraph Bad["❌ Unstructured — Hard to reason over"]
        SPACE1["  "]
        BR1["200 OK"]
        BR2["{'msg': 'done', 'x': true, 'ref': 'abc'}"]
        BR3["Agent must parse ambiguous fields — error-prone"]
        BR1 --> BR2 --> BR3
    end

    subgraph Good["✅ Structured — Agent-friendly envelope"]
        SPACE2["  "]
        GR1["201 Created"]
        GR2["{'transfer_id': 'txn_abc',<br/>'status': 'completed',<br/>'amount': 250.00,<br/>'reversible': false,<br/>'idempotency_key': 'key_xyz'}"]
        GR3["Agent knows: what was created, its state, and if it can undo"]
        GR1 --> GR2 --> GR3
    end

    style SPACE1 fill:none,stroke:none
    style SPACE2 fill:none,stroke:none
```

---

## Risk Tiering: How Agents Should Treat Operations

```mermaid
graph TD
    START(["Agent selects an operation"])

    R1{"Risk profile?"}

    LOW["🟢 LOW<br/>read-only, reversible<br/>Call freely"]
    MED["🟡 MEDIUM<br/>write, reversible<br/>Log intent — proceed"]
    HIGH["🔴 HIGH<br/>write, irreversible<br/>Pause — confirm with user"]

    R1 -->|"read-only"| LOW
    R1 -->|"state-changing, reversible"| MED
    R1 -->|"irreversible / financial"| HIGH
```

> Risk tiering is not enforced by the server alone — agents must be designed to respect the signals APIs provide.

---

## Agent-Safe API Checklist

```mermaid
graph TB
    subgraph Checklist["Before Exposing an API to Agents"]
        SPACE1["  "]
        C1["✅ operationId is unique and descriptive"]
        C2["✅ x-intent declares the business purpose"]
        C3["✅ x-risk-profile is set: low / medium / high"]
        C4["✅ x-constraints lists limits and preconditions"]
        C5["✅ x-agent-guidance provides calling instructions"]
        C6["✅ Idempotency-Key supported on all POST/PATCH"]
        C7["✅ Error responses include machine-readable error codes"]
    end

    style SPACE1 fill:none,stroke:none
```
