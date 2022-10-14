# k6utils

A bunch of random functions for k6 performance testing that I found missing, but useful in our work.

# Use

Just import:
```ts
import k6utils from 'k6/x/k6utils';
```

## sleepMilliseconds
```ts
  k6utils.sleepMilliseconds(666);
```
## CSV operator

A native CSV operator that loads all CSV records to a map in memory. Empty lines are skipped. Header in the CSV file is mandatory for mapping purposes.

```ts
const data = k6utils.load('data.csv', ',');
data[0].csvColumnName;
```

### CSV random record
Returns random record from the CSV map. Load
```ts
const allRows = k6utils.load('data.csv', ',');
const oneRandomRow = k6utils.takeRandomRow();
```