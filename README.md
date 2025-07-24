# Go Cryptography

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/crypto.svg)](https://pkg.go.dev/golang.org/x/crypto)

This repository holds supplementary Go cryptography packages.

## Report Issues / Send Patches

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://go.dev/doc/contribute.

The git repository is https://go.googlesource.com/crypto.

The main issue tracker for the crypto repository is located at
https://go.dev/issues. Prefix your issue with "x/crypto:" in the
subject line, so it is easy to find.

Note that contributions to the cryptography package receive additional scrutiny
due to their sensitive nature. Patches may take longer than normal to receive
feedback.

# Stateful SHA3 Keccak-256 Usage

This document describes the enhanced `NewLegacyKeccak256()` function that supports setting and exporting internal state.

## Overview

The SHA3 Keccak-256 implementation now supports:
1. **Setting initial state** when creating a hash instance
2. **Exporting internal state** from an existing hash instance
3. **Importing state** into a hash instance

## API Changes

### Original Function (Unchanged)

```go
// NewLegacyKeccak256 creates a new Keccak-256 hash.
func NewLegacyKeccak256() hash.Hash
```

### New Function

```go
// NewLegacyKeccak256WithState creates a new Keccak-256 hash with optional initial state.
// Returns a StatefulHash interface that supports state export/import operations.
func NewLegacyKeccak256WithState(initialState []byte) (StatefulHash, error)
```

### New Interface

```go
// StatefulHash extends hash.Hash with methods to export and import internal state
type StatefulHash interface {
    hash.Hash
    // ExportState returns the internal state of the hash function
    ExportState() []byte
    // ImportState sets the internal state of the hash function
    ImportState(state []byte) error
}
```

## Usage Examples

### Basic Usage (Backward Compatible)

```go
// Create a hash the traditional way (no changes needed)
h := sha3.NewLegacyKeccak256()
h.Write([]byte("data"))
result := h.Sum(nil)
```

### Setting Initial State

```go
// Create a hash with exported state from another hash
h1 := sha3.NewLegacyKeccak256()
h1.Write([]byte("partial data"))

// Export the state using type assertion
state := h1.(interface {
    ExportState() []byte
}).ExportState()

// Create a new hash with the exported state
h2, err := sha3.NewLegacyKeccak256WithState(state)
if err != nil {
    log.Fatal(err)
}
h2.Write([]byte("more data"))
result := h2.Sum(nil)
```

### Using StatefulHash Interface

```go
// Create a stateful hash
h1, err := sha3.NewLegacyKeccak256WithState(nil)
if err != nil {
    log.Fatal(err)
}

h1.Write([]byte("data"))

// Export state using interface method
state := h1.ExportState()

// Create another stateful hash and import state
h2, err := sha3.NewLegacyKeccak256WithState(nil)
if err != nil {
    log.Fatal(err)
}

err = h2.ImportState(state)
if err != nil {
    log.Fatal(err)
}

// Both hashes now have the same internal state
```

## Use Cases

1. **Checkpointing**: Save hash state at specific points for later resumption
2. **Parallel Processing**: Distribute hash computation across multiple goroutines
3. **State Persistence**: Save hash state to disk and restore later
4. **Testing**: Create reproducible hash states for testing purposes

## Backward Compatibility

All existing code using `NewLegacyKeccak256()` will continue to work without any changes. The function signature uses variadic parameters to maintain compatibility.

## State Format

The exported state is a binary representation of the internal hash state, including:
- The sponge state (1600 bits / 200 bytes)
- Buffer position and rate information
- Domain separation byte
- Sponge direction (absorbing/squeezing)

The state format is compatible with the existing `MarshalBinary`/`UnmarshalBinary` methods.