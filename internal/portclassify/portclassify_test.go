package portclassify_test

import (
	"testing"

	"github.com/example/portwatch/internal/portclassify"
)

func TestClassify_WebPort(t *testing.T) {
	c := portclassify.New()
	for _, p := range []int{80, 443, 8080} {
		if got := c.Classify(p); got != portclassify.CategoryWeb {
			t.Errorf("port %d: want %s, got %s", p, portclassify.CategoryWeb, got)
		}
	}
}

func TestClassify_DatabasePort(t *testing.T) {
	c := portclassify.New()
	for _, p := range []int{3306, 5432, 6379, 27017} {
		if got := c.Classify(p); got != portclassify.CategoryDatabase {
			t.Errorf("port %d: want %s, got %s", p, portclassify.CategoryDatabase, got)
		}
	}
}

func TestClassify_SystemPort(t *testing.T) {
	c := portclassify.New()
	for _, p := range []int{22, 25, 53, 1023} {
		if got := c.Classify(p); got != portclassify.CategorySystem {
			t.Errorf("port %d: want %s, got %s", p, portclassify.CategorySystem, got)
		}
	}
}

func TestClassify_UserPort(t *testing.T) {
	c := portclassify.New()
	if got := c.Classify(1024); got != portclassify.CategoryUser {
		t.Errorf("want user, got %s", got)
	}
	if got := c.Classify(49151); got != portclassify.CategoryUser {
		t.Errorf("want user, got %s", got)
	}
}

func TestClassify_DynamicPort(t *testing.T) {
	c := portclassify.New()
	if got := c.Classify(49152); got != portclassify.CategoryDynamic {
		t.Errorf("want dynamic, got %s", got)
	}
	if got := c.Classify(65535); got != portclassify.CategoryDynamic {
		t.Errorf("want dynamic, got %s", got)
	}
}

func TestClassify_InvalidPort_ReturnsUnknown(t *testing.T) {
	c := portclassify.New()
	for _, p := range []int{-1, 65536, 99999} {
		if got := c.Classify(p); got != portclassify.CategoryUnknown {
			t.Errorf("port %d: want unknown, got %s", p, got)
		}
	}
}

func TestClassifyAll_ReturnsMappedCategories(t *testing.T) {
	c := portclassify.New()
	result := c.ClassifyAll([]int{80, 22, 3306})
	if result[80] != portclassify.CategoryWeb {
		t.Errorf("80: want web, got %s", result[80])
	}
	if result[22] != portclassify.CategorySystem {
		t.Errorf("22: want system, got %s", result[22])
	}
	if result[3306] != portclassify.CategoryDatabase {
		t.Errorf("3306: want database, got %s", result[3306])
	}
}
