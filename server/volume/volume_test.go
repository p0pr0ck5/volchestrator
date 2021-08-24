package volume

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type args struct {
		v *Volume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid volume",
			args{
				v: &Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bar",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVolume_ValidateTransition(t *testing.T) {
	type args struct {
		currentVolume *Volume
		newVolume     *Volume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"valid no change",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
			},
			false,
		},
		{
			"valid region change",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-1",
					Tag:    "foo",
					Status: Available,
				},
			},
			false,
		},
		{
			"valid tag change",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bar",
					Status: Available,
				},
			},
			false,
		},
		{
			"invalid region change",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Attached,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-1",
					Tag:    "foo",
					Status: Attached,
				},
			},
			true,
		},
		{
			"invalid tag change",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Attached,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "bar",
					Status: Attached,
				},
			},
			true,
		},
		{
			"valid state change - available",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Unavailable,
				},
			},
			false,
		},
		{
			"valid state change - unavailable",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Unavailable,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Deleting,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.args.currentVolume
			v.Init()
			if err := v.ValidateTransition(tt.args.newVolume); (err != nil) != tt.wantErr {
				t.Errorf("Volume.ValidateTransition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
