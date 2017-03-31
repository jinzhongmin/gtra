# gtra 谷歌翻译api

golang实现的谷歌翻译api

参考并使用了https://github.com/matheuss/google-translate-api的部分代码

## 安装

```bash
go get github.com/jinzhongmin/gtra
```

## 例子

```go
package main

import (
	"fmt"
	"github.com/jinzhongmin/gtra"
)

func main() {
	str := gtra.TranslateT("hello world!", "en", "zh-cn")
	fmt.Println(str)
}

```



## 说明

返回数据没有处理，请自行处理
