package tg

import (
	"reflect"
	"strings"
)

// redactToken walks val (mirroring defaultsInternal's traversal shape) and
// rewrites every string it finds that contains token, replacing it with
// redactedTokenValue(token). It's called on the request struct right before
// encoding, so a leaked token (e.g. from a careless OnError handler echoing
// err.Error() back into a message) never reaches the wire.
func redactToken(val any, token string) {
	if token == "" {
		return
	}
	redactTokenInternal(reflect.ValueOf(val), token)
}

func redactTokenInternal(reflectVal reflect.Value, token string) {
	for reflectVal.Kind() == reflect.Ptr {
		if reflectVal.IsNil() {
			return
		}
		reflectVal = reflectVal.Elem()
	}

	switch reflectVal.Kind() {
	case reflect.String:
		if s := reflectVal.String(); strings.Contains(s, token) && reflectVal.CanSet() {
			reflectVal.SetString(strings.ReplaceAll(s, token, redactedTokenValue(token)))
		}
	case reflect.Struct:
		for i := 0; i < reflectVal.NumField(); i++ {
			field := reflectVal.Field(i)
			switch field.Kind() {
			case reflect.String, reflect.Struct:
				if field.CanAddr() {
					redactTokenInternal(field.Addr(), token)
				}
			case reflect.Array, reflect.Slice, reflect.Ptr, reflect.Interface:
				redactTokenInternal(field, token)
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < reflectVal.Len(); i++ {
			redactTokenInternal(reflectVal.Index(i), token)
		}
	case reflect.Interface:
		if reflectVal.IsNil() {
			return
		}
		elem := reflectVal.Elem()
		// A pointer/struct held in the interface (the common case: every RichText/
		// RichBlock/... variant implements its interface via a pointer) is reached
		// through real, separately-addressable memory, so recursing into it can
		// mutate it in place. A bare scalar held directly in an interface (e.g. a
		// RichText field holding a plain string, since RichText is Go `any`) is
		// NOT addressable through .Elem() - the only way to rewrite it is to
		// reassign the whole interface field.
		if elem.Kind() == reflect.String {
			if s := elem.String(); strings.Contains(s, token) && reflectVal.CanSet() {
				reflectVal.Set(reflect.ValueOf(strings.ReplaceAll(s, token, redactedTokenValue(token))))
			}
			return
		}
		redactTokenInternal(elem, token)
	}
}

// redactedTokenValue keeps the bot id, replaces the secret 1:1 by length, e.g. "1234:asdf" -> "1234:****".
func redactedTokenValue(token string) string {
	if id, secret, ok := strings.Cut(token, ":"); ok {
		return id + ":" + strings.Repeat("*", len(secret))
	}
	return strings.Repeat("*", len(token))
}
