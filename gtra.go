package gtra

import (
	"errors"
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
func New() *Translate {
	t := new(Translate)

	l := NewConstLang()
	t.sl = l.AUTO
	t.tl = l.ZHCN

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
func (t *Translate) To(tl string, fn func(e error) string) string {
	r := get(t.url(tl, []string{"t"}))

	if r[0] != '[' {
		return fn(errors.New("DENY"))
	}

	p := gjson.Parse(r)

	return p.Get("0.0.0").String()
}

//From ..
func (t *Translate) From(sl string) *Translate {
	t.sl = sl

	return t
}

//Translate ..
func (t *Translate) Translate(str string) *Translate {
	t.str = str
	return t
}

//BySelf ..
func (t *Translate) BySelf(tl string, dt []string, fn func(e error) gjson.Result) gjson.Result {
	r := get(t.url(tl, dt))

	if r[0] != '[' {
		return fn(errors.New("DENY"))
	}

	return gjson.Parse(r)
}
