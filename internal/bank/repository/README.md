# Repository Layer

Welcome to the **Repository Layer**—our data vault! This package handles all communication with the database, cleanly separating data storage from our business logic.

## What happens here?

- **Database Isolation:** The rest of the app doesn't know *how* data is stored, it just uses the [repository.Repository](repository.go) interface.
- **Type-Safe SQL:** We use a Postgres implementation ([postgres/repository.go](postgres/repository.go)) powered by [`go-jet`](https://github.com/go-jet/jet). This generates type-safe SQL queries, saving us from typos! (Check out the generated models in [postgres/gen/](postgres/gen/)).
- **Data Mapping:** This layer translates raw database rows into our clean `domain` entities.

## Your Next Step

We've already set up the Postgres interactions for you!

Feel free to explore how `go-jet` queries are written. Now that you know where the data is stored, let's see how business rules are enforced by moving to the **[Service Layer](../service/README.md)**.