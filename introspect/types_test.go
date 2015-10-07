package introspect

import "testing"

func TestParseParenthess(t *testing.T) {
	var data = []struct {
		signautre            string
		lSep, rSep           rune
		rLeft, rMide, rRight string
	}{
		{"(nnqq)", '(', ')', "", "nnqq", ""},
		{"(nn(ss)qq)", '(', ')', "", "nn(ss)qq", ""},
		{"nn(ss)qq", '(', ')', "nn", "ss", "qq"},
		{"()", '(', ')', "", "", ""},
		{"i(ii)a(s)", '(', ')', "i", "ii", "a(s)"},
	}
	for _, d := range data {
		l, m, r := ParseParenthess(d.signautre, d.lSep, d.rSep)
		if d.rLeft != l || d.rMide != m || d.rRight != r {
			t.Fatalf("%q split to %q %q %q\n", d.signautre, l, m, r)
		}
	}
}
func TestSplit(t *testing.T) {
	var data = []struct {
		signature string
		items     []string
	}{
		{"(nqqq)", []string{"(nqqq)"}},
		{"nqqq", []string{"n", "q", "q", "q"}},
		{"aa{sv}", []string{"aa{sv}"}},
		{"a(ai)(si)", []string{"a(ai)", "(si)"}},
		{"asaia{sd}i", []string{"as", "ai", "a{sd}", "i"}},
	}
	for _, d := range data {
		for i, v := range Split(d.signature) {
			if d.items[i] != v {
				t.Fatalf("%q split to %v\n", d.signature, Split(d.signature))
			}
		}
	}
}
