package golib

import (
	"errors"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func ParseExcel2(sheet *xlsx.Sheet, headerLine int, obj interface{}) error {
	indecColumnMap := make(map[int]string)
	arrobj := reflect.ValueOf(obj).Elem()

	t := reflect.TypeOf(obj).Elem().Elem().Elem()
	// fmt.Println("Type:", t.Name())
	// fmt.Println("Kind:", t.Kind())
	m := make(map[string]string)
	musthave := make(map[string]bool)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i) //获取结构体的每一个字段
		tag := field.Tag.Get("h")
		if tag == "" {
			tag = field.Name
		}
		must := field.Tag.Get("must")
		if must == "true" {
			musthave[tag] = true
		}
		// fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		m[tag] = field.Name
	}

	for j, row := range sheet.Rows {

		if j < headerLine {
			continue
		}
		if j == headerLine {
			tempmap := make(map[string]bool)
			for i, cell := range row.Cells {

				header := cell.String()
				// log.Println(header)
				if v, has := m[header]; has {
					indecColumnMap[i] = v
					musthave[header] = false
					if tempmap[header] {
						log.Println("Sheet:" + sheet.Name + " Header name:" + header + " duplicated.")
					}
					tempmap[header] = true
				}

			}
			// log.Println(indecColumnMap)
			// spew.Dump(m)
			for k, v := range musthave {
				if v {
					log.Println("Mandatory header", k, "not found")
					return errors.New("must have")
				}
			}

		} else {
			if len(row.Cells) < 0 {
				continue
			}
			// log.Println(indecColumnMap)
			newobj := reflect.New(t)
			for index, fieldName := range indecColumnMap {
				if index >= len(row.Cells) {
					continue
				}
				field := newobj.Elem().FieldByName(fieldName)
				switch field.Kind() {
				case reflect.Int:
					tempint := strings.TrimSpace(row.Cells[index].String())
					if tempint == "" {
						field.SetInt(0)
					} else {
						i, err := strconv.Atoi(tempint)
						if err != nil {
							log.Println("line："+strconv.Itoa(j+1)+" column:"+strconv.Itoa(index+1)+" int error:"+tempint, err)
						}
						field.SetInt(int64(i))
					}
				case reflect.String:
					field.SetString(row.Cells[index].String())
				default:
					// log.Println(field.())
					vv, _ := new(big.Rat).SetString(row.Cells[index].String())
					// field.Set(vv)
					// log.Println(vv)
					field.Set(reflect.ValueOf(vv))
				}

			}
			arrobj.Set(reflect.Append(arrobj, (newobj)))

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
func ParseExcel(filename, sheetname string, headerLine int, obj interface{}) error {

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// m := make(map[string]bool)
	for _, sheet := range xlFile.Sheets {
		if sheet.Name != sheetname {
			// log.Println(sheet.Name)
			continue
		}
		return ParseExcel2(sheet, headerLine, obj)
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

func ParseExcel3(filename string, sheetId int, headerLine int, obj interface{}) error {

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// m := make(map[string]bool)
	for i, sheet := range xlFile.Sheets {
		if i != sheetId {
			// log.Println(sheet.Name)
			continue
		}
		return ParseExcel2(sheet, headerLine, obj)
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
