package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"neomatica/neosync/infra/constants"
	"time"

	"github.com/redis/go-redis/v9"
)

type ConfigurationDraftInsert struct {
	DeviceImei string         `json:"device_imei"`
	CfgHash    uint32         `json:"cfg_hash"`
	Section    string         `json:"section"`
	Changes    map[string]any `json:"changes"`
	Timestamp  int64          `json:"timestamp"`
}

func draftKey(uuid, imei, section string) string {
	return fmt.Sprintf("draft:%s:%s:%s", uuid, imei, section)
}

func SaveOrUpdateDraft(uuid, imei, section string, newDraft ConfigurationDraftInsert) error {
	key := draftKey(uuid, imei, section)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := Client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	if err == redis.Nil {
		newDraft.Timestamp = time.Now().Unix()

		data, _ := json.Marshal(newDraft)

		return Client.Set(ctx, key, data, constants.Redis_DraftTTL).Err()
	}

	var existing ConfigurationDraftInsert
	if err := json.Unmarshal([]byte(val), &existing); err != nil {
		return err
	}

	for k, v := range newDraft.Changes {
		existing.Changes[k] = v
	}

	if newDraft.CfgHash != 0 && newDraft.CfgHash != existing.CfgHash {
		existing.CfgHash = newDraft.CfgHash
	}

	existing.Timestamp = time.Now().Unix()
	data, _ := json.Marshal(existing)

	if err := Client.Set(ctx, key, data, constants.Redis_DraftTTL); err != nil {
		return err.Err()
	}

	return nil
}

func GetDraft(uuid, imei, section string) (*ConfigurationDraftInsert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := draftKey(uuid, imei, section)

	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var draft ConfigurationDraftInsert
	if err := json.Unmarshal([]byte(val), &draft); err != nil {
		return nil, err
	}

	return &draft, nil
}

func HasDraft(uuid, imei string) (bool, error) {
	ctx := context.Background()
	pattern := fmt.Sprintf("draft:%s:%s:*", uuid, imei)
	keys, err := Client.Keys(ctx, pattern).Result()
	return len(keys) > 0, err
}

func GetMergedDraft(uuid, imei string) (*ConfigurationDraftInsert, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pattern := fmt.Sprintf("draft:%s:%s:*", uuid, imei)

	var cursor uint64
	merged := ConfigurationDraftInsert{
		DeviceImei: imei,
		CfgHash:    0,
		Changes:    make(map[string]any),
		Timestamp:  0,
	}

	for {
		keys, nextCursor, err := Client.Scan(ctx, cursor, pattern, 58).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := Client.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}

				return nil, err
			}

			var draft ConfigurationDraftInsert
			if err := json.Unmarshal([]byte(val), &draft); err != nil {
				return nil, err
			}

			for k, v := range draft.Changes {
				merged.Changes[k] = v
			}

			merged.CfgHash = draft.CfgHash

			if draft.Timestamp > merged.Timestamp {
				merged.Timestamp = draft.Timestamp
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(merged.Changes) == 0 {
		return nil, nil
	}

	return &merged, nil
}

func DeleteAllDrafts(uuid, imei string) error {
	pattern := fmt.Sprintf("draft:%s:%s:*", uuid, imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cursor uint64
	deletedCount := 0
	for {
		keys, nextCursor, err := Client.Scan(ctx, cursor, pattern, 58).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := Client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			deletedCount += len(keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	var verifyCursor uint64
	for {
		verifyKeys, verifyNextCursor, err := Client.Scan(ctx, verifyCursor, pattern, 58).Result()
		if err != nil {
			break
		}

		if len(verifyKeys) > 0 {
			_ = Client.Del(ctx, verifyKeys...).Err()
		}

		verifyCursor = verifyNextCursor
		if verifyCursor == 0 {
			break
		}
	}

	return nil
}

func DeleteDraft(uuid, imei, section string) error {
	key := draftKey(uuid, imei, section)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return Client.Del(ctx, key).Err()
}

// UpdateAllDraftsHash обновляет CfgHash во всех draft для указанного пользователя и устройства
func UpdateAllDraftsHash(uuid, imei string, newHash uint32) error {
	pattern := fmt.Sprintf("draft:%s:%s:*", uuid, imei)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cursor uint64
	for {
		keys, nextCursor, err := Client.Scan(ctx, cursor, pattern, 58).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			val, err := Client.Get(ctx, key).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}

				return err
			}

			var draft ConfigurationDraftInsert
			if err := json.Unmarshal([]byte(val), &draft); err != nil {
				continue
			}

			if draft.CfgHash != newHash {
				draft.CfgHash = newHash
				draft.Timestamp = time.Now().Unix()

				data, _ := json.Marshal(draft)

				if err := Client.Set(ctx, key, data, constants.Redis_DraftTTL).Err(); err != nil {
					return err
				}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
