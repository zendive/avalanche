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
$ make build
$ ./avalanche https://localhost:3001/users
Version: go1.22.2, CPUs: 12
Press Ctrl+C to stop...
718 750 737 790 767 770 775 775 790 760 779 779 762 818 792 807 806 798 768 756 801 800 806 796 786 790 795 801 803 800 790 803 803 802 815 799 788 782 798 810 800 796 784 783 801 789 788 774 800 800 804 808 791 793 793 776 808 799 790 775
Test complete in:       60.0(s) for url: https://localhost:3001/users
Successfull (rps):      avg=788.683,    min=718.000,    max=818.000,    stddev=18.542
Response time (s):      avg=0.679,      min=0.003,      max=23.326,     stddev=0.895
Fetch rate (fps):       avg=818.533,    min=778.000,    max=838.000,    stddev=10.390
Σ of status codes:      map[200:47321 error:477 total:47798]
---

```

##### Requirements

- Go: 1.22.2

##### Install
```bash
go install github.com/zendive/avalanche@latest
```

##### Links

- <https://go.dev/>
