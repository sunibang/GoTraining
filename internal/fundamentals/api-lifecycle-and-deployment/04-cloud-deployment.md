# Cloud Deployment

---

## ECS Needs Explicit Signals to Manage Traffic

```mermaid
graph TB
    subgraph Signals["🚦 Health Signals"]
        direction TB
        LIV["/healthz — LIVENESS<br/>Is the process alive?<br/>ECS restarts container on failure"]
        RDY["/readyz — READINESS<br/>Is the app ready for traffic?<br/>ALB removes from rotation on failure"]
        LIV --> RDY
    end

    subgraph ECS["☁️ AWS ECS + ALB"]
        direction TB
        ALB["⚖️ Application Load Balancer<br/>Polls /readyz every 10s"]
        TASK["📦 ECS Task<br/>Running Go container"]
        ALB --> TASK
    end

    RDY ~~~ ALB
    ALB -->|"GET /readyz → 200 OK"| TASK
    ALB -->|"503 → drain from target group"| TASK
    LIV -.->|"failing → ECS restarts"| TASK
    RDY -.->|"failing → ALB stops routing"| TASK
```

> Liveness tells ECS to restart. Readiness tells the ALB to stop sending traffic. They serve different purposes.

---

## Rolling Deployment with Health Checks

```mermaid
sequenceDiagram
    autonumber
    participant ECS as ☁️ ECS Scheduler
    participant ALB as ⚖️ ALB
    participant Old as 📦 Old Task (v1.4.1)
    participant New as 📦 New Task (v1.4.2)

    Note over ECS,New: Rolling deployment — minimum healthy % maintained

    ECS->>New: Start new task with image v1.4.2
    New->>New: App initialising — readyState = false

    ALB->>New: GET /readyz
    New-->>ALB: 503 Service Unavailable — not ready yet

    Note over ALB,Old: ALB keeps routing to old task during startup

    New->>New: App ready — readyState = true
    ALB->>New: GET /readyz
    New-->>ALB: 200 OK

    ALB->>ALB: Add new task to target group
    ECS->>Old: Send SIGTERM — begin graceful shutdown

    Old->>Old: readyState = false — stop accepting new requests
    Old->>Old: Drain in-flight requests (30s)
    Old->>ECS: Exit 0 — clean shutdown

    ECS->>Old: Remove task
    Note over ECS,New: Deployment complete — zero dropped requests
```

---

## Graceful Shutdown: SIGTERM to Exit 0

```mermaid
graph TB
    SIGTERM["📡 SIGTERM received<br/>ECS stops task during deployment<br/>or scale-in event"]

    RDY["🔴 /readyz → 503<br/>ALB stops routing new requests<br/>immediately"]

    DRAIN["⏳ Drain in-flight requests<br/>context with 30s deadline<br/>server.Shutdown(ctx)"]

    EXIT["✅ os.Exit(0)<br/>Clean exit — ECS marks task<br/>as stopped successfully"]

    SIGTERM --> RDY --> DRAIN --> EXIT
```

> `server.Shutdown(ctx)` stops accepting new connections but waits for active requests to complete before returning.

---

## Liveness vs Readiness: What Each Check Does

```mermaid
graph TB

    subgraph NeverDo["❌ Never Do This"]
        N1["Do NOT return 200 on /readyz during shutdown<br/>ALB will route to a dying container"]
        N2["Do NOT exit immediately on SIGTERM<br/>In-flight requests will be dropped mid-response"]
    end
    
    subgraph Liveness["/healthz — Liveness"]
        L1["✅ Return 200 if the process is running"]
        L2["❌ Return 500 if the app is deadlocked<br/>or in an unrecoverable state"]
        L3["ECS action: restart the container"]
    end

    subgraph Readiness["/readyz — Readiness"]
        R1["✅ Return 200 when app is initialised<br/>and ready for traffic"]
        R2["❌ Return 503 during startup<br/>or graceful shutdown"]
        R3["ALB action: remove from target group"]
    end
```

> Set readyState to false on SIGTERM immediately — before draining. This stops ALB routing before the drain begins.

---

## Blue/Green vs Rolling: Choosing a Strategy

```mermaid
graph TB
    subgraph Rolling["🔄 Rolling Deployment"]
        RO1["Replace tasks one by one<br/>Old and new run simultaneously"]
        RO2["✅ No extra infrastructure cost"]
        RO3["✅ Built into ECS natively"]
        RO4["⚠️ Mixed versions in flight during rollout"]
        RO1 --> RO2 --> RO3 --> RO4
    end

    subgraph BlueGreen["🟦🟩 Blue/Green Deployment"]
        BG1["Spin up full new environment (green)<br/>Switch ALB target group when ready"]
        BG2["✅ Instant rollback — flip back to blue"]
        BG3["✅ Zero mixed-version traffic"]
        BG4["⚠️ Double the infrastructure during cutover"]
        BG1 --> BG2 --> BG3 --> BG4
    end
```

> Rolling suits stateless APIs with backwards-compatible changes. Blue/Green suits high-risk releases where instant rollback matters.

---

## ECS Task Definition: Key Fields

```mermaid
graph TB
    subgraph TaskDef["📋 ECS Task Definition"]
        T1["image: account.dkr.ecr.region.amazonaws.com/myapp:1.4.2<br/>← immutable semver tag"]
        T2["cpu: 256 / memory: 512<br/>← right-size for your workload"]
        T3["healthCheck: GET /healthz<br/>interval: 10s, threshold: 3"]
        T4["environment: [{PORT: 8080}]<br/>secrets: [{DB_URL: secretsmanager/prod/db}]"]
        T5["logConfiguration: awslogs<br/>group: /ecs/myapp, region: ap-southeast-2"]
    end

    subgraph Ops["✅ Operational Outcomes"]
        O1["ECS replaces unhealthy tasks automatically"]
        O2["Secrets injected at runtime — not baked into image"]
        O3["Logs stream to CloudWatch for every task"]
    end

    TaskDef --> Ops
```

> Never bake secrets into the image. Inject them at runtime via AWS Secrets Manager. The image is immutable and shareable — secrets are not.
