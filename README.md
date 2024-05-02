# avalanche

Perform stress load test on specific URL in avalanche-fetching manner

##### Motivation

- Explore server boundaries with maximum requests
  - trigger errors like `Error: ENFILE: file table overflow`
  - memory usage
- Measure response per second (rps) of specific URL
- Testing recoverability after denial of service simulation

##### Example

```bash
$ make run
Commencing avalanche fetch on: http://localhost:8282/
Press Ctrl+C to stop...
27 40 0 0 229 42 56 57 126 34 17 101 0 81 44 0 126 51 0 81 124 0 122 22 84 33 29 ^C
Test complete in:       27.389391021s for url: http://localhost:8282/
Successfull (rps):      min=0, avg=56.851851851851855, max=229, stddev=53.82461484522962
Response time (s):      min=0.655551337, avg=6.319355257750636, max=20.723939862, stddev=3.0739266088686144
Fetch rate (fps):       min=6580, avg=56119.14814814815, max=87586, stddev=17952.485274568746
Σ of status codes:      map[total:1513670 error:1511308 200:1578 500:784]
---
```

##### Requirements

- Go: 1.22.2

```bash
make run
make test
```

##### Links

- <https://go.dev/>
