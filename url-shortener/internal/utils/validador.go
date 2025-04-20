package utils

import (
	"net/url"
)

// IsValidURL verifica si una cadena es una URL válida.
func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false // No se pudo parsear
	}
	// Verifica que tenga un esquema (http, https) y un host
	return u.Scheme != "" && u.Host != ""
}