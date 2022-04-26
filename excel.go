package golib

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com1/tealeg/xlsx"
)

func ParseExcel(filename, sheetname string, headerLine int, obj interface{}) error {
	arrobj := reflect.ValueOf(obj).Elem()

	t := reflect.TypeOf(obj).Elem().Elem().Elem()
	fmt.Println("Type:", t.Name())
	fmt.Println("Kind:", t.Kind())
	m := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i) //获取结构体的每一个字段
		tag := field.Tag.Get("h")
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

					header := strings.TrimSpace(cell.String())
					if v, has := m[header]; has {
						indecColumnMap[i] = v
					}

				}
			} else {
				newobj := reflect.New(t)
				for index, fieldName := range indecColumnMap {

					field := newobj.Elem().FieldByName(fieldName)
					switch field.Kind() {
					case reflect.Int:
						tempint := row.Cells[index].String()
						i, err := strconv.Atoi(tempint)
						if err != nil {
							log.Println("line：" + strconv.Itoa(i+1) + " column:" + strconv.Itoa(index+1) + " int error:" + tempint)
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
