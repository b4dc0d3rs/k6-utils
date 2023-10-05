# k6utils

A bunch of random functions for k6 performance testing that I found missing, but useful in our work.

# Compile for development
```sh
xk6 build v0.46.0 \
  --with github.com/b4dc0d3rs/k6-utils=.

./k6 run k6utils.js
```

# Use

Just import:
```js
import k6utils from 'k6/x/k6utils';
```

## sleepMilliseconds
```js
k6utils.sleepMilliseconds(666);
```
## CSV operator

A native CSV operator that loads all CSV records to a map in memory. Empty lines are skipped. Header in the CSV file is mandatory for mapping purposes.

```js
const data = k6utils.load('data.csv', ',');
data[0].csvColumnName;
```

### CSV random record
Returns random record from the CSV map. Load
```js
const allRows = k6utils.load('data.csv', ',');

// a row can be returned many times
const oneRandomRow = k6utils.takeRandomRow();

const row5 = k6utils.takeRowByIndex(4)

// this method removes polled row from in-memory cache
// each row is returned only once. The item is removed from in-memory cache before returning.
const uniqueRandomRow = k6utils.pollRandomRow();
```

### Expiring cache
There is one global in-memory cache that evicts KV set after pre-configured number of seconds since insertion passed.

```js
// configure it in setup method
k6utils.createCacheWithExpiryInSeconds(1)

// insert anything anytime
k6utils.putToCache('key', 'value')

// get anywhere anytime, even in a different method
k6utils.getFromCache('key')
```