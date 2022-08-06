package envfile

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yakumo-saki/go-envconfig/internal/util"
)

// ENVファイルを読み込んでmapにして返す。
// この時点では単純にkey, value ともにstring
func ReadEnvfile(path string) (map[string]string, error) {
	ret := make(map[string]string)
	filebyte, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}

	lines := strings.Split(string(filebyte), "\n")
	for _, lineStr := range lines {
		if len(lineStr) == 0 {
			continue
		}

		line := util.StringTrim(lineStr)
		if util.IsCommentLine(line) {
			continue
		}

		// fmt.Println(line)
		sep := strings.SplitN(line, "=", 2)
		if len(sep) != 2 {
			return ret, fmt.Errorf("separeted count %d != 2. line=%s", len(sep), line)
		}
		ret[sep[0]] = sep[1]
	}

	return ret, nil
}
