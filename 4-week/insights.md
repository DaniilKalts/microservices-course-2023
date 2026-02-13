# Test Types:
- **Unit Tests:** Test business logic in isolation using mocks for dependencies.
- **Integration Tests:** Test concrete infrastructure implementations (DB, cache) with external systems.
- **E2E Tests:** Verify complete user workflows against the running system.
- **Load Tests:** Measure performance and stability under high traffic.
- **Fuzzing:** Inject random data to find crashes and edge cases.

## Stub vs Mock:
- **Stub:** Returns predefined data (State Verification). Used to simulate dependency behavior.
- **Mock:** Expects specific calls/parameters (Behavior Verification). Used to verify interactions with dependencies.

### Table Driven Tests:
- A pattern where test cases are defined as data (structs) in a slice/array.
- Iterates over the test cases, running the same test logic with different inputs and expected outputs.
