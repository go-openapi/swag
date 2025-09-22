# benchmarks

Disclaimer: `go-openapi` has no affiliation with the providers of the tested libraries.

Credits: the definition of the small, medium and large payloads has been taken from the benchmarks
published at `github.com/goccy/go-json`. Thanks @goccy.  We'd love to include your library as one
of our supported adapters.

This benchmark is not a competitive benchmark, but merely a way for us to ensure that
our "JSON adapter" feature runs with good performances.

> Expected behaviour: the `Adapter` layer is relatively thin, yet it involves an unavoidable
> allocattion to materialized the returned interface type. So this slightly degrades
> performance when compared to calling the concrete type directly.

**TODO**: publish another run on `go1.25` and try out go's new garbage collection algorithm.

```sh
go version go1.24.5 linux/amd64
go test -v -bench . -run Bench -benchtime 30s -benchmem
```


```
goos: linux
goarch: amd64
pkg: github.com/go-openapi/swag/jsonutils/adapters/testintegration/benchmarks
cpu: Intel(R) Core(TM) i5-6200U CPU @ 2.30GHz
BenchmarkJSON
BenchmarkJSON/with_standard_library
BenchmarkJSON/with_standard_library/standard_ReadJSON_-_small
BenchmarkJSON/with_standard_library/standard_ReadJSON_-_small-4         	 9776988	      3561 ns/op	  40.44 MB/s	     432 B/op	      10 allocs/op
BenchmarkJSON/with_standard_library/standard_WriteJSON_-_small
BenchmarkJSON/with_standard_library/standard_WriteJSON_-_small-4        	35532856	       880.4 ns/op	 149.93 MB/s	     160 B/op	       2 allocs/op
BenchmarkJSON/with_standard_library/standard_FromDynamicJSON_-_small
BenchmarkJSON/with_standard_library/standard_FromDynamicJSON_-_small-4  	 7411022	      4742 ns/op	    1528 B/op	      32 allocs/op
BenchmarkJSON/with_standard_library/standard_ReadJSON_-_medium
BenchmarkJSON/with_standard_library/standard_ReadJSON_-_medium-4        	 1514880	     23488 ns/op	  88.13 MB/s	     640 B/op	      19 allocs/op
BenchmarkJSON/with_standard_library/standard_WriteJSON_-_medium
BenchmarkJSON/with_standard_library/standard_WriteJSON_-_medium-4       	18548271	      1801 ns/op	 175.48 MB/s	     336 B/op	       2 allocs/op
BenchmarkJSON/with_standard_library/standard_FromDynamicJSON_-_medium
BenchmarkJSON/with_standard_library/standard_FromDynamicJSON_-_medium-4 	 3349936	     10512 ns/op	    5675 B/op	      78 allocs/op
BenchmarkJSON/with_standard_library/standard_ReadJSON_-_large
BenchmarkJSON/with_standard_library/standard_ReadJSON_-_large-4         	  105868	    322861 ns/op	  87.09 MB/s	    4448 B/op	     148 allocs/op
BenchmarkJSON/with_standard_library/standard_WriteJSON_-_large
BenchmarkJSON/with_standard_library/standard_WriteJSON_-_large-4        	 1516041	     23693 ns/op	 204.15 MB/s	    4882 B/op	       2 allocs/op
BenchmarkJSON/with_standard_library/standard_FromDynamicJSON_-_large
BenchmarkJSON/with_standard_library/standard_FromDynamicJSON_-_large-4  	  200446	    165414 ns/op	   90246 B/op	    1325 allocs/op
BenchmarkJSON/with_easyjson_library
BenchmarkJSON/with_easyjson_library/easyjson_ReadJSON_-_small
BenchmarkJSON/with_easyjson_library/easyjson_ReadJSON_-_small-4         	25455423	      1218 ns/op	 118.21 MB/s	     208 B/op	       6 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_WriteJSON_-_small
BenchmarkJSON/with_easyjson_library/easyjson_WriteJSON_-_small-4        	40651999	       881.9 ns/op	 149.67 MB/s	     736 B/op	       5 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_FromDynamicJSON_-_small
BenchmarkJSON/with_easyjson_library/easyjson_FromDynamicJSON_-_small-4  	 7403109	      4788 ns/op	    2088 B/op	      34 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_ReadJSON_-_medium
BenchmarkJSON/with_easyjson_library/easyjson_ReadJSON_-_medium-4        	 3000784	     11842 ns/op	 174.80 MB/s	     272 B/op	      10 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_WriteJSON_-_medium
BenchmarkJSON/with_easyjson_library/easyjson_WriteJSON_-_medium-4       	24347379	      1365 ns/op	 231.47 MB/s	     912 B/op	       5 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_FromDynamicJSON_-_medium
BenchmarkJSON/with_easyjson_library/easyjson_FromDynamicJSON_-_medium-4 	 3389486	     10448 ns/op	    6234 B/op	      80 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_ReadJSON_-_large
BenchmarkJSON/with_easyjson_library/easyjson_ReadJSON_-_large-4         	  283608	    114098 ns/op	 246.44 MB/s	    3921 B/op	     133 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_WriteJSON_-_large
BenchmarkJSON/with_easyjson_library/easyjson_WriteJSON_-_large-4        	 2948630	     12013 ns/op	 402.64 MB/s	    5564 B/op	       9 allocs/op
BenchmarkJSON/with_easyjson_library/easyjson_FromDynamicJSON_-_large
BenchmarkJSON/with_easyjson_library/easyjson_FromDynamicJSON_-_large-4  	  210428	    154746 ns/op	   90817 B/op	    1331 allocs/op
PASS
ok  	github.com/go-openapi/swag/jsonutils/adapters/testintegration/benchmarks	615.559s
```
