# Where should we place the interface?

There's a debate about where should we place the interface. Two options on the table.

## Define on Producer Side

### Advantages

- Producer can provide a clear picture of what functionality it provides.
- It's easier to distribute widely and reuse.

### Use cases

- When an interface need to be widely used and distributed. For instance, the built-in `context`, `io.Writer` 
interfaces are defined on the producer side. 
- When a producer provides multiple implementations for the same interface. In this case, it makes sense to place 
interface on the producer side. For instance, the `image.Image` package.

## Define on Consumer Side

### Advantages

- Consumer have a say on what functionalities they need. The interface doesn't have to be one-fit-for-all.
- When writing unit tests, mocking can be simplified.
- If the interface only contains primitive types, we can swap implementation without changing any code on the consumer side.

### Use cases

- When you want to decouple packages. For instance, an API server could depend on a database access layer (let's assume that 
the access logic is implemented in the `store` package). The store interface can be defined in the API server package to
avoid strong coupling between these two packages. This creates better abstraction, and it will make the package more modular 
and self-contained.
- When the interface provided by the producer is heavy, and you don't need all the exposed functions.
A great man has said:
> Clients should not be forced to depend on methods they do not use.
> –Robert C. Martin
- Reference: https://github.com/golang/go/wiki/CodeReviewComments#interfaces

## Verdict

I know you've been waiting for the answer. So what should we do? 
Unfortunately, there's no universal rule that tells us where to put the interfaces. 
I hate to say this but 
> It depends.