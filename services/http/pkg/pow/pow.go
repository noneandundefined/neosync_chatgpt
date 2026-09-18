package pow

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"sync"
	"time"
)

const TTL = 2 * time.Minute
const MaxDifficulty = 6
const maxUsed = 100000

// Each process has its own key, so a proof accepted here cannot be replayed on another instance.
type Guard struct {
	key         [32]byte
	mu          sync.Mutex
	used        map[string]int64
	lastCleanup int64
}

func New() (*Guard, error) {
	guard := &Guard{used: make(map[string]int64)}
	_, err := rand.Read(guard.key[:])
	return guard, err
}

func (g *Guard) sign(value string) string {
	mac := hmac.New(sha256.New, g.key[:])
	mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (g *Guard) Challenge(binding string, difficulty int) (string, error) {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	payload := strconv.FormatInt(time.Now().Add(TTL).Unix(), 10) + ":" + base64.RawURLEncoding.EncodeToString(random)
	signature := g.sign(payload + ":" + strconv.Itoa(difficulty) + ":" + binding)
	return payload + ":" + signature, nil
}

func Matches(input string, difficulty int) bool {
	if difficulty < 1 || difficulty > MaxDifficulty {
		return false
	}
	hash := sha256.Sum256([]byte(input))
	for i := 0; i < difficulty; i++ {
		if i%2 == 0 && hash[i/2]>>4 != 0 {
			return false
		}
		if i%2 == 1 && hash[i/2]&15 != 0 {
			return false
		}
	}
	return true
}

func (g *Guard) Verify(challenge, nonce, binding string, difficulty int) bool {
	if len(challenge) > 256 || len(nonce) == 0 || len(nonce) > 128 {
		return false
	}
	parts := strings.Split(challenge, ":")
	if len(parts) != 3 {
		return false
	}
	expires, err := strconv.ParseInt(parts[0], 10, 64)
	now := time.Now().Unix()
	if err != nil || expires <= now || expires > now+int64(TTL/time.Second) {
		return false
	}
	expected := g.sign(parts[0] + ":" + parts[1] + ":" + strconv.Itoa(difficulty) + ":" + binding)
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) || !Matches(challenge+nonce, difficulty) {
		return false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if now != g.lastCleanup {
		for key, expiry := range g.used {
			if expiry <= now {
				delete(g.used, key)
			}
		}
		g.lastCleanup = now
	}
	if _, exists := g.used[challenge]; exists {
		return false
	}
	if len(g.used) >= maxUsed {
		return false
	}
	g.used[challenge] = expires
	return true
}
