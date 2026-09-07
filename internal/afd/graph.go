package afd

import (
	"fmt"
	"strings"
)

// GenerarDOT genera la representación en Graphviz del AFD, mostrando estados iniciales, de aceptación y transiciones.
func GenerarDOT(afd *AFD) string {
	var builder strings.Builder
	builder.WriteString("digraph AFD {\n")
	builder.WriteString("\trankdir=LR;\n")
	builder.WriteString("\tnode [shape=circle, fontname=\"Arial\"];\n")

	// 1. Definir los estados de aceptación (Doble círculo)[cite: 6]
	for _, estado := range afd.Estados {
		if estado.Aceptacion {
			builder.WriteString(fmt.Sprintf("\tnode%d [shape=doublecircle];\n", estado.ID))
		}
	}

	// 2. Definir el estado inicial (Flecha de entrada)[cite: 6]
	builder.WriteString("\tinicio [shape=point];\n")
	builder.WriteString(fmt.Sprintf("\tinicio -> node%d;\n", afd.Inicial.ID))

	// 3. Escribir todas las transiciones con sus símbolos correspondientes[cite: 6]
	for _, estado := range afd.Estados {
		for _, t := range estado.Transiciones {
			label := t.Simbolo
			if label == "|" {
				label = "\\|" // Escape seguro para Graphviz
			}
			builder.WriteString(fmt.Sprintf("\tnode%d -> node%d [label=\"%s\"];\n", estado.ID, t.Destino.ID, label))
		}
	}

	builder.WriteString("}\n")
	return builder.String()
}