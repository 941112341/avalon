package inline

import (
	"math/rand"
	"os"
)

func GetEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return val
}

func GetEnvs(key string, defaultValues ...string) string {
	defaultValues = append(defaultValues, "")
	d := defaultValues[rand.Intn(len(defaultValues))]
	return GetEnv(key, d)
}
