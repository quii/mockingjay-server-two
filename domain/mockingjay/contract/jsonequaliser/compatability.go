package jsonequaliser

import (
	"fmt"
	"reflect"
)

type jsonNode map[string]interface{}

var (
	msgFieldMissing       = "Missing field"
	msgNotString          = "Field is not a string in other JSON"
	msgNotBool            = "Field is not a boolean in other JSON"
	msgNotFloat           = "Field is not a float in other JSON"
	msgEmptyArray         = "Array in other JSON is empty so I cant check"
	msgNotMap             = "Not a map in other JSON"
	msgDifferentArrayType = "Type of array is different"
	msgEmptyRootArray     = "Empty arrays are not suitable for comparison"
)

func emptyJSONArrayHandler(parseError error) (errorMessages map[string]string, err error) {
	if _, ok := parseError.(*emptyJSONArrayError); ok {
		errorMessages = make(map[string]string)
		errorMessages["rootArray"] = msgEmptyRootArray
		return
	}

	return errorMessages, parseError
}

func IsCompatible(a, b string) (errorMessages map[string]string, err error) {
	aMap, err := getJSONNodeFromString(a)

	if err != nil {
		return emptyJSONArrayHandler(err)
	}

	bMap, err := getJSONNodeFromString(b)

	if err != nil {
		return emptyJSONArrayHandler(err)
	}

	return isStructurallyTheSame(aMap, bMap, make(map[string]string), "")
}

func isStructurallyTheSame(a, b jsonNode, messages map[string]string, baseNode string) (map[string]string, error) {
	for jsonFieldName, v := range a {
		messageNodeName := jsonFieldName
		if baseNode != "" {
			messageNodeName = baseNode + "->" + jsonFieldName
		}

		if fieldMissingIn(b, jsonFieldName) {
			messages[messageNodeName] = msgFieldMissing
			continue
		}

		if a[jsonFieldName] == nil {
			continue
		}

		switch v.(type) {
		case string:
			if !isString(b, jsonFieldName) {
				messages[messageNodeName] = msgNotString
				continue
			}
		case bool:
			if !isBool(b, jsonFieldName) {
				messages[messageNodeName] = msgNotBool
				continue
			}
		case float64:
			if !isFloat(b, jsonFieldName) {
				messages[messageNodeName] = msgNotFloat
				continue
			}

		case interface{}:

			aArr, aIsArray := a[jsonFieldName].([]interface{})

			bArr, bIsArray := b[jsonFieldName].([]interface{})

			if aIsArray && len(aArr) == 0 {
				continue
			}

			if !bIsArray && aIsArray || aIsArray && len(bArr) == 0 {
				messages[messageNodeName] = msgEmptyArray
				continue
			}

			var aLeaf, bLeaf jsonNode
			var aIsMap, bIsMap bool

			if aIsArray && bIsArray {
				aLeaf, aIsMap = aArr[0].(map[string]interface{})
				bLeaf, bIsMap = bArr[0].(map[string]interface{})
			} else {
				aLeaf, aIsMap = a[jsonFieldName].(map[string]interface{})
				bLeaf, bIsMap = b[jsonFieldName].(map[string]interface{})
			}

			if aIsMap && bIsMap {
				messages, err := isStructurallyTheSame(aLeaf, bLeaf, messages, messageNodeName)
				if err != nil {
					return messages, err
				}
				continue
			} else if aIsMap && !bIsMap {
				messages[messageNodeName] = msgNotMap
				continue
			} else if reflect.TypeOf(aArr[0]) != reflect.TypeOf(bArr[0]) {
				messages[messageNodeName] = msgDifferentArrayType
				continue
			}
		default:
			return messages, fmt.Errorf("Unmatched type of json found, got a %v", reflect.TypeOf(v))
		}
	}

	return messages, nil
}
