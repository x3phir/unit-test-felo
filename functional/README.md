# FELO Functional Tests

This folder contains the functional-test scaffold for FELO.

Current status:

- real infra is defined in `docker-compose.functional.yml`
- scenario tests are written for the five priority services
- tests compile under the `functional` build tag
- tests skip until real gRPC adapters are connected

Read:

- `docs/functional-testing-strategy.md`
- `docs/run-functional-testing-guide.md`
