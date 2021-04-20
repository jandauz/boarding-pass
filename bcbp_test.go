package bcbp

import (
    "encoding/json"
    "flag"
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/google/go-cmp/cmp"
)

// update is a flag to regenerate .golden files. It should be call any time
// a new test input file is added or there is a change in logic.
var update = flag.Bool("update", false, "update .golden files") //nolint

func TestFromStr(t *testing.T) {
    testFromStr(t, "testdata/*.input", false)
}

func TestFromStr_Errors(t *testing.T) {
    testFromStr(t, "testdata/errors/*.input", true)
}

func testFromStr(t *testing.T, in string, wantErr bool) {
    t.Helper()

    match, err := filepath.Glob(in)
    if err != nil {
        t.Fatal(err)
    }

    for _, in := range match {
        t.Run(in, func(t *testing.T) {
            data, err := os.ReadFile(in)
            if err != nil {
                t.Errorf("failed reading .input file: %v", err)
                return
            }

            // Special case for malformed spec scenario where spec needs
            // to be modified
            if strings.Contains(in, "malformed_spec") {
                defer func() {
                    spec[0].items = nil
                }()
                spec[0].items = []item{}
            }

            var got []byte
            b, err := FromStr(string(data))
            switch {
            case wantErr && err == nil:
                t.Error("FromStr() = nil: expected error")
                return
            case wantErr && err != nil:
                got = []byte(err.Error())
            case err != nil:
                t.Errorf("FromStr(%s) returned unexpected error: %+v", data, err)
                return
            default:
                got, err = json.MarshalIndent(b, "", "  ")
                if err != nil {
                    t.Errorf("json.MarshalIndent() returned unexpected error: %+v", err)
                }
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
