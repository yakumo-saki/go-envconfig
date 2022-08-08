# go-envconfig

設定ファイルと環境変数から値を取得してstructにセットします。  

## 特徴

* 設定なしでもつかえる
* 設定ファイルの位置、名前は自由
* 設定をsliceに格納することができる
* structが階層化されていても大丈夫

## インストール

```
go get https://github.com/yakumo-saki/go-envconfig
```

## Getting started

（オプション）設定ファイルを作成する
```
STR_CONF=value_from_config_file
INT_CONF=1234
```

（オプション）環境変数をセットする
```
STR_CONF=value_from_env
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

## 基本動作

* envconfig.AddPath() で追加された順番に設定ファイルを読み込みます。
* なお、指定された設定ファイルが存在しない場合は無視します
* envconfig.AddPath()されたファイルを読み終わった後に環境変数を読み込みます。
* 同一の設定が存在する場合、後に出現したものが優先されます。

## Advanced

### デフォルト値を入れたい

envconfig.LoadConfigに渡すstructに値をセットしてから渡すことで設定可能。  

```golang
func main() {
    cfg := Conf{}
    cfg.StrConf = "default value"

    err := envconfig.LoadConfig(&cfg)
    fmt.Println(cfg.StrConf) // "default value"
}
```


### 設定のキーを指定したい

デフォルトでは、フィールド名をUPPER_SNAKE_CASEにしたものが設定のキーとして使用されます。  
以下の例であれば、`STR_CONF` が使用されます。これを変更したい場合は以下のように指定します。

```golang
type Conf struct {
    StrConf string `cfg:"MY_CONFIG"`
}
```

### 連番で入力した設定をSliceで欲しい

cfgタグの二番目のオプションに `slice` を追加することで可能です。


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

cfgタグの二番目のオプションに `mergeslice` を追加することで可能です。

```golang
type Conf struct {
    SliceConfig []string `cfg:"SLICE_CONFIG,mergeslice"`
}
```


#### 実行例

```config1
SLICE_CONFIG00=abc
```

```config2
SLICE_CONFIG00=xyz
```

→ SLICE_CONFIG = []string{"abc", "xyz"}

### 指定したフィールド以外に値をセットしない

```
TEST=this is not read
```

```golang
type SimpleConfig struct {
	Test string
}

func main() {
    cfg := SimpleConfig{}
    
    // Use Strict mode.
    UseStrict() 
    
    envconfig.LoadConfig(&cfg)
    fmt.Println(cfg.Test)  // empty. not "this is not read"
}
```

## 制限

### 一度定義されたsliceを空にすることはできません

```config1
SLICE_0=abc
SLICE_1=def
```

```config2
# これで上書きしても、[]string = {""} となり空にはなりません。
SLICE_0=
```