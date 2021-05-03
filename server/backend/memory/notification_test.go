package memory

import (
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/notification"
)

func slurpNotifications(ch chan *notification.Notification) {
	if ch == nil {
		return
	}

	for range ch {
		<-ch
	}
}

func TestMemory_WriteNotification(t *testing.T) {
	type fields struct {
		notificationChMap map[string]chan *notification.Notification
	}
	type args struct {
		n *notification.Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid write",
			fields{
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
			},
			args{
				&notification.Notification{
					ClientID: "foo",
					Message:  "in a bottle",
				},
			},
			false,
		},
		{
			"write to nil channel",
			fields{
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
			},
			args{
				&notification.Notification{
					ClientID: "bar",
					Message:  "in a bottle",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				notificationChMap: tt.fields.notificationChMap,
			}
			for _, ch := range tt.fields.notificationChMap {
				go slurpNotifications(ch)
			}
			if err := m.WriteNotification(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Memory.WriteNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_GetNotifications(t *testing.T) {
	type fields struct {
		notificationChMap map[string]chan *notification.Notification
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		exists bool
	}{
		{
			"noexistent client",
			fields{
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
			},
			args{
				"foo",
			},
			true,
		},
		{
			"valid client",
			fields{
				notificationChMap: map[string]chan *notification.Notification{
					"foo": make(chan *notification.Notification),
				},
			},
			args{
				"bar",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				notificationChMap: tt.fields.notificationChMap,
			}
			if got, _ := m.GetNotifications(tt.args.id); (got == nil) == tt.exists {
				t.Errorf("Memory.GetNotifications() = %v, want %v", got, tt.exists)
			}
		})
	}
}
