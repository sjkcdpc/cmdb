package tlog

import (
	"net"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    *Tlogger
		wantErr bool
	}{
		{"", args{"10.1.16.201:6667"}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Dial(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want {
				t.Errorf("Dial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTlogger_Log(t *testing.T) {
	type fields struct {
		conn net.Conn
	}
	type args struct {
		args []interface{}
	}

	var names []interface{}
	names = append(names, "first")
	names = append(names, "second")
	names = append(names, "")
	names = append(names, 234)
	names = append(names, 23465.324)

	//strs := []string{"first", "second"}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{args: args{names}},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.Log(PublicConfig{"PlayerLogin", "2.2.5.1", time.Now().Format("2006-01-02 15:04:05"), "100695782", 1, 20001}, tt.args.args...)
		})
	}

	r.Close()
}

func TestLog(t *testing.T) {
	type args struct {
		PlatID int
		args   []interface{}
	}

	var names []interface{}
	names = append(names, "first")
	names = append(names, "second")
	names = append(names, "")
	names = append(names, 234)
	names = append(names, 23465.324)

	tests := []struct {
		name string
		args args
	}{
		{args: args{1, names}},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(PublicConfig{"PlayerLogin", "2.2.5.1", time.Now().Format("2006-01-02 15:04:05"), "100695782", tt.args.PlatID, 20001}, tt.args.args...)
		})
	}

	r.Close()
}

func TestTlogger_CreateFormat(t *testing.T) {
	type fields struct {
		conn        net.Conn
		required    []string
		format      string
		GameSvrId   string
		GameAppID   string
		IZoneAreaID int
	}

	var names []interface{}
	names = append(names, "first")
	names = append(names, "second")
	names = append(names, "")
	names = append(names, 234)
	names = append(names, 23465.324)

	type args struct {
		PlatID int
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{args: args{1, names}, want: "PlayerLogin|2.2.5.1|%v|100695782|1|20001|first|second|NULL|234|23465.324\n"},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := r.CreateFormat("PlayerLogin", tt.args.PlatID, tt.args.args...); got != tt.want {
				t.Errorf("Tlogger.CreateFormat() = %v, want %v", got, tt.want)
			}
		})
	}

	r.Close()
}

func TestTlogger_LogRaw(t *testing.T) {
	type fields struct {
		conn        net.Conn
		required    []string
		format      string
		GameSvrId   string
		GameAppID   string
		IZoneAreaID int
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{args: args{"2.2.5.1|2018-02-05 17:29:35|100695782|1|20001|first|second|NULL|234|23465.324\n"}},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.LogRaw(tt.args.message)
		})
	}

	r.Close()
}

func TestTlogger_Close(t *testing.T) {
	type fields struct {
		conn        net.Conn
		required    []string
		format      string
		GameSvrId   string
		GameAppID   string
		IZoneAreaID int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Tlogger{
				conn:        tt.fields.conn,
				required:    tt.fields.required,
				format:      tt.fields.format,
				GameSvrId:   tt.fields.GameSvrId,
				GameAppID:   tt.fields.GameAppID,
				IZoneAreaID: tt.fields.IZoneAreaID,
			}
			if err := r.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Tlogger.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTlogger_SetRequired(t *testing.T) {
	type fields struct {
		conn        net.Conn
		required    []string
		format      string
		GameSvrId   string
		GameAppID   string
		IZoneAreaID int
	}
	type args struct {
		GameSvrId   string
		GameAppID   string
		iZoneAreaID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.SetRequired(tt.args.GameSvrId, tt.args.GameAppID, tt.args.iZoneAreaID); (err != nil) != tt.wantErr {
				t.Errorf("Tlogger.SetRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	r.Close()
}

func TestCreateFormat(t *testing.T) {
	type args struct {
		PlatID int
		args   []interface{}
	}

	var names []interface{}
	names = append(names, "first")
	names = append(names, "second")
	names = append(names, "%s")
	names = append(names, 234)
	names = append(names, 23465.324)

	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{1, names}, want: "PlayerLogin|2.2.5.1|%v|100695782|1|20001|first|second|%s|234|23465.324\n"},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateFormat("PlayerLogin", tt.args.PlatID, tt.args.args...); got != tt.want {
				t.Errorf("CreateFormat() = %v, want %v", got, tt.want)
			}
		})
	}

	r.Close()
}

func TestLogRaw(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{args: args{"2.2.5.1|2018-02-05 17:29:35|100695782|1|20001|first|second|NULL|234|23465.324\n"}},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogRaw(tt.args.message)
		})
	}

	r.Close()
}

func TestTlogger_LogFormat(t *testing.T) {
	type fields struct {
		conn        net.Conn
		required    []string
		format      string
		GameSvrId   string
		GameAppID   string
		IZoneAreaID int
	}
	type args struct {
		format string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{args: args{"2.2.5.1|%v|100695782|1|20001|first|second|%v|234|23465.324\n"}},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.LogFormat(tt.args.format, 444)
		})
	}

	r.Close()
}

func TestLogFormat(t *testing.T) {
	type args struct {
		format string
	}
	tests := []struct {
		name string
		args args
	}{
		{args: args{"2.2.5.1|%v|100695782|1|20001|first|second|NULL|234|23465.324\n"}},
	}

	r, _ := Dial("10.1.16.201:6667")
	r.SetRequired("2.2.5.1", "100695782", 20001)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogFormat(tt.args.format)
		})
	}

	r.Close()
}
