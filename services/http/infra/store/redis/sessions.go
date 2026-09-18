package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	redisdb "github.com/redis/go-redis/v9"
)

type AuthSession struct {
	ID          string
	UserUUID    string
	RefreshHash string
	IPHash      string
	ExpiresAt   time.Time
}

func sessionKey(id string) string {
	return "auth:session:" + id
}

func userSessionsKey(user string) string {
	return "auth:user_sessions:" + user
}

var createSessionScript = redisdb.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 1 then return 0 end
local now = redis.call('TIME')
local nowms = tonumber(now[1]) * 1000 + math.floor(tonumber(now[2]) / 1000)
local expires = tonumber(ARGV[5])
if expires <= nowms then return 0 end
redis.call('HSET', KEYS[1], 'id', ARGV[1], 'user_uuid', ARGV[2],
    'refresh_hash', ARGV[3], 'ip_hash', ARGV[4], 'expires_at', ARGV[5])
redis.call('PEXPIREAT', KEYS[1], ARGV[5])
redis.call('ZREMRANGEBYSCORE', KEYS[2], '-inf', nowms)
redis.call('ZADD', KEYS[2], expires, ARGV[1])
local latest = redis.call('ZREVRANGE', KEYS[2], 0, 0, 'WITHSCORES')
redis.call('PEXPIREAT', KEYS[2], latest[2])
return 1
`)

func CreateAuthSession(ctx context.Context, session *AuthSession) error {
	if session == nil || session.ID == "" || session.UserUUID == "" || session.IPHash == "" {
		return errors.New("invalid Redis session")
	}

	keys := []string{
		sessionKey(session.ID),
		userSessionsKey(session.UserUUID),
	}

	created, err := createSessionScript.Run(ctx, Client, keys,
		session.ID, session.UserUUID, session.RefreshHash, session.IPHash, session.ExpiresAt.UnixMilli()).Int()
	if err != nil {
		return err
	}

	if created != 1 {
		return errors.New("session already exists or expired")
	}

	return nil
}

func GetAuthSession(ctx context.Context, id string) (*AuthSession, error) {
	values, err := Client.HGetAll(ctx, sessionKey(id)).Result()
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, nil
	}

	expires, err := strconv.ParseInt(values["expires_at"], 10, 64)
	if err != nil || values["id"] != id || values["user_uuid"] == "" || values["ip_hash"] == "" {
		return nil, errors.New("invalid Redis session")
	}

	// Sessions without a user index cannot be revoked together on password reset.
	score, err := Client.ZScore(ctx, userSessionsKey(values["user_uuid"]), id).Result()
	if errors.Is(err, redisdb.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if score != float64(expires) || time.Now().UnixMilli() >= expires {
		return nil, nil
	}

	return &AuthSession{
		ID:          id,
		UserUUID:    values["user_uuid"],
		RefreshHash: values["refresh_hash"],
		IPHash:      values["ip_hash"],
		ExpiresAt:   time.UnixMilli(expires),
	}, nil
}

var rotateRefreshScript = redisdb.NewScript(`
if ARGV[2] == '' or ARGV[3] == '' then return 0 end
if redis.call('HGET', KEYS[1], 'user_uuid') ~= ARGV[1] then return 0 end
if redis.call('HGET', KEYS[1], 'refresh_hash') ~= ARGV[2] then return 0 end
local expires = tonumber(redis.call('HGET', KEYS[1], 'expires_at'))
local score = tonumber(redis.call('ZSCORE', KEYS[2], ARGV[4]))
local now = redis.call('TIME')
local nowms = tonumber(now[1]) * 1000 + math.floor(tonumber(now[2]) / 1000)
if not expires or score ~= expires or expires <= nowms then return 0 end
redis.call('HSET', KEYS[1], 'refresh_hash', ARGV[3])
return 1
`)

func RotateAuthSessionRefresh(ctx context.Context, id, userUUID, oldHash, newHash string) (bool, error) {
	keys := []string{
		sessionKey(id),
		userSessionsKey(userUUID),
	}

	rotated, err := rotateRefreshScript.Run(ctx, Client, keys, userUUID, oldHash, newHash, id).Int()

	return rotated == 1, err
}

var revokeSessionScript = redisdb.NewScript(`
if redis.call('HGET', KEYS[1], 'user_uuid') == ARGV[1] then
    redis.call('DEL', KEYS[1])
    redis.call('ZREM', KEYS[2], ARGV[2])
end
return 1
`)

func DeleteAuthSession(ctx context.Context, id, userUUID string) error {
	keys := []string{
		sessionKey(id),
		userSessionsKey(userUUID),
	}

	return revokeSessionScript.Run(ctx, Client, keys, userUUID, id).Err()
}

var revokeUserSessionsScript = redisdb.NewScript(`
local ids = redis.call('ZRANGE', KEYS[1], 0, -1)
for _, id in ipairs(ids) do redis.call('DEL', ARGV[1] .. id) end
redis.call('DEL', KEYS[1])
return 1
`)

func DeleteUserAuthSessions(ctx context.Context, userUUID string) error {
	keys := []string{
		userSessionsKey(userUUID),
	}

	return revokeUserSessionsScript.Run(ctx, Client, keys, sessionKey("")).Err()
}
