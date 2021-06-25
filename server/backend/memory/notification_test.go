package memory

import (
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/notification"
)

func slurpNotifications(ch <-chan *notification.Notification) {
	if ch == nil {
		return
	}

	for range ch {
		<-ch
	}
}

func TestMemory_WriteNotification(t *testing.T) {
	type fields struct {
		notificationMap map[string]*ChQueue
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
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
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
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
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
				notificationMap: tt.fields.notificationMap,
			}
			for _, queue := range tt.fields.notificationMap {
				notificationCh, _ := queue.Read()
				go slurpNotifications(notificationCh)
			}
			if err := m.WriteNotification(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Memory.WriteNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_GetNotifications(t *testing.T) {
	type fields struct {
		notificationMap map[string]*ChQueue
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
			"valid client",
			fields{
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				"bar",
			},
			false,
		},
		{
			"noexistent client",
			fields{
				notificationMap: map[string]*ChQueue{
					"foo": MustNewChQueue(),
				},
			},
			args{
				"foo",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				notificationMap: tt.fields.notificationMap,
			}
			if got, _ := m.GetNotifications(tt.args.id); (got == nil) == tt.exists {
				t.Errorf("Memory.GetNotifications() = %v, want %v", got, tt.exists)
			}
		})
	}
}
