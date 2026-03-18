# Contract-First vs Code-First Design

---

## The Core Question

```mermaid
graph TB
    Q["**Which comes first?**"]

    CF["📝 **CONTRACT-FIRST**<br/>Write the spec (OpenAPI YAML)<br/>Then generate + implement code"]
    COF["💻 **CODE-FIRST**<br/>Write the Go code<br/>Spec is generated from it"]

    Q --> CF
    Q --> COF
```

> Both approaches produce an OpenAPI spec and a working API. The difference is **which artifact drives the other**.

---

## Contract-First: Spec Before Code

```mermaid
graph TB
    subgraph Step1["1. DESIGN"]
        S1["📝 Write openapi.yaml<br/>Endpoints, request/response shapes, errors"]
        S2["👥 Review with all stakeholders<br/>Frontend, Mobile, QA, Backend"]
        S1 --> S2
    end

    subgraph Step2["2. GENERATE"]
        G1["⚙️ oapi-codegen / openapi-generator"]
        G2["📦 types.gen.go — structs"]
        G3["📋 server.gen.go — ServerInterface{}"]
        G1 --> G2
        G1 --> G3
    end

    subgraph Step3["3. IMPLEMENT"]
        I1["✍️ Implement ServerInterface in Go"]
        I2["🔨 Compiler enforces the contract<br/>Missing method = build error"]
        I1 --> I2
    end

    subgraph Step4["4. PARALLEL WORK"]
        direction LR
        FE["🖥️ Frontend<br/>Build UI against mock"]
        BE["⚙️ Backend<br/>Implement handlers"]
        QA["🧪 QA<br/>Write tests from spec"]
        MOB["📱 Mobile<br/>Generate SDK"]
    end

    S2 --> G1
    G3 --> I1
    S2 --> FE
    S2 --> BE
    S2 --> QA
    S2 --> MOB
```

> All teams work in parallel from the day the spec is agreed. No waiting for a working backend.

---

## Code-First: Code Before Spec

```mermaid
graph TB
    subgraph Step3["3. ZERO DRIFT"]
        Z1["✏️ Add a field to a struct"]
        Z2["🔄 Spec updates automatically on restart"]
        Z3["0️⃣ Docs can never be out of date"]
        Z1 --> Z2 --> Z3
    end

    subgraph Step2["2. AUTO-GENERATION"]
        F1["⚙️ Huma / Fuego<br/>Reads Go types via reflection"]
        F2["📄 /openapi.json<br/>Full OpenAPI 3.1 spec"]
        F3["📖 /docs<br/>Swagger UI"]
        F1 --> F2
        F1 --> F3
    end

    subgraph Step1["1. WRITE CODE"]
        C1["📦 Go structs with JSON tags"]
        C2["✍️ Handler functions<br/>Register with huma.Get / huma.Post"]
        C1 --> C2
    end

    C2 --> F1
```

> The code IS the source of truth. The spec is always a live reflection of what's deployed.

---

## Parallel Development: Why Contract-First Enables It

```mermaid
sequenceDiagram
    autonumber
    participant PM as 📋 Product
    participant BE as ⚙️ Backend
    participant FE as 🖥️ Frontend
    participant QA as 🧪 QA
    participant MOB as 📱 Mobile

    PM->>BE: Define API requirements
    BE->>BE: Write openapi.yaml
    BE->>PM: Review spec
    PM->>FE: Spec approved — start building
    PM->>QA: Spec approved — write contract tests
    PM->>MOB: Spec approved — generate SDK

    Note over BE,MOB: All teams unblocked. Work in parallel.

    BE->>BE: Implement Go handlers
    FE->>FE: Build UI against mock server
    QA->>QA: Write Dredd / Schemathesis tests
    MOB->>MOB: Auto-generate SDK from spec

    BE->>FE: Backend ready — swap mock for real
    QA->>QA: Run spec conformance tests → ✅
```

---

## Spec Drift: The Code-First Advantage

```mermaid
sequenceDiagram
    autonumber
    participant Dev as 👨‍💻 Developer
    participant Code as 📦 Go Code
    participant Huma as ⚙️ Huma
    participant Spec as 📄 openapi.json

    Note over Dev,Spec: Code-First — spec cannot drift

    Dev->>Code: Add `CurrencyCode string` to PaymentResponse
    Code->>Huma: Server restarts
    Huma->>Spec: Spec auto-updated — new field visible

    Note over Dev,Spec: Contract-First — drift is possible

    Dev->>Code: Add field in Go struct
    Note over Code,Spec: ⚠️ Developer forgot to update openapi.yaml
    Note over Code,Spec: Spec now WRONG — clients break, tests fail
```

> Code-first makes drift structurally impossible. Contract-first requires discipline (and tooling like Dredd) to detect it.

---

## Decision Guide

```mermaid
graph TD
    START(["New API?"])

    Q1{"Multiple teams or<br/>external consumers?"}
    Q2{"Formal review / change<br/>control process required?"}
    Q3{"Public or partner-facing API?"}
    Q4{"Small Go team,<br/>internal service?"}

    CF["✅ **CONTRACT-FIRST**<br/>oapi-codegen<br/>Spec is the contract"]
    COF["✅ **CODE-FIRST**<br/>Huma / Fuego<br/>Code is the truth"]

    START --> Q1
    Q1 -->|"Yes — multi-team / external"| CF
    Q1 -->|"No — single team"| Q2
    Q2 -->|"Yes — compliance"| CF
    Q2 -->|"No — agile"| Q3
    Q3 -->|"Yes — public"| CF
    Q3 -->|"No — internal"| Q4
    Q4 -->|"Yes — ship fast"| COF
    Q4 -->|"No — complex domain"| CF

    CF --- CFD["Best for:<br/>Public APIs<br/>Multi-team<br/>Compliance"]
    COF --- COFD["Best for:<br/>Internal APIs<br/>Single team<br/>Rapid iteration"]

    style CFD fill:none,stroke:#d2a8ff,color:#d2a8ff
    style COFD fill:none,stroke:#3fb950,color:#3fb950
```

---

## Trade-offs at a Glance

```mermaid
graph TB

    subgraph COF2["💻 CODE-FIRST"]
        COF_A["✅ Spec can never drift from code"]
        COF_B["✅ No YAML — just Go structs"]
        COF_C["✅ Fast to iterate"]
        COF_D["⚠️ API shape emerges during coding"]
        COF_E["⚠️ Harder to align multi-team upfront"]
    end

    subgraph CF2["📝 CONTRACT-FIRST"]
        CF_A["✅ Teams work in parallel from day 1"]
        CF_B["✅ Compiler enforces spec conformance"]
        CF_C["✅ Spec is the single source of truth"]
        CF_D["⚠️ More upfront design effort"]
        CF_E["⚠️ Spec can drift if not enforced"]
    end
```

> Example: public-facing APIs use contract-first — external consumers need a stable, reviewed spec. Internal services may use code-first for speed.
