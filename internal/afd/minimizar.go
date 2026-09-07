package afd

import (
	"fmt"
)

// encontrarGrupo busca a qué partición pertenece un estado actual
func encontrarGrupo(estado *EstadoAFD, particiones [][]*EstadoAFD) int {
	if estado == nil {
		return -1
	}
	for i, grupo := range particiones {
		for _, e := range grupo {
			if e.ID == estado.ID {
				return i
			}
		}
	}
	return -1
}

// obtenerDestino devuelve el estado destino al consumir un símbolo
func obtenerDestino(estado *EstadoAFD, simbolo string) *EstadoAFD {
	for _, t := range estado.Transiciones {
		if t.Simbolo == simbolo {
			return t.Destino
		}
	}
	return nil
}

// Minimizar implementa el algoritmo de particiones para optimizar el AFD.
func Minimizar(afdOriginal *AFD) *AFD {
	// 1. Partición inicial: Separamos en estados de Aceptación y de No Aceptación
	var aceptacion []*EstadoAFD
	var normales []*EstadoAFD

	for _, e := range afdOriginal.Estados {
		if e.Aceptacion {
			aceptacion = append(aceptacion, e)
		} else {
			normales = append(normales, e)
		}
	}

	var particiones [][]*EstadoAFD
	if len(normales) > 0 {
		particiones = append(particiones, normales)
	}
	if len(aceptacion) > 0 {
		particiones = append(particiones, aceptacion)
	}

	// 2. Refinamiento iterativo de las particiones
	cambio := true
	for cambio {
		cambio = false
		var nuevasParticiones [][]*EstadoAFD

		for _, grupo := range particiones {
			// Si el grupo tiene 1 o 0 elementos, no se puede subdividir más
			if len(grupo) <= 1 {
				nuevasParticiones = append(nuevasParticiones, grupo)
				continue
			}

			// Agrupar los estados por su "firma" de comportamiento (hacia qué grupos transicionan)
			comportamientos := make(map[string][]*EstadoAFD)
			for _, estado := range grupo {
				firma := ""
				for _, sim := range afdOriginal.Alfabeto {
					destino := obtenerDestino(estado, sim)
					idGrupoDestino := encontrarGrupo(destino, particiones)
					firma += fmt.Sprintf("%s:%d|", sim, idGrupoDestino)
				}
				comportamientos[firma] = append(comportamientos[firma], estado)
			}

			// Convertir el mapa de comportamientos en nuevas particiones
			for _, subgrupo := range comportamientos {
				nuevasParticiones = append(nuevasParticiones, subgrupo)
			}
			// Si un grupo se dividió en más de un comportamiento, hubo cambios
			if len(comportamientos) > 1 {
				cambio = true
			}
		}
		particiones = nuevasParticiones
	}

	// 3. Reconstruir el nuevo AFD Minimizado a partir de las particiones finales
	estadosMinimizados := make(map[int]*EstadoAFD)
	var nuevoInicial *EstadoAFD
	var listaEstados []*EstadoAFD

	// Crear un nuevo estado por cada partición final
	for i, grupo := range particiones {
		representante := grupo[0]
		nuevoEstado := &EstadoAFD{
			ID:         i,
			Aceptacion: representante.Aceptacion,
		}
		estadosMinimizados[i] = nuevoEstado
		listaEstados = append(listaEstados, nuevoEstado)

		// Identificar si este grupo contiene al estado inicial original
		for _, e := range grupo {
			if e.ID == afdOriginal.Inicial.ID {
				nuevoInicial = nuevoEstado
			}
		}
	}

	// Enlazar las transiciones para los nuevos estados
	for i, grupo := range particiones {
		representante := grupo[0]
		estadoActual := estadosMinimizados[i]

		for _, sim := range afdOriginal.Alfabeto {
			destinoOriginal := obtenerDestino(representante, sim)
			if destinoOriginal != nil {
				idGrupoDestino := encontrarGrupo(destinoOriginal, particiones)
				if idGrupoDestino != -1 {
					estadoActual.Transiciones = append(estadoActual.Transiciones, TransicionAFD{
						Simbolo: sim,
						Destino: estadosMinimizados[idGrupoDestino],
					})
				}
			}
		}
	}

	return &AFD{
		Inicial:  nuevoInicial,
		Estados:  listaEstados,
		Alfabeto: afdOriginal.Alfabeto,
	}
}