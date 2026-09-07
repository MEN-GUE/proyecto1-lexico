package afd

import "strings"

// ObtenerAlfabeto extrae los símbolos únicos de la expresión regular
func ObtenerAlfabeto(regex string) []string {
	alfabetoMap := make(map[string]bool)
	var alfabeto []string
	
	// Operadores reservados de nuestra sintaxis, incluyendo 'E' como épsilon
	operadores := "|*+?~()E"

	for _, char := range regex {
		simbolo := string(char)
		// Ignoramos operadores y espacios
		if !strings.Contains(operadores, simbolo) && simbolo != " " {
			if !alfabetoMap[simbolo] {
				alfabetoMap[simbolo] = true
				alfabeto = append(alfabeto, simbolo)
			}
		}
	}
	return alfabeto
}