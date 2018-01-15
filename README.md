# gtra 谷歌翻译api

golang实现的谷歌翻译api

参考并使用了 https://github.com/matheuss/google-translate-api 的部分代码

bug：有时会出现403错误网页

## 安装

```bash
go get -u -v github.com/jinzhongmin/gtra
```

## 例子

```go
package main

import (
	"fmt"

	"github.com/jinzhongmin/gtra"
)

func main() {
	t := gtra.New()
	l := gtra.NewConstLang()
	fmt.Println(t.Translate("Hello World！").To(l.ZHCN, func(e error) string {
		return "IF ERROR ,DO IN THIS FUNCTION"
	}))

}

```



## 说明

func (t *Translate) BySelf(tl string, dt []string) (gjson.Result, error)

TranslateBySelf提供自助翻译服务，dt参数自行google

## License

MIT
