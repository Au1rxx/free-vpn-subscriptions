package store

import (
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestShouldRetryParseBatch(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		attempt int
		want    bool
	}{
		{name: "deadlock", err: fmt.Errorf("persist: %w", &mysql.MySQLError{Number: 1213}), attempt: 1, want: true},
		{name: "lock wait timeout", err: &mysql.MySQLError{Number: 1205}, attempt: 2, want: true},
		{name: "retry budget exhausted", err: &mysql.MySQLError{Number: 1213}, attempt: parseBatchMaxAttempts, want: false},
		{name: "ordinary error", err: fmt.Errorf("invalid node"), attempt: 1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetryParseBatch(tt.err, tt.attempt); got != tt.want {
				t.Fatalf("shouldRetryParseBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
