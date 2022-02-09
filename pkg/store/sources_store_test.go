package store

import (
	"testing"
	"time"
)

func TestGetSourceIDsForRefreshSequential(t *testing.T) {

	s := &SourcesStore{
		refreshParallelism: 1,
		status: map[string]*Status{
			"src1": {
				SourceID: "src1",
			},
			"src2": {
				SourceID: "src2",
			},
		},
	}

	ids := s.GetSourceIDsForRefresh()
	if len(ids) != 1 {
		t.Error("expected 1 id")
	}
	if err := s.LockSource(ids[0]); err != nil {
		t.Error(err)
	}
	lastID := ids[0]

	ids = s.GetSourceIDsForRefresh()
	if len(ids) != 0 {
		t.Error("all concurrent refresh slots should be taken")
	}

	if err := s.RefreshSuccess(lastID); err != nil {
		t.Error(err)
	}

	ids = s.GetSourceIDsForRefresh()
	if len(ids) != 1 {
		t.Error("expected 1 id")
	}

	if ids[0] == lastID {
		t.Error("the next source should have been returned")
	}

}

func TestGetSourceIDsForRefreshParallel(t *testing.T) {

	s := &SourcesStore{
		refreshParallelism: 2,
		status: map[string]*Status{
			"src1": {
				SourceID: "src1",
			},
			"src2": {
				SourceID: "src2",
			},
			"src3": {
				SourceID:         "src3",
				lastRefreshStart: time.Now().UTC(),
			},
		},
	}

	ids := s.GetSourceIDsForRefresh()
	if len(ids) != 2 {
		t.Error("expected 2 ids")
	}
	for _, id := range ids {
		if err := s.LockSource(id); err != nil {
			t.Error(err)
		}

		if id == "src3" {
			t.Error("unexpected src3")
		}
	}

	nextIds := s.GetSourceIDsForRefresh()
	if len(nextIds) != 0 {
		t.Error("all concurrent refresh slots should be taken")
	}

	for _, id := range ids {
		if err := s.RefreshSuccess(id); err != nil {
			t.Error(err)
		}
	}

	ids = s.GetSourceIDsForRefresh()
	t.Log(ids)
	t.Log(s.status["src1"])
	if len(ids) != 2 {
		t.Error("expected 2 id")
	}

	if ids[0] != "src3" {
		t.Error("expected src3 to be least refreshed")
	}
}
