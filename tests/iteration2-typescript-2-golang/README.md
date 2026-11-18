# Iteration 2: TypeScript → `.sets` → Go

| Step | Knows | Files |
| ---- | ----- | ----- |
| 1. Input | — | [input/types.ts](input/types.ts) (same as iteration 1) |
| 2. TS → `.sets` standard (2 passes) | value model + TypeScript | [ts-to-sets/STANDARD.md](ts-to-sets/STANDARD.md), [REVIEW.md](ts-to-sets/REVIEW.md) |
| 3. Transpiled IR (the only hand-off) | — | [sets/ts.sets](sets/ts.sets) (prelude), [sets/types.sets](sets/types.sets) |
| 4. `.sets` → Go standard (2 passes) | value model + `.sets` only | [sets-to-go/STANDARD.md](sets-to-go/STANDARD.md), [REVIEW.md](sets-to-go/REVIEW.md) |
| 5. Generated Go | — | [go/](go/): `prelude.go`, `model.go`, `sets_runtime.go`, `model_test.go` |
| 6. Results | — | [FINDINGS.md](FINDINGS.md), [VECTORS.md](VECTORS.md) |

Run the Go vectors with `cd go && go test ./...`.
