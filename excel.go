package golib

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com1/tealeg/xlsx"
)

func ParseExcel(filename, sheetname string, headerLine int, obj interface{}) error {
	arrobj := reflect.ValueOf(obj).Elem()

	t := reflect.TypeOf(obj).Elem().Elem().Elem()
	// fmt.Println("Type:", t.Name())
	// fmt.Println("Kind:", t.Kind())
	m := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i) //获取结构体的每一个字段
		tag := field.Tag.Get("h")
		if tag == "" {
			tag = field.Name
		}
		// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		m[tag] = field.Name
	}

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// m := make(map[string]bool)
	indecColumnMap := make(map[int]string)
	for _, sheet := range xlFile.Sheets {
		if sheet.Name != sheetname {
			// log.Println(sheet.Name)
			continue
		}
		for j, row := range sheet.Rows {

			if j < headerLine {
				continue
			}
			if j == headerLine {
				for i, cell := range row.Cells {

					header := cell.String()
					log.Println(header)
					if v, has := m[header]; has {
						indecColumnMap[i] = v
					}

				}
				// log.Println(indecColumnMap)
				// spew.Dump(m)
			} else {
				if len(row.Cells) < 0 {
					continue
				}
				newobj := reflect.New(t)
				for index, fieldName := range indecColumnMap {
					if index >= len(row.Cells) {
						continue
					}
					field := newobj.Elem().FieldByName(fieldName)
					switch field.Kind() {
					case reflect.Int:
						tempint := strings.TrimSpace(row.Cells[index].String())
						i, err := strconv.Atoi(tempint)
						if err != nil {
							log.Println("line："+strconv.Itoa(j+1)+" column:"+strconv.Itoa(index+1)+" int error:"+tempint, err)
						}
						field.SetInt(int64(i))
					case reflect.String:
						field.SetString(row.Cells[index].String())
					}

				}
				arrobj.Set(reflect.Append(arrobj, (newobj)))

			}
		}
	}

	// newobj := reflect.New(t)
	// for _, v := range m {

	// 	field := newobj.Elem().FieldByName(v)
	// 	switch field.Kind() {
	// 	case reflect.Int:
	// 		field.SetInt(8)
	// 	case reflect.String:
	// 		field.SetString("yes")
	// 	}

	// }
	// arrobj.Set(reflect.Append(arrobj, (newobj)))
	// arr := reflect.ArrayOf(0, t)

	return nil
}
