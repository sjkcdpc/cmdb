package log

import (
	"github.com/sirupsen/logrus"
	"github.com/mds1455975151/cmdb/settings"
	"github.com/heirko/go-contrib/logrusHelper"

)

func InitLogger() {

	var logSetting = settings.Get("logger")

	// Read configuration
	var c = logrusHelper.UnmarshalConfiguration(logSetting) // Unmarshal configuration from Viper
	logrusHelper.SetConfig(logrus.StandardLogger(), c)      // for e.g. apply it to logrus default instance

	// ### End Read Configuration

	//// ### Use logrus as normal
	//logrus.WithFields(logrus.Fields{
	//	"animal": "walrus",
	//}).Error("A walrus appears")
	//
	//Infof(map[string]interface{}{
	//	"aaa": 1,
	//}, "aaa", nil)
}

