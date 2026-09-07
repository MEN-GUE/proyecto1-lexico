package afn

// Adaptación para el Proyecto 1: Actualización de la ruta de importación para coincidir con la nueva inicialización del módulo[cite: 6].
import "proyecto1/internal/ast"

type Transicion struct {
	Simbolo string
	Destino *Estado
}

type Estado struct {
	ID           int
	Transiciones []Transicion
}

type AFN struct {
	Inicial *Estado
	Final   *Estado
}

var contadorEstados int

func ReiniciarContador() {
	contadorEstados = 0
}

func nuevoEstado() *Estado {
	estado := &Estado{ID: contadorEstados}
	contadorEstados++
	return estado
}

func ConstruirThompson(nodo *ast.Nodo) *AFN {
	if nodo == nil {
		return nil
	}

	switch nodo.Valor {
	case "~": 
		afnIzq := ConstruirThompson(nodo.Izq)
		afnDer := ConstruirThompson(nodo.Der)
		afnIzq.Final.Transiciones = append(afnIzq.Final.Transiciones, Transicion{Simbolo: "#", Destino: afnDer.Inicial})
		return &AFN{Inicial: afnIzq.Inicial, Final: afnDer.Final}

	case "|": 
		afnIzq := ConstruirThompson(nodo.Izq)
		afnDer := ConstruirThompson(nodo.Der)
		inicial := nuevoEstado()
		final := nuevoEstado()

		inicial.Transiciones = append(inicial.Transiciones, Transicion{Simbolo: "#", Destino: afnIzq.Inicial})
		inicial.Transiciones = append(inicial.Transiciones, Transicion{Simbolo: "#", Destino: afnDer.Inicial})
		afnIzq.Final.Transiciones = append(afnIzq.Final.Transiciones, Transicion{Simbolo: "#", Destino: final})
		afnDer.Final.Transiciones = append(afnDer.Final.Transiciones, Transicion{Simbolo: "#", Destino: final})

		return &AFN{Inicial: inicial, Final: final}

	case "*": 
		afnBase := ConstruirThompson(nodo.Izq)
		inicial := nuevoEstado()
		final := nuevoEstado()

		inicial.Transiciones = append(inicial.Transiciones, Transicion{Simbolo: "#", Destino: afnBase.Inicial})
		inicial.Transiciones = append(inicial.Transiciones, Transicion{Simbolo: "#", Destino: final})
		afnBase.Final.Transiciones = append(afnBase.Final.Transiciones, Transicion{Simbolo: "#", Destino: afnBase.Inicial})
		afnBase.Final.Transiciones = append(afnBase.Final.Transiciones, Transicion{Simbolo: "#", Destino: final})

		return &AFN{Inicial: inicial, Final: final}

	default: 
		inicial := nuevoEstado()
		final := nuevoEstado()
		inicial.Transiciones = append(inicial.Transiciones, Transicion{Simbolo: nodo.Valor, Destino: final})
		return &AFN{Inicial: inicial, Final: final}
	}
}