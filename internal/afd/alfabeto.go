package afd

import "strings"

// ObtenerAlfabeto extrae de forma segura los símbolos únicos de la expresión regular
func ObtenerAlfabeto(regex string) []string {
	alfabetoMap := make(map[string]bool)
	var alfabeto []string
	
	// El símbolo '#' se reserva a nivel de sistema para representar transiciones épsilon[cite: 7].
	// Los operadores estructurales de la notación postfix e infix quedan excluidos del mapa.
	operadores := "|*+?~()#"

	for _, char := range regex {
		simbolo := string(char)
		// La validación ignora tanto los operadores lógicos como los espacios en blanco
		if !strings.Contains(operadores, simbolo) && simbolo != " " {
			if !alfabetoMap[simbolo] {
				alfabetoMap[simbolo] = true
				alfabeto = append(alfabeto, simbolo)
			}
		}
	}
	return alfabeto
}