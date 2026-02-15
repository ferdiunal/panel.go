package core

import "reflect"

// ElementCloner is an optional interface that allows field implementations
// to provide their own clone logic.
type ElementCloner interface {
	Clone() Element
}

// CloneElement returns an isolated clone for request/item scoped mutations.
// If the element does not implement ElementCloner, a reflection-based deep clone
// fallback is used.
func CloneElement(element Element) Element {
	if element == nil {
		return nil
	}

	if cloner, ok := element.(ElementCloner); ok {
		if cloned := cloner.Clone(); cloned != nil {
			return cloned
		}
	}

	value := reflect.ValueOf(element)
	clonedValue := deepCloneValue(value, make(map[uintptr]reflect.Value))
	if !clonedValue.IsValid() {
		return element
	}

	clonedElement, ok := clonedValue.Interface().(Element)
	if !ok || clonedElement == nil {
		return element
	}

	return clonedElement
}

func deepCloneValue(value reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}

	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}

		ptr := value.Pointer()
		if cached, ok := visited[ptr]; ok {
			return cached
		}

		cloned := reflect.New(value.Elem().Type())
		visited[ptr] = cloned
		cloned.Elem().Set(value.Elem())
		return deepCloneStructPointerFields(cloned, visited)

	case reflect.Interface:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}

		clonedElem := deepCloneValue(value.Elem(), visited)
		if !clonedElem.IsValid() {
			return reflect.Zero(value.Type())
		}
		out := reflect.New(value.Type()).Elem()
		if clonedElem.Type().AssignableTo(value.Type()) {
			out.Set(clonedElem)
			return out
		}
		if clonedElem.Type().Implements(value.Type()) {
			out.Set(clonedElem)
			return out
		}
		out.Set(value)
		return out

	case reflect.Struct:
		cloned := reflect.New(value.Type()).Elem()
		cloned.Set(value)
		for i := 0; i < cloned.NumField(); i++ {
			field := cloned.Field(i)
			if !field.CanSet() {
				continue
			}
			clonedField := deepCloneValue(field, visited)
			if clonedField.IsValid() && clonedField.Type().AssignableTo(field.Type()) {
				field.Set(clonedField)
			}
		}
		return cloned

	case reflect.Slice:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}

		cloned := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		for i := 0; i < value.Len(); i++ {
			elem := deepCloneValue(value.Index(i), visited)
			if elem.IsValid() && elem.Type().AssignableTo(value.Type().Elem()) {
				cloned.Index(i).Set(elem)
			} else {
				cloned.Index(i).Set(value.Index(i))
			}
		}
		return cloned

	case reflect.Array:
		cloned := reflect.New(value.Type()).Elem()
		for i := 0; i < value.Len(); i++ {
			elem := deepCloneValue(value.Index(i), visited)
			if elem.IsValid() && elem.Type().AssignableTo(value.Type().Elem()) {
				cloned.Index(i).Set(elem)
			} else {
				cloned.Index(i).Set(value.Index(i))
			}
		}
		return cloned

	case reflect.Map:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}

		cloned := reflect.MakeMapWithSize(value.Type(), value.Len())
		iter := value.MapRange()
		for iter.Next() {
			key := deepCloneValue(iter.Key(), visited)
			val := deepCloneValue(iter.Value(), visited)
			if !key.IsValid() || !val.IsValid() {
				continue
			}

			if !key.Type().AssignableTo(value.Type().Key()) {
				key = iter.Key()
			}
			if !val.Type().AssignableTo(value.Type().Elem()) {
				val = iter.Value()
			}
			cloned.SetMapIndex(key, val)
		}
		return cloned

	default:
		return value
	}
}

func deepCloneStructPointerFields(value reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if !value.IsValid() || value.Kind() != reflect.Pointer || value.IsNil() {
		return value
	}

	elem := value.Elem()
	if elem.Kind() != reflect.Struct {
		return value
	}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if !field.CanSet() {
			continue
		}

		clonedField := deepCloneValue(field, visited)
		if clonedField.IsValid() && clonedField.Type().AssignableTo(field.Type()) {
			field.Set(clonedField)
		}
	}

	return value
}
