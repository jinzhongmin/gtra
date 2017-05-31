package gtra

import (
	"errors"
	"github.com/tidwall/gjson"
	"github.com/jinzhongmin/gtra/lang"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Translate struct {
	str string

	sl string
	tl string
	dt []string
}

func NewTra(str string) *Translate {
	t := new(Translate)
	t.str = str

	t.sl = lang.AUTO
	t.tl = lang.ZHCN
	t.dt = []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"}

	return t

}
func (t *Translate) tk() string {
	tkki, tkkf := gettkk()
	tki, tkf := gettk(int32(tkki), int32(tkkf), t.str)

	return (strconv.Itoa(int(tki)) + "." + strconv.Itoa(int(tkf)))
}
func (t *Translate) url(tl string) string {
	urls := "https://translate.google.cn/translate_a/single?"

	urls = urls + "&client=gtx"
	urls = urls + "&sl=" + t.sl
	urls = urls + "&tl=" + tl
	urls = urls + "&hl=" + tl
	for i := range t.dt {
		urls = urls + "&dt=" + t.dt[i]
	}
	//urls = urls + "&ie=UTF-8"
	//urls = urls + "&oe=UTF-8"
	//urls = urls + "&otf=1"
	//urls = urls + "&ssel=0"
	//urls = urls + "&tsel=0"
	//urls = urls + "&kc=7"
	urls = urls + "&q=" + url.QueryEscape(t.str)
	urls = urls + "&tk=" + t.tk()

	return urls
}
func (t *Translate) To(tl string) (string, error) {
	r := get(t.url(tl))

	var e error
	if r[0] != '[' {
		log.Panicln("server deny!")
		e = errors.New("DENY")
	}

	j, _ := toJson(r, 0)

	j = strings.Replace(strings.Replace(j, "]", "}", -1), "[", "{", -1)
	p := gjson.Parse(j)

	return p.Get("0.0.0").String(), e
}
func (t *Translate) SetSl(sl string) *Translate {
	t.sl = sl

	return t
}

/*func Translate(str string, from string, to string, dt []string) string {
	toJson(gtra(str, from, to, dt))
}*/

/**************************\
 * Base for get remote src*
\**************************/
func get(url string) string {

	c := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	ret, err := c.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	return string(dataLoad(ret.Body))
}
func dataLoad(body io.ReadCloser) []byte {

	data := make([]byte, 0)

	l := 0
	_, err := body.Read(data)

	if err == io.EOF {
		data = data[:l]
	} else if err != nil {
		log.Panicln(err)
	} else {
		for err == nil {

			_data := make([]byte, 100)
			l, err = body.Read(_data)

			tmp := make([]byte, len(data)+l)
			tmp = append(data, _data...)
			data = tmp
		}
	}

	return data
}

/**************************\
 *       translate        *
\**************************/
/*func gtra(tra string, from string, to string, dt []string) string {

	tkka, tkkb := gettkk()
	tka, tkb := gettk(int32(tkka), int32(tkkb), tra)

	_url := "https://translate.google.cn/translate_a/single?client=t&ie=UTF-8&oe=UTF-8&otf=1&ssel=0&tsel=0&kc=7"
	_url += ("&q=" + url.QueryEscape(tra))
	_url += ("&sl=" + from + "&tl=" + to + "&hl=" + to)

	for i := range dt {
		_url = _url + "&dt=" + dt[i]
	}

	_url += ("&tk=" + strconv.FormatInt(int64(tka), 10) + "." + strconv.FormatInt(int64(tkb), 10))

	return get(_url)
}*/

/**************************\
 * 	   calculate tk       *
\**************************/
func gettkk() (int64, int64) {
	s := get("https://translate.google.cn/")
	reg := regexp.MustCompile(`TKK=.*\)\(\)\)'\)`)
	tkkc := reg.FindAllString(s, -1)

	reg = regexp.MustCompile(`\-*\d{3,}`)
	tkkv := reg.FindAllString(tkkc[0], -1)

	a, _ := strconv.ParseInt(tkkv[0], 10, 64)
	b, _ := strconv.ParseInt(tkkv[1], 10, 64)
	c, _ := strconv.ParseInt(tkkv[2], 10, 64)

	return c, a + b
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

/**************************\
 * 	        parse         *
\**************************/
func toJson(str string, i int) (string, int) {
	_str := ""

	if str[i] == '[' {
		i++
	}
	n := 0
	l := i

	for ; i < len(str); i++ {
		if str[i] == ',' {
			_str = _str + "\"" + strconv.Itoa(n) + "\":" + str[l:i+1]
			n++
			l = i + 1
		} else if str[i] == ']' {
			_str = _str + "\"" + strconv.Itoa(n) + "\":" + str[l:i+1] + ","
			break
		} else if str[i] == '[' {
			_s, _i := toJson(str, i)

			i = _i + 1
			l = i + 1

			_str = _str + "\"" + strconv.Itoa(n) + "\":" + _s
			n++
		}
	}
	_str = "[" + _str

	return _str, i
}

/*func toJson(list string) string {
	var l string
	for i := 0; i < len(list); i++ {
		if list[i] == '[' {
			l = list[0:i+1] + "\"0\":" + list[i+1:]
			list = l
		}
	}
	n := 0
	for i := 0; i < len(list); i++ {

		if list[i] == ',' {
			n++
			l = list[0:i+1] + "\"" + strconv.Itoa(n) + "\":" + list[i+1:]
			list = l
		}

		if i < len(list)-1 && list[i+1] == '[' {
			n = 0
		}
	}
	return strings.Replace(strings.Replace(l, "]", "}", -1), "[", "{", -1)

}*/
