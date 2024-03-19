package kutimeconv

import (
	"testing"
	"time"
)

func Test_parseUptime(t *testing.T) {
	buf := []byte("1696518.17 1696518.17\n")

	uptime, err := parseUptime(buf)
	if err != nil {
		t.Errorf("parseUptime() error = %v", err)
		return
	}

	if uptime != 1696518170000000 {
		t.Errorf("parseUptime() uptime = %v, want 1696518170000000", uptime)
	}
}

func Test_uptimeDiff(t *testing.T) {
	tests := []struct {
		name    string
		uptime1 uint64
		uptime2 uint64
		want    time.Duration
	}{
		{
			name:    "uptime1 < uptime2",
			uptime1: 100,
			uptime2: 200,
			want:    100,
		},
		{
			name:    "uptime1 = uptime2",
			uptime1: 100,
			uptime2: 100,
			want:    0,
		},
		{
			name:    "uptime1 > uptime2",
			uptime1: 200,
			uptime2: 10,
			want:    -190,
		},
		{
			name:    "uptime1 > uptime2 error",
			uptime1: 200,
			uptime2: 220,
			want:    20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uptimeDiff(tt.uptime1, tt.uptime2)
			if got != tt.want {
				t.Errorf("uptimeDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uptimeToTime(t *testing.T) {
	tests := []struct {
		name      string
		nowUptime uint64
		uptime    uint64
		wantErr   bool
	}{
		{
			name:      "nowUptime < uptime",
			nowUptime: 100,
			uptime:    200,
			wantErr:   false,
		},
		{
			name:      "nowUptime = uptime",
			nowUptime: 100,
			uptime:    100,
			wantErr:   false,
		},
		{
			name:      "nowUptime > uptime",
			nowUptime: 200,
			uptime:    100,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			tm, err := uptimeToTime(tt.nowUptime, tt.uptime, now)
			if (err != nil) != tt.wantErr {
				t.Errorf("uptimeToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if now.Add(-time.Duration(tt.nowUptime)).Equal(tm) {
				t.Errorf("uptimeToTime() = %v, want %v", tm, now.Add(-time.Duration(tt.nowUptime)))
			}
		})
	}
}

func Test_timeToUptime(t *testing.T) {
	now := time.Now()
	nowUptime, err := GetKernelUptime()
	if err != nil {
		t.Errorf("GetKernelUptime() error = %v", err)
		return
	}
	tests := []struct {
		name    string
		t       time.Time
		addNano int64
		wantErr bool
	}{
		{
			name:    "t > now",
			t:       now.Add(time.Duration(time.Hour)),
			addNano: 1000000000 * 60 * 60,
			wantErr: false,
		},
		{
			name:    "t = now",
			t:       now,
			addNano: 0,
			wantErr: false,
		},
		{
			name:    "t < now",
			t:       now.Add(-time.Duration(time.Hour)),
			addNano: -1000000000 * 60 * 60,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uptime, err := timeToUptime(tt.t, now, nowUptime)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeToUptime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if uptime != uint64(int64(nowUptime)+tt.addNano) {
				t.Errorf("timeToUptime() = %v, want %v", uptime, uint64(int64(nowUptime)+tt.addNano))
			}
		})
	}
}
