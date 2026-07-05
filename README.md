# Fehlerbaum (Fault Tree Analysis) in Go

A fault tree analysis tool written in Go that computes Minimal Cut Sets (MCS), failure probabilities, and reliability values — and exports the tree as a Graphviz visualization.

## Structs

### `EVENT`
A basic event (leaf node) with a known failure probability and reliability.

| Field | Type | Description |
|---|---|---|
| `Title` | string | Name of the event |
| `Reliability` | float64 | Probability of working: R |
| `Failure` | float64 | Probability of failing: F = 1 - R |

### `AND_NODE`
Gate that fails only when **all** children fail.

- **Failure:** `F = F₁ × F₂ × ...`
- **Reliability:** `R = 1 - (1-R₁)(1-R₂)...`

### `OR_NODE`
Gate that fails when **any** child fails.

- **Failure:** `F = 1 - (1-F₁)(1-F₂)...`
- **Reliability:** `R = R₁ × R₂ × ...`

### `NODE` Interface
All node types implement:
- `getCutSets() [][]EVENT` — returns all Minimal Cut Sets
- `getReliability() float64`
- `getFailure() float64`
- `toDot() string` — Graphviz DOT representation

## Generate PNG

```bash
go run Fehlerbaum.go
dot -Tpng Fehlerbaum.dot -o Fehlerbaum.png
```
