/**
 * Read the configuration file
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package config

import (
    "strings"
    "sync"

    "github.com/astaxie/beego/logs"
    "github.com/owner888/kaligo/util"
)

var (
    configMaps sync.Map
)

type configMap interface {
    Load(key any) (value any, ok bool)
}

// StrMap is use for string -> map
type StrMap map[string]interface{}

func (m StrMap) Load(key any) (value any, ok bool) {
    value = m[key.(string)]
    ok = value != nil
    return
}

// Env 读取环境变量(configMaps存入的值)，支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
    if len(defaultValue) > 0 {
        return Get(envName, defaultValue[0])
    }
    return Get(envName)
}

// Add 新增配置项
func Add(key string, value map[string]interface{}) {
    logs.Debug("Add", key)
    logs.Debug("Add", value)
    configMaps.Store(key, value)
}

// Get 获取配置项，允许使用点式获取，如：core.name
func Get(key string, defaultValue ...interface{}) interface{} {
    var keys []string = strings.Split(key, ".")
    lastIndex := len(keys) - 1
    var maps configMap = &configMaps
    for i, k := range keys {
        if val, ok := maps.Load(k); ok {
            if i == lastIndex {
                return val
            } else if m, ok := val.(map[string]interface{}); ok {
                maps = StrMap(m)
            }
        }
    }
    if len(defaultValue) > 0 {
        return defaultValue[0]
    }
    return nil
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
    return util.ToString(Get(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
    return util.ToInt(Get(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
    return util.ToInt64(Get(path, defaultValue...))
}

// GetInt32 获取 Int64 类型的配置信息
func GetInt32(path string, defaultValue ...interface{}) int32 {
    return util.ToInt32(Get(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
    return util.ToUint(Get(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
    return util.ToBool(Get(path, defaultValue...))
}
