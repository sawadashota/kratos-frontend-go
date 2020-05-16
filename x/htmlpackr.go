package x

import (
	"html/template"
	"io"

	"github.com/gobuffalo/packr/v2"
)

type Box struct {
	box *packr.Box
}

func NewBox(box *packr.Box) *Box {
	return &Box{
		box: box,
	}
}

const (
	layoutHTMLName  = "layout"
	contentHTMLName = "content"
)

func (b *Box) ParseHTML(name, layout, content string) (*HTMLTemplate, error) {
	t := template.New(name)

	layoutHTML, err := b.box.FindString(layout)
	if err != nil {
		return nil, err
	}
	t, err = t.New(layoutHTMLName).Parse(layoutHTML)
	if err != nil {
		return nil, err
	}

	contentHTML, err := b.box.FindString(content)
	if err != nil {
		return nil, err
	}
	t, err = t.New(contentHTMLName).Parse(contentHTML)
	if err != nil {
		return nil, err
	}
	return &HTMLTemplate{
		template: t,
	}, nil
}

func (b *Box) MustParseHTML(name, layout, content string) *HTMLTemplate {
	t, err := b.ParseHTML(name, layout, content)
	if err != nil {
		panic(err)
	}
	return t
}

type HTMLTemplate struct {
	template *template.Template
}

func (t *HTMLTemplate) Render(w io.Writer, data interface{}) error {
	return t.template.ExecuteTemplate(w, layoutHTMLName, data)
}
