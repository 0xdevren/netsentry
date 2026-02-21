# NetSentry: Next-Generation Network Configuration Validator

NetSentry is an enterprise-grade, deterministic configuration validation suite designed for immutable infrastructure deployments, automated compliance verification, and rigorous security posture enforcement across heterogeneous network environments.

## Architecture and Capabilities

NetSentry parses raw, proprietary device configurations into canonical, vendor-agnostic infrastructure models. These models are subsequently evaluated against declarative, YAML-based policies utilizing a highly concurrent, multi-threaded execution engine.

### Core Features

*   **Vendor-Agnostic Parsing**: Native lexing and parsing engines for Cisco IOS, Cisco NX-OS, Juniper JunOS, and Arista EOS, translating vendor-specific syntax into canonical `ConfigModel` representations.
*   **Deterministic Evaluation Engine**: A robust policy execution runtime utilizing a channel-based worker pool for concurrent rule evaluation, ensuring predictable and rapid validation across extensive rule sets.
*   **Declarative Policy DSL**: A structured Domain Specific Language (DSL) defining validation requirements through explicit match strategies (`contains`, `not_contains`, `regex`, `required_block`) and remediation actions.
*   **Topology Graph Analysis**: Multi-device adjacency inference and Depth-First Search (DFS) topology traversal to definitively identify routing loops, asymmetric paths, subnet overlaps, and duplicate addressing.
*   **Cryptographic Drift Detection**: SHA-256 state hashing and differential analysis to detect and quantify configuration drift against established baselines, facilitating strict change control.
*   **Comprehensive Reporting Pipeline**: Pluggable output generation supporting human-readable formatted tables, structured JSON/YAML for programmatic consumption, and static HTML reports.
*   **Extensible Operator Interfaces**: A unified Command Line Interface (CLI) application and an embedded HTTP REST API server, supporting diverse execution environments from local terminals to automated CI/CD pipelines.
*   **Robust Observability**: System integration with Prometheus metrics and OpenTelemetry tracing, providing deep insight into validation durations, rule execution frequencies, and subsystem latencies.

## Operational Paradigms

NetSentry accommodates multiple execution vectors:

1.  **Continuous Integration / Continuous Deployment (CI/CD)**: Validating proposed configuration changes prior to deployment.
2.  **Audit and Compliance**: Generating point-in-time compliance reports against established framework baselines (e.g., CIS Benchmarks).
3.  **Post-Deployment Verification**: Ensuring operational device states accurately reflect the intended infrastructure-as-code definitions.

## Project Navigation

Comprehensive technical documentation is localized within the `docs/` directory:

| Document | Purpose and Scope |
| :--- | :--- |
| [Architecture Reference](docs/architecture.md) | Component-level architectural design, dependency graphs, and temporal data flow diagrams. |
| [Algorithm Specifications](docs/algorithms.md) | Pseudocode implementations of the core validation pipeline, graph traversal, and drift computation mechanics. |
| [API Integration Guide](docs/api.md) | Exhaustive documentation of the HTTP REST API endpoints, request models, and response structures. |
| [CLI Usage Reference](docs/cli_reference.md) | Exhaustive documentation of command-line interfaces, subcommands, flags, exit codes, and robust troubleshooting matrices. |
| [Policy DSL Specification](docs/policy_dsl.md) | Syntax and semantic rules governing the YAML Policy DSL, including match evaluators, declarative loops, and operational consequences. |
| [Plugin SDK Development](docs/plugin_guide.md) | Interfaces and methodologies for extending parsing capabilities or integrating bespoke programmatic validation logic. |
| [Developer Guidelines](docs/developer_guide.md) | Procedures for establishing local development environments, executing test suites, and adhering to stylistic conventions. |
| [Security Posture](docs/security.md) | Declarations on secret handling, runtime security contexts, and vulnerability disclosure protocols. |
| [Contribution Protocol](docs/contributing.md) | Requirements for formulating pull requests, conventional commit guidelines, and integration testing mandates. |

## Quick Start Subroutine Execution

```bash
# Compile the static binary for the host architecture
make build

# Execute validation against a configuration block
./bin/netsentry validate --config test_router.conf --policy standard.yaml --format table

# Initiate topological relationship inference
./bin/netsentry topology --config router_A.conf --config router_B.conf
```

## Compilation and Deployment

Initialize the build process using the included Makefile directives. Detailed deployment configurations for Docker and Kubernetes environments are available within the `deployments/` directory.

## Licensing Considerations

This software is distributed under the terms of the Apache License, version 2.0. Refer to the `LICENSE` file for full legal stipulations.
