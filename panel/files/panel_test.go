package files

import "testing"

func TestFilesView(t *testing.T) {
	m := New(".")
	v := m.View()
	if v == "" { t.Fatal("empty") }
}
