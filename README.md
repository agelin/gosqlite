#### Intro ####

The SQLite Go library as provided by Russ Cox is capable, but very spartan.

This version will expand on the existing code by providing extra functionality that may improve ease of use.

#### License ####

The New BSD License covers work by Russ Cox, with the MIT License covering the additional work by Richard B. Lyman.

#### Benchmarking ####

To run the memprof...
... first: `sqlite.test.exe -test.memprofile=mem.prof -test.bench=.* -test.memprofilerate=1`
... then: `go tool pprof sqlite.test.exe mem.prof --svg > out.svg`
