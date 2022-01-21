package env

import (
	"log"
	"os"
	"strconv"
)

func RequiredString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("you must define %s environment variable", key)
	}
	return val
}

func OptionalString(key, defult string) string {
	val := os.Getenv(key)
	if val == "" {
		return defult
	} else {
		return val
	}
}

func RequiredInt(key string) int {
	val := RequiredString(key)
	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("failed to parse %s as int; %s", val, err)
	}
	return n
}

func OptionalInt(key string, defult int) int {
	val := os.Getenv(key)
	n, err := strconv.Atoi(val)
	if err != nil {
		return defult
	}
	return n
}
