# GoN Verifier

## Verifier

`Verfier` verifier each participant's task evidence and output a `taskpoint.xlsx` file
- `taskpoint1.xlsx` Stage 1
- `taskpoint2.xlsx` Stage 2
- `taskpoint2b.xlsx` Stage 2 (with new chan/port pair)
- `taskpoint3.xlsx` Stage 3

## Ranker

`Ranker` read `taskpoint3.xlsx` under each participant's directory.
- If a rankable task has reason start with `race/xxx/yyy/zzz`, then it will be ranked.
- Ranker will struct the reason and output a corresponding `rank.xlsx` file.
- Ranker will append the top 10 back to `taskpoint3.xlsx` file under each participant's directory.

## Scorecard

`Scorecard` reads `taskpoint{1,2,3}.xlsx` under each participant's directory and output a `scorecard.xlsx` file.

