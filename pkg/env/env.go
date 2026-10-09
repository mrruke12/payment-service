package env

import (
	"fmt"
	"os"
)

func KeysMustExist(keys []string) ([]string, error) {
	res := make([]string, 0, len(keys))

	for _, key := range keys {
		value, exist := os.LookupEnv(key)

		if !exist {
			return nil, fmt.Errorf("key \"%s\" is not present in the env", key)
		}

		res = append(res, value)
	}

	return res, nil
}
