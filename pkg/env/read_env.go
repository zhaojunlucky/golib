package env

import (
	"os"
	"strings"
)

var OSEnv = NewOSEnv()

type ReadEnv struct {
	Parent *Env
	envs   map[string]string
}

func NewOSEnv() Env {
	env := ReadEnv{
		Parent: nil,
		envs:   make(map[string]string),
	}
	env.initOSEnv()
	return &env
}

func (env *ReadEnv) initOSEnv() {
	for _, envStr := range os.Environ() {
		sepIndex := strings.Index(envStr, "=")
		if sepIndex < 0 {
			continue
		}
		env.envs[envStr[:sepIndex]] = envStr[sepIndex+1:]
	}
}

func (env *ReadEnv) Get(key string) string {
	if val, ok := env.envs[key]; ok {
		return val
	}
	return os.Getenv(key)
}

func (env *ReadEnv) GetAll() map[string]string {
	envs := make(map[string]string)

	for key, value := range env.envs {
		envs[key] = value
	}
	for _, envStr := range os.Environ() {
		key, value, found := strings.Cut(envStr, "=")
		if !found {
			continue
		}

		if _, ok := env.envs[key]; !ok {
			envs[key] = value
		}
	}

	return envs

}

func (env *ReadEnv) Set(key, value string) {

}

func (env *ReadEnv) SetAll(envs map[string]string) {

}

func (env *ReadEnv) Expand(s string) string {
	return os.Expand(s, func(s string) string {
		return env.Get(s)
	})
}
