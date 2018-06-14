package settings

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

var settings map[string]*viper.Viper

func Initialize() {
	// initalize cmdb config
	cmdb := viper.New()
	cmdb.SetConfigFile("cmdb")
	cmdb.AddConfigPath("configs/")   // 可能的配置文件路径
	cmdb.AddConfigPath("/etc/cmdb/") // 可能的配置文件路径
	cmdb.AddConfigPath("$HOME/.cmdb/")
	cmdb.AddConfigPath(".")
	err := cmdb.ReadInConfig()
	if err != nil {
		logrus.Errorf("Config file load failed: %s\n", err)
	}

	logrus.Infof("Application Setting: %#v args:%#v", viper.AllSettings(), os.Args)
	logrus.Infof("cmdb.json Settings: %#v args:%#v", cmdb.AllSettings(), os.Args)
}

func Get(name string) *viper.Viper {

	if settings == nil {
		settings = make(map[string]*viper.Viper)
	}

	if settings[name] != nil {
		return settings[name]
	}

	// initialize setting config
	setting := viper.New()
	setting.SetConfigName(name + ".default")     // name of config file (without extension)
	setting.AddConfigPath("configs/")            // path to look for the config file in
	setting.AddConfigPath("./../configs/")       // path to look for the config file in
	setting.AddConfigPath("./../../configs/")    // path to look for the config file in
	setting.AddConfigPath("./../../../configs/") // path to look for the config file in
	setting.AddConfigPath("/etc/cmdb/")          // path to look for the config file in
	setting.AddConfigPath("$HOME/.cmdb")         // call multiple times to add many search paths
	setting.AddConfigPath(".")                   // optionally look for config in the working directory

	errDefault := setting.ReadInConfig() // Find and read the config file

	setting.SetConfigName(name)           // name of config file (without extension)
	errCurrent := setting.MergeInConfig() // Find and read the config file

	if errDefault != nil && errCurrent != nil {

		logrus.WithFields(logrus.Fields{
			"name":       name,
			"errDefault": errDefault,
			"errCurrent": errCurrent,
		}).Warn("Config file load failed.")
	}

	settings[name] = setting

	return setting
}

//用于获取不同的配置文件，调用settings.Get("cmdb") 即：cmdb.yaml
func Set(name string, key string, value interface{}) {

	if settings == nil {
		settings = make(map[string]*viper.Viper)
	}

	if settings[name] == nil {
		settings[name] = viper.New()
	}

	var setting = settings[name]

	setting.Set(key, value)
}

func Flush() {
	settings = nil
}

// GetString returns the value associated with the key as a string.
func GetString(name string, key string) string {
	v := Get(name)
	return v.GetString(key)
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(name string, key string) bool {
	v := Get(name)
	return v.GetBool(key)
}

// GetInt returns the value associated with the key as an integer.
func GetInt(name string, key string) int {
	v := Get(name)
	return v.GetInt(key)
}

// GetInt64 returns the value associated with the key as an integer.
func GetInt64(name string, key string) int64 {
	v := Get(name)
	return v.GetInt64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func GetFloat64(name string, key string) float64 {
	v := Get(name)
	return v.GetFloat64(key)
}

// GetTime returns the value associated with the key as time.
func GetTime(name string, key string) time.Time {
	v := Get(name)
	return v.GetTime(key)
}

// GetDuration returns the value associated with the key as a duration.
func GetDuration(name string, key string) time.Duration {
	v := Get(name)
	return v.GetDuration(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func GetStringSlice(name string, key string) []string {
	v := Get(name)
	return v.GetStringSlice(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func GetStringMap(name string, key string) map[string]interface{} {
	v := Get(name)
	return v.GetStringMap(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func GetStringMapString(name string, key string) map[string]string {
	v := Get(name)
	return v.GetStringMapString(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func GetStringMapStringSlice(name string, key string) map[string][]string {
	v := Get(name)
	return v.GetStringMapStringSlice(key)
}

// GetSizeInBytes returns the size of the value associated with the given key
// in bytes.
func GetSizeInBytes(name string, key string) uint {
	v := Get(name)
	return v.GetSizeInBytes(key)
}
