package eventlog

import "time"

// QueryOptions controls which entries are returned by Query.
type QueryOptions struct {
	// Since filters entries to those at or after this time.
	// Zero value means no lower bound.
	Since time.Time
	// Limit caps the number of returned entries (0 = unlimited).
	Limit int
}

// Query returns entries from the log that match the given options.
func (l *EventLog) Query(opts QueryOptions) ([]Entry, error) {
	all, err := l.Load()
	if err != nil {
		return nil, err
	}

	var result []Entry
	for _, e := range all {
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		result = append(result, e)
		if opts.Limit > 0 && len(result) >= opts.Limit {
			break
		}
	}
	return result, nil
}

// Latest returns the most recent n entries from the log.
func (l *EventLog) Latest(n int) ([]Entry, error) {
	all, err := l.Load()
	if err != nil {
		return nil, err
	}
	if n <= 0 || len(all) == 0 {
		return all, nil
	}
	if n >= len(all) {
		return all, nil
	}
	return all[len(all)-n:], nil
}
