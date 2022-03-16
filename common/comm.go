package common

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// DataToStructByTagSql 根据结构体中sql标签映射到结构体中并且转换类型
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()    //利用反射获得obj的值
	for i := 0; i < objValue.NumField(); i++ { //遍历该对象
		value := data[objValue.Type().Field(i).Tag.Get("sql")] //获取sql对应的值
		name := objValue.Type().Field(i).Name                  //获取对应字段的名称
		structFieldType := objValue.Field(i).Type()            //获取字段的类型
		val := reflect.ValueOf(value)                          //获取变量类型,也可以直接写string类型
		var err error
		if structFieldType != val.Type() {
			//类型转换
			val, err = TypeConversion(value, structFieldType.Name())
			if err != nil {
				fmt.Println("maybe there are type couldn't transfer,", err)
			}
		}
		//设置类型值
		objValue.FieldByName(name).Set(val)
	}
}

// TypeConversion 类型转换
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
