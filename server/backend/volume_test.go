package backend

import (
	"reflect"
	"testing"
	"time"

	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func TestBackend_UpdateVolume(t *testing.T) {
	mockVolumes := []*volume.Volume{
		{
			ID:     "foo",
			Region: "us-west-2",
			Tag:    "bar",
			Status: volume.Available,
		},
		{
			ID:     "bar",
			Region: "us-west-1",
			Tag:    "baz",
			Status: volume.Unavailable,
		},
	}

	type fields struct {
		b backend
	}
	type args struct {
		v *volume.Volume
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *volume.Volume
		wantErr bool
	}{
		{
			"full update",
			fields{
				b: NewMemoryBackend(WithVolumes(mockVolumes)),
			},
			args{
				v: &volume.Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "baz",
					Status: volume.Available,
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-2",
				Tag:    "baz",
				Status: volume.Available,
			},
			false,
		},
		{
			"partial update",
			fields{
				b: NewMemoryBackend(WithVolumes(mockVolumes)),
			},
			args{
				v: &volume.Volume{
					ID:     "foo",
					Region: "us-west-1",
				},
			},
			&volume.Volume{
				ID:     "foo",
				Region: "us-west-1",
				Tag:    "bar",
				Status: volume.Available,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				b: tt.fields.b,
			}
			if err := b.UpdateVolume(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Backend.UpdateVolume() error = %v, wantErr %v", err, tt.wantErr)
			}

			v, _ := b.ReadVolume(tt.want.ID)
			v.UpdatedAt = time.Time{} // mock

			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("Backend.UpdateVolume() = %v, want %v", v, tt.want)
			}
		})
	}
}
