package redis

import (
	"context"
	"encoding/json"
	"neomatica/neosync-tcp/infra/logger"
	"neomatica/neosync-tcp/pkg/telemetry"
	"time"

	"github.com/redis/go-redis/v9"
)

func telemetryKey(imei string) string {
	return cacheKey("telemetry:imei:" + imei)
}

func SaveTelemetry(imei string, value any) error {
	key := telemetryKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("SaveTelemetry imei={%s}: Failed save telemetry in redis: %s", imei, err.Error())
		return err
	}

	return client.Set(ctx, key, data, 0).Err()
}

func GetTelemetry(imei string) (*telemetry.Telemetry, error) {
	key := telemetryKey(imei)

	tl := &telemetry.Telemetry{
		BLESENSORINFO: make(map[string]telemetry.BLESensor),
		FUELINFO:      make(map[string]telemetry.BLESensor),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return tl, err
	}

	if len(val) == 0 {
		return tl, nil
	}

	if err := json.Unmarshal([]byte(val), tl); err != nil {
		logger.Error("GetTelemetry imei={%s}: Failed get telemetry in redis: %s", imei, err.Error())
		return tl, err
	}

	return tl, nil
}
