package bcbp

import "testing"

func TestErrorType_String(t *testing.T) {
    defer func() {
        r := recover()
        if r == nil {
            t.Errorf("errorType.String() = nil: expected error")
        }
    }()

    et := ErrorType("test")
    _ = et.String()
}
