package parser

import (
	"log"
	"strconv"
	"strings"
)

// ParseArgs parses arguments like key:value key:key:value key:val,val,val
func ParseArgs(args []string) interface{} {

	if len(args) == 0 {
		return nil
	}

	if len(args) == 1 {
		return parseItem(args[0])
	}

	data := map[string]interface{}{}

	for _, item := range args {

		result := parseItem(item)

		stringResult, ok := result.(string)
		if ok {
			data[stringResult] = true
			continue
		}

		newData, ok := result.(map[string]interface{})
		if !ok {
			log.Fatalf("item not like \"name:...\": %s", item)
		}

		data = mergeMap(data, newData)
	}

	return data
}

func mergeMap(one, two map[string]interface{}) map[string]interface{} {

	for k, v := range two {
		one[k] = merge(one[k], v)
	}

	return one
}

func merge(one, two interface{}) interface{} {

	if one == nil {
		return two
	}

	if two == nil {
		return one
	}

	oneString, oneStringOK := one.(string)
	twoString, twoStringOK := two.(string)

	if oneStringOK && twoStringOK {
		return []interface{}{
			oneString,
			twoString,
		}
	}

	oneSlice, oneSliceOK := one.([]interface{})
	twoSlice, twoSliceOK := two.([]interface{})

	if oneSliceOK {
		if twoSliceOK {
			return append(oneSlice, twoSlice...)
		}
		return append(oneSlice, two)
	}

	if twoSliceOK {
		newSlice := []interface{}{one}
		return append(newSlice, twoSlice...)
	}

	oneMap, oneMapOK := one.(map[string]interface{})
	twoMap, twoMapOK := two.(map[string]interface{})

	if oneMapOK || twoMapOK {
		if !oneMapOK || !twoMapOK {
			log.Fatalf("cannot merge map with non-map: %s, %s", one, two)
		}

		newMap := map[string]interface{}{}
		for k1, v1 := range oneMap {
			v2 := twoMap[k1]
			newMap[k1] = merge(v1, v2)
		}

		return newMap
	}

	return []interface{}{
		one,
		two,
	}
}

func parseItem(item string) interface{} {

	p := strings.Index(item, ":")
	if p != -1 {
		name := item[0:p]
		other := item[p+1:]
		var value interface{}
		if other == "" {
			value = true
		} else {
			value = parseItem(other)
		}

		m := map[string]interface{}{}
		m[name] = value

		return m
	}

	p = strings.Index(item, ",")
	if p != -1 {
		items := strings.Split(item, ",")
		newItems := []interface{}{}
		for _, item := range items {
			newItems = append(newItems, parseItem(item))
		}
		return newItems
	}

	i64, err := strconv.ParseInt(item, 10, 0)
	if err == nil {
		return i64
	}

	f64, err := strconv.ParseFloat(item, 0)
	if err == nil {
		return f64
	}

	if strings.EqualFold(item, "true") {
		return true
	}

	if strings.EqualFold(item, "false") {
		return false
	}

	return item
}
