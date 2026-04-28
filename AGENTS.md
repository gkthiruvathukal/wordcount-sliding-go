# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Overview

A Go implementation of a sliding-window word frequency counter. Given a text stream on stdin, it tracks the top-N most frequent words in the last-K words seen, printing a summary every M steps. The program has two orthogonal axes of variation — pipeline style and sliding-window implementation — selectable at runtime via flags.

## Build and Run

```bash
# Build all apps and run help + tests
./do-app.sh

# Test libraries only
./do-libs.sh

# Run acceptance/performance benchmarks against hamlet data (requires build/ to exist first)
./do-hamlet.sh
```

Build output goes to `build/sliding_wordcount`. The Go workspace (`go.work`) ties together the two modules — run `go` commands from the repo root or use the scripts above.

### Run a single test

```bash
# Library tests
cd libs/cq && go test -v

# App tests
cd app/sliding_wordcount && go test -v
```

### Run the binary directly

```bash
cat data/hamlet.txt | ./build/sliding_wordcount [flags]
```

Key flags (all optional):

| Flag | Default | Meaning |
|---|---|---|
| `-last-n-words` | 1000 | Sliding window size |
| `-show-top` | 10 | Top-N words to print |
| `-min-word-length` | 5 | Ignore shorter words |
| `-every-steps` | 1000 | Print frequency (every N words) |
| `-ignore-case` | false | Fold to lowercase |
| `-go-routines` | false | Use goroutine pipeline (default: imperative loop) |
| `-circular-queue` | true | Use `CircularQueue` for sliding window (default) |
| `-channel` | false | Use buffered channel for sliding window |
| `-channel-size` | 10 | Buffered channel size for inter-stage communication (goroutine mode only) |

`-circular-queue` and `-channel` are mutually exclusive. If neither is given, `-circular-queue` is used.

## Architecture

### Module layout (Go workspace)

```
go.work                          # workspace linking both modules
libs/cq/                         # module "cq" — reusable generic circular queue
  cqueue_generic.go              # CircularQueue[T comparable] — the active implementation
  cqueue.go                      # CQueueString — original non-generic version, kept for reference
  cqueue_test.go                 # unit + property-based tests (testing/quick)
app/sliding_wordcount/           # module "sliding_wordcount" — main binary
  sliding_wordcount.go           # all application logic
data/
  hamlet.txt                     # source text
  hamlet-10-copies.txt           # 10x copy used for benchmarking
```

### Two orthogonal axes

**Pipeline style** (`-go-routines`):

- **Imperative** (default): single-goroutine loop — read line → regex match → update sliding window → print on step boundary.
- **Goroutine pipeline**: three stages connected by buffered channels: `generateWords` → `filterBasedOnCommandLine` → sliding analysis. Each stage runs in its own goroutine.

**Sliding window implementation** (`-circular-queue` / `-channel`):

- **`CircularQueue`** (default): explicit ring buffer from `libs/cq`. `IsFull()` / `Enqueue()` / `Dequeue()` map directly onto the algorithm. Semantically transparent and the preferred implementation.
- **Buffered channel**: uses a `chan string` of capacity `lastNWords` as a FIFO. Fullness is detected via `len(window) == cap(window)`. Functionally equivalent but unconventional — channels are designed for goroutine communication, not local queues. Included as an educational comparison.

This gives four selectable combinations, all in `driver()`:

| `-go-routines` | `-circular-queue` / `-channel` | Function called |
|---|---|---|
| false | circular-queue | `imperativePipelineCQ` |
| false | channel | `imperativePipelineChannel` |
| true | circular-queue | `generateWords` → `filterBasedOnCommandLine` → `slidingAnalysisCQ` |
| true | channel | `generateWords` → `filterBasedOnCommandLine` → `slidingAnalysisChannel` |

### `CircularQueue[T]` (`libs/cq/cqueue_generic.go`)

Fixed-capacity ring buffer with `storePos` / `retrievePos` indices and a `count`. `IsFull()` is checked before every `Enqueue`; the sliding-window logic always calls `Dequeue` first when full, then `Enqueue` the new word. `cqueue.go` retains the original string-only non-generic version for reference but is not used by the app.

### Design note: why `CircularQueue` over the channel approach

Both implementations produce identical output and have near-identical performance (~0.09s imperative, ~0.15s goroutine pipeline on 10 copies of Hamlet). The `CircularQueue` is preferred because:

- Its API (`IsFull`, `Enqueue`, `Dequeue`) reads as a direct description of the sliding-window algorithm.
- The channel approach relies on `len(ch) == cap(ch)` as a fullness check, which is only safe in single-goroutine contexts and would be a silent race if the function were ever made concurrent.
- Channels are semantically for goroutine communication; using one as a local queue surprises readers, especially when the goroutine pipeline already uses channels for inter-stage communication.
