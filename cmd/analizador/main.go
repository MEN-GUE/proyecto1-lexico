package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"proyecto1/internal/afd"
	"proyecto1/internal/afn"
	"proyecto1/internal/ast"
	"proyecto1/internal/shuntingyard"
)

// procesarLinea garantiza un análisis léxico seguro sin dependencias frágiles de formato
func procesarLinea(linea string, index int) {
	// fields abstrae la complejidad de los múltiples espacios o tabulaciones
	partes := strings.Fields(linea)
	if len(partes) < 2 {
		fmt.Printf("Advertencia: La línea %d no posee suficientes elementos para conformar una expresión 'r' y una cadena 'w'[cite: 7]. Se omite.\n", index)
		return
	}

	// El formato asume que el último bloque de caracteres ininterrumpidos corresponde a la cadena a evaluar w[cite: 7].
	cadena := partes[len(partes)-1]
	// Todo el contenido que antecede a la cadena final se agrupa para conformar la expresión regular r[cite: 7].
	regex := strings.Join(partes[:len(partes)-1], "")

	fmt.Printf("\n================================================================================\n")
	fmt.Printf("Expresión Regular (r): %s\n", regex)
	fmt.Printf("Cadena a evaluar (w): %s\n", cadena)

	tokens := shuntingyard.Tokenizar(regex)
	tokensSimples := shuntingyard.SimplificarExtensiones(tokens)
	tokensConConcat := shuntingyard.AgregarConcatenacionExplicita(tokensSimples)
	postfix := shuntingyard.Convertir(tokensConConcat)

	raizAST := ast.ConstruirAST(postfix)

	afn.ReiniciarContador()
	automataNoDeterminista := afn.ConstruirThompson(raizAST)
	graficarAFN(automataNoDeterminista, index)
	
	resultadoAFN := afn.Simular(automataNoDeterminista, cadena)
	imprimirResultado("AFN", cadena, resultadoAFN)
	
	alfabeto := afd.ObtenerAlfabeto(regex)
	
	automataDeterminista := afd.ConstruirSubconjuntos(automataNoDeterminista, alfabeto)
	graficarAFD(automataDeterminista, index, "afd")
	
	resultadoAFD := afd.Simular(automataDeterminista, cadena)
	imprimirResultado("AFD", cadena, resultadoAFD)
	
	automataMinimizado := afd.Minimizar(automataDeterminista)
	graficarAFD(automataMinimizado, index, "min")
	
	resultadoMin := afd.Simular(automataMinimizado, cadena)
	imprimirResultado("AFD Minimizado", cadena, resultadoMin)
}

func graficarAFN(automata *afn.AFN, index int) {
	dotSource := afn.GenerarDOT(automata)
	dotFilename := fmt.Sprintf("afn_%d.dot", index)
	pngFilename := fmt.Sprintf("afn_%d.png", index)

	os.WriteFile(dotFilename, []byte(dotSource), 0644)
	exec.Command("dot", "-Tpng", dotFilename, "-o", pngFilename).Run()
	// Llamada directa al subsistema X Window System de la distribución Linux para desplegar las imágenes renderizadas al vuelo
	exec.Command("xdg-open", pngFilename).Start()
}

func graficarAFD(automata *afd.AFD, index int, prefijo string) {
	dotSource := afd.GenerarDOT(automata)
	dotFilename := fmt.Sprintf("%s_%d.dot", prefijo, index)
	pngFilename := fmt.Sprintf("%s_%d.png", prefijo, index)

	os.WriteFile(dotFilename, []byte(dotSource), 0644)
	err := exec.Command("dot", "-Tpng", dotFilename, "-o", pngFilename).Run()
	if err != nil {
		fmt.Printf("Error a nivel de sistema operativo al procesar Graphviz: %v\n", err)
		return
	}
	
	fmt.Printf("✅ %s generado de forma íntegra en el nodo: %s\n", strings.ToUpper(prefijo), pngFilename)
	exec.Command("xdg-open", pngFilename).Start()
}

// imprimirResultado estandariza la salida exigida por la rúbrica indicando de forma binaria si la cadena pertenece a L(r)[cite: 7].
func imprimirResultado(tipo string, cadena string, aceptada bool) {
	if aceptada {
		fmt.Printf("Simulación %s: sí, '%s' pertenece a L(r)\n", tipo, cadena)
	} else {
		fmt.Printf("Simulación %s: no, '%s' NO pertenece a L(r)\n", tipo, cadena)
	}
}

func main() {
	fmt.Println("=== ANALIZADOR LÉXICO - MOTOR DE EVALUACIÓN DE AUTÓMATAS ===")

	archivo, err := os.Open("expresiones.txt")
	if err != nil {
		fmt.Println("Excepción de I/O: Fallo de puntero al invocar 'expresiones.txt'.")
		return
	}
	defer archivo.Close()

	escaner := bufio.NewScanner(archivo)
	index := 1

	for escaner.Scan() {
		linea := strings.TrimSpace(escaner.Text())
		if len(linea) > 0 {
			procesarLinea(linea, index)
			index++
		}
	}
}