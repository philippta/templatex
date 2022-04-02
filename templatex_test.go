package templatex_test

import (
	"strings"
	"testing"

	"github.com/philippta/templatex"
	"github.com/philippta/templatex/test/templates"
)

func TestExecuteTemplate(t *testing.T) {
	tmpl := templatex.New()
	if err := tmpl.ParseDir("test/templates"); err != nil {
		t.Errorf("error executing template: %v", err)
	}

	{
		var b strings.Builder
		if err := tmpl.ExecuteTemplate(&b, "profile/payments/methods", "testdata"); err != nil {
			t.Errorf("error executing template: %v\n", err)
		}
		expect := "header layout profile payments_layout methods method footer \n"
		got := b.String()
		if expect != got {
			t.Errorf("expected:\n%v\ngot: %v\n", expect, got)
		}
	}

	{
		var b strings.Builder
		if err := tmpl.ExecuteTemplate(&b, "profile/edit", "testdata"); err != nil {

			t.Errorf("error executing template: %v\n", err)
		}

		expect := "header layout profile edit testdatafooter \n"
		got := b.String()
		if expect != got {
			t.Errorf("expected:\n%v\ngot: %v\n", expect, got)
		}
	}

	{
		var b strings.Builder
		if err := tmpl.ExecuteTemplate(&b, "profile/view", "testdata"); err != nil {
			t.Errorf("error executing template: %v\n", err)
		}
		expect := "header layout profile view testdatafooter \n"
		got := b.String()
		if expect != got {
			t.Errorf("expected:\n%v\ngot: %v\n", expect, got)
		}
	}
}

func TestExecuteTemplateFS(t *testing.T) {
	tmpl := templatex.New()

	if err := tmpl.ParseFS(templates.FS, "."); err != nil {
		t.Errorf("error executing template: %v", err)
	}

	{
		var b strings.Builder
		if err := tmpl.ExecuteTemplate(&b, "profile/payments/methods", "testdata"); err != nil {
			t.Errorf("error executing template: %v\n", err)
		}
		expect := "header layout profile payments_layout methods method footer \n"
		got := b.String()
		if expect != got {
			t.Errorf("expected:\n%v\ngot: %v\n", expect, got)
		}
	}

	{
		var b strings.Builder
		if err := tmpl.ExecuteTemplate(&b, "profile/edit", "testdata"); err != nil {

			t.Errorf("error executing template: %v\n", err)
		}

		expect := "header layout profile edit testdatafooter \n"
		got := b.String()
		if expect != got {
			t.Errorf("expected:\n%v\ngot: %v\n", expect, got)
		}
	}

	{
		var b strings.Builder
		if err := tmpl.ExecuteTemplate(&b, "profile/view", "testdata"); err != nil {
			t.Errorf("error executing template: %v\n", err)
		}
		expect := "header layout profile view testdatafooter \n"
		got := b.String()
		if expect != got {
			t.Errorf("expected:\n%v\ngot: %v\n", expect, got)
		}
	}
}
