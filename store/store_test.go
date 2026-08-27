package store

import "testing"

func TestOpenAndUpsertRepos(t *testing.T) {
	s, err := Open(t.TempDir() + "/txs.db")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer s.Close()
	if err := s.UpsertRepos([]Repo{{Name: "a", Stars: 10}}); err != nil {
		t.Fatal(err)
	}
	got, err := s.ListRepos("stars")
	if err != nil {
		t.Fatalf("ListRepos failed: %v", err)
	}
	if len(got) != 1 || got[0].Name != "a" {
		t.Fatal(got)
	}
}

func TestPushAndListJobs(t *testing.T) {
	s, err := Open(t.TempDir() + "/txs.db")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer s.Close()
	id, err := s.PushJob("echo hi")
	if err != nil {
		t.Fatalf("PushJob failed: %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id 1 got %d", id)
	}
	jobs, err := s.ListJobs("pending")
	if err != nil {
		t.Fatalf("ListJobs failed: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatal("pending")
	}
}
