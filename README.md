# go-envconfig

設定ファイルと環境変数から値を取得してstructにセットします。  

## branching

main branch is newest release.
edge branch for development.
semver tag is available.

## 特徴

* 設定なしでもつかえる
* 設定ファイルの位置、名前は自由
* 設定をsliceに格納することができる
* structが階層化されていても大丈夫

## インストール

```
go get -u https://github.com/yakumo-saki/go-envconfig
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
    ec := envconfig.New()
    ec.AddPath("path/to/your/config")        
    ec.AddPath("path/to/your/another/config")  // if not found, thats ok. simply ignored
    
    cfg := Conf{}
    err := ec.LoadConfig(&cfg)
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

### TAG

細かい動作を指定したい場合は、受け側のstructにtagを記述することで制御可能。
通常の変数、slice、mapの3パターンで動作が変わる。
以下、環境変数もしくは、設定ファイルを合わせて設定と呼ぶ。

```golang
type Conf struct {
    StrConf string `cfg:"ENV_OR_CFGKEY_NAME,option"`
    IntConf int
}
```

#### 通常の変数時

##### ENV_OR_CFGKEY_NAME

設定のキー。指定された名前の設定の値がセットされる。

##### option

なし。指定するとエラーになる。

#### Slice時

##### ENV_OR_CFGKEY_NAME

設定のキーのプレフィックス。
例えば、 `MY_CFG` を指定すると、`MY_CFG_0` `MY_CFG_1` ... `MY_CFG_99` の値がsliceにセットされる。

##### `merge` (default)

環境変数と設定ファイルにプレフィックスに該当する設定があった場合、内容をマージする。

```
環境変数 
MY_CFG_0=env_0

ファイル CFG_FILE
MY_CFG_0=file_0

結果
MY_CFG=["env_0", "file_0"]
```

##### `overwrite` 

環境変数と設定ファイルにプレフィックスに該当する設定があった場合、最後に見つけた内容のみをsliceにセットする。

```
ファイル CFG_FILE
MY_CFG_0=file_0

環境変数 
MY_CFG_0=env_0

結果
MY_CFG=["env_0", "file_0"]
```

#### map (value is not slice) 

eg) map[key] is string, int, float...


##### `keymerge` (default)

##### `overwrite`

##### `valuemerge` 

causes panic. (v0.3.x) 
値をマージすることができないため。

#### map (value is slice) 

eg) map[key] is []string, []int ...

実行結果は以下の入力があったときのものである。

```
file1.env
MYMAP_ALICE_1=alice1_file1
MYMAP_ALICE_2=alice2_file1
MYMAP_BOB_1=bob_file1

file2.env
MYMAP_ALICE_1=alice1_fileTWO
```


##### `keymerge` (default)

同一マップが複数箇所で定義された場合、キー単位で最後に定義された箇所のみが有効になる。

```
MYMAP["ALICE"] = []string{"alice1_fileTWO"}
MYMAP["BOB"] = []string{"bob_file1"}
```

##### `valuemerge`

同一マップが複数箇所で定義された場合、キー単位で値をマージする 


```
MYMAP["ALICE"] = []string{"alice1_file1","alice2_file1","alice1_fileTWO"}
MYMAP["BOB"] = []string{"bob_file1"}
```


##### `overwrite`

同一マップが複数箇所で定義された場合、最後に定義された箇所のみのマップを生成する

```
MYMAP["ALICE"] = []string{"alice1_fileTWO"}
```


## Examples

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
    
    ec := envconfig.New()
    ec.LoadConfig(&cfg)
    fmt.Println(cfg.Test)  // empty. NOT "this is not read"
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
