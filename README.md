# go-envconfig

設定ファイルと環境変数から

## インストール

```
go get https://github.com/yakumo-saki/go-envconfig
```

## Getting started

（オプション）設定ファイルを作成する
```
StrConf=value_from_config_file
IntConf=1234
```

（オプション）環境変数をセットする
```
StrConf=value_from_env
```

```golang
package main

type Conf struct {
    StrConf string
    IntConf int
}

func main() {
    envconfig.AddPath("path/to/your/config")
    envconfig.AddPath("path/to/your/another/config")
    
    cfg := Conf{}
    err := envconfig.LoadConfig(&cfg)
    if err != nil {
        panic(err)
    }

    fmt.Println(cfg.StrConf) // "value_from_config_file" or "value_from_env"
    fmt.Println(cfg.IntConf) // 1234
}
```

## Advanced usage

### デフォルト値を入れたい

envconfig.LoadConfigに渡すstructに値をセットしてから渡すことで設定可能

```golang
func main() {
    cfg := Conf{}
    cfg.StrConf = "default value"

    err := envconfig.LoadConfig(&cfg)

}
```


### キーを指定したい

```golang
type Conf struct {
    StrConf string `cfg:"MY_CONFIG"`
}
```

### Sliceで欲しい

```golang
type Conf struct {
    SliceConfig []string `cfg:"SLICE_CONFIG,slice"`
}
```

`SLICE_CONFIG0`,`SLICE_CONFIG00`,
`SLICE_CONFIG1`,`SLICE_CONFIG01`,
...`SLICE_CONFIG99` までを探してSliceにして読み込む。連番でなくてもOK。
この設定の場合は、複数のファイルに `SLICE_CONFIGnn`が一つでもある場合上書きされる。

note: 上記の例の場合、`SLICE_CONFIG` というキーは読み込まれないことに注意

#### sliceの例

```config1
SLICE_CONFIG00=abc
```

```config2
SLICE_CONFIG00=xyz
```

→ SLICE_CONFIG = []string{"xyz"}  
"abc" は上書きされる。

### 複数ファイルにあるsliceを上書きではなくマージしたい

```golang
type Conf struct {
    SliceConfig []string `cfg:"SLICE_CONFIG,mergeslice"`
}
```

#### sliceの例

```config1
SLICE_CONFIG00=abc
```

```config2
SLICE_CONFIG00=xyz
```

→ SLICE_CONFIG = []string{"abc", "xyz"}
