# gtra 谷歌翻译api

golang实现的谷歌翻译api

参考并使用了 https://github.com/matheuss/google-translate-api 的部分代码

bug：有时会出现403错误网页

## 安装

```bash
go get -u -v github.com/jinzhongmin/gtra
```

## 例子

```golang
package main

import (
	"fmt"

	"github.com/jinzhongmin/gtra"
	"github.com/jinzhongmin/gtra/lang"
)

func main() {
	t := gtra.NewTranslater()
	fmt.Println(t.Translate("hello"))
	fmt.Println(gtra.Translate("world", lang.ZHCN))

	fmt.Println(t.Vector(lang.ZHCN, lang.EN).Translate("你好世界"))
	fmt.Println(t.To(lang.JA).Translate("你好世界"))

	_, j := t.Vector(lang.EN, lang.ZHCN).Dt("like", "t", "at")
	fmt.Println(j.String())
}


```



## api

[api](https://godoc.org/github.com/jinzhongmin/gtra)

func (t *Translater) Dt(src string, dt ...string) (error, gjson.Result)

dt参数请参考 [&](https://stackoverflow.com/questions/26714426/what-is-the-meaning-of-google-translate-query-params#answers)

## License

MIT © [jinzhongmin] (https://github.com/jinzhongmin)
