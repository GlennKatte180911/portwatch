package eventlog_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/eventlog"
)

func seedLog(t *testing.T) *eventlog.EventLog {
	t.Helper()
	log := eventlog.New(tempLogPath(t))
	_ = log.Append([]int{80}, nil)
	_ = log.Append([]int{443}, nil)
	_ = log.Append([]int{8080}, nil)
	return log
}

func TestQuery_Since_FiltersOldEntries(t *testing.T) {
	log := eventlog.New(tempLogPath(t))
	_ = log.Append([]int{80}, nil)
	time.Sleep(2 * time.Millisecond)
	cutoff := time.Now().UTC()
	_ = log.Append([]int{443}, nil)

	entries, err := log.Query(eventlog.QueryOptions{Since: cutoff})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after cutoff, got %d", len(entries))
	}
	if entries[0].Added[0] != 443 {
		t.Errorf("expected port 443, got %d", entries[0].Added[0])
	}
}

func TestQuery_Limit_CapsResults(t *testing.T) {
	log := seedLog(t)
	entries, err := log.Query(eventlog.QueryOptions{Limit: 2})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLatest_ReturnsLastN(t *testing.T) {
	log := seedLog(t)
	entries, err := log.Latest(2)
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].Added[0] != 8080 {
		t.Errorf("expected last entry port 8080, got %d", entries[1].Added[0])
	}
}

func TestLatest_NGreaterThanTotal_ReturnsAll(t *testing.T) {
	log := seedLog(t)
	entries, err := log.Latest(100)
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}
