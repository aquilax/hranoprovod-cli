package reporter

import "testing"

func Test_shorten(t *testing.T) {
	type args struct {
		t   string
		max int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Returns short string as it is",
			args{
				"this is a long string",
				50,
			},
			"this is a long string",
		},
		{
			"Shortens long string",
			args{
				"0123456789",
				7,
			},
			"012…789",
		},
		{
			"Shortens long string",
			args{
				"sweets/pasencia white/sm bonus/100g",
				20,
			},
			"sweets/pa…bonus/100g",
		},
	}
	shorten := getShorten(Options{
		ShortenStrings: true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shorten(tt.args.t, tt.args.max)
			if got != tt.want {
				t.Errorf("shorten() = %v, want %v", got, tt.want)
			}
			l := len([]rune(got))
			if l > tt.args.max {
				t.Errorf("shorten() length = %v, want less than %v", l, tt.args.max)
			}
		})
	}
}
