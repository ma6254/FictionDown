package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	n := 0
	nn := 5
	err := Retry(nn, 500*time.Millisecond, func() error {
		n++
		t.Logf("retry %d/%d", n, nn)
		return fmt.Errorf("Test Error")
	})
	if err == nil {
		t.Fatalf("not return error")
	}
	if n != nn {
		t.Fatalf("retry count error %d/%d", n, nn)
	}
}
