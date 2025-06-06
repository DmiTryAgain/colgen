package colgen

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseRules(t *testing.T) {
	type args struct {
		lines         []string
		useListSuffix bool
	}
	tests := []struct {
		name    string
		args    args
		want    []Rule
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				lines: []string{
					"News,Tag,Category",
					"News:MapP(db)",
					"Tag:Index(OrderNumber),OrderNumber",
				},
				useListSuffix: false,
			},
			want: []Rule{
				{
					EntityName:    "Category",
					BaseGen:       true,
					UseListSuffix: false,
				},
				{
					EntityName:    "News",
					BaseGen:       true,
					UseListSuffix: false,
					CustomRules: []CustomRule{{
						Name: "MapP", Field: "", Arg: "db"},
					},
				},

				{
					EntityName:    "Tag",
					BaseGen:       true,
					UseListSuffix: false,
					CustomRules: []CustomRule{
						{Name: "Index", Field: "OrderNumber", Arg: ""},
						{Name: "", Field: "OrderNumber", Arg: ""},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "simple",
			args: args{
				lines: []string{
					"News",
				},
				useListSuffix: false,
			},
			want: []Rule{
				{
					EntityName:    "News",
					BaseGen:       true,
					UseListSuffix: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRules(tt.args.lines, tt.args.useListSuffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRules() \n\tgot  = %v, \n\twant = %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	const want = `// Code generated by colgen devel; DO NOT EDIT.
package newsportal

import (
	"pkg/db"
	"pkg/newsportal"
)

type Categories []Category

func (ll Categories) IDs() []int {
	r := make([]int, len(ll))
	for i := range ll {
		r[i] = ll[i].ID
	}
	return r
}

func (ll Categories) Index() map[int]Category {
	r := make(map[int]Category, len(ll))
	for i := range ll {
		r[ll[i].ID] = ll[i]
	}
	return r
}

type NewsList []News

func (ll NewsList) IDs() []int {
	r := make([]int, len(ll))
	for i := range ll {
		r[i] = ll[i].ID
	}
	return r
}

func (ll NewsList) Index() map[int]News {
	r := make(map[int]News, len(ll))
	for i := range ll {
		r[ll[i].ID] = ll[i]
	}
	return r
}

func NewNewsList(in []db.News) NewsList { return MapP(in, NewNews) }

type Tags []Tag

func (ll Tags) IDs() []int {
	r := make([]int, len(ll))
	for i := range ll {
		r[i] = ll[i].ID
	}
	return r
}

func (ll Tags) Index() map[int]Tag {
	r := make(map[int]Tag, len(ll))
	for i := range ll {
		r[ll[i].ID] = ll[i]
	}
	return r
}

func (ll Tags) IndexByOrderNumber() map[int64]Tag {
	r := make(map[int64]Tag, len(ll))
	for i := range ll {
		r[ll[i].OrderNumber] = ll[i]
	}
	return r
}

func (ll Tags) OrderNumbers() []int64 {
	r := make([]int64, len(ll))
	for i := range ll {
		r[i] = ll[i].OrderNumber
	}
	return r
}

func (ll Tags) UniqueOrderNumbers() []int64 {
	idx := make(map[int64]struct{})
	for i := range ll {
		if _, ok := idx[ll[i].OrderNumber]; !ok {
			idx[ll[i].OrderNumber] = struct{}{}
		}
	}

	r, i := make([]int64, len(idx)), 0
	for k := range idx {
		r[i] = k
		i++
	}
	return r
}
`

	// start
	raw := `
News,Tag,Category
News:MapP(db)
Tag:Index(OrderNumber),OrderNumber,UniqueOrderNumber
`
	lines := strings.Split(raw, "\n")
	imports := "pkg/db,pkg/newsportal"

	g := NewGenerator("newsportal", imports, "", "devel")
	rules, err := ParseRules(lines, false)
	if err != nil {
		t.Fatal(err)
	}

	err = g.UsePackageDir(".")
	if err != nil {
		t.Fatal(err)
	}

	data, err := g.Generate(rules)
	if err != nil {
		t.Fatal(err)
	}

	dataF, err := g.Format()
	if err != nil {
		t.Fatal(err)
	}

	if string(dataF) != want {
		t.Errorf("Generate() different result:\ngot  = %q\nwant = %q", data, want)

		os.Stdout.Write(dataF)
	}
}
