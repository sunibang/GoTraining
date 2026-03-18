# Tool Discovery

---

## What Is Tool Discovery?

```mermaid
graph TB
    subgraph Without["❌ Without Discovery — Hardcoded Knowledge"]
        W1["Agent knows tools at build time"]
        W2["New tools require redeployment"]
        W3["No runtime capability negotiation"]
        W1 --> W2 --> W3
    end

    subgraph With["✅ With Discovery — Dynamic Capability Negotiation"]
        SPACE1[" "]
        D1["Agent fetches manifest at runtime"]
        D2["New tools appear without agent changes"]
        D3["Agent adapts to available capabilities"]
        D1 --> D2 --> D3
    end

    style SPACE1 fill:none,stroke:none

    Without -->|"move to"| With
```

> Tool discovery lets agents **ask** what a server can do, rather than having that knowledge baked in at build time. Servers evolve; agents adapt automatically.

---

## The Manifest: A Server's Capability Advertisement

```mermaid
graph TB
    MANIFEST["📋 **GET /mcp/manifest**<br/>The capability advertisement"]

    subgraph Contents["Manifest Contents"]
        M1["server_info<br/>name, version, description"]
        M2["resources[]<br/>URI templates + intent + risk_profile"]
        M3["tools[]<br/>name, description, inputSchema, risk_profile, constraints"]
        M4["prompts[]<br/>name, description, arguments"]
    end

    MANIFEST --> M1
    MANIFEST --> M2
    MANIFEST --> M3
    MANIFEST --> M4
```

---

## Tool Schema: The Contract an Agent Reads

```mermaid
graph TB
    subgraph Schema["🔧 Tool Definition"]
        S1["name: transfer-funds"]
        S2["description: Move money between two accounts"]
        S3["inputSchema:<br/>  from_account: string (required)<br/>  to_account: string (required)<br/>  amount: number, min: 0.01<br/>  idempotency_key: string (required)"]
        S4["x-intent: move-funds-between-accounts"]
        S5["x-risk-profile: high — irreversible financial operation"]
        S6["x-constraints:<br/>  max_amount: daily_limit<br/>  mfa_required: true"]
        S7["x-agent-guidance: always confirm with user before calling"]
        S1 --> S2 --> S3 --> S4 --> S5 --> S6 --> S7
    end
```

> The schema is not just validation — it is the **complete contract** the agent uses to decide whether and how to call the tool.

---

## Agent Discovery Flow

```mermaid
sequenceDiagram
    autonumber
    participant Agent as 🤖 AI Agent
    participant MCP as 🔌 MCP Server

    Agent->>MCP: GET /mcp/manifest
    MCP-->>Agent: {tools: [...], resources: [...], prompts: [...]}

    Agent->>Agent: Index tools by name and intent
    Agent->>Agent: Build risk map: high / medium / low

    Note over Agent: User asks: "Transfer $200 to savings"

    Agent->>Agent: Match intent → tool: transfer-funds
    Agent->>Agent: Check risk_profile: HIGH
    Agent->>Agent: Check constraints: mfa_required=true
    Agent->>Agent: Present summary to user — await confirmation

    Agent->>MCP: POST /mcp/rpc tools/call transfer-funds {...}
    MCP-->>Agent: {"result": {"transfer_id": "txn_abc", "status": "completed"}}
```

---

## Trust Levels: Not All Tools Are Equal

```mermaid
graph TD
    TOOL(["Agent evaluates a tool"])

    R{"Risk profile?"}

    LOW["🟢 LOW<br/>read-only or reversible<br/>Call autonomously"]
    MED["🟡 MEDIUM<br/>state-changing, reversible<br/>Log intent — proceed"]
    HIGH["🔴 HIGH<br/>irreversible / financial<br/>Show summary — require explicit user approval"]
    BLOCKED["⛔ BLOCKED<br/>outside policy scope<br/>Refuse — log attempt"]

    R -->|"read-only"| LOW
    R -->|"write, reversible"| MED
    R -->|"write, irreversible"| HIGH
    R -->|"out of scope"| BLOCKED
```

> The agent does not decide trust on its own — it reads the risk signal the server declared and acts accordingly.

---

## Versioned Manifests: Tools Evolve Safely

```mermaid
graph TB
    subgraph V1["v1 Manifest"]
        V1T1["tool: get-account"]
        V1T2["tool: transfer-funds"]
    end

    subgraph V2["v2 Manifest"]
        V2T1["tool: get-account"]
        V2T2["tool: transfer-funds"]
        V2T3["tool: schedule-payment  ← NEW"]
        V2T4["tool: get-account-v1 deprecated: true  ← SUNSET SIGNAL"]
    end

    V1 -->|"server deploys new version"| V2
    V2 -->|"agent re-fetches manifest"| Agent["🤖 Agent auto-discovers schedule-payment"]
```

> Manifests should be versioned. Agents that re-fetch the manifest on each session automatically gain new capabilities and respect deprecation signals — no redeployment needed.

---

## Discovery Failure Modes

```mermaid
graph TB
    subgraph Failures["Common Tool Discovery Failures"]
        F1["❌ Missing inputSchema<br/>Agent cannot validate arguments<br/>→ always provide full JSON Schema"]
        F2["❌ No risk_profile<br/>Agent assumes safe — may call destructive tools freely<br/>→ always declare risk level"]
        F3["❌ Ambiguous description<br/>Agent picks the wrong tool<br/>→ use precise, action-oriented names and descriptions"]
        F4["❌ No deprecation signal<br/>Agent keeps calling removed tools<br/>→ use deprecated: true with migration hints"]
    end
```

---

## Full Discovery + Execution Flow

```mermaid
sequenceDiagram
    autonumber
    participant Agent as 🤖 AI Agent
    participant MCP as 🔌 MCP Server (Go)
    participant API as 🚪 Bank API

    Agent->>MCP: GET /mcp/manifest
    MCP-->>Agent: Full tool + resource + prompt catalogue

    Note over Agent: Session begins — manifest cached

    Agent->>MCP: POST /mcp/rpc-resources/read-bank://accounts/acc_123
    MCP->>API: GET /api/v1/accounts/acc_123
    API-->>MCP: {"id":"acc_123","balance":1500.00}
    MCP-->>Agent: {"result":{"id":"acc_123","balance":1500.00}}

    Note over Agent: User confirms transfer

    Agent->>MCP: POST /mcp/rpc-tools/call-transfer-funds
    MCP->>API: POST /api/v1/transfers
    API-->>MCP: 201 Created {"transfer_id":"txn_abc"}
    MCP-->>Agent: {"result":{"transfer_id":"txn_abc","status":"completed"}}

    Note over Agent,API: Discovery → Read → Act. The full agentic loop.
```
