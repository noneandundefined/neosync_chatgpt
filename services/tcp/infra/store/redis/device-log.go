package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"neomatica/neosync-tcp/infra/constants"
	"neomatica/neosync-tcp/util"
	"time"

	"github.com/redis/go-redis/v9"
)

func deviceStreamKey(imei string) string {
	return fmt.Sprintf("adm:device:%s:events", imei)
}

func WriteAdmLog(event string, imei string, payload any) error {
	imeiPrepare := util.PrepareImei(imei)
	timeutc := time.Now().UTC().Unix()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var payloadStr string
	switch t := payload.(type) {
	case nil:
		payloadStr = ""

	case string:
		payloadStr = t

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		payloadStr = fmt.Sprint(t)

	default:
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		payloadStr = string(b)
	}

	key := deviceStreamKey(imeiPrepare)

	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		MaxLen: 200,
		Approx: true,
		Values: map[string]any{
			"event":   event,
			"ts":      timeutc,
			"payload": payloadStr,
		},
	}).Result()
	if err != nil {
		return err
	}

	return client.Expire(ctx, key, constants.Redis_LogTTL).Err()
}
