package portcache_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portcache"
)

func TestSet_Get_HitWithinTTL(t *testing.T) {
	c := portcache.New(5 * time.Second)
	c.Set("localhost", []int{80, 443})

	e, ok := c.Get("localhost")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(e.Ports) != 2 || e.Ports[0] != 80 || e.Ports[1] != 443 {
		t.Fatalf("unexpected ports: %v", e.Ports)
	}
}

func TestGet_MissWhenExpired(t *testing.T) {
	c := portcache.New(10 * time.Millisecond)
	c.Set("localhost", []int{22})

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("localhost")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestGet_MissWhenZeroTTL(t *testing.T) {
	c := portcache.New(0)
	c.Set("localhost", []int{8080})

	_, ok := c.Get("localhost")
	if ok {
		t.Fatal("expected cache miss with zero TTL")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := portcache.New(5 * time.Second)
	c.Set("host-a", []int{3000})
	c.Invalidate("host-a")

	_, ok := c.Get("host-a")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	c := portcache.New(10 * time.Millisecond)
	c.Set("host-a", []int{80})
	c.Set("host-b", []int{443})

	time.Sleep(20 * time.Millisecond)
	c.Purge()

	if c.Len() != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", c.Len())
	}
}

func TestPurge_KeepsValidEntries(t *testing.T) {
	c := portcache.New(5 * time.Second)
	c.Set("host-a", []int{80})
	c.Set("host-b", []int{443})

	c.Purge()

	if c.Len() != 2 {
		t.Fatalf("expected 2 entries after purge, got %d", c.Len())
	}
}

func TestSet_ScannedAt_IsRecent(t *testing.T) {
	before := time.Now()
	c := portcache.New(5 * time.Second)
	c.Set("localhost", []int{9090})
	after := time.Now()

	e, ok := c.Get("localhost")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if e.ScannedAt.Before(before) || e.ScannedAt.After(after) {
		t.Errorf("ScannedAt %v not within expected range [%v, %v]", e.ScannedAt, before, after)
	}
}
