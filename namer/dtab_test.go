package namer

import "testing"

type dtabtest struct {
	text   string
	ok     bool
	dtab   Dtab
	pretty string
}

var testdtabs = []dtabtest{
	dtabtest{"", true, []*Dentry{}, ""},
	dtabtest{
		`/foo=>/bar;/foo=>/bah#word`,
		true,
		[]*Dentry{&Dentry{"/foo", "/bar"}, &Dentry{"/foo", "/bah#word"}},
		"/foo  => /bar ;\n/foo  => /bah#word ;\n",
	},
	dtabtest{
		`/foo=>/bar;/foo/bar/baz=>/bah#word;`,
		true,
		[]*Dentry{&Dentry{"/foo", "/bar"}, &Dentry{"/foo/bar/baz", "/bah#word"}},
		"/foo          => /bar ;\n/foo/bar/baz  => /bah#word ;\n",
	},
}

func eqDtabs(dtab0, dtab1 Dtab) bool {
	if len(dtab0) != len(dtab1) {
		return false
	}
	for i, dentry0 := range dtab0 {
		if dentry0.Prefix != dtab1[i].Prefix || dentry0.Destination != dtab1[i].Destination {
			return false
		}
	}
	return true
}

func TestDtab(t *testing.T) {
	for _, test := range testdtabs {
		dtab, err := parseDtab(test.text)
		if test.ok {
			if err != nil {
				t.Error("unexpected parse error", err)
			} else {
				if !eqDtabs(dtab, test.dtab) {
					t.Errorf("expected dtab: '%s', got '%s'", test.dtab, dtab)
				}
				if pretty := dtab.Pretty(); pretty != test.pretty {
					t.Errorf("expected pretty: '%s', got '%s'",
						test.pretty, pretty)
				}
			}
		} else {
			if err == nil {
				t.Error("expected parser error got", dtab)
			}
		}
	}
}
