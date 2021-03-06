package tpl_test

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	. "github.com/wwgberlin/go-weather-widget/tpl"
)

const layoutTemplateName = "layout"
const contentTemplateName = "content"
const headTemplateName = "head"

func TestTemplateLayout(t *testing.T) {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Fatalf("Template layout.tmpl was expected to parse without any errors. %v", err)
	}

	if err = tmpl.ExecuteTemplate(&b, layoutTemplateName, "some data"); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&b)

	if body := doc.Find("html body"); body.Length() == 0 {
		t.Error("Expected to render html and body elements")
	}

	if tmpl.Lookup(layoutTemplateName) == nil {
		t.Error("Template layout.tmpl was expected to define template layout")
	}

	if tmpl.Lookup(contentTemplateName) == nil {
		t.Error("Template layout.tmpl was expected to define empty template content")
	}

	if tmpl.Lookup(headTemplateName) == nil {
		t.Error("Template layout.tmpl was expected to define empty template head")
	}
}

func TestLayoutWithHead(t *testing.T) {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Fatalf("Template head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(fmt.Sprintf(
		`{{define "%s"}}
			<head><title>{{.}}</title></head>
		{{end}}`, headTemplateName))

	if err = tmpl.ExecuteTemplate(&b, layoutTemplateName, "title"); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(&b)

	if title := doc.Find("title"); title.Length() == 0 {
		t.Error("Expected to render title element")
	} else if title.Text() != "title" {
		t.Errorf("Template head.tmpl was expected to be rendered with title 'title' but got '%s'", title.Text())
	}
}

func TestLayoutWithContent(t *testing.T) {
	var b bytes.Buffer
	const expected = "some data"
	tmpl, err := template.ParseFiles("./templates/layouts/layout.tmpl")

	if err != nil {
		t.Fatalf("head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "content"}}{{.}}{{end}}`)

	if err = tmpl.ExecuteTemplate(&b, "layout", expected); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(&b)

	if body := doc.Find("body"); body.Length() == 0 {
		t.Error("Expected to render body element")
	} else {
		txt := strings.TrimSpace(body.Text())
		if txt != expected {
			t.Errorf("Template layout.tmpl was expected to be rendered with '%s' but got '%s'", expected, txt)
		}
	}
}

func TestTemplateHead(t *testing.T) {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")
	if err != nil {
		t.Fatalf("Template head.tmpl was expected to parse without any errors. %v", err)
	}

	if tmpl.Lookup("styles") == nil {
		t.Error("Template head.tmpl was expected to define empty template styles")
	}

	if tmpl.Lookup("title") == nil {
		t.Error("Template head.tmpl was expected to define empty template title")
	}

	if err = tmpl.ExecuteTemplate(&b, "head", nil); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&b)

	if head := doc.Find("head"); head.Length() == 0 {
		t.Error("Expected to render head element")
	}
}

func TestTemplatesHeadWithTitle(t *testing.T) {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Fatalf("Template head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "title"}}<title>{{.}}</title>{{end}}`)

	if tmpl.ExecuteTemplate(&b, "head", "some_title"); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&b)

	if title := doc.Find("head title"); title.Length() == 0 {
		t.Error("Expected to render title element")
	} else {
		txt := strings.TrimSpace(title.Text())
		if txt != "some_title" {
			t.Errorf("Head element was expected to have title 'some_title', but got '%s'", txt)
		}
	}
}

func TestTemplatesHeadWithStyles(t *testing.T) {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/layouts/head.tmpl")

	if err != nil {
		t.Fatalf("Template head.tmpl was expected to parse without any errors. %v", err)
	}

	tmpl, err = tmpl.Parse(`{{define "styles"}}<link rel={{.}}/>{{end}}`)

	if err = tmpl.ExecuteTemplate(&b, "head", "some_link"); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(&b)

	if link := doc.Find("link"); link.Length() == 0 {
		t.Error("Expected to render link element")
	} else {
		rel, ok := link.Attr("rel")
		if !ok || rel != "some_link/" {
			t.Errorf("Head element expected to have link 'some_link', but got '%s'", rel)
		}
	}
}

func TestTemplateWidget(t *testing.T) {
	var b bytes.Buffer
	h := copyFuncMap(DefaultHelpers)
	if _, ok := h["clothes"]; !ok {
		t.Fatal("Function clothes was expected to be added to DefaultHelpers")
	}

	h["clothes"] = myClothes("crown", "cape")

	tmpl := template.New("widget").Funcs(h)
	tmpl, err := tmpl.ParseFiles("./templates/widget.tmpl")

	if err != nil {
		t.Fatalf("widget.tmpl was expected to parse without any errors. %v", err)
	}

	if tmpl.Lookup("content") == nil {
		t.Error("widget.tmpl was expected to define template content")
	}

	if tmpl.Lookup("styles") == nil {
		t.Error("widget.tmpl was expected to define template styles")
	}

	if err = tmpl.ExecuteTemplate(&b, "content", map[string]interface{}{
		"location":    "Berlin",
		"description": "It's spring time",
		"celsius":     25,
	}); err != nil {
		t.Fatalf("Template was expected to execute without errors. %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&b)
	if doc.Find("div.gopher").Length() == 0 {
		t.Error("expected to render div with class 'gopher'")
	} else {
		gopherDiv := doc.Find("div.gopher")
		if gopherDiv.Find(".crown").Length() == 0 {
			t.Error("gopher was expected to wear a crown")
		}
		if gopherDiv.Find(".cape").Length() == 0 {
			t.Error("gopher was expected to wear a cape")
		}
	}
}

func myClothes(ret ...string) func(args ...interface{}) ([]string, error) {
	return func(args ...interface{}) ([]string, error) {
		if len(args) < 2 {
			return nil, errors.New("clothe expects 2 arguments to be passed (description, celsius)")
		}
		if desc, ok := args[0].(string); !ok {
			return nil, errors.New("first argument in clothes was expected to be a string (description)")
		} else if desc != "It's spring time" {
			return nil, errors.New("first argument in clothes was expected to be the weather description")
		}
		if celsius, ok := args[1].(int); !ok {
			return nil, errors.New("second argument in clothes was expected to be an integer (celsius)")
		} else if celsius != 25 {
			return nil, errors.New("first argument in clothes was expected to be the weather celsius")
		}
		return ret, nil
	}
}

func copyFuncMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
