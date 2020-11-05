package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

//TypeConversion 类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}
	//else if .......增加其他一些类型的转换
	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

//DataToStructByTagSql 根据tag将数据映射到结构体中
func DataToStructByTagSql(data map[string]string, obj interface{}){
	objvalue := reflect.ValueOf(obj).Elem()
	//fmt.Println(objvalue.Type().Field(3).Tag.Get("sql"))
	for i:=0 ; i<objvalue.NumField(); i++{
		value:= data[objvalue.Type().Field(i).Tag.Get("sql")]
		name:= objvalue.Type().Field(i).Name
		structFieldType := objvalue.Field(i).Type()
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type(){
			val , err = TypeConversion(value , structFieldType.Name())
			if err != nil {

			}
		}
		objvalue.FieldByName(name).Set(val)
	}
}