# gx - Modern Go Script Runner

Run Go files like scripts with smart caching and auto-reload.

# gx is fast because it avoids the Go toolchain entirely on repeated runs. This removes ~170ms of fixed overhead per execution at the cost of weaker dependency invalidation guarantees.

## Performance

Benchmarks show significant speedup over `go run`:

```bash
slyt3000@slyt3000:~/go-run-bench$ strace -c gx run main.go
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ------------------
 33,78    0,006128        2042         3         1 waitid
 18,66    0,003385          29       114           rt_sigaction
  7,80    0,001415         128        11           rt_sigprocmask
  6,16    0,001117        1117         1           execve
  5,62    0,001019          92        11           read
  4,78    0,000867          39        22           mmap
  4,37    0,000793          36        22           prctl
  3,73    0,000676         135         5           clone
  3,18    0,000576          52        11           close
  2,99    0,000542          77         7         1 openat
  1,62    0,000294          58         5           newfstatat
  1,56    0,000283          20        14           fcntl
  1,41    0,000255          51         5           futex
  0,57    0,000104          52         2           prlimit64
  0,56    0,000101          33         3         2 epoll_ctl
  0,52    0,000095          47         2           madvise
  0,52    0,000094          47         2           rt_sigreturn
  0,29    0,000053          26         2           getpid
  0,29    0,000053          53         1           epoll_create1
  0,24    0,000043          43         1           eventfd2
  0,21    0,000038          38         1           arch_prctl
  0,21    0,000038          38         1           pipe2
  0,20    0,000037          37         1           sched_getaffinity
  0,20    0,000036          36         1           pidfd_open
  0,18    0,000032          16         2           sigaltstack
  0,17    0,000031          31         1           fstat
  0,12    0,000021          21         1           gettid
  0,08    0,000014          14         1           pidfd_send_signal
------ ----------- ----------- --------- --------- ------------------
100,00    0,018140          71       253         4 total
```
```bash
slyt3000@slyt3000:~/go-run-bench$ strace -c go run main.go
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ------------------
 42,73    0,047695         324       147        43 futex
 18,98    0,021179          49       426        47 newfstatat
 10,16    0,011341        2268         5         1 waitid
  5,40    0,006022           6       893           fcntl
  4,64    0,005179          12       427       191 openat
  4,08    0,004551         119        38           mmap
  3,71    0,004139          16       251           close
  2,71    0,003027         432         7           clone
  2,38    0,002661           6       422           read
  1,27    0,001414          45        31           getdents64
  1,23    0,001369          12       114           rt_sigaction
  1,02    0,001137          35        32         6 rt_sigreturn
  0,62    0,000697           3       232       217 epoll_ctl
  0,19    0,000210           8        26           fstat
  0,15    0,000167          11        15           getpid
  0,14    0,000161           8        19           rt_sigprocmask
  0,11    0,000126         126         1           mkdirat
  0,10    0,000117          58         2           faccessat2
  0,07    0,000075           6        12           tgkill
  0,06    0,000068          13         5           pipe2
  0,04    0,000047          23         2           sigaltstack
  0,03    0,000035           1        24           prctl
  0,03    0,000033           1        32           sched_yield
  0,03    0,000032          16         2           prlimit64
  0,02    0,000027           3         7           pread64
  0,02    0,000021          10         2           flock
  0,02    0,000017          17         1           nanosleep
  0,01    0,000016           4         4           epoll_pwait
  0,01    0,000011           5         2           readlinkat
  0,01    0,000010          10         1           gettid
  0,01    0,000010          10         1           pidfd_open
  0,01    0,000008           8         1           uname
  0,01    0,000007           7         1           pidfd_send_signal
  0,00    0,000000           0         2           madvise
  0,00    0,000000           0         1           execve
  0,00    0,000000           0         1           arch_prctl
  0,00    0,000000           0         1           sched_getaffinity
  0,00    0,000000           0         1           eventfd2
  0,00    0,000000           0         1           epoll_create1
------ ----------- ----------- --------- --------- ------------------
100,00    0,111609          34      3192       505 total

```
```bash
gx:
real	0m0,151s
user	0m0,059s
sys	0m0,123s

go run:
real	0m3,511s
user	0m4,795s
sys	0m2,656s

```





## Features

-  Smart caching - compile once, run instantly
-  Auto-reload - watch mode for development  
-  Cache management
-  Proper exit codes

## Commands
```bash
gx run <script> [args]    # Run a Go script (with caching)
gx watch <script> [args]  # Auto-reload on file changes
gx clean                  # Clear cache
gx version                # Show version
```

## Installation
```bash
go install github.com/slyt3/gx@latest
```

## Examples

### Run a script
```bash
gx run server.go --port 8080
```

### Development with auto-reload
```bash
gx watch server.go --port 8080
# Edit server.go and save → automatically restarts!
```

### Clear cache
```bash
gx clean
```

## How it works

1. **First run:** Compiles and caches the binary
2. **Subsequent runs:** Uses cached binary (validated by file size + modTime)
3. **Watch mode:** Monitors file changes and auto-restarts

## Why gx?

- Faster than `go run` for repeated executions
- No manual recompilation during development
- Works great for scripts, servers, and CLI tools

## License

MIT
