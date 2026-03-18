# Containerisation

---

## Why Tiny, Static Binaries Matter

```mermaid
graph TB
    subgraph Problem["❌ Typical Language Runtime Images"]
        SPACE1["  "]
        P1["🐍 Python app: ~900 MB<br/>Full interpreter + pip packages"]
        P2["☕ Java app: ~500 MB<br/>Full JRE + dependencies"]
        P3["🟨 Node app: ~400 MB<br/>Full Node.js + node_modules"]
    end

    style SPACE1 fill:none,stroke:none
    subgraph Solution["✅ Go Scratch Image"]
        S1["⚡ Go binary: ~12 MB<br/>Statically compiled — zero runtime"]
        S2["🔒 No shell, no package manager<br/>Minimal attack surface"]
        S3["🚀 Fast registry pull times<br/>Faster container startup"]
    end

    Problem -.->|"replace with"| Solution
```

> A static Go binary needs no runtime. The container image is just the binary — nothing else.

---

## Multi-Stage Docker Build

```mermaid
graph TB
    subgraph Stage1["🏗️ Stage 1: Builder (golang:1.22-alpine)"]
        B1["FROM golang:1.22-alpine AS builder"]
        B2["COPY go.mod go.sum ./"]
        B3["RUN go mod download"]
        B4["COPY . ."]
        B5["RUN CGO_ENABLED=0 GOOS=linux go build -o /app ."]
        B1 --> B2 --> B3 --> B4 --> B5
    end

    subgraph Stage2["📦 Stage 2: Runner (scratch)"]
        R1["FROM scratch"]
        R2["COPY --from=builder /app /app"]
        R3["EXPOSE 8080"]
        R4["ENTRYPOINT /app"]
        R1 --> R2 --> R3 --> R4
    end

    Stage1 -->|"copy binary only<br/>SDK discarded"| Stage2

    RESULT["🎯 Final image: ~12 MB<br/>No Go SDK, no shell, no OS packages"]
    Stage2 --> RESULT
```

> The builder stage has the full Go SDK. The runner stage has only the compiled binary. The SDK is thrown away.

---

## Image Tagging Strategy

```mermaid
graph TB
    subgraph Tags["🏷️ Tag Every Image Two Ways"]
        T1["Immutable semver tag<br/>myapp:1.4.2<br/>← pinned in deployment config"]
        T2["Mutable latest tag<br/>myapp:latest<br/>← for local dev and quick pulls"]
        T1 ~~~ T2
    end

    subgraph Why["✅ Why Both?"]
        W1["Semver: reproducible deployments<br/>know exactly what is running"]
        W2["Latest: convenience for developers<br/>never use latest in production"]
        W1 ~~~ W2
    end

    Tags --> Why

```

> Production deployments must reference immutable tags. `latest` is a moving target — unsafe for rollbacks.

---

## Build-Time Version Metadata

```mermaid
graph TB
    subgraph BuildTime["🏗️ Build Time: go build -ldflags"]
        LD1["VERSION=1.4.2"]
        LD2["COMMIT=abc1234"]
        LD3["BUILD_TIME=2025-07-01T10:00:00Z"]
    end

    subgraph Runtime["⚙️ Runtime — GET /buildinfo"]
        RT1["version: 1.4.2"]
        RT2["commit: abc1234"]
        RT3["built: 2025-07-01T10:00:00Z"]
    end

    subgraph Benefits["✅ Benefits"]
        BN1["🔍 Instantly know which version is running"]
        BN2["🐛 Correlate production errors to exact Git commit"]
        BN3["📊 Dashboards show deployment progress across tasks"]
    end

    BuildTime --> Runtime --> Benefits
```

> Embed version metadata at build time. Expose it via `/buildinfo`. Never guess what is running in production.

---

## Container Security: The Scratch Advantage

```mermaid
graph TB
    subgraph UbuntuBased["❌ ubuntu:22.04 Base"]
        U1["Shell (bash, sh)"]
        U2["Package manager (apt)"]
        U3["System utilities (curl, wget)"]
        U4["~80 MB of OS packages"]
        U5["⚠️ Each package = potential CVE surface"]
    end

    subgraph ScratchBased["✅ scratch Base"]
        SPACE1["  "]
        SC1["No shell"]
        SC2["No package manager"]
        SC3["No system utilities"]
        SC4["Zero OS packages"]
        SC5["🔒 Attacker has nothing to pivot from"]
    end

    style SPACE1 fill:none,stroke:none
```

> With a scratch image, a compromised container has no tools to escalate with. The attack surface is the application code itself — nothing more.

---

## CI/CD Image Build Pipeline

```mermaid
sequenceDiagram
    autonumber
    participant Dev as 👩‍💻 Developer
    participant CI as ⚙️ CI Pipeline
    participant Registry as 🗄️ Container Registry

    Dev->>CI: git push → tag v1.4.2
    CI->>CI: Run tests — fail fast if broken
    CI->>CI: docker build -t myapp:1.4.2 .
    CI->>CI: docker scan myapp:1.4.2 — check CVEs

    CI->>Registry: docker push myapp:1.4.2
    CI->>Registry: docker push myapp:latest

    Note over CI,Registry: Tag with immutable semver AND latest<br/>Deployment config references the semver tag
```
