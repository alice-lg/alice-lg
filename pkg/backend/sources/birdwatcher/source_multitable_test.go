package birdwatcher

import (
	"testing"
)

func TestGetMasterPipeName(t *testing.T) {
	config := Config{
		PipeProtocolPrefix: "pp",
		PeerTablePrefix:    "pb",
	}

	bw := &MultiTableBirdwatcher{
		GenericBirdwatcher: GenericBirdwatcher{
			config: config,
		},
	}

	peerProto := "pb_0200_as123456"
	expected := "pp_0200_as123456"
	if res := bw.getMasterPipeName(peerProto); res != expected {
		t.Error("Expected:", peerProto, "but got:", res)
	}

}
