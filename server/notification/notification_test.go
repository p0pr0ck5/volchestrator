package notification

import (
	"reflect"
	"testing"
)

func TestNotification_SetMessageID(t *testing.T) {
	type fields struct {
		ClientID  string
		Message   string
		messageID uint64
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Notification
	}{
		{
			"no existing ID",
			fields{
				ClientID: "foo",
				Message:  "bar",
			},
			args{
				id: 1,
			},
			&Notification{
				ClientID:  "foo",
				Message:   "bar",
				messageID: 1,
			},
		},
		{
			"existing ID",
			fields{
				ClientID:  "foo",
				Message:   "bar",
				messageID: 2,
			},
			args{
				id: 1,
			},
			&Notification{
				ClientID:  "foo",
				Message:   "bar",
				messageID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Notification{
				ClientID:  tt.fields.ClientID,
				Message:   tt.fields.Message,
				messageID: tt.fields.messageID,
			}
			n.SetMessageID(tt.args.id)
			if !reflect.DeepEqual(n, tt.want) {
				t.Errorf("NotificationSetMessageID() = %v, want %v", n, tt.want)
			}
		})
	}
}
