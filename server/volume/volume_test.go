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
		{
			"invalid volume - missing id",
			args{
				v: &Volume{
					Region: "us-west-2",
					Tag:    "bar",
				},
			},
			true,
		},
		{
			"invalid volume - missing region",
			args{
				v: &Volume{
					ID:  "foo",
					Tag: "bar",
				},
			},
			true,
		},
		{
			"invalid volume - missing tag",
			args{
				v: &Volume{
					ID:     "foo",
					Region: "us-west-2",
				},
			},
			true,
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
			"invalid id change",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
				&Volume{
					ID:     "bar",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Available,
				},
			},
			true,
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
			"invalid state change - available",
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
					Status: Attached,
				},
			},
			true,
		},
		{
			"invalid state change - attaching",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Attaching,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Unavailable,
				},
			},
			true,
		},
		{
			"invalid state change - attached",
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
					Tag:    "foo",
					Status: Attaching,
				},
			},
			true,
		},
		{
			"invalid state change - detaching",
			args{
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Detaching,
				},
				&Volume{
					ID:     "foo",
					Region: "us-west-2",
					Tag:    "foo",
					Status: Attaching,
				},
			},
			true,
		},
		{
			"invalid state change - unavailable",
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
					Status: Attaching,
				},
			},
			true,
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
