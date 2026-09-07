package afd

// Simular evalúa la cadena caracter por caracter de forma determinista[cite: 6].
func Simular(afd *AFD, cadena string) bool {
	actual := afd.Inicial

	for _, char := range cadena {
		simbolo := string(char)
		encontrado := false

		// Buscamos la transición exacta para el símbolo actual
		for _, t := range actual.Transiciones {
			if t.Simbolo == simbolo {
				actual = t.Destino
				encontrado = true
				break
			}
		}

		// Si no hay transición para un símbolo, caemos en un estado de error implícito (sumidero) y se rechaza.
		if !encontrado {
			return false
		}
	}

	// La cadena es aceptada únicamente si el estado donde terminó tiene la bandera de aceptación[cite: 6].
	return actual.Aceptacion
}