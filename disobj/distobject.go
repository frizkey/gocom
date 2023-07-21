package distobj

import (
	"encoding/json"
	"time"

	"github.com/adlindo/gocom/pubsub"
)

func Invoke(prefix, className, methodName string, params ...interface{}) ([]interface{}, error) {

	targetName := className + ">>" + methodName

	retStr, err := pubsub.Get().Request(targetName, params, 5*time.Minute)

	if err != nil {
		return nil, err
	}

	ret := []interface{}{}
	err = json.Unmarshal([]byte(retStr), &ret)

	return ret, err
}
