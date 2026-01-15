1. One service per `.proto` file (shared messages in separate common `.proto` files).
2. `minikube` runs a local Kubernetes cluster to test microservices + infrastructure.
3. `Miro` is for system design diagrams/whiteboarding.
4. `context.WithTimeout` caps DB connect/query time; on timeout `ctx.Done()` closes and the call should error.
5. `context.WithValue` stores request-scoped metadata (e.g., request ID), not general parameters.
6. Context cancellation stops goroutines on events (signals/cancel) and prevents leaks.
7. For complex SQL, prefer raw SQL over `squirrel`.
8. Triggers/procedures reduce transparency; prefer business logic in Go, DB for constraints/transactions.
