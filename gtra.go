package gtra

import (
	"errors"
	"log"
	"net/url"
	"strconv"

	"github.com/tidwall/gjson"
)

//Translate ..
type Translate struct {
	str string

	sl string
	tl string
}

//New ..
func New(str string) *Translate {
	t := new(Translate)
	t.str = str

	t.sl = LangAUTO
	t.tl = LangZHCN

	return t

}
func (t *Translate) tk() string {
	tkki, tkkf := gettkk()
	tki, tkf := gettk(int32(tkki), int32(tkkf), t.str)

	return (strconv.Itoa(int(tki)) + "." + strconv.Itoa(int(tkf)))
}
func (t *Translate) url(tl string, dt []string) string {
	urls := "https://translate.google.cn/translate_a/single?"

	urls = urls + "&client=gtx"
	urls = urls + "&sl=" + t.sl
	urls = urls + "&tl=" + tl
	urls = urls + "&hl=" + tl
	for i := range dt {
		urls = urls + "&dt=" + dt[i]
	}

	urls = urls + "&q=" + url.QueryEscape(t.str)
	urls = urls + "&tk=" + t.tk()

	return urls
}

//To ..
func (t *Translate) To(tl string) (string, error) {
	r := get(t.url(tl, []string{"t"}))

	if r[0] != '[' {
		log.Panicln("server deny!")
		e := errors.New("DENY")
		return "", e
	}

	p := gjson.Parse(r)

	return p.Get("0.0.0").String(), nil
}

//From ..
func (t *Translate) From(sl string) *Translate {
	t.sl = sl

	return t
}

//TranslateBySelf ..
func (t *Translate) TranslateBySelf(tl string, dt []string) (gjson.Result, error) {
	r := get(t.url(tl, dt))

	if r[0] != '[' {
		log.Panicln("server deny!")
		e := errors.New("DENY")
		return gjson.Parse(""), e
	}

	return gjson.Parse(r), nil
}
