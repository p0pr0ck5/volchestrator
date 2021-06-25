package memory

import (
	"testing"

	"github.com/p0pr0ck5/volchestrator/server/notification"
)

func TestNewChQueue(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"valid new",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewChQueue()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChQueue_Read(t *testing.T) {
	type fields struct {
		ch         chan *notification.Notification
		shutdownCh chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantNil bool
		wantErr bool
	}{
		{
			"open queue",
			fields{
				ch:         make(chan *notification.Notification),
				shutdownCh: make(chan struct{}),
			},
			false,
			false,
		},
		{
			"closed queue",
			fields{
				ch:         make(chan *notification.Notification),
				shutdownCh: make(chan struct{}),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChQueue{
				ch:         tt.fields.ch,
				shutdownCh: tt.fields.shutdownCh,
			}
			if tt.wantNil {
				close(c.shutdownCh)
			}
			got, err := c.Read()
			if (got == nil) != tt.wantNil {
				t.Errorf("ChQueue.Read() chan = %v, wantNil %v", err, tt.wantNil)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ChQueue.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChQueue_Write(t *testing.T) {
	type fields struct {
		ch         chan *notification.Notification
		shutdownCh chan struct{}
	}
	type args struct {
		n *notification.Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		closed  bool
		wantErr bool
	}{
		{
			"open queue",
			fields{
				ch:         make(chan *notification.Notification),
				shutdownCh: make(chan struct{}),
			},
			args{
				&notification.Notification{},
			},
			false,
			false,
		},
		{
			"closed queue",
			fields{
				ch:         make(chan *notification.Notification),
				shutdownCh: make(chan struct{}),
			},
			args{
				&notification.Notification{},
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChQueue{
				ch:         tt.fields.ch,
				shutdownCh: tt.fields.shutdownCh,
			}
			if tt.closed {
				close(c.shutdownCh)
			}
			if err := c.Write(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("ChQueue.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChQueue_Close(t *testing.T) {
	type fields struct {
		ch         chan *notification.Notification
		shutdownCh chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		closed  bool
		wantErr bool
	}{
		{
			"open queue",
			fields{
				ch:         make(chan *notification.Notification),
				shutdownCh: make(chan struct{}),
			},
			false,
			false,
		},
		{
			"closed queue",
			fields{
				ch:         make(chan *notification.Notification),
				shutdownCh: make(chan struct{}),
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChQueue{
				ch:         tt.fields.ch,
				shutdownCh: tt.fields.shutdownCh,
			}
			if tt.closed {
				close(c.shutdownCh)
			}
			if err := c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("ChQueue.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
