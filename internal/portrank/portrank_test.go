package portrank

import (
	"testing"
)

func TestNew_ContainsDefaultPorts(t *testing.T) {
	r := New()
	if got := r.Get(22); got != RankHigh {
		t.Errorf("expected RankHigh for port 22, got %v", got)
	}
	if got := r.Get(443); got != RankLow {
		t.Errorf("expected RankLow for port 443, got %v", got)
	}
}

func TestGet_UnknownPort_ReturnsLow(t *testing.T) {
	r := New()
	if got := r.Get(9999); got != RankLow {
		t.Errorf("expected RankLow for unknown port, got %v", got)
	}
}

func TestSet_OverwritesRank(t *testing.T) {
	r := New()
	r.Set(80, RankHigh)
	if got := r.Get(80); got != RankHigh {
		t.Errorf("expected RankHigh after Set, got %v", got)
	}
}

func TestSet_AddsNewPort(t *testing.T) {
	r := New()
	r.Set(12345, RankMedium)
	if got := r.Get(12345); got != RankMedium {
		t.Errorf("expected RankMedium for newly added port, got %v", got)
	}
}

func TestAnnotate_ReturnsMappedRanks(t *testing.T) {
	r := New()
	annotated := r.Annotate([]int{22, 80, 9999})
	if annotated[22] != RankHigh {
		t.Errorf("expected RankHigh for 22, got %v", annotated[22])
	}
	if annotated[80] != RankLow {
		t.Errorf("expected RankLow for 80, got %v", annotated[80])
	}
	if annotated[9999] != RankLow {
		t.Errorf("expected RankLow for unknown 9999, got %v", annotated[9999])
	}
}

func TestAnnotate_EmptySlice(t *testing.T) {
	r := New()
	if got := r.Annotate([]int{}); len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMaxRank_ReturnsHighest(t *testing.T) {
	r := New()
	// 22 = High, 80 = Low
	if got := r.MaxRank([]int{80, 22}); got != RankHigh {
		t.Errorf("expected RankHigh, got %v", got)
	}
}

func TestMaxRank_EmptySlice_ReturnsLow(t *testing.T) {
	r := New()
	if got := r.MaxRank([]int{}); got != RankLow {
		t.Errorf("expected RankLow for empty slice, got %v", got)
	}
}

func TestRank_String(t *testing.T) {
	cases := []struct {
		rank Rank
		want string
	}{
		{RankLow, "low"},
		{RankMedium, "medium"},
		{RankHigh, "high"},
	}
	for _, tc := range cases {
		if got := tc.rank.String(); got != tc.want {
			t.Errorf("Rank(%d).String() = %q, want %q", tc.rank, got, tc.want)
		}
	}
}
