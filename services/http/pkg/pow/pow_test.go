package pow

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func solve(challenge string) string {
	for n := 0; ; n++ {
		nonce := strconv.Itoa(n)
		if Matches(challenge+nonce, 1) {
			return nonce
		}
	}
}

func TestProofCannotBeForgedReboundOrReplayed(t *testing.T) {
	guard, err := New()
	if err != nil {
		t.Fatal(err)
	}
	challenge, err := guard.Challenge("POST|/signin|192.0.2.1", 1)
	if err != nil {
		t.Fatal(err)
	}
	nonce := solve(challenge)
	if guard.Verify(challenge, nonce, "POST|/signin|192.0.2.2", 1) {
		t.Fatal("accepted another IP")
	}
	if guard.Verify(challenge, nonce, "POST|/reset|192.0.2.1", 1) {
		t.Fatal("accepted another route")
	}
	if guard.Verify(challenge, nonce, "POST|/signin|192.0.2.1", 2) {
		t.Fatal("accepted another difficulty")
	}
	if guard.Verify("chosen", solve("chosen"), "POST|/signin|192.0.2.1", 1) {
		t.Fatal("accepted client challenge")
	}
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if guard.Verify(challenge, nonce, "POST|/signin|192.0.2.1", 1) {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()
	if accepted.Load() != 1 {
		t.Fatalf("proof accepted %d times", accepted.Load())
	}
}

func TestExpiredChallengeAndInvalidDifficulty(t *testing.T) {
	guard, _ := New()
	payload := strconv.FormatInt(time.Now().Add(-time.Second).Unix(), 10) + ":random"
	challenge := payload + ":" + guard.sign(payload+":1:binding")
	if guard.Verify(challenge, solve(challenge), "binding", 1) {
		t.Fatal("expired proof accepted")
	}
	for _, difficulty := range []int{-1, 0, MaxDifficulty + 1, 100} {
		if Matches("test", difficulty) {
			t.Fatal("invalid difficulty accepted")
		}
	}
}
