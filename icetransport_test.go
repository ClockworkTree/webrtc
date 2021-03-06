// +build !js

package webrtc

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/pion/transport/test"
	"github.com/stretchr/testify/assert"
)

func TestICETransport_OnSelectedCandidatePairChange(t *testing.T) {
	report := test.CheckRoutines(t)
	defer report()

	lim := test.TimeOut(time.Second * 30)
	defer lim.Stop()

	pcOffer, pcAnswer, err := newPair()
	if err != nil {
		t.Fatal(err)
	}

	iceComplete := make(chan bool)
	pcAnswer.OnICEConnectionStateChange(func(iceState ICEConnectionState) {
		if iceState == ICEConnectionStateConnected {
			time.Sleep(3 * time.Second)
			close(iceComplete)
		}
	})

	senderCalledCandidateChange := int32(0)
	pcOffer.SCTP().Transport().ICETransport().OnSelectedCandidatePairChange(func(pair *ICECandidatePair) {
		atomic.StoreInt32(&senderCalledCandidateChange, 1)
	})

	assert.NoError(t, signalPair(pcOffer, pcAnswer))
	<-iceComplete

	if atomic.LoadInt32(&senderCalledCandidateChange) == 0 {
		t.Fatalf("Sender ICETransport OnSelectedCandidateChange was never called")
	}

	closePairNow(t, pcOffer, pcAnswer)
}
