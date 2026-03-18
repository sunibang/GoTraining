# Policy as Code

---

## The Problem: Hardcoded Rules Everywhere

```mermaid
graph TB
    subgraph Nightmare["❌ HARDCODED POLICY: No"]
        direction LR
        S1["⚙️ Accounts Service<br/>if role == 'admin' { allow }"]
        S2["⚙️ Payments Service<br/>if role == 'admin' { allow }"]
        S3["⚙️ Transfers Service<br/>if role == 'admin' { allow }"]
        S4["⚙️ Reports Service<br/>if role == 'admin' { allow }"]
    end

    PROBLEM["⚠️ Add a new role? Change a rule?<br/>Deploy ALL services.<br/>Rules drift. Audits fail."]

    S1 --> PROBLEM
    S2 --> PROBLEM
    S3 --> PROBLEM
    S4 --> PROBLEM
```

---

## Decouple Decision from Enforcement

```mermaid
graph TB
    subgraph Services["Your Microservices — Policy Enforcement Points"]
        direction LR
        A["⚙️ Accounts"]
        P["⚙️ Payments"]
        T["⚙️ Transfers"]
    end

    subgraph PaC["Policy as Code — Policy Decision Point"]
        POLICY["📋 policy.rego<br/>Versioned in Git<br/>Reviewed · Tested · Audited"]
        OPA["🔍 Open Policy Agent (OPA)"]
        POLICY --> OPA
    end

    A -->|"input: {user, action, resource}"| OPA
    P -->|"input: {user, action, resource}"| OPA
    T -->|"input: {user, action, resource}"| OPA

    OPA -->|"allow: true / false"| A
    OPA -->|"allow: true / false"| P
    OPA -->|"allow: true / false"| T
```

> Services **enforce** policy (PEP). OPA **decides** policy (PDP). Rules live in one place, versioned, auditable.

---

## Authorisation Middleware Flow

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant AuthN as 🔐 JWT Middleware
    participant AuthZ as 🛡️ Policy Middleware
    participant OPA as 🔍 OPA Engine
    participant Handler as ⚙️ Handler

    Client->>AuthN: DELETE /api/v1/accounts/acc_001<br/>Authorization: Bearer eyJ...

    AuthN->>AuthN: Validate JWT → UserIdentity{roles:["viewer"]}
    AuthN->>AuthZ: Forward with identity in context

    AuthZ->>AuthZ: Build PolicyRequest:<br/>{user:"alice", roles:["viewer"],<br/>action:"delete", resource:"accounts"}
    AuthZ->>OPA: POST /v1/data/bankx/authz {"input": {...}}

    OPA->>OPA: Evaluate policy.rego
    OPA-->>AuthZ: {"result": {"allow": false}}
    AuthZ-->>Client: 403 Forbidden {"code":"FORBIDDEN"}

    Note over Client,Handler: Handler is NEVER called.<br/>Policy evaluated at the middleware layer.
```

---

## OPA Rego Policy

```mermaid
graph TB
    subgraph Rego["policy.rego — Evaluated by OPA"]
        R1["package bankx.authz"]
        R2["default allow := false"]
        R3["allow if {<br/>  input.roles[_] == 'admin'<br/>}"]
        R4["allow if {<br/>  input.action == 'read'<br/>  input.roles[_] == 'viewer'<br/>}"]
        R1 --> R2 --> R3
        R2 --> R4
    end

    subgraph Benefits["Why Policy as Code?"]
        B1["✅ Version controlled in Git"]
        B2["✅ Unit-testable with opa test"]
        B3["✅ Single source of truth across services"]
        B4["✅ No service redeployment to change rules"]
        B5["✅ Auditable change history for compliance"]
    end
```

---

## Enforcement: Inside vs Outside the Application

```mermaid
graph TB
    subgraph Outside["🏗️ OUTSIDE — Infrastructure / Gateway"]
        GW["API Gateway / Service Mesh<br/>(Kong · Envoy · Istio)"]
        GW_OPA["OPA sidecar or plugin"]
        COARSE["Coarse-grained:<br/>Can service A call service B's endpoint?"]
        GW --> GW_OPA --> COARSE
    end

    subgraph Inside["⚙️ INSIDE — Application Middleware"]
        APP["Go Middleware"]
        APP_OPA["OPA HTTP call or embedded wasm"]
        FINE["Fine-grained:<br/>Can alice read account acc_001?"]
        APP --> APP_OPA --> FINE
    end

    UNIFIED["📋 Unified policy.rego<br/>Same rules, different enforcement points"]

    Outside --> Inside
    Inside --> UNIFIED
```

> Both enforcement points share **the same policy files**. Consistency across the entire stack.

---

## RBAC vs ABAC — When to Use Each

```mermaid
graph TD
    START(["Define access rule"])
    Q1{"Based only<br/>on role?"}
    Q2{"Contextual factors?<br/>time · department · classification"}
    Q3{"Cross-resource<br/>ownership rules?"}

    RBAC["✅ RBAC<br/>Simple · Fast<br/>role → permission map<br/>Easy to reason and audit"]
    ABAC["✅ ABAC<br/>Attribute-rich policies<br/>user.dept == resource.dept<br/>Powerful but complex"]
    EITHER["✅ Hybrid<br/>RBAC for coarse access<br/>ABAC for fine-grained rows"]

    START --> Q1
    Q1 -->|"Yes"| RBAC
    Q1 -->|"No"| Q2
    Q2 -->|"Yes"| Q3
    Q2 -->|"No"| RBAC
    Q3 -->|"Yes"| ABAC
    Q3 -->|"No"| EITHER
```

---

## Policy Lifecycle: From Commit to Decision

```mermaid
sequenceDiagram
    autonumber
    participant Dev as 👩‍💻 Developer
    participant Git as 📦 Git
    participant CI as ⚙️ CI Pipeline
    participant OPA as 🔍 OPA Server
    participant API as 🚪 API

    Dev->>Git: git push policy.rego
    Git->>CI: Trigger pipeline
    CI->>CI: opa test ./policies/...
    CI->>CI: opa fmt --fail
    CI->>OPA: Deploy bundle (if tests pass)
    OPA->>OPA: Hot-reload new policy — zero downtime

    API->>OPA: POST /v1/data/bankx/authz {input}
    OPA-->>API: {allow: true}

    Note over Dev,API: No API service redeployment.<br/>Policy change is live in seconds.
```

> Policy changes ship independently. No service restarts. No coordination with 10 teams.

---

## What Goes in a Policy Request

```mermaid
graph TB
    subgraph Input["PolicyRequest — Sent to OPA"]
        I1["user: 'alice'"]
        I2["roles: ['viewer', 'payments-reader']"]
        I3["action: 'delete'"]
        I4["resource: 'accounts'"]
        I5["resource_id: 'acc_001'"]
        I6["tenant: 'bankx'"]
        I7["time: '2024-01-15T09:30:00Z'"]
    end

    subgraph Output["PolicyResponse — Returned by OPA"]
        O1["allow: false"]
        O2["reason: 'viewer role cannot delete'"]
        O3["required_role: 'admin'"]
    end

    Input -->|"POST /v1/data/bankx/authz"| Output
```

> Include enough context for the policy to make a meaningful decision. Sparse inputs produce blunt policies.
