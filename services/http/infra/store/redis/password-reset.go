package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	redisdb "github.com/redis/go-redis/v9"
)

func SetPasswordResetToken(ctx context.Context, userUUID, token string, expiresAt time.Time) error {
	hash := sha256.Sum256([]byte(token))
	key := "auth:password_reset:" + userUUID

	return Client.SetArgs(ctx, key, hex.EncodeToString(hash[:]), redisdb.SetArgs{
		ExpireAt: expiresAt,
	}).Err()
}

var consumePasswordResetScript = redisdb.NewScript(`
if redis.call('GET', KEYS[1]) ~= ARGV[1] then return 0 end

local ids = redis.call('ZRANGE', KEYS[2], 0, -1)

for _, id in ipairs(ids) do
    redis.call('DEL', ARGV[2] .. id)
end

redis.call('DEL', KEYS[2], KEYS[1])
return 1
`)

/* Погашение ссылки и отзыв сессий выполняются одной операцией Redis */
func ConsumePasswordResetToken(ctx context.Context, userUUID, token string) (bool, error) {
	if token == "" {
		return false, nil
	}

	hash := sha256.Sum256([]byte(token))
	keys := []string{
		"auth:password_reset:" + userUUID,
		userSessionsKey(userUUID),
	}

	consumed, err := consumePasswordResetScript.Run(ctx, Client, keys, hex.EncodeToString(hash[:]), sessionKey("")).Int()

	return consumed == 1, err
}
