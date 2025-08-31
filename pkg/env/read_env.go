package env

import (
	"os"
	"strings"
)

var OSEnv = NewOSEnv()

type ReadEnv struct {
	Parent Env
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

func NewEmptyReadEnv() Env {
	return &ReadEnv{
		Parent: nil,
		envs:   make(map[string]string),
	}
}

func NewReadEnv(parent Env, envs map[string]string) *ReadEnv {
	if parent == nil {
		parent = OSEnv
	}
	env := &ReadEnv{
		Parent: parent,
		envs:   make(map[string]string, max(16, len(envs))),
	}

	for k, v := range envs {
		env.envs[k] = env.Expand(v)
	}

	return env
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
	if env.Parent != nil {
		return env.Parent.Get(key)
	}
	return ""
}

func (env *ReadEnv) Contains(key string) bool {
	if _, ok := env.envs[key]; ok {
		return true
	}
	if env.Parent != nil {
		return env.Parent.Contains(key)
	}
	return false
}

func (env *ReadEnv) GetAll() map[string]string {
	newEnv := make(map[string]string)

	if env.Parent != nil {
		for key, value := range env.Parent.GetAll() {
			newEnv[key] = value
		}
	}

	for key, value := range env.envs {
		newEnv[key] = value
	}
	return newEnv
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
