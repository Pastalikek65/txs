package sys

import "testing"

func TestSysView(t *testing.T) {
	m := New()
	v := m.View()
	if v == "" { t.Fatal("empty") }
}
