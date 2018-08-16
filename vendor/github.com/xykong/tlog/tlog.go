/*
2、Q：是否可以不是使用tlog api写日志？
A:可以使用UDP协议直接往tlogd指定端口发送文本日志数据包，需要在日志内容后加"\n",一条日志一个udp报
B.由于使用udp传输，为了降低丢包率，建议一个包的小于1024个字节。
C. 日志格式和命名参考经分《日志内容规范》。
*/

package tlog

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Tlogger struct {
	conn     net.Conn
	required []string
	format   string

	GameSvrId   string
	GameAppID   string
	IZoneAreaID int
}

var lastTlogger *Tlogger

func Dial(address string) (*Tlogger, error) {

	lastTlogger = &Tlogger{}

	conn, err := net.Dial("udp", address)
	if err != nil {
		logrus.Errorf("Tlogger Dial failed: %s", err.Error())
		return nil, err
	}

	lastTlogger.conn = conn

	logrus.WithFields(logrus.Fields{
		"address":    address,
		"remoteAddr": conn.RemoteAddr(),
	}).Info("Tlogger Dial succeed.")

	return lastTlogger, nil
}

func (r *Tlogger) Log(conf PublicConfig, args ...interface{}) {

	if r == nil || r.conn == nil {
		logrus.Errorf("Tlogger is not initialized.")
		return
	}

	if r.format == "" {
		logrus.Errorf("Tlogger SetRequired is required before send log.")
		return
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s|%s|%v|%s|%d|%d", conf.Table, conf.GameSvrId, conf.DtEventTime, conf.GameAppId, conf.PlatId, conf.IZoneAreaId))

	for _, param := range args {
		buffer.WriteString("|")

		val := fmt.Sprintf("%v", param)
		if val == "" {
			buffer.WriteString("NULL")
		} else {
			if len(val) > 64 {
				val = val[0:64]
			}

			buffer.WriteString(val)
		}
	}

	buffer.WriteString("\n")

	_, err := r.conn.Write(buffer.Bytes())
	if err != nil {
		logrus.Errorf("Tlogger Log send failed: %s", err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"table":  conf.Table,
		"PlatID": conf.PlatId,
		"args":   args,
		"tlog":   buffer.String(),
	}).Info("Tlog Log")
}

func (r *Tlogger) CreateFormat(table string, PlatID int, args ...interface{}) string {

	if r.format == "" {
		logrus.Errorf("Tlogger SetRequired is required before send log.")
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf(r.format, table, "%v", PlatID))

	for _, param := range args {
		buffer.WriteString("|")

		val := fmt.Sprintf("%v", param)
		if val == "" {
			buffer.WriteString("NULL")
		} else {
			if len(val) > 64 {
				val = val[0:64]
			}

			buffer.WriteString(val)
		}
	}

	buffer.WriteString("\n")

	return buffer.String()
}

func (r *Tlogger) LogFormat(format string, args ...interface{}) {

	params := append([]interface{}{time.Now().Format("2006-01-02 15:04:05")}, args...)

	message := fmt.Sprintf(format, params...)

	_, err := r.conn.Write([]byte(message))
	if err != nil {
		logrus.Errorf("Tlogger Log send failed: %s", err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"format": format,
		"args":   args,
		"tlog":   message,
	}).Info("Tlog LogFormat")
}

func (r *Tlogger) LogRaw(message string) {

	if r == nil || r.conn == nil {
		logrus.Errorf("Tlogger is not initialized.")
		return
	}

	_, err := r.conn.Write([]byte(message))
	if err != nil {
		logrus.Errorf("Tlogger Log send failed: %s", err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"tlog": message,
	}).Info("Tlog LogRaw")
}

func (r *Tlogger) Close() error {

	if r == nil || r.conn == nil {
		return nil
	}

	return r.conn.Close()
}

/*
3	TLOG-表公共字段说明
3.1	公共字段
TLOG每张表有以下6个必填字段（包括特性需求的表结构）：
<entry name="GameSvrId" type="string" size="25" desc="登录的游戏服务器编" />
      <entry name="dtEventTime" type="datetime" desc="游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS" />
       <entry name="GameAppID" type="string" size="32" desc="游戏APPID" />
       <entry name="PlatID" type="int" defaultvalue="0" desc="ios 0/android 1" />
       <entry name="iZoneAreaID" type="int" defaultvalue="0" desc="(必填)针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0" />
       <entry name="OpenID" type="string" size="64" desc="用户OPENID号" />

示例内容：
2.2.5.1|2014-01-16 10:00:48|100695782|0|20001|0FEB999268EF9FEA26D4CB219C37910D
2.1.5.1|2014-01-16 10:04:30|100695782|1|10001|1ABC1C01947A64079B544426652F4C14
2.2.5.1|2014-01-16 10:03:45|100695782|0|20001|56A3FDD286A3A327D4C3267835A55557
2.1.5.1|2014-01-16 10:04:41|100695782|1|10001|1E2EC6FE7C3CA63EBD760B02FF2B6DC6
2.2.5.1|2014-01-16 10:07:27|100695782|0|20001|56A3FDD286A3A327D4C3267835A55557
1.1.5.1|2014-01-16 10:09:21|wxce852959d18822db|0|40001|o2ebPjtSsPY3DYZOO6b0Mx_KVlv8

以下对每个字段进行说明：
GameSvrId     服务器编号，游戏侧自行开发实现。用于定位数据来源于那台服务器。
dtEventTime    格式 YYYY-MM-DD HH:MM:SS
GameAppID    游戏APPID,通过MSDK获取
PlatID          IOS填0，安卓填1，不能写其它值。
iZoneAreaID    针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0
OpenID        玩家唯一标识，通过MSDK获取。

*/
func (r *Tlogger) SetRequired(GameSvrId, GameAppID string, iZoneAreaID int) error {

	r.GameSvrId = GameSvrId
	r.GameAppID = GameAppID
	r.IZoneAreaID = iZoneAreaID

	r.format = fmt.Sprintf("%%s|%s|%%v|%s|%%d|%d", GameSvrId, GameAppID, iZoneAreaID)

	return nil
}

func Log(conf PublicConfig, args ...interface{}) {

	lastTlogger.Log(conf, args...)
}

func CreateFormat(table string, PlatID int, args ...interface{}) string {

	return lastTlogger.CreateFormat(table, PlatID, args...)
}

func LogRaw(message string) {

	lastTlogger.LogRaw(message)
}

func LogFormat(format string, args ...interface{}) {

	lastTlogger.LogFormat(format, args...)
}

type PublicConfig struct {
	Table       string
	GameSvrId   string
	DtEventTime string
	GameAppId   string
	PlatId      int
	IZoneAreaId int
}
