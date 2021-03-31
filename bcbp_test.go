package bcbp

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// update is a flag to regenerate .golden files. It should be call any time
// a new test input file is added or there is a change in logic.
var update = flag.Bool("update", false, "update .golden files") //nolint

func TestFromStr_Errors(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want error
	}{
		{
			name: "Insufficient data",
			s:    "",
			want: ErrInsufficientData,
		},
		{
			name: "Non ASCII characters",
			s:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789世界",
			want: ErrNonASCII,
		},
		{
			name: "Unsupported format",
			s:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			want: ErrUnsupportedFormat,
		},
		{
			name: "Invalid field format (Number Of Legs Encoded)",
			s:    "MbcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			want: ErrInvalidFieldFormat,
		},
		{
			name: "Invalid field format (Beginning of Version Number)",
			s:    "M1DESMARAIS/LUC       EABC123 YULFRAAC 0834 326J001A0025 1010",
			want: ErrInvalidFieldFormat,
		},
		{
			name: "Invalid field format (Beginning of Security Data)",
			s:    "M1DESMARAIS/LUC       EABC123 YULFRAAC 0834 326J001A0025 1000",
			want: ErrInvalidFieldFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := FromStr(tt.s)
			if got == nil {
				t.Errorf("FromStr() = %+v: want %+v", got, tt.want)
			}

			if !errors.Is(got, tt.want) {
				t.Errorf("FromStr() = %+v: want %+v", got, tt.want)
			}
		})
	}
}

func TestFromStr(t *testing.T) {
	match, err := filepath.Glob("testdata/*.input")
	if err != nil {
		t.Fatal(err)
	}

	for _, in := range match {
		t.Run(in, func(t *testing.T) {
			data, err := os.ReadFile(in)
			if err != nil {
				t.Errorf("failed reading .input file: %v", err)
			}

			b, err := FromStr(string(data))
			if err != nil {
				t.Errorf("FromStr(%v) returned unexpected error: %+v", data, err)
			}

			got, err := json.MarshalIndent(b, "", "  ")
			if err != nil {
				t.Errorf("json.MarshalIndent() returned unexpected error: %+v", err)
			}

			runFromStrTest(t, got, in)
		})
	}
}

func runFromStrTest(t *testing.T, got []byte, in string) {
	t.Helper()

	out := in[:len(in)-len(".input")] + ".golden"
	if *update {
		t.Logf("update %s golden file", out)
		if err := os.WriteFile(out, got, 0666); err != nil {
			t.Fatalf("failed to update golden file: %v", err)
		}
	}

	want, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read golden file %v: %v", out, err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("output mismatch (-want +got):\n%s", diff)
	}
}

func benchmarkFromStr(in string, b *testing.B) {
	data, err := os.ReadFile(in)
	if err != nil {
		b.Errorf("failed reading .input file: %v", err)
	}

	s := string(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromStr(s)
	}
}

func BenchmarkFromStr_Mandatory_No_Security_Single(b *testing.B) {
	benchmarkFromStr("testdata/mandatory_no_security_single.input", b)
}

func BenchmarkFromStr_Mandatory_Single(b *testing.B) {
	benchmarkFromStr("testdata/mandatory_single.input", b)
}

func BenchmarkFromStr_Full_No_Security_Single(b *testing.B) {
	benchmarkFromStr("testdata/full_no_security_single.input", b)
}

func BenchmarkFromStr_Full_Single(b *testing.B) {
	benchmarkFromStr("testdata/full_single.input", b)
}

func BenchmarkFromStr_Full_No_Security_Multi(b *testing.B) {
	benchmarkFromStr("testdata/full_no_security_multi.input", b)
}

func BenchmarkFromStr_Full_Multi(b *testing.B) {
	benchmarkFromStr("testdata/full_multi.input", b)
}
