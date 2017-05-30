# gtra 谷歌翻译api

golang实现的谷歌翻译api

旧版本(tra.go_old)参考并使用了 https://github.com/matheuss/google-translate-api 的部分代码

旧版本中使用了javascript代码，代码冗余

新版本(gtra.go) 参考了 https://github.com/matheuss/google-translate-api 的代码，精简了代码，全部换成了go

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
	"github.com/jinzhongmin/gtra/lang"
)

func main() {
	t := gtra.NewTra("wide")
	s, err := t.To(lang.ZHCN)
	if err != nil {
	} else {
		fmt.Println(s)
	}
}

```



## 说明
func NewTra(str string) *Translate

func (t *Translate) To(tl string) (string, error)

翻译到指定语言,有时服务器拒绝就会有错误

func (t *Translate) SetSl(sl string) *Translate

指定要翻译的文本的语言，可以这样用
```go
	s , _ := t.SetSl(lang.EN).To(lang.ZHCN)
	fmt.Println(s)
```
 
## License

MIT © Matheus Fernandes
