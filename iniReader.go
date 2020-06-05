package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type MysqlConfig struct {
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type RedisConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
}
type Config struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}

func LoadIni(filename string, data interface{}) {
	if !isDataOK(data) {
		return
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("檔案讀取問題", err)
		return
	}
	// fmt.Println(string(file))
	//切割文件
	strs := cutIni(file)
	var section string
	var configField string
	var key string
	var val string
	for _, str := range strs {
		//跳過注釋
		if IsAnnotate(str) {
			continue
		}
		//如果是[mysql]或[redis]這樣的字段 就當作節
		if IsSection(str) {
			section = string(Section(str)[1])
			//找出configField
			configField = SetConfigField(data, section)
			fmt.Println("找到configField:", configField)
		} else {
			key, val = KeyAndValue(str)
			if key != "" {
				// fmt.Println(key, val)

				structVal := reflect.ValueOf(data).Elem().FieldByName(configField)
				sType := structVal.Type()
				for i := 0; i < sType.NumField(); i++ {
					field := sType.Field(i).Tag.Get("ini")
					if field == key {
						switch structVal.Field(i).Type().Kind() {
						case reflect.String:
							structVal.Field(i).SetString(val)
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							valint, _ := strconv.ParseInt(val, 10, 64)
							structVal.Field(i).SetInt(valint)
						}
					}
				}
			}
		}
	}
}
func KeyAndValue(str string) (key, val string) {
	temp := strings.Split(str, "=")
	if len(temp) != 2 {
		fmt.Println(str, "格式錯誤")
		return
	}
	if len(temp[0]) == 0 || len(temp[1]) == 0 {
		fmt.Println(str, "有值為空")
		return
	}
	key = strings.TrimSpace(temp[0])
	val = strings.TrimSpace(temp[1])
	return
}

//SetConfigField 判斷configField是什麼
func SetConfigField(data interface{}, section string) string {
	confObj := reflect.TypeOf(data).Elem()
	for i := 0; i < confObj.NumField(); i++ {
		tag := confObj.Field(i).Tag.Get("ini")
		if tag == section {
			//取得Config struct內含字段的名稱
			return confObj.Field(i).Name
		}
	}
	return ""
}

//判斷data為struct且為指針
func isDataOK(data interface{}) bool {
	obj := reflect.TypeOf(data)
	if obj.Kind() == reflect.Ptr && obj.Elem().Kind() == reflect.Struct {
		return true
	}
	if obj.Kind() != reflect.Ptr {
		fmt.Println("data is not pointer")
	}

	if obj.Elem().Kind() != reflect.Struct {
		fmt.Println("data is not struct")
	}
	return false
}

func cutIni(file []byte) []string {
	strs := strings.Split(string(file), "\r\n")

	for i, str := range strs {
		strs[i] = strings.TrimSpace(str)
	}
	return strs
}

func IsAnnotate(str string) bool {
	return strings.HasPrefix(str, "#") || strings.HasPrefix(str, ";")
}

func IsSection(str string) bool {
	if m, _ := regexp.MatchString(`\[\s*(\S+)\s*\]`, str); m {
		return true
	}
	return false
}
func Section(str string) (s [][]byte) {
	b := []byte(str)
	// fmt.Println(str)
	re := regexp.MustCompile(`\[\s*(\S+)\s*\]`)
	s = re.FindSubmatch(b)
	return s
}

func main() {
	var cfg Config
	LoadIni("MyConfig.ini", &cfg)
	fmt.Println(cfg)
}
