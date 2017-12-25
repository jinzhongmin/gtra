package gtra

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

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

	d, _ := ioutil.ReadAll(ret.Body)
	return string(d)
}
