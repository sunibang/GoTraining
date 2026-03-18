# Domain Layer

Welcome to the **Domain Layer**, the absolute center of our Go Bank architecture. This package defines the core business entities, types, and the specific errors our bank cares about.

## Meet the Entities

- **[Account](account.go)**: A bank account tracking identity, balance, and status (e.g., Active or Locked).
- **[Transaction](transaction.go)**: A record of money moving in (deposit) or out (withdrawal).
- **[Transfer](transfer.go)**: A helper type for logging and API mapping (actual transfers are just paired transactions!).

## Why keep it isolated?

Notice how this package imports nothing from the rest of our app? That's intentional! 
- It has **zero outward dependencies**—no databases, no HTTP handlers.
- All business errors (like `ErrInsufficientFunds`) live here so the rest of the app can handle them consistently using `errors.Is`.

## Your Next Step

The good news? This layer is **fully pre-built** for you. 

Take a quick peek at the code to get familiar with the structures. Once you're done, let's see where this data is actually saved by exploring the **[Repository Layer](../repository/README.md)**.