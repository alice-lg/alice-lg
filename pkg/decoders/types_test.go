package decoders

import (
	"testing"
	"time"
)

func TestDurationTimeframe(t *testing.T) {
	d := DurationTimeframe("01w1d5h", 0)
	if d != 197*time.Hour {
		t.Error("unexpected", d)
	}

	d = DurationTimeframe("5d20h05m", 0)
	if d != 140*time.Hour+5*time.Minute {
		t.Error("unexpected", d)
	}

	d = DurationTimeframe("01:02:03", 0)
	if d != 1*time.Hour+2*time.Minute+3*time.Second {
		t.Error("unexpected", d)
	}

}
