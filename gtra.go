package gtra

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhongmin/errs"

	"github.com/jinzhongmin/gtra/lang"
	"github.com/tidwall/gjson"
)

var ERR_SERVER_DENY = errors.New("server deny.")
var ERR_LANG_CODE_NOT_EXIST = errors.New("Language code does not exist")

type Translater struct {
	sl string
	tl string
}

func NewTranslater() *Translater {
	t := new(Translater)
	t.sl = lang.AUTO
	t.tl = lang.ZHCN
	return t
}

func (t *Translater) Vector(sl string, tl string) *Translater {
	if lang.Verify(sl) {
		t.sl = sl
	} else {
		errs.Warn(ERR_LANG_CODE_NOT_EXIST)
		t.sl = lang.AUTO
	}

	if lang.Verify(tl) {
		t.tl = tl
		return t
	}
	errs.Warn(ERR_LANG_CODE_NOT_EXIST)
	return t
}

func (t *Translater) To(tl string) *Translater {
	if lang.Verify(tl) {
		t.tl = tl
		return t
	}
	errs.Warn(ERR_LANG_CODE_NOT_EXIST)
	return t
}

func (t *Translater) Dt(src string, dt ...string) (error, gjson.Result) {
	r := getResult(src, t.sl, t.tl, dt...)
	if r[0] != '[' {
		return ERR_SERVER_DENY, gjson.Result{}
	}
	return nil, gjson.Parse(r)
}

func (t *Translater) Translate(src string) (error, string) {
	r := getResult(src, t.sl, t.tl, "t")
	if r[0] != '[' {
		return ERR_SERVER_DENY, ""
	}
	p := gjson.Parse(r)
	return nil, p.Get("0.0.0").String()
}

func Translate(src string, tl string, sl ...string) (error, string) {
	t := NewTranslater()
	if len(sl) == 1 {
		if lang.Verify(sl[0]) {
			t.sl = sl[0]
		} else {
			errs.Warn(ERR_LANG_CODE_NOT_EXIST)
			t.sl = lang.AUTO
		}
	}
	if lang.Verify(tl) {
		t.tl = tl
	} else {
		errs.Warn(ERR_LANG_CODE_NOT_EXIST)
	}

	r := getResult(src, t.sl, t.tl, "t")
	if r[0] != '[' {
		return ERR_SERVER_DENY, ""
	}

	p := gjson.Parse(r)
	return nil, p.Get("0.0.0").String()
}

//tk get
func tk(s string) string {
	tkki, tkkf := gettkk()
	tki, tkf := gettk(int32(tkki), int32(tkkf), s)

	return (strconv.Itoa(int(tki)) + "." + strconv.Itoa(int(tkf)))
}
func gettkk() (int64, int64) {
	s := get("https://translate.google.cn/")
	reg := regexp.MustCompile(`tkk:'.*?'`)
	tkkc := reg.FindString(s)
	tkkc = strings.Replace(tkkc, "tkk:", "", -1)
	tkkv := strings.Split(tkkc, ".")
	a, _ := strconv.ParseInt(tkkv[0], 10, 64)
	b, _ := strconv.ParseInt(tkkv[1], 10, 64)

	return a, b
}
func gettk(tkka int32, tkkb int32, str string) (int32, int32) {
	b := make([]byte, len(str))
	b = []byte(str)
	a := tkka
	for i := 0; i < len(b); i++ {
		a += int32((b[i]))
		a = xr(a, "+-a^+6")
	}
	a = xr(a, "+-3^+b+-f")
	a ^= tkkb
	if a < 0 {
		a = -a
	}
	a %= 1E6

	return a, a ^ tkka
}
func xr(a int32, b string) int32 {
	for c := 0; c < len(b)-2; c = c + 3 {
		var d int32
		d = int32(b[c+2])

		if d >= int32('a') {
			d = d - 87
		} else {
			d = d - 48
		}

		if int32('+') == int32(b[c+1]) {
			if a >= 0 {
				d = a >> uint32(d)
			} else {
				t := 4294967295 + int(a)
				t = t + 1
				d = int32(t >> uint(d))
			}
		} else {
			d = a << uint32(d)
		}

		if int32('+') == int32(b[c]) {
			a = a + d
		} else {
			a = a ^ d
		}
	}
	return a
}

//get content
func getResult(src string, sl string, tl string, dt ...string) string {
	l := "https://translate.google.cn/translate_a/single?&client=gtx"

	l = l + "&sl=" + sl
	l = l + "&tl=" + tl
	l = l + "&hl=" + tl
	for _, i := range dt {
		l = l + "&dt=" + i
	}
	l = l + "&q=" + url.QueryEscape(src)
	l = l + "&tk=" + tk(src)
	return get(l)
}
func get(l string) string {
	c := &http.Client{}
	req, err := http.NewRequest("GET", l, nil)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	ret, err := c.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	d, _ := ioutil.ReadAll(ret.Body)
	return string(d)
}
