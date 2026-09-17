package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"neomatica/neosync/types"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func deviceStreamKey(imei string) string {
	return fmt.Sprintf("adm:device:%s:events", imei)
}

func ReadAdmLog(imei string) ([]types.AdmLogItem, error) {
	key := deviceStreamKey(imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msgs, err := Client.XRevRange(ctx, key, "+", "-").Result()
	if err != nil {
		if err == redis.Nil {
			return []types.AdmLogItem{}, nil
		}

		return nil, err
	}

	items := make([]types.AdmLogItem, 0, len(msgs))

	for _, msg := range msgs {
		var payload any
		if p, ok := msg.Values["payload"].(string); ok && p != "" {
			if json.Unmarshal([]byte(p), &payload) != nil {
				payload = p
			}
		}

		ts, _ := strconv.ParseInt(fmt.Sprint(msg.Values["ts"]), 10, 64)

		items = append(items, types.AdmLogItem{
			ID:      msg.ID,
			TS:      ts,
			Event:   fmt.Sprint(msg.Values["event"]),
			Payload: payload,
		})
	}

	return items, nil
}
