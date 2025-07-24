// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha3

import (
	"bytes"
	"testing"
)

// TestStateExportImport tests the ExportState and ImportState functionality
func TestStateExportImport(t *testing.T) {
	// Create a hash and write some data
	h1 := NewLegacyKeccak256()
	data := []byte("test data for state export/import")
	h1.Write(data)

	// Export the state
	state1 := h1.(*state).ExportState()

	// Create a new hash with the exported state
	h2, err := NewLegacyKeccak256WithState(state1)
	if err != nil {
		t.Fatalf("Failed to create hash with state: %v", err)
	}

	// Both hashes should produce the same result
	sum1 := h1.Sum(nil)
	sum2 := h2.Sum(nil)

	if !bytes.Equal(sum1, sum2) {
		t.Errorf("Hashes with same state produced different results: %x vs %x", sum1, sum2)
	}

	// Write additional data to both hashes
	additionalData := []byte("additional data after state export/import")
	h1.Write(additionalData)
	h2.Write(additionalData)

	// Compare results after additional writes
	sum1After := h1.Sum(nil)
	sum2After := h2.Sum(nil)

	if !bytes.Equal(sum1After, sum2After) {
		t.Errorf("Hashes with same state produced different results after additional writes: %x vs %x", sum1After, sum2After)
	}
}

// TestStatefulHashInterface tests the StatefulHash interface methods
func TestStatefulHashInterface(t *testing.T) {
	// Create a stateful hash
	h1, err := NewLegacyKeccak256WithState(nil)
	if err != nil {
		t.Fatalf("Failed to create stateful hash: %v", err)
	}

	// Write some data
	data := []byte("interface test data")
	h1.Write(data)

	// Export state using the interface method
	exportedState := h1.ExportState()

	// Create another stateful hash and import the state
	h2, err := NewLegacyKeccak256WithState(nil)
	if err != nil {
		t.Fatalf("Failed to create second stateful hash: %v", err)
	}

	err = h2.ImportState(exportedState)
	if err != nil {
		t.Fatalf("Failed to import state: %v", err)
	}

	// Both hashes should produce the same result
	sum1 := h1.Sum(nil)
	sum2 := h2.Sum(nil)

	if !bytes.Equal(sum1, sum2) {
		t.Errorf("Hashes with imported state produced different results: %x vs %x", sum1, sum2)
	}
}

// TestBackwardCompatibility ensures the original NewLegacyKeccak256() still works
func TestBackwardCompatibility(t *testing.T) {
	// Test that calling without arguments still works
	h1 := NewLegacyKeccak256()
	h2 := NewLegacyKeccak256()

	data := []byte("backward compatibility test")
	h1.Write(data)
	h2.Write(data)

	sum1 := h1.Sum(nil)
	sum2 := h2.Sum(nil)

	if !bytes.Equal(sum1, sum2) {
		t.Errorf("Backward compatibility broken: %x vs %x", sum1, sum2)
	}
}
