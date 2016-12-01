package logging

import (
	"sort"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

// CreateHook creates a hook based on a viper config
type CreateHook func(*viper.Viper) logrus.Hook

var (
	knownHooks map[string]CreateHook
	hooksLock  *sync.Mutex
)

func init() { // using init avoids a race
	hooksLock = new(sync.Mutex)
	knownHooks = make(map[string]CreateHook, 50)
}

// RegisterHook for use through configuration system
func RegisterHook(name string, factory CreateHook) {
	hooksLock.Lock()
	defer hooksLock.Unlock()
	knownHooks[strings.ToLower(name)] = factory

}

// KnownHooks returns the list of keys for the registered hooks
func KnownHooks() []string {
	var hooks []string
	for k := range knownHooks {
		hooks = append(hooks, k)
	}
	sort.Strings(hooks)
	return hooks
}

func parseHooks(v *viper.Viper) []logrus.Hook {
	if v.IsSet("hooks") {
		hs := v.Get("hooks")
		var res []logrus.Hook
		switch ele := hs.(type) {
		case []interface{}:
			for _, v := range ele {
				mm, err := cast.ToStringMapE(v)
				if err != nil {
					continue
				}
				h := parseHook(mm)
				if h == nil {
					continue
				}
				res = append(res, h)
			}
			return res
		case map[interface{}]interface{}:
			h := parseHook(v.GetStringMap("hooks"))
			if h == nil {
				return nil
			}
			res = append(res, h)
			return res
		}
	}
	return nil
}

func parseHook(v map[string]interface{}) logrus.Hook {
	if nme, ok := v["name"]; ok {
		name, err := cast.ToStringE(nme)
		if err != nil {
			return nil
		}

		hooksLock.Lock()
		defer hooksLock.Unlock()
		if create, ok := knownHooks[strings.ToLower(name)]; ok {
			vv := viper.New()
			vv.Set("nested", v)
			h := create(vv.Sub("nested"))

			return h
		}
	}
	return nil
}

// hooks are "special" they get merged for real
// so if you define hooks then the hooks from the parent logger trickle down
func mergeHooks(child, parent *viper.Viper) {
	var result []interface{}
	known := make(map[string]int, 20)

	if parent.IsSet("hooks") {
		switch hc := parent.Get("hooks").(type) {
		case []interface{}:
			for _, v := range hc {
				if mp, ok := v.(map[interface{}]interface{}); ok {
					if nmi, ok := mp["name"]; ok {
						if nm, ok := nmi.(string); ok {
							known[nm] = len(result)
							result = append(result, v)
						}
					}
				}
			}
		case map[interface{}]interface{}:
			if nmi, ok := hc["name"]; ok {
				if nm, ok := nmi.(string); ok {
					known[nm] = len(result)
					result = append(result, hc)
				}
			}
		}
	}

	if child.IsSet("hooks") {
		switch hc := child.Get("hooks").(type) {
		case []interface{}:
			if len(result) == 0 {
				for _, v := range hc {
					if mp, ok := v.(map[interface{}]interface{}); ok {
						if nmi, ok := mp["name"]; ok {
							if nm, ok := nmi.(string); ok {
								known[nm] = len(result)
								result = append(result, v)
							}
						}
					}
				}
			} else {
				for _, v := range hc {
					if mp, ok := v.(map[interface{}]interface{}); ok {
						if nmi, ok := mp["name"]; ok {
							if nm, ok := nmi.(string); ok {
								idx, k := known[nm]
								if k {
									result[idx] = v
								} else {
									known[nm] = len(result)
									result = append(result, v)
								}
							}
						}
					}
				}
			}
		case map[interface{}]interface{}:
			if nmi, ok := hc["name"]; ok {
				if nm, ok := nmi.(string); ok {
					idx, k := known[nm]
					if k {
						result[idx] = hc
					} else {
						known[nm] = len(result)
						result = append(result, hc)
					}
				}
			}
		}
	}

	child.Set("hooks", result)
}
