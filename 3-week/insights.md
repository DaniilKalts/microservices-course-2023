1. Putting all the logic for a single handler (database access, business logic, and HTTP/gRPC handling) in one file is a bad practice because it quickly becomes messy and hard to navigate.
2. When everything is in one place, testing becomes harder: concerns are tightly coupled, you can’t test parts in isolation, and it’s less clear which layer is responsible when a test fails.
3. Large, mixed-responsibility handler files increase cognitive load: it’s harder to keep context in your head, reason about changes, and add new code safely.
4. It’s usually better not to use protobuf-generated types as your service/domain layer models. Transport types can have constraints or awkward mappings (enums, optional fields, timestamps, `oneof`, etc.), and if the transport changes (e.g., gRPC → HTTP/JSON), you don’t want your core logic to depend on those types. Prefer defining your own internal models/DTOs and use mappers/converters at the boundary to translate to/from protobuf messages.
5. Layers should depend on abstractions (interfaces/contracts), not on concrete implementations. Each layer should only know what methods it can call, not how the lower layer implements them—keeping dependencies flowing from higher-level code down to lower-level details.
6. Where to store interfaces: where they’re used vs. where they’re implemented.
   - Store interfaces where they’re used: each consumer defines an interface with only the methods it needs. This keeps interfaces small and focused, but it can become messy if you rely heavily on a DI container, because the container must wire many different interfaces.
   - Store interfaces where they’re implemented: interfaces are centralized (e.g., under a `repository` package), which can be convenient to reference and reduces duplication. In this approach, the `repository` package contains only interfaces, while the actual database CRUD implementations live in entity-specific packages.
7. Use constants for table names and column names so it’s easy to update queries if a name changes.
8. The idea behind a pool is to limit the number of resources and reuse them once they become available.
9. When you have no idea about naming a method - brainstorm it with the others.
10. If external package contains a single utility function - it may be better to copy and paste to the codebase, rather than downloading the whole dependency.
