package afd

import (
	"fmt"
	"sort"

	"proyecto1/internal/afn"
)

// Estructuras nativas para el AFD
type TransicionAFD struct {
	Simbolo string
	Destino *EstadoAFD
}

type EstadoAFD struct {
	ID           int
	Transiciones []TransicionAFD
	Aceptacion   bool
}

type AFD struct {
	Inicial  *EstadoAFD
	Estados  []*EstadoAFD
	Alfabeto []string
}

// tempState nos sirve para el control de la tabla de subconjuntos
type tempState struct {
	afdState *EstadoAFD
	afnSet   []*afn.Estado
	marcado  bool
}

// hashEstados crea un identificador único para un subconjunto de estados (ej. "[1 2 4]")
func hashEstados(estados []*afn.Estado) string {
	var ids []int
	for _, e := range estados {
		ids = append(ids, e.ID)
	}
	sort.Ints(ids)
	return fmt.Sprintf("%v", ids)
}

// contieneAceptacion verifica si en el subconjunto actual existe el estado de aceptación del AFN
func contieneAceptacion(estados []*afn.Estado, idAceptacion int) bool {
	for _, e := range estados {
		if e.ID == idAceptacion {
			return true
		}
	}
	return false
}

// ConstruirSubconjuntos aplica el algoritmo para convertir un AFN a AFD[cite: 6].
func ConstruirSubconjuntos(automataAFN *afn.AFN, alfabeto []string) *AFD {
	estadosMapa := make(map[string]*tempState)
	var estadosAFD []*EstadoAFD
	contador := 0

	// 1. Estado inicial: Cierre-épsilon del estado inicial del AFN
	inicialAFN := []*afn.Estado{automataAFN.Inicial}
	cierreInicial := afn.CerraduraEpsilon(inicialAFN)
	hashInicial := hashEstados(cierreInicial)

	estadoInicial := &EstadoAFD{
		ID:         contador,
		Aceptacion: contieneAceptacion(cierreInicial, automataAFN.Final.ID),
	}
	contador++

	tempIni := &tempState{
		afdState: estadoInicial,
		afnSet:   cierreInicial,
		marcado:  false,
	}

	estadosMapa[hashInicial] = tempIni
	estadosAFD = append(estadosAFD, estadoInicial)

	// 2. Ciclo principal del algoritmo (hasta que todos los estados estén marcados)
	hayNoMarcados := true
	for hayNoMarcados {
		hayNoMarcados = false
		var T *tempState

		// Buscar un estado en el Dstates que no esté marcado
		for _, state := range estadosMapa {
			if !state.marcado {
				T = state
				hayNoMarcados = true
				break
			}
		}

		if !hayNoMarcados {
			break
		}

		T.marcado = true // Marcamos T

		// Por cada símbolo en el alfabeto, calcular el movimiento
		for _, simbolo := range alfabeto {
			// U = CerraduraEpsilon(Mover(T, símbolo))
			movimiento := afn.Mover(T.afnSet, simbolo)
			if len(movimiento) == 0 {
				continue // Va a un estado de error/muerto, lo omitimos en la representación
			}

			cierreMov := afn.CerraduraEpsilon(movimiento)
			hashU := hashEstados(cierreMov)

			// Si U no está en Dstates, agregarlo como estado no marcado
			if estadosMapa[hashU] == nil {
				nuevoEstado := &EstadoAFD{
					ID:         contador,
					Aceptacion: contieneAceptacion(cierreMov, automataAFN.Final.ID),
				}
				contador++

				estadosMapa[hashU] = &tempState{
					afdState: nuevoEstado,
					afnSet:   cierreMov,
					marcado:  false,
				}
				estadosAFD = append(estadosAFD, nuevoEstado)
			}

			// Agregar la transición de T a U con el símbolo
			T.afdState.Transiciones = append(T.afdState.Transiciones, TransicionAFD{
				Simbolo: simbolo,
				Destino: estadosMapa[hashU].afdState,
			})
		}
	}

	return &AFD{
		Inicial:  estadoInicial,
		Estados:  estadosAFD,
		Alfabeto: alfabeto,
	}
}