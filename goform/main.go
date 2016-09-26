//generate struct from form data
//such as ie=utf-8&kw=%E7%94%9F%E4%B8%AA%E5%A5%B3%E5%AD%A9&fid=820625&tid=4797520575&floor_num=24&quote_id=98339446889&rich_text=1&tbs=a71f7e74bb9fd3621474863536&content=tako~&lp_type=0&lp_sub_type=0&new_vcode=1&tag=11&repostid=98339446889&anonymous=0&vcode=&vcode_md5=
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strings"
)

var data = `ie=utf-8&kw=%E7%94%9F%E4%B8%AA%E5%A5%B3%E5%AD%A9&fid=820625&tid=4797520575&floor_num=24&quote_id=98339446889&rich_text=1&tbs=a71f7e74bb9fd3621474863536&content=tako~&lp_type=0&lp_sub_type=0&new_vcode=1&tag=11&repostid=98339446889&anonymous=0&vcode=&vcode_md5=`

var file = flag.String("i", "", "input file")

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	data = string(b)
	values, err := url.ParseQuery(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(`type Foo struct{`)
	for k, v := range values {
		str := ""
		if len(v) > 0 {
			str = v[0]
		}
		fmt.Println(line(k, str))
	}
	fmt.Println(`}`)
}

var numRegexp = regexp.MustCompile(`[0-9]+`)

const fmtStr = "    %s %s `form:\"%s\"`"

func line(k, v string) string {
	typ := "string"
	if numRegexp.MatchString(v) {
		typ = "int"
	}
	return fmt.Sprintf(fmtStr, goStyleKey(k), typ, k)
}

func goStyleKey(str string) string {
	result := ""
	flag := true
	for i := 0; i < len(str); i++ {
		v := string(str[i])
		if flag {
			result += strings.ToUpper(v)
			flag = false
			continue
		}
		if v == "_" {
			flag = true
			continue
		}
		result += v
	}
	return result
}
