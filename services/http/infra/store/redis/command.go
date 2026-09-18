package redis

import (
	"context"
	"neomatica/neosync/infra/constants"
	"time"
)

func commandKey(imei string) string {
	return cacheKey("command:imei:" + imei)
}

func SetCommandResponse(imei, command, response string) error {
	key := commandKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.HSet(ctx, key, command, response).Err(); err != nil {
		return err
	}

	return Client.Expire(ctx, key, constants.Redis_CommandTTL).Err()
}

func GetCommandResponse(imei, command string) (string, error) {
	key := commandKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.HGet(ctx, key, command).Result()
}

func CommandExists(imei, command string) (bool, error) {
	key := commandKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.HExists(ctx, key, command).Result()
}
