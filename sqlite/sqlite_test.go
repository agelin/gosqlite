package sqlite

import (
	"testing"
)

func TestOpenMemory(t *testing.T) {
	_, open_err := Open(":memory:")
  if open_err != nil {
    t.Error("Failed to open in memory only database")
  }
}

