package gocom

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

type DistObjReqMsg struct {
	Params []interface{}
}

type DistObjResMsg struct {
	IsErr  bool
	ErrMsg string
	Result []interface{}
}

type DistObjClient struct {
	pubSub   PubSubClient
	implList map[string]interface{}
}

func (o *DistObjClient) Invoke(className string, methodName string, param ...interface{}) ([]interface{}, error) {

	req := DistObjReqMsg{
		Params: param,
	}

	ret, err := o.pubSub.Request("distobj__"+className+"__"+methodName, req)

	if err != nil {
		return nil, err
	}

	res := DistObjResMsg{}
	err = json.Unmarshal([]byte(ret), &res)

	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (o *DistObjClient) fixType(param interface{}, paramType, argType reflect.Type) (interface{}, error) {

	// if got float64 convert to correct type
	if paramType.Kind() == reflect.Float64 {

		if argType.Kind() != reflect.Interface {
			paramReal := param.(float64)

			switch argType.Kind() {
			case reflect.Int:
				param = int(paramReal)
			case reflect.Int8:
				param = int8(paramReal)
			case reflect.Int16:
				param = int16(paramReal)
			case reflect.Int32:
				param = int32(paramReal)
			case reflect.Int64:
				param = int64(paramReal)
			case reflect.Float32:
				param = float32(paramReal)
			case reflect.Float64:
			default:
				return nil, fmt.Errorf("Want %s, got %s", argType.Name(), paramType.Name())
			}
		}
	} else if paramType.Kind() == reflect.Array {

		ret := []interface{}{}
		paramArr := param.([]interface{})

		for _, elm := range paramArr {
			retElm, err := o.fixType(elm, reflect.TypeOf(elm), paramType.Elem())

			if err != nil {

				return nil, err
			}

			ret = append(ret, retElm)
		}

		param = ret
	} else if paramType.Kind() == reflect.Map {

		ret := map[string]interface{}{}
		paramMap := param.(map[string]interface{})

		for key, elm := range paramMap {
			retElm, err := o.fixType(elm, reflect.TypeOf(elm), paramType.Elem())

			if err != nil {

				return nil, err
			}

			ret[key] = retElm
		}

		param = ret
	}

	return param, nil
}

func (o *DistObjClient) proxy(className string, impl interface{}, method reflect.Method) {

	o.pubSub.RequestSubscribe("distobj__"+className+"__"+method.Name,
		func(name, msg string) string {

			req := DistObjReqMsg{}
			res := DistObjResMsg{}

			err := json.Unmarshal([]byte(msg), &req)

			if err != nil {
				res.IsErr = true
				res.ErrMsg = err.Error()
			}

			args := make([]reflect.Value, method.Type.NumIn())
			args[0] = reflect.ValueOf(impl)

			for i, param := range req.Params {

				argType := method.Type.In(i + 1)

				param, err = o.fixType(param, reflect.TypeOf(param), argType)

				if err != nil {
					res.IsErr = true
					res.ErrMsg = fmt.Sprintf("Param %d : %s", i, err.Error())
					break
				}

				paramType := reflect.TypeOf(param)

				if !paramType.ConvertibleTo(argType) {
					res.IsErr = true
					res.ErrMsg = fmt.Sprintf("Param %d must type %s, got %s", i, argType.Name(), paramType.Name())
					break
				}

				args[i+1] = reflect.ValueOf(param)
			}

			if !res.IsErr {
				retVal := method.Func.Call(args)
				res.Result = []interface{}{}

				for _, elm := range retVal {
					res.Result = append(res.Result, elm.Interface())
				}
			}

			byteRet, _ := json.Marshal(res)
			return string(byteRet)
		})
}

func (o *DistObjClient) AddImpl(className string, impl interface{}) error {

	var err error

	implType := reflect.TypeOf(impl)

	for i := 0; i < implType.NumMethod(); i++ {

		o.proxy(className, impl, implType.Method(i))
	}

	if err == nil {

		o.implList[className] = impl
	}

	return err
}

var distObjMap map[string]*DistObjClient = map[string]*DistObjClient{}
var distObjMutex sync.Mutex

func DistObj(name ...string) *DistObjClient {

	targetName := ""

	if len(name) > 0 {
		targetName = name[0]
	}

	ret, ok := distObjMap[targetName]

	if !ok {

		distObjMutex.Lock()
		defer distObjMutex.Unlock()

		// check if prev lock already create
		ret, ok = distObjMap[targetName]

		if !ok {

			pubSub := PubSub(targetName)

			if pubSub != nil {

				ret = &DistObjClient{
					pubSub:   pubSub,
					implList: map[string]interface{}{},
				}

				distObjMap[targetName] = ret
			}
		}
	}

	return ret
}
