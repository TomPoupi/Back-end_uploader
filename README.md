# uploader
https://github.com/TomPoupi/uploader

je m'apprête à faire des choses, je travail sur intégration


exemple : testunitaire :

import "testing"
func TestFooerTableDriven(t *testing.T) {
      // Defining the columns of the table
        var tests = []struct {
        name string
            input int
            want  string
        }{
            // the table itself
            {"9 should be Foo", 9, "Foo"},
            {"3 should be Foo", 3, "Foo"},
            {"1 is not Foo", 1, "1"},
            {"0 should be Foo", 0, "Foo"},
        }
      // The execution loop
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                ans := Fooer(tt.input)
                if ans != tt.want {
                    t.Errorf("got %s, want %s", ans, tt.want)
                }
            })
        }
    }