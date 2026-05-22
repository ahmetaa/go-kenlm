# go-kenlm

Go bindings for [KenLM](https://github.com/kpu/kenlm) — load an n-gram language
model (`.binary` or `.arpa`) and score text under it.

Supported platforms: **linux/amd64**, **darwin/arm64**.

## Install

```bash
go get github.com/ahmetaa/go-kenlm
```

You also need the KenLM static libraries on disk. Two options:

### Option A — download prebuilt libs (fastest)

Grab the tarball for your platform from the
[latest release](https://github.com/ahmetaa/go-kenlm/releases/latest) and unpack
it so that `third_party/kenlm/build/lib/libkenlm.a` exists relative to the
go-kenlm source:

```bash
# inside your clone of go-kenlm (or the module cache copy)
mkdir -p third_party/kenlm/build
tar -xzf go-kenlm-libs-linux-amd64.tar.gz \
    --strip-components=1 -C third_party/kenlm/build
```

### Option B — build from source

System prerequisites:

- **macOS:** `brew install boost cmake`
- **Linux (Debian/Ubuntu):** `sudo apt install cmake build-essential libboost-program-options-dev libboost-thread-dev zlib1g-dev libbz2-dev liblzma-dev`

Then, from your clone of this repo:

```bash
git submodule update --init --recursive
make build-kenlm
```

This produces `third_party/kenlm/build/lib/libkenlm{,_util}.a`. The build is
one-shot — subsequent `go build`/`go test` calls just link against those
files.

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/ahmetaa/go-kenlm"
)

func main() {
    m, err := kenlm.LoadModel("char-6gram-kenlm.binary")
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    // Whitespace-separated tokens — for a char-LM, split by characters.
    fmt.Println(m.Score("h e l l o"))     // log10 P("h e l l o")
    fmt.Println(m.ScoreChars("hello"))    // convenience: rune-split + join
    fmt.Println(m.Order())                // n-gram order, e.g. 6
}
```

Scores are **log10 probabilities** and include the implicit `<s>` / `</s>`
sentence markers. A `Model` is safe for concurrent reads.

## API

| Function                          | Returns          |
|-----------------------------------|------------------|
| `LoadModel(path string)`          | `*Model, error`  |
| `(*Model).Score(tokens string)`   | `float64` (log10 P) |
| `(*Model).ScoreChars(s string)`   | `float64` (log10 P, rune-split) |
| `(*Model).Order()`                | `int`            |
| `(*Model).Close()`                | —                |

That's the whole surface.

## Troubleshooting

- **`fatal error: 'lm/model.hh' file not found`** — the submodule isn't
  checked out. Run `git submodule update --init --recursive`.
- **`ld: library 'kenlm' not found`** — you skipped the build step. Run
  `make build-kenlm` (Option B) or unpack the prebuilt tarball (Option A).
- **`Could NOT find Boost`** during `make build-kenlm` — install Boost
  (see prerequisites above).

## License

This binding is MIT-licensed. Upstream KenLM is LGPL — see
`third_party/kenlm/LICENSE`.
