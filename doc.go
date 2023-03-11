package verifier

// Validation Logic for Stage2 and Stage2
// Stage 2
// - never-go-back:
//   1. query ibc class on iris and get its trace
//   2. compare the trace with flow
//   3. query the owner of ibc asset
// - other types:
//   1. convert to flow to query target, like i--(1)-->s
//   2. retrieve tx hash on sender chain and validate port, channel, recipient, sender, ...
//   3. query ibc asset on last chain and its owner
