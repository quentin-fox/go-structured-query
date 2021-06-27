package postgres

import (
	"strings"
	"testing"

	"github.com/matryer/is"

	"go/parser"
	"go/token"
)

func TestTablesTemplate(t *testing.T) {
	is := is.New(t)

	template, err := getTablesTemplate()
	is.NoErr(err)

	var writer strings.Builder

	data := TablesTemplateData{
		PackageName: "tables",
		Imports: []string{
			`sq "github.com/bokwoon95/go-structured-query"`,
		},
		Tables: []Table{
			{
				Name:        "users",
				Schema:      "public",
				StructName:  "USERS",
				RawType:     "BASE TABLE",
				Constructor: "TABLE_USERS",
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

	expected := `// Code generated by 'sqgen-postgres tables'; DO NOT EDIT.
package tables

import (
	sq "github.com/bokwoon95/go-structured-query"
)

// USERS references the public.users table.
type USERS struct {
	*sq.TableInfo
	ID sq.NumberField
	FIRST_NAME sq.StringField
	DATE_CREATED sq.TimeField
}

// TABLE_USERS creates an instance of the public.users table.
func TABLE_USERS() USERS {
	tbl := USERS{TableInfo: &sq.TableInfo{
		Schema: "public",
		Name: "users",
	},}
	tbl.ID = sq.NewNumberField("id", tbl.TableInfo)
	tbl.FIRST_NAME = sq.NewStringField("first_name", tbl.TableInfo)
	tbl.DATE_CREATED = sq.NewTimeField("date_created", tbl.TableInfo)
	return tbl
}

// As modifies the alias of the underlying table.
func (tbl USERS) As(alias string) USERS {
	tbl.TableInfo.Alias = alias
	return tbl
}`
	is.Equal(out, expected)

	// checks that the go parser can parse the contents of out to an AST
	fs := token.NewFileSet()
	_, err = parser.ParseFile(fs, "", out, parser.AllErrors)
	is.NoErr(err)
}

func TestFunctionsTemplate(t *testing.T) {
	is := is.New(t)

	template, err := getFunctionsTemplate()
	is.NoErr(err)

	var writer strings.Builder

	data := FunctionsTemplateData{
		PackageName: "tables",
		Imports: []string{
			`sq "github.com/bokwoon95/go-structured-query"`,
		},
		Functions: []Function{
			{
				Name:        "insert_user",
				Schema:      "public",
				StructName:  "FUNCTION_INSERT_USER",
				Constructor: "INSERT_USER",
				Arguments: []FunctionField{
					{
						Name:        "first_name",
						GoType:      GoTypeString,
						FieldType:   FieldTypeString,
						Constructor: FieldConstructorString,
					},
					{
						Name:        "date_created",
						GoType:      GoTypeTime,
						FieldType:   FieldTypeTime,
						Constructor: FieldConstructorTime,
					},
				},
				Results: []FunctionField{
					{
						Name:        "user_id",
						GoType:      GoTypeInt,
						FieldType:   FieldTypeNumber,
						Constructor: FieldConstructorNumber,
					},
				},
			},
		},
	}

	err = template.Execute(&writer, data)
	is.NoErr(err)

	out := writer.String()

	expected := `// Code generated by 'sqgen-postgres functions'; DO NOT EDIT.
package tables

import (
	sq "github.com/bokwoon95/go-structured-query"
)

// FUNCTION_INSERT_USER references the public.insert_user function.
type FUNCTION_INSERT_USER struct {
	*sq.FunctionInfo
	USER_ID sq.NumberField
}

// INSERT_USER creates an instance of the public.insert_user function.
func INSERT_USER(
	first_name string,
	date_created time.Time,
	) FUNCTION_INSERT_USER {
	return INSERT_USER_(first_name, date_created)
}

// INSERT_USER_ creates an instance of the public.insert_user function.
func INSERT_USER_(
	first_name interface{},
	date_created interface{},
	) FUNCTION_INSERT_USER {
	f := FUNCTION_INSERT_USER{FunctionInfo: &sq.FunctionInfo{
		Schema: "public",
		Name: "insert_user",
		Arguments: []interface{}{first_name, date_created},
	},}
	f.USER_ID = sq.NewNumberField("user_id", f.FunctionInfo)
	return f
}

// As modifies the alias of the underlying function.
func (f FUNCTION_INSERT_USER) As(alias string) FUNCTION_INSERT_USER {
	f.FunctionInfo.Alias = alias
	return f
}`

	is.Equal(len(out), len(expected))
	is.Equal(out, expected)

	// TODO once function template is fully stable, test exact output vs. expected string

	fs := token.NewFileSet()
	_, err = parser.ParseFile(fs, "", out, parser.AllErrors)
	is.NoErr(err)
}
