package portannounce_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/portannounce"
	"github.com/user/portwatch/internal/portdiff"
)

func sampleDiff() portdiff.Diff {
	return portdiff.Diff{Added: []int{8080}, Removed: []int{9090}}
}

func TestAnnounce_CallsRegisteredHandlers(t *testing.T) {
	a := portannounce.New()

	var got portdiff.Diff
	a.Subscribe(func(d portdiff.Diff) { got = d })

	a.Announce(sampleDiff())

	if len(got.Added) != 1 || got.Added[0] != 8080 {
		t.Errorf("expected Added=[8080], got %v", got.Added)
	}
}

func TestAnnounce_SkipsEmptyDiff(t *testing.T) {
	a := portannounce.New()
	called := false
	a.Subscribe(func(d portdiff.Diff) { called = true })

	a.Announce(portdiff.Diff{})

	if called {
		t.Error("handler should not be called for empty diff")
	}
}

func TestUnsubscribe_RemovesHandler(t *testing.T) {
	a := portannounce.New()
	called := false

	unsub := a.Subscribe(func(d portdiff.Diff) { called = true })
	unsub()

	a.Announce(sampleDiff())

	if called {
		t.Error("unsubscribed handler should not be called")
	}
}

func TestCount_ReflectsActiveSubscribers(t *testing.T) {
	a := portannounce.New()

	if a.Count() != 0 {
		t.Fatalf("expected 0 subscribers, got %d", a.Count())
	}

	unsub1 := a.Subscribe(func(portdiff.Diff) {})
	a.Subscribe(func(portdiff.Diff) {})

	if a.Count() != 2 {
		t.Fatalf("expected 2 subscribers, got %d", a.Count())
	}

	unsub1()

	if a.Count() != 1 {
		t.Errorf("expected 1 subscriber after unsub, got %d", a.Count())
	}
}

func TestAnnounce_ConcurrentSafe(t *testing.T) {
	a := portannounce.New()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsub := a.Subscribe(func(portdiff.Diff) {})
			a.Announce(sampleDiff())
			unsub()
		}()
	}

	wg.Wait()
}
