package util

import (
	"errors"
	"fmt"
	"neomatica/neosync/infra/locale"
	"strconv"
	"strings"
)

func ParseIdsQuery(tr locale.Translator, raw string) ([]uint64, error) {
	if raw == "" {
		return nil, errors.New(tr.TErr("ids-query-param-empty"))
	}

	parts := strings.Split(raw, ",")
	result := make([]uint64, 0, len(parts))

	for _, part := range parts {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, fmt.Errorf(tr.TErr("invalid-id"), part)
		}

		result = append(result, id)
	}

	return result, nil
}
