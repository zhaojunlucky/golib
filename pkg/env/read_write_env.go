package env

import (
	"os"
)

type ReadWriteEnv struct {
	Parent Env
	envs   map[string]string
}

func NewReadWriteEnv(parent Env, envs map[string]string) *ReadWriteEnv {
	if parent == nil {
		parent = OSEnv
	}
	env := &ReadWriteEnv{
		Parent: parent,
		envs:   make(map[string]string, max(16, len(envs))),
	}
	if len(envs) > 0 {
		env.SetAll(envs)
	}
	return env
}

func (env *ReadWriteEnv) Get(key string) string {
	if val, ok := env.envs[key]; ok {
		return val
	}
	return env.Parent.Get(key)
}

func (env *ReadWriteEnv) Set(key, value string) {
	env.envs[key] = env.Expand(value)
}

func (env *ReadWriteEnv) SetAll(envs map[string]string) {
	for key, value := range envs {
		env.envs[key] = env.Expand(value)
	}
}

func (env *ReadWriteEnv) GetAll() map[string]string {
	envs := make(map[string]string)

	for key, value := range env.envs {
		envs[key] = value
	}
	if env.Parent != nil {
		for key, value := range env.Parent.GetAll() {
			if _, ok := envs[key]; ok {
				continue
			}
			envs[key] = value
		}
	}

	return envs

}

func (env *ReadWriteEnv) Expand(s string) string {
	return os.Expand(s, func(s string) string {
		return env.Get(s)
	})
}
