# REST vs RPC

---

## Two Philosophies of API Design

```mermaid
graph TB
    subgraph RPC["RPC — Remote Procedure Call"]
        SPACE1["  "]
        P1["⚙️ **Actions are central**<br/>Call a function remotely"]
        P2["📡 **Transport-agnostic**<br/>HTTP, TCP, gRPC"]
        P3["📋 **Strong typing**<br/>Protobuf, Thrift, JSON-RPC"]
        P1 --> P2 --> P3
    end

    subgraph REST["REST — Representational State Transfer"]
        SPACE2["  "]
        R1["📦 **Resources are central**<br/>Nouns, not verbs"]
        R2["🔀 **HTTP verbs carry intent**<br/>GET, POST, PUT, DELETE"]
        R3["🔗 **Uniform interface**<br/>Standard URL patterns"]
        R1 --> R2 --> R3
    end

    style SPACE1 fill:none,stroke:none
    style SPACE2 fill:none,stroke:none
```

> REST models the **world as resources**. RPC models it as **function calls**.

---

## REST: Resource-Oriented Design

```mermaid
graph TB
    subgraph Bad["❌ RPC-style in REST URLs"]
        B1["POST /createPayment"]
        B2["GET  /getPaymentById?id=123"]
        B3["POST /cancelPayment/123"]
        B4["POST /processRefund"]
    end

    subgraph Good["✅ REST — Nouns + HTTP Verbs"]
        SPACE1["  "]
        G1["GET    /payments          → list payments"]
        G2["POST   /payments          → create payment"]
        G3["GET    /payments/{id}     → get one"]
        G4["PUT    /payments/{id}     → update"]
        G5["DELETE /payments/{id}     → cancel"]
    end

    style SPACE1 fill:none,stroke:none
```

> The URL identifies the resource. The HTTP method identifies the operation.

---

## RPC: Action-Oriented Design

```mermaid
graph TB
    subgraph GRPC["gRPC — Google's RPC Framework"]
        PROTO["📋 payments.proto<br/>Define service + messages"]
        GEN["⚙️ protoc code generation"]
        CLIENT["📱 Generated client stub"]
        SERVER["⚙️ Generated server interface"]
        PROTO --> GEN --> CLIENT
        GEN --> SERVER
    end

    subgraph Methods["Service Methods"]
        M1["CreatePayment(CreatePaymentRequest)"]
        M2["GetPayment(GetPaymentRequest)"]
        M3["CancelPayment(CancelPaymentRequest)"]
        M4["ListPayments(ListPaymentsRequest)"]
    end

    SERVER --> M1
    SERVER --> M2
    SERVER --> M3
    SERVER --> M4
```

> gRPC: strongly typed, binary protocol, HTTP/2, bidirectional streaming — fast and efficient.

---

## REST vs gRPC: Wire Format Comparison

```mermaid
sequenceDiagram
    autonumber
    participant Client as 📱 Client
    participant REST as 🌐 REST API
    participant gRPC as ⚡ gRPC API

    Note over Client,gRPC: Same operation: fetch payment

    Client->>REST: GET /payments/pay_123<br/>Accept: application/json
    REST-->>Client: 200 OK<br/>{"id":"pay_123","amount":250.00,...}<br/>~200 bytes JSON text

    Client->>gRPC: GetPayment(id: "pay_123")<br/>Binary Protobuf over HTTP/2
    gRPC-->>Client: Payment{id:"pay_123",amount:25000,...}<br/>~40 bytes binary
```

> gRPC payloads are typically **5–10× smaller** than equivalent JSON. Critical for high-throughput systems.

---

## Streaming: Where RPC Wins

```mermaid
graph TB
    subgraph GRPC_STREAM["gRPC — Server-Side Streaming"]
        SPACE1["  "]
        GC["📱 Client"]
        GC -->|"StreamTransactions(accountId)"| GS["⚡ Open stream"]
        GS -->|"event 1 →"| GC
        GS -->|"event 2 →"| GC
        GS -->|"event N →"| GC
    end

    subgraph REST_STREAM["REST — Polling"]
        RC["📱 Client"]
        RC -->|"GET /transactions?after=t1"| RS1["📦 Batch 1"]
        RC -->|"GET /transactions?after=t2"| RS2["📦 Batch 2"]
        RC -->|"GET /transactions?after=t3"| RS3["📦 Batch 3"]
    end

    style SPACE1 fill:none,stroke:none
```

> gRPC streaming keeps a single connection open. No polling, no overhead. Ideal for real-time feeds.

---

## When to Use REST vs RPC

```mermaid
graph TD
    START(["New API?"])

    Q1{"Public-facing<br/>or external consumers?"}
    Q2{"Need streaming<br/>or low-latency IPC?"}
    Q3{"Multiple language<br/>clients?"}
    Q4{"Simple CRUD<br/>over resources?"}

    REST["✅ **REST**<br/>HTTP + JSON<br/>Universal, cacheable<br/>Easy to explore"]
    GRPC["✅ **gRPC**<br/>Protobuf + HTTP/2<br/>Fast, typed<br/>Streaming support"]
    EITHER["✅ **Either works**<br/>Pick team preference"]

    START --> Q1
    Q1 -->|"Yes — public"| REST
    Q1 -->|"No — internal"| Q2
    Q2 -->|"Yes — streaming / speed"| GRPC
    Q2 -->|"No — request/response"| Q3
    Q3 -->|"Yes — polyglot"| GRPC
    Q3 -->|"No — Go only"| Q4
    Q4 -->|"Yes"| REST
    Q4 -->|"No — complex actions"| EITHER
```

---

## Side-by-Side Summary

```mermaid
graph TB
    subgraph G["⚡ gRPC / RPC"]
        GA["Methods + messages"]
        GB["Protobuf over HTTP/2"]
        GC2["Binary — not human-readable"]
        GD["Needs gRPC client"]
        GE["No native caching"]
        GF["Streaming built-in"]
        GG["🚀 High performance"]
    end

    subgraph R["🌐 REST"]
        RA["Nouns + HTTP verbs"]
        RB["JSON over HTTP/1.1"]
        RC2["Human-readable"]
        RD["Browser-friendly"]
        RE["Cacheable (GET)"]
        RF["Stateless by design"]
        RG["🏆 Most common API style"]
    end
```

> Possibility of both: public APIs use REST. Internal microservice-to-microservice calls may use gRPC for throughput.
