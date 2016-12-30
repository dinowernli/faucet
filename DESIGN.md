# Some notes about the Faucet design

The distributed system consists of three components: coordinators, workers, and clients.

## Coordinator

The coordinator is a service which receives incoming requests for checks. It acts as a frontend for clients and dispatches the incoming requests to workers. The coordinator has a configuration describing where to find workers and uses persistent storage to keep track of all builds (present and past).

All operations on the coordinator are designed to be cheap, anything expensive is offloaded to workers.

## Worker

A worker is a service capable of performing expensive operations such as executing builds. Operations running on workers are designed to be simple (i.e., as stateless as possible) and the system can cope with workers joining, leaving, or failing at any time.

Workers know how to checkout a repository, reason about changesets, execute builds, run tests, etc.

## Client

A client is a binary used by the end-user and executed from within a repository. It talks exclusively to the coordinator. The client's main purpose is to kick off builds by shipping the current commit hash to the coordinator.
