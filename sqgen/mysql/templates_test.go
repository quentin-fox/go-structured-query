package mysql

import (
	"strings"
	"testing"

	"go/parser"
	"go/token"

	"github.com/matryer/is"
)

func TestTablesTemplate(t *testing.T) {
	is := is.New(t)

	template, err := getTablesTemplate()
	is.NoErr(err)

	var writer strings.Builder

	data := TablesTemplateData{
		PackageName: "tables",
		Imports: []string{
			`sq "github.com/quentin-fox/go-structured-query"`,
		},
		Tables: []Table{
			{
				Name:        "users",
				Schema:      "public",
				StructName:  "TABLE_USERS",
				RawType:     "BASE TABLE",
				Constructor: "USERS",
				Fields: []TableField{
					{
						Name:        "id",
						RawType:     "integer",
						Type:        FieldTypeNumber,
						Constructor: FieldConstructorNumber,
					},
					{
						Name:        "first_name",
						RawType:     "text",
						Type:        FieldTypeString,
						Constructor: FieldConstructorString,
					},
					{
						Name:        "date_created",
						RawType:     "timestamp",
						Type:        FieldTypeTime,
						Constructor: FieldConstructorTime,
					},
				},
			},
		},
	}

	err = template.Execute(&writer, data)
	is.NoErr(err)

	out := writer.String()

	expected := `// Code generated by 'sqgen-mysql tables'; DO NOT EDIT.
package tables

import (
	sq "github.com/quentin-fox/go-structured-query"
)

// TABLE_USERS references the public.users table.
type TABLE_USERS struct {
	*sq.TableInfo
	ID sq.NumberField
	FIRST_NAME sq.StringField
	DATE_CREATED sq.TimeField
}

// USERS creates an instance of the public.users table.
func USERS() TABLE_USERS {
	tbl := TABLE_USERS{TableInfo: &sq.TableInfo{
		Schema: "public",
		Name: "users",
	},}
	tbl.ID = sq.NewNumberField("id", tbl.TableInfo)
	tbl.FIRST_NAME = sq.NewStringField("first_name", tbl.TableInfo)
	tbl.DATE_CREATED = sq.NewTimeField("date_created", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl TABLE_USERS) As(alias string) TABLE_USERS {
	tbl.TableInfo.Alias = alias
	return tbl
}`

	is.Equal(out, expected)

	// checks that the go parser can parse the contents of out to an AST
	fs := token.NewFileSet()
	_, err = parser.ParseFile(fs, "", out, parser.AllErrors)
	is.NoErr(err)
}
