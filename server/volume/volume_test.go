package volume

import "testing"

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
			if err := Validate(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
