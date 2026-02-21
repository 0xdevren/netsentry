# Programmatic Application Programming Interface (API) Reference

The NetSentry application embeds a fully functional HTTP REST API listener capable of seamlessly integrating into external orchestration workflows, automated provisioning environments, and specialized compliance web portals.

## Global Communication Parameters

*   **Initialization Sequence**: Execute `--serve --addr <interface>:<port>`. Defaults to `:8080`.
*   **Encapsulation**: All payload representations transmitted must strictly adhere to `application/json` formatting. 
*   **Error Topology**: Exceptions return structured JSON dictionaries delineating precise failure modalities accompanying relative HTTP status integers.

### Mock Error Response

```json
{
  "error": "validation constraint error: invalid device config payload definition"
}
```

## 1. Subsystem Viability Check (`/healthz`)

Responds exclusively with architectural metadata verifying proper thread initiation parameters and active process retention. Identifies versioning to ensure operational predictability.

**HTTP Method**: `GET`
**URI Target**: `/healthz`

### Request Implementation

```bash
curl -i -X GET http://localhost:8080/healthz
```

### Anticipated Response Output

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 21 Feb 2026 09:00:00 GMT

{
  "commit": "unknown",
  "status": "ok",
  "version": "dev"
}
```

## 2. Comprehensive Compliance Validation (`/api/v1/validate`)

The primary functional endpoint. Triggers full synchronous execution of the analysis engine against explicitly provided string configurations and concurrent policy definitions.

**HTTP Method**: `POST`
**URI Target**: `/api/v1/validate`

### Payload Configuration Requirements

| Property Definition | Specification Variable | Mandatory Indicator | Description |
| :--- | :--- | :--- | :--- |
| `config` | `string` | Yes | Escaped ASCII string containing the proprietary hardware configuration directives. |
| `policy_yaml` | `string` | Yes | Escaped ASCII string reflecting the declarative YAML policy rulesets governing the compliance evaluation. |
| `strict` | `boolean` | No | Overrides warning definitions forcing subsequent exit calculations to signify strict operational failure constraints. |
| `concurrency` | `integer` | No | Dictates logical thread utilization. Recommended threshold remains `4` to prevent arbitrary CPU exhaustion. |

### Execution Syntax

```bash
curl -X POST http://localhost:8080/api/v1/validate \
-H "Content-Type: application/json" \
-d '{
  "config": "hostname Core-Router\nbgp as 65000\n",
  "policy_yaml": "name: Example\nrules:\n  - id: BGP-CHECK\n    match:\n      contains: \"bgp as\"\n    action:\n      deny: false\n    severity: HIGH",
  "strict": true
}'
```

### Deterministic Output Structure

The output structurally mirrors internal `policy.Report` data definitions, enumerating device topologies and granular array listings of validation anomalies.

```json
{
  "metadata": {
    "device_name": "Core-Router",
    "device_type": "unknown",
    "policy_name": "Example"
  },
  "summary": {
    "passed": 1,
    "failed": 0,
    "warnings": 0,
    "skipped": 0,
    "total": 1,
    "score": 100
  },
  "results": [
    {
      "rule_id": "BGP-CHECK",
      "status": "PASS",
      "message": "rule BGP-CHECK passed",
      "severity": "HIGH",
      "remediation": ""
    }
  ]
}
```

## Troubleshooting & Failure Categorization

Integration frameworks must anticipate distinct HTTP status error mappings reflecting explicit failure sequences within the computational logic limits.

| Evaluated Status Parameter | Technical Definition | Rectification Vector |
| :--- | :--- | :--- |
| `400 Bad Request` | Structural anomalies exist within the supplied JSON data frame | Verify payload serialization and JSON structural boundaries. Confirm all mandatory fields contain applicable non-empty values. |
| `422 Unprocessable Entity` | String parameters could not undergo parsing algorithms successfully | Ensure the provided `config` structure represents an identifiable hardware structure supported by registered lexers. |
| `500 Internal Server Error` | Unrecoverable thread starvation or execution routing failure | Investigate internal system telemetry limitations and execution duration restrictions limiting process execution functionality. |
