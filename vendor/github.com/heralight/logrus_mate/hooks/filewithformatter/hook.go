package logrus_file

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/heralight/logrus_mate"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus_mate.RegisterHook("filewithformatter", NewFileHook)
}

func NewFileHook(options logrus_mate.Options) (hook logrus.Hook, err error) {

	conf := FileLogConifg{}

	if err = options.ToObject(&conf); err != nil {
		return
	}

	path := strings.Split(conf.Filename, "/")
	if len(path) > 1 {
		exec.Command("mkdir", path[0]).Run()
	}

	w := NewFileWriter()

	if err = w.Init(conf); err != nil {
		return
	}

	//w.SetPrefix("[-] ")

	hook = &FileHook{W: w}

	return
}

type FileHook struct {
	W *FileLogWriter
}

func (p *FileHook) Fire(entry *logrus.Entry) (err error) {
	if p.W.Level < entry.Level {
		return nil
	}
	message, err := entry.String()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	return p.W.WriteMsg(message)

}

func (p *FileHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
