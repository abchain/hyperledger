package utils

import (
	"strings"
)

func IsNullOrEmpty(str string) bool {
	str = strings.Trim(str, " ")
	if str == "" {
		return true
	}
	return false
}

//ConvertToAbsPath 将相对路径转化为绝对路径
// func ConvertToAbsPath(p string) string {
// 	if filepath.IsAbs(p) {
// 		return p
// 	}
// 	base := filepath.Dir(os.Args[0])
// 	return filepath.Join(base, p)
// }

// func GetReplaceAbsPath(raw, rep string) string {
// 	if IsNullOrEmpty(raw) {
// 		return ConvertToAbsPath(rep)
// 	}
// 	return ConvertToAbsPath(raw)
// }

// func ReadCert(path string) ([]byte, error) {
// 	var (
// 		b   []byte
// 		err error
// 	)
// 	if !IsNullOrEmpty(path) {
// 		path = ConvertToAbsPath(path)
// 		b, err = ioutil.ReadFile(path)
// 	}
// 	if err != nil {
// 		return b, errors.New("can not load cert file ")
// 	}
// 	return b, err
// }
