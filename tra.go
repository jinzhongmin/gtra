package gtra

import (
	"log"
	"regexp"

	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

func getTKK() otto.Value {
	doc, err := goquery.NewDocument("https://translate.google.cn/")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(doc.Text())
	reg := regexp.MustCompile(`TKK=.*\)\(\)\)'\)`)
	TKK := reg.FindAllString(doc.Text(), -1)[0]

	vm := otto.New()
	vm.Set("TKK", TKK)
	vm.Run(TKK)
	TKKV, err := vm.Get("TKK")
	if err != nil {
		log.Fatal(err)
	}

	return TKKV
}
func gettk(str string) string {
	TKK := getTKK()

	tkjs := `
function sM(a) {
    var b;
    if (null !== yr)
        b = yr;
    else {
        b = wr(String.fromCharCode(84));
        var c = wr(String.fromCharCode(75));
        b = [b(), b()];
        b[1] = c();
        b = (yr = window[b.join(c())] || "") || ""
    }
    var d = wr(String.fromCharCode(116))
        , c = wr(String.fromCharCode(107))
        , d = [d(), d()];
    d[1] = c();
    c = "&" + d.join("") + "=";
    d = b.toString().split(".");
    b = Number(d[0]) || 0;
    for (var e = [], f = 0, g = 0; g < a.length; g++) {
        var l = a.charCodeAt(g);
        128 > l ? e[f++] = l : (2048 > l ? e[f++] = l >> 6 | 192 : (55296 == (l & 64512) && g + 1 < a.length && 56320 == (a.charCodeAt(g + 1) & 64512) ? (l = 65536 + ((l & 1023) << 10) + (a.charCodeAt(++g) & 1023),
            e[f++] = l >> 18 | 240,
            e[f++] = l >> 12 & 63 | 128) : e[f++] = l >> 12 | 224,
            e[f++] = l >> 6 & 63 | 128),
            e[f++] = l & 63 | 128)
    }
    a = b;
    for (f = 0; f < e.length; f++)
        a += e[f],
            a = xr(a, "+-a^+6");
    a = xr(a, "+-3^+b+-f");
    a ^= Number(d[1]) || 0;
    0 > a && (a = (a & 2147483647) + 2147483648);
    a %= 1E6;
    return c + (a.toString() + "." + (a ^ b))
}

var yr = null;
var wr = function (a) {
    return function () {
        return a
    }
}
var xr = function (a, b) {
    for (var c = 0; c < b.length - 2; c += 3) {
        var d = b.charAt(c + 2)
            , d = "a" <= d ? d.charCodeAt(0) - 87 : Number(d)
            , d = "+" == b.charAt(c + 1) ? a >>> d : a << d;
        a = "+" == b.charAt(c) ? a + d & 4294967295 : a ^ d
    }
    return a
};
var window = {
    TKK: TKKV   //TKKV
};

tk = sM(str).replace('&tk=', '')

//在golang中设置str和TKKV的值，最后获取tk的值
`
	scriptStr := tkjs

	vm := otto.New()
	vm.Set("TKKV", TKK)
	vm.Set("str", str)

	vm.Run(scriptStr)

	tk, err := vm.Get("tk")
	if err != nil {
		log.Fatal(err)
	}
	tkStr, err := tk.ToString()
	if err != nil {
		log.Fatal(err)
	}

	return tkStr
}
func getUrl(str string, from string, to string) string {
	Static := "https://translate.google.cn/translate_a/single?client=t&dt=at&dt=bd&dt=ex&dt=ld&dt=md&dt=qca&dt=rw&dt=rm&dt=ss&dt=t&ie=UTF-8&oe=UTF-8&otf=1&ssel=0&tsel=0&kc=7"
	urls := Static + "&q=" + url.QueryEscape(str)
	urls = urls + "&sl=" + from
	urls = urls + "&tl=" + to
	urls = urls + "&hl=" + to
	urls = urls + "&tk=" + gettk(str)

	return urls
}
func TranslateT(str string, from string, to string) string {
	urls := getUrl(str, from, to)
	doc, err := goquery.NewDocument(urls)
	if err != nil {
		log.Fatal(err)
	}
	res := doc.Text()

	return res
}
