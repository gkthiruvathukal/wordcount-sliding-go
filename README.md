# wordcount-sliding-go

A Go implementation of a sliding-window word frequency counter. Given a text stream on stdin, the program tracks the most frequent words within a fixed-size window of recent words, printing a ranked word cloud at regular intervals.

The project exists as an exploration of two orthogonal design choices in Go:

1. **Pipeline style** — imperative loop vs. goroutine/channel pipeline
2. **Sliding window implementation** — explicit circular queue vs. buffered channel

All four combinations are supported and selectable at runtime via flags.

## Background

A sliding window word count works as follows: as words are read from the input stream, a fixed-capacity window of the most recent N words is maintained. When a new word enters and the window is full, the oldest word is evicted and its count is decremented. The result at any point is a word frequency map over only the last N words — not the entire document — which makes it useful for tracking local trends in text.

## Build

Requires Go 1.22+. The project uses a Go workspace (`go.work`) linking two modules: the `cq` library and the `sliding_wordcount` app.

```bash
# Build binary and run tests
./do-app.sh

# Test the cq library only
./do-libs.sh
```

The binary is written to `build/sliding_wordcount`.

## Usage

```bash
cat data/hamlet.txt | ./build/sliding_wordcount [flags]
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `-last-n-words N` | 1000 | Sliding window size (number of words tracked) |
| `-show-top N` | 10 | Number of top words to display |
| `-min-word-length N` | 5 | Ignore words shorter than N characters |
| `-every-steps N` | 1000 | Print word cloud every N words |
| `-ignore-case` | false | Fold all words to lowercase |
| `-go-routines` | false | Use goroutine/channel pipeline (default: imperative loop) |
| `-circular-queue` | true | Use `CircularQueue` for sliding window (default) |
| `-channel` | false | Use buffered channel for sliding window |
| `-channel-size N` | 10 | Channel buffer size for inter-stage communication (goroutine mode) |

`-circular-queue` and `-channel` are mutually exclusive. If neither is given, `-circular-queue` is the default.

### Example

```bash
# Default: imperative loop + circular queue
cat data/hamlet.txt | ./build/sliding_wordcount

# Goroutine pipeline + circular queue, larger window, top 5 words every 500
cat data/hamlet.txt | ./build/sliding_wordcount -go-routines -last-n-words 2000 -show-top 5 -every-steps 500

# Compare channel-based window
cat data/hamlet.txt | ./build/sliding_wordcount -channel
```

## Benchmarking

`do-hamlet.sh` runs the binary against a 10-copy concatenation of Hamlet, sweeping the `-channel-size` parameter from 1 to 10,000,000 for the goroutine pipeline. Timing results are written to `*.txt` files in the project root.

```bash
./do-app.sh   # build first
./do-hamlet.sh
```

## The Two Pipeline Styles

### Imperative loop

The default. A single goroutine reads stdin line by line, applies a regex to extract words, and updates the sliding window and word cloud directly.

**Strengths:**
- Simple, easy to follow top-to-bottom
- No goroutine or channel overhead
- Fastest wall-clock time (~0.09s on 10× Hamlet)

**Weaknesses:**
- All stages (read, filter, analyse) are entangled in one function, making individual stages harder to test or swap independently

### Goroutine pipeline (`-go-routines`)

Three stages connected by buffered channels, each running in its own goroutine:

1. `generateWords` — reads stdin, emits raw words
2. `filterBasedOnCommandLine` — drops short words, optionally lowercases
3. `slidingAnalysis` — maintains the window and word cloud, prints results

**Strengths:**
- Stages are decoupled and independently readable
- Models a functional-style pipeline naturally in Go
- `-channel-size` lets you tune backpressure between stages

**Weaknesses:**
- Goroutine and channel coordination adds overhead (~0.14–0.15s on 10× Hamlet)
- The analysis stage is still single-goroutine, so there is no parallelism in the compute-heavy part

## The Two Sliding Window Implementations

### Circular queue (`-circular-queue`, default)

Implemented in `libs/cq` as `CircularQueue[T comparable]` — a generic, fixed-capacity ring buffer. The API (`IsFull`, `Enqueue`, `Dequeue`) maps directly onto the sliding-window algorithm: check if full, evict the oldest word, enqueue the new one.

**Strengths:**
- Semantically transparent — the data structure name and API match the algorithm's intent
- Well-tested with both unit tests and property-based tests (`testing/quick`)
- Safe to refactor into concurrent contexts without hidden assumptions

**Weaknesses:**
- Requires an explicit library (`libs/cq`) rather than a built-in type

### Buffered channel (`-channel`)

Uses a `chan string` of capacity `lastNWords` as a local FIFO. Fullness is detected via `len(window) == cap(window)`; receiving from the channel evicts the oldest word.

**Strengths:**
- No library dependency — channels are a built-in Go type
- Concise

**Weaknesses:**
- Channels are designed for goroutine communication, not local queues. Using one this way surprises readers, especially in code that already uses channels for inter-stage communication
- `len(ch) == cap(ch)` as a fullness check is only safe in single-goroutine contexts. It is a TOCTOU pattern that would silently race if the function were ever made concurrent
- The circularity and eviction semantics are implicit rather than named

In practice both implementations produce identical output and near-identical performance. The circular queue is the preferred implementation.

## Project Structure

```
go.work                        # Go workspace linking both modules
libs/cq/                       # module "cq" — generic circular queue library
  cqueue_generic.go            # CircularQueue[T] — active implementation
  cqueue.go                    # original string-only version, kept for reference
  cqueue_test.go               # unit + property-based tests
app/sliding_wordcount/         # module "sliding_wordcount" — main binary
  sliding_wordcount.go
data/
  hamlet.txt
  hamlet-10-copies.txt         # used for benchmarking
```

## License

See [LICENSE](LICENSE).
