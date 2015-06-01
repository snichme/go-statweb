package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/russross/blackfriday"
)

type PageName string

type Page struct {
	name PageName
}

type PageSettings struct {
	Layout string `json:"layout"`
}

func (p *PageSettings) GetLayout() string {
	if p.Layout != "" {
		return p.Layout
	}
	return "default.html"
}

type PageData struct {
	Settings PageSettings           `json:"settings"`
	Data     map[string]interface{} `json:"data"`
}

func (p *Page) path(filetype string) string {
	return fmt.Sprintf("page/%s.%s", p.name, filetype)
}

func (p *Page) renderMd() ([]byte, error) {
	md, err := ioutil.ReadFile(p.path("md"))
	if err != nil {
		return []byte{}, err
	}
	return blackfriday.MarkdownCommon(md), nil

}

// TODO: Needs to return defalt page data or give error
func (p *Page) getPageData() *PageData {
	var data = &PageData{}
	path := p.path("json")
	if _, err := os.Stat(path); err == nil {
		js, _ := ioutil.ReadFile(path)
		json.Unmarshal(js, &data)
	}
	return data
}

func (p *Page) Render() Response {
	body, err := p.renderMd()
	data := p.getPageData()
	layout := data.Settings.GetLayout()
	if err != nil {
		if os.IsNotExist(err) {
			return NotFound
		}
		return Error("Problems", err)
	}
	data.Data["Body"] = template.HTML(body)
	return Template(200, layout, data.Data)
}
