# Model Context Protocol (MCP)

---

## The Problem MCP Solves

```mermaid
graph TB
    subgraph Before["Before MCP — N × M Integration Problem"]
        C1["🤖 Claude"]
        C2["🤖 GPT"]
        C3["🤖 Gemini"]
        S1["🗄️ Database"]
        S2["📁 File System"]
        S3["🌐 Web API"]
        C1 --> S1
        C1 --> S2
        C1 --> S3
        C2 --> S1
        C2 --> S2
        C2 --> S3
        C3 --> S1
        C3 --> S2
        C3 --> S3
    end

    subgraph After["After MCP — One Standard Protocol"]
        MC1["🤖 Claude"]
        MC2["🤖 GPT"]
        MC3["🤖 Gemini"]
        MCP["🔌 Model Context Protocol"]
        MS1["🗄️ Database Server"]
        MS2["📁 File System Server"]
        MS3["🌐 Web API Server"]
        MC1 --> MCP
        MC2 --> MCP
        MC3 --> MCP
        MCP --> MS1
        MCP --> MS2
        MCP --> MS3
    end
```

> Just as USB-C standardises how devices connect to peripherals, **MCP standardises how AI agents connect to application context** — regardless of the underlying LLM or backend.

---

## MCP: The USB-C Port for AI

```mermaid
graph TB
    subgraph Clients["AI Clients — any LLM"]
        CLAUDE["🤖 Claude"]
        GPT["🤖 GPT"]
        GEM["🤖 Gemini"]
        CUSTOM["🤖 Custom Agent"]
    end

    MCP["🔌 **Model Context Protocol**<br/>One standard. Any client. Any server."]

    subgraph Primitives["MCP Server Primitives"]
        RES["📄 Resources<br/>Read-only data the agent can query"]
        TOOLS["🔧 Tools<br/>Side-effecting actions the agent can invoke"]
        PROMPTS["💬 Prompts<br/>Reusable prompt templates"]
    end

    CLAUDE --> MCP
    GPT --> MCP
    GEM --> MCP
    CUSTOM --> MCP

    MCP --> RES
    MCP --> TOOLS
    MCP --> PROMPTS
```

---

## MCP Primitives in Detail

```mermaid
graph TB
    subgraph Resources["📄 Resources — READ ONLY"]
        R1["URI-addressed data"]
        R2["bank://accounts/{id}"]
        R3["risk_profile: read-only<br/>No side effects — safe to call freely"]
        R1 --> R2 --> R3
    end

    subgraph Tools["🔧 Tools — ACT"]
        T1["Named functions with typed parameters"]
        T2["transfer-funds, get-account, cancel-payment"]
        T3["risk_profile: low / medium / high<br/>May have irreversible side effects"]
        T1 --> T2 --> T3
    end

    subgraph Prompts["💬 Prompts — GUIDE"]
        P1["Parameterised prompt templates"]
        P2["summarise-transactions, diagnose-failed-payment"]
        P3["Reusable reasoning patterns<br/>Consistent agent behaviour across calls"]
        P1 --> P2 --> P3
    end
```

---

## Anatomy of an MCP Server

```mermaid
graph TB
    SERVER["🔌 **MCP Server**<br/>Your Go API"]

    subgraph R["📄 Resources — READ"]
        R1["bank://accounts/{id}<br/>intent: retrieve-account-balance<br/>risk_profile: read-only"]
        R2["bank://transactions/{id}<br/>intent: audit-payment-trail<br/>risk_profile: read-only"]
    end

    subgraph T["🔧 Tools — ACT"]
        T1["get-account<br/>intent: fetch-account-details<br/>risk_profile: low"]
        T2["transfer-funds<br/>intent: move-money-between-accounts<br/>risk_profile: high<br/>constraints: mfa-required, within-daily-limit"]
    end

    subgraph P["💬 Prompts — GUIDE"]
        P1["summarise-transactions<br/>intent: generate-spending-summary"]
        P2["diagnose-failed-payment<br/>intent: root-cause-analysis"]
    end

    SERVER --> R
    SERVER --> T
    SERVER --> P
```

---

## MCP Wire Protocol: JSON-RPC 2.0

```mermaid
graph TB
    subgraph JSONRPC["JSON-RPC 2.0 — The MCP Wire Protocol"]
        REQ["📨 Request<br/>{jsonrpc: '2.0', id: 1,<br/>method: 'tools/call',<br/>params: {name, arguments}}"]
        RES_OK["✅ Success Response<br/>{jsonrpc: '2.0', id: 1,<br/>result: {...}}"]
        RES_ERR["❌ Error Response<br/>{jsonrpc: '2.0', id: 1,<br/>error: {code, message}}"]
        REQ -->|"success path"| RES_OK
        REQ -->|"failure path"| RES_ERR
    end

    subgraph Stateless["Stateless by Design"]
        SL1["Each request is fully self-contained"]
        SL2["No session. No server memory between calls."]
        SL3["Any instance handles any request"]
        SL4["Scales horizontally — just like REST"]
        SL1 --> SL2 --> SL3 --> SL4
    end

    JSONRPC --> Stateless
```

> MCP tool calls are **stateless** — each JSON-RPC request carries all the context it needs. A natural fit for Go's concurrent, stateless HTTP server model.

---

## MCP Request Lifecycle

```mermaid
sequenceDiagram
    autonumber
    participant Agent as 🤖 AI Agent
    participant MCP as 🔌 MCP Server (Go)
    participant API as 🚪 Backend API

    Agent->>MCP: GET /mcp/manifest
    MCP-->>Agent: Resources + Tools + Prompts<br/>with full semantic metadata

    Agent->>Agent: Selects tool: transfer-funds
    Agent->>Agent: Reads: risk_profile=high, mfa_required=true
    Agent->>Agent: Presents details to user — awaits explicit approval

    Agent->>MCP: POST /mcp/rpc<br/>{"method":"tools/call","params":{"name":"transfer-funds","arguments":{...}}}
    MCP->>API: POST /api/v1/transfers
    API-->>MCP: 201 Created {"transfer_id":"txn_abc"}
    MCP-->>Agent: {"result":{"transfer_id":"txn_abc","status":"completed"}}

    Note over Agent,API: The agent never guessed. It read intent, checked constraints, confirmed risk.
```

---

## MCP vs Traditional REST: What Changes

```mermaid
graph TB
    subgraph REST["🌐 Traditional REST API"]
        RA["Designed for human developers"]
        RB["Docs live in Swagger UI / Confluence"]
        RC["No built-in risk signals"]
        RD["Agent must reverse-engineer intent"]
        RA --> RB --> RC --> RD
    end

    subgraph MCP["🔌 MCP Server"]
        MA["Designed for AI agents"]
        MB["Manifest is the docs — machine-readable"]
        MC["Risk profiles and constraints are first-class"]
        MD["Agent reads intent directly — no guessing"]
        MA --> MB --> MC --> MD
    end
```

> MCP does not replace REST. It **wraps** your existing REST API with a machine-readable capability layer that agents can safely discover and invoke.
