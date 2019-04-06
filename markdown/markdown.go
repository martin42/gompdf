package markdown

import (
	"strings"
)

const asterisk = "*"
const asterisks = "**"
const underscore = "_"
const underscores = "__"
const backtick = "`"
const backslash = `\`

type Item struct {
	Text    string
	Italic  bool
	Bold    bool
	Code    bool
	Newline bool
}

type Items []Item

func (is *Items) add(text string, italic, bold, code bool, nl bool) {
	*is = append(*is, Item{
		Text:    text,
		Italic:  italic,
		Bold:    bold,
		Code:    code,
		Newline: nl,
	})
}

func (i Item) WordItems() Items {
	is := Items{}
	words := []string{}
	currWord := ""
	for _, r := range i.Text {
		if r == ' ' {
			currWord += string(r)
			words = append(words, currWord)
			currWord = ""
		} else {
			currWord += string(r)
		}
	}
	if currWord != "" {
		words = append(words, currWord)
	}
	for _, word := range words {
		is.add(word, i.Italic, i.Bold, i.Code, i.Newline)
	}
	return is
}

func (is Items) WordItems() Items {
	wis := Items{}
	for _, i := range is {
		wis = append(wis, i.WordItems()...)
	}
	return wis
}

func NewProcessor() *Processor {
	return &Processor{}
}

type Processor struct {
}

func (p *Processor) Process(s string) Items {
	if s == "" {
		return Items{}
	}
	items := Items{}
	bs := []byte(s)
	bold := false
	italic := false
	code := false
	i := 0
	for {
		if i >= len(s) {
			return items
		}
		b := bs[i]
		if strings.HasPrefix(s[i:], asterisks) || strings.HasPrefix(s[i:], underscores) {
			bold = !bold
			items.add("", italic, bold, code, false)
			i += 2
		} else if strings.HasPrefix(s[i:], asterisk) || strings.HasPrefix(s[i:], underscore) {
			italic = !italic
			items.add("", italic, bold, code, false)
			i += 1
		} else if strings.HasPrefix(s[i:], backtick) {
			code = !code
			items.add("", italic, bold, code, false)
			i += 1
		} else if strings.HasPrefix(s[i:], backslash) {
			items.add(string(b), italic, bold, code, true)
			items.add("", italic, bold, code, false)
			i += 1
		} else {
			if len(items) == 0 {
				items.add(string(b), italic, bold, code, false)
			} else {
				ibs := []byte(items[len(items)-1].Text)
				ibs = append(ibs, b)
				items[len(items)-1].Text = string(ibs)
			}
			i += 1
		}
	}
}
