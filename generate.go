package main

import (
	"encoding/json"
	"net/url"
	"os"
	"text/template"
)

type ReadmeData struct {
	KnownLanguages []Language `json:"known_languages"`
	WishLanguages  []Language `json:"wish_languages"`
	Love           []Badge    `json:"love"`
	Using          []Badge    `json:"using"`
	Interest       []Badge    `json:"interest"`
	DontAsk        []Badge    `json:"dont_ask"`
	Tools          []Badge    `json:"tools"`
}

// Uses https://simpleicons.org and https://shields.io
type Badge struct {
	Text      string `json:"text"`
	BGColor   string `json:"color"`
	Logo      string `json:"logo"`
	LogoColor string `json:"logo_color"`
}

func (b Badge) Url() string {
	base := "https://img.shields.io/badge/-" + url.PathEscape(b.Text) + "-" + b.BGColor
	if b.Logo != "" {
		base += "?style=flat&logo=" + b.Logo + "&logoColor=" + b.LogoColor
	}

	return base
}

type Language struct {
	Name        string `json:"name"`
	GithubEmoji string `json:"github_emoji"`
}

func (l Language) String() string {
	return l.Name + " " + l.GithubEmoji
}

var funcs = template.FuncMap{
	// From a list, do "a, b, c, and d"
	"join": func(a []Language) string {
		if len(a) == 0 {
			return ""
		}
		if len(a) == 1 {
			return a[0].String()
		}
		if len(a) == 2 {
			return a[0].String() + " and " + a[1].String()
		}
		s := ""
		for i, v := range a {
			if i == len(a)-1 {
				s += "and " + v.String()
			} else {
				s += v.String() + ", "
			}
		}
		return s
	},
	// Integrate image with markdown syntax
	"badge": func(b Badge) string {
		return "![" + b.Text + "](" + b.Url() + ")"
	},
}

func main() {
	f, err := os.Open("data.json")
	die(err)
	defer f.Close()

	var data ReadmeData
	die(json.NewDecoder(f).Decode(&data))

	tmpl, err := template.New("README.md.tmpl").Funcs(funcs).ParseFiles("README.md.tmpl")
	die(err)
	tmpl.Execute(os.Stdout, data)
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}
