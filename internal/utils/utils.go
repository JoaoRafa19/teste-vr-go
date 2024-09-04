package utils

import "reflect"

// Contains verifica se o valor está presente no array ou slice e se ambos têm o mesmo tipo
func Contains(arr any, value any) bool {
	// Verifica se arr é um slice ou array
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Slice && arrValue.Kind() != reflect.Array {
		return false
	}

	// Verifica se o array/slice está vazio
	if arrValue.Len() == 0 {
		return false
	}

	// Verifica se o tipo do primeiro elemento do array/slice é compatível com o valor
	if reflect.TypeOf(arrValue.Index(0).Interface()) != reflect.TypeOf(value) {
		return false
	}

	// Percorre o array/slice e compara os valores
	for i := 0; i < arrValue.Len(); i++ {
		if reflect.DeepEqual(arrValue.Index(i).Interface(), value) {
			return true
		}
	}

	return false
}
