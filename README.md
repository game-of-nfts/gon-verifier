# Gon Verifier

## Install

```bash
go install ./cmd/gon-verifier 
```

## Usage

```bash
gon-verifier <evidence.xlsx>
```

It will output four evidence results, including point lost reasons:

- `taskpoint1.xlxs` Stage one result.
- `taskpoint2.xlsx` Stage two result.
- `taskpoint2b.xlsx` Some use new chan/pair of stars and juno, this will generate points for thess task.
- `taskpoint3.xlsx` Stage three result, except for quiz game.