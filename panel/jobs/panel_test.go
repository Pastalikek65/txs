package jobs

import (
	"testing"

	"github.com/Pastalikek65/txs/store"
)

func TestPushAndRunOnce(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/db")
	defer s.Close()
	id, err := Push(s, "echo hi")
	if err != nil { t.Fatal(err) }
	if id == 0 { t.Fatal("no id") }
	rc, err := RunOnce(s, "")
	if err != nil { t.Fatal(err) }
	if rc != 0 { t.Fatalf("want 0 got %d", rc) }
	jobs, _ := s.ListJobs("done")
	if len(jobs) != 1 { t.Fatalf("want done 1 got %d", len(jobs)) }
}

func TestPushEmpty(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/db")
	defer s.Close()
	if _, err := Push(s, ""); err == nil { t.Fatal("want error for empty") }
	if _, err := Push(s, "   "); err == nil { t.Fatal("want error for whitespace") }
}
