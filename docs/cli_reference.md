# Definitive Command Line Execution Reference

This highly rigid technical reference enumerates deterministic behaviors, execution commands, mandatory logical arguments, and rigorous troubleshooting matrices accompanying all operational limits of the NetSentry core engine CLI implementations.

## Global Parameter Modifications

Specific logical overrides passed globally alter explicit baseline configurations fundamentally guiding diagnostic mechanisms universally across nested operations.

| Instruction Flag | Value Type | Pre-configured Parameter | Operational Implication |
| :--- | :--- | :--- | :--- |
| `--log-level` | `String` | `info` | Adjusts minimal telemetry visibility barriers preventing standard-state metrics capturing processing power dynamically. Accepts `debug`, `info`, `warn`, `error`. |
| `--log-json` | `Boolean` | `false` | Commands programmatic JSON serialization formatting specifically allowing log aggregation ingestion directly replacing pure aesthetic textual structures. |

## 1. Compliance Certification Protocol (`validate`)

Initiates the holistic synchronous evaluation sequence converting proprietary payload definitions progressively against targeted YAML structural conditions, identifying anomalous conditions globally.

**Invocation Construct**: `$ netsentry validate --config <filepath> --policy <filepath> [modifiers]`

### Argument Directives

| Instruction Flag | Mandatory Assertion | Procedural Implication |
| :--- | :--- | :--- |
| `--config` | Yes | Points toward concrete temporal definitions describing active infrastructure state. Requires specific explicit string logic targeting recognized text structures. |
| `--policy` | Yes | Local disk path addressing explicit YAML declarative files defining enforcement limitations dynamically. |
| `--format` | No | Overrides terminal visualization matrices. Accepts deterministic models (`table`, `json`, `yaml`, `html`). Identifies `table` by default rendering colorized ascii output directly. |
| `--output` | No | Aborts standard terminal writing procedures re-routing entire structured text output components specifying precise file storage logic explicitly defined by operational string limits. |
| `--strict` | No | Escalates specific evaluation warnings (`WARN`) strictly elevating overall pipeline result codes towards full structural failures effectively terminating integrated CI/CD chains unceremoniously. |
| `--timeout` | No | Commands deterministic temporal termination metrics utilizing sequence mapping sequences avoiding continuous execution traps natively (e.g., `45s`, `2m`). |
| `--concurrency` | No | Instructs precise limitation models targeting simultaneous multithreaded computation vectors calculating regular extensions globally limiting total system memory ingestion bounds. |

### Anticipated Formatted Visualization (Mockup)

Execution invoking generic configurations utilizing the `table` rendering parameter structurally outputs deterministic visual representations explicitly:

```text
DEVICE : Aggregation-Switch-02
POLICY : Zero-Trust Architecture Controls (v3.2)
------------------------------------------------------------------------
  RULE-ID      STATUS   SEVERITY   MESSAGE
-------------+--------+----------+--------------------------------------
  SSH-CIPHER   PASS     HIGH       rule SSH-CIPHER passed
  SNMP-V3      FAIL     CRITICAL   rule SNMP-V3 violated: Legacy SNMP v2
                                   enabled. Require authPriv encryption.
  BANNER-MOTD  WARN     LOW        rule BANNER-MOTD warning: Unauth
                                   access statement missing exact match.
------------------------------------------------------------------------
SUMMARY:
  Passed   : 1
  Failed   : 1
  Warnings : 1
  Skipped  : 0
  Score    : 33% (Action Required)
```

### Deterministic Output Escalations (Exit Codes)

| State Vector | Functional Designation | Remediation Context |
| :--- | :--- | :--- |
| `0` | Absolute Compliance | Payload fulfills declared operational bounds successfully. |
| `1` | Integrity Violation | Active string structures triggered evaluation faults defining failure paths distinctly. |
| `2` | Execution Fault | An internal architectural constraint failure triggered uncontrolled runtime halting parameters specifically identifying code failures rather than logical constraints. |
| `3` | Parsing Failure | Configurable inputs prevented operational initialization limiting functionality significantly mapping missing objects explicitly. |
| `4` | Temporal Restriction | Configurable duration limit restrictions forcefully aborted background thread synchronization routines prematurely. |

## 2. Infrastructure Deviation Analysis (`drift`)

Initiates distinct symmetric comparisons executing differential tracking algorithms mapping progressive operational states discovering unintended structural configuration adjustments dynamically globally.

**Invocation Construct**: `$ netsentry drift --baseline <filepath> --current <filepath> [modifiers]`

### Argument Directives

| Instruction Flag | Functional Designation |
| :--- | :--- |
| `--baseline` | Path explicit defining prior functional states accurately denoting control parameters globally. |
| `--current` | Path explicit describing recent acquisition strings investigating possible logical shifts. |
| `--threshold` | Float variable explicitly specifying acceptable absolute variation metrics defining deviation failures strictly bypassing minor temporal sequence rearrangements globally. |

## 3. Topographical Integrity Verification (`topology`)

Generates structured internal directed node maps compiling relationships defined intrinsically within distinct configuration definitions computing specific global network layout assertions avoiding explicit operational interactions exclusively through abstract string decoding alone.

**Invocation Construct**: `$ netsentry topology --config <filepath> [--config <filepath> ...]`

## 4. Policy Execution Testing (`policy lint`)

Analyzes offline configuration boundaries executing absolute logic constraint mechanisms verifying formal definitions ensuring DSL mappings avoid fatal failures specifically when executed within live operational boundaries.

**Invocation Construct**: `$ netsentry policy lint <filepath>`

## Operational Anomaly Remediation (Troubleshooting)

Operational limitations occasionally manifest during structural interactions.

| Identifiable Failure Signature | Direct Diagnostic Remediation Pathway |
| :--- | :--- |
| `parser: failed to classify structure` | The configuration file exhibits completely unidentified header lines circumventing heuristic rules globally. Execute explicitly targeted manual cleanup removing unidentifiable logging string metadata lines positioned strictly preceding authentic configuration sequences limits explicitly. |
| `policy engine: invalid regex` | RE2 parsing algorithms automatically reject syntactically unverified parameters executing computational vectors inherently exposing Denial-Of-Service possibilities. Refactor boundaries targeting simplistic character combinations globally without lookaheads. |
| `timeout constraint exhausted` | Multi-threaded computations consumed entire allocatable operational processing periods explicitly avoiding returning finalized aggregation arrays. Minimize global payload limits investigating possible infinite-loop syntax structures implicitly. Reduce regex string mapping criteria recursively limiting required evaluation cycles heavily. |
