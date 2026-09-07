package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"proyecto1/internal/afn"
	"proyecto1/internal/ast"
	"proyecto1/internal/shuntingyard"
	// "proyecto1/internal/afd" // TODO: Descomentar al implementar Subconjuntos
)

func procesarLinea(linea string, index int) {
	partes := strings.Split(linea, " ")
	if len(partes) != 2 {
		fmt.Printf("Error: La línea %d no cumple con el formato 'regex cadena'\n", index)
		return
	}

	regex := partes[0]
	cadena := partes[1]

	fmt.Printf("\n================================================================================\n")
	fmt.Printf("Expresión Regular (r): %s\n", regex)
	fmt.Printf("Cadena a evaluar (w): %s\n", cadena)

	// 1. Conversión Infix a Postfix (Shunting Yard)
	tokens := shuntingyard.Tokenizar(regex)
	tokensSimples := shuntingyard.SimplificarExtensiones(tokens)
	tokensConConcat := shuntingyard.AgregarConcatenacionExplicita(tokensSimples)
	postfix := shuntingyard.Convertir(tokensConConcat)

	// 2. Construcción del Árbol Sintáctico Abstracto
	raizAST := ast.ConstruirAST(postfix)

	// 3. Construcción del AFN (Algoritmo de Thompson)
	afn.ReiniciarContador()
	automataNoDeterminista := afn.ConstruirThompson(raizAST)

	// 4. Graficación del AFN
	graficarAFN(automataNoDeterminista, index)

	// 5. Simulación del AFN
	resultadoAFN := afn.Simular(automataNoDeterminista, cadena)
	imprimirResultado("AFN", cadena, resultadoAFN)

	// ====================================================================
	// FASE 2: CONSTRUCCIÓN DE SUBCONJUNTOS Y MINIMIZACIÓN
	// ====================================================================
	
	// TODO: Generación de AFD con Subconjuntos
	// automataDeterminista := afd.ConstruirSubconjuntos(automataNoDeterminista)
	// graficarAFD(automataDeterminista, index, "afd")
	// resultadoAFD := afd.Simular(automataDeterminista, cadena)
	// imprimirResultado("AFD", cadena, resultadoAFD)

	// TODO: Minimización de AFD
	// automataMinimizado := afd.Minimizar(automataDeterminista)
	// graficarAFD(automataMinimizado, index, "min")
	// resultadoMin := afd.Simular(automataMinimizado, cadena)
	// imprimirResultado("AFD Minimizado", cadena, resultadoMin)
}

func graficarAFN(automata *afn.AFN, index int) {
	dotSource := afn.GenerarDOT(automata)
	dotFilename := fmt.Sprintf("afn_%d.dot", index)
	pngFilename := fmt.Sprintf("afn_%d.png", index)

	os.WriteFile(dotFilename, []byte(dotSource), 0644)
	err := exec.Command("dot", "-Tpng", dotFilename, "-o", pngFilename).Run()
	if err != nil {
		fmt.Printf("Error al generar la imagen %s. Verifica tu instalación de Graphviz.\n", pngFilename)
		return
	}
	
	fmt.Printf("✅ AFN generado exitosamente en: %s\n", pngFilename)
	exec.Command("xdg-open", pngFilename).Start()
}

// Función auxiliar para imprimir las métricas de aceptación requeridas[cite: 6]
func imprimirResultado(tipo string, cadena string, aceptada bool) {
	if aceptada {
		fmt.Printf("Simulación %s: sí, '%s' pertenece a L(r)\n", tipo, cadena)
	} else {
		fmt.Printf("Simulación %s: no, '%s' NO pertenece a L(r)\n", tipo, cadena)
	}
}

func main() {
	fmt.Println("=== ANALIZADOR LÉXICO - PROYECTO 1 ===")

	archivo, err := os.Open("expresiones.txt")
	if err != nil {
		fmt.Println("Error: No se encontró 'expresiones.txt' en la raíz del proyecto.")
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

	if err := escaner.Err(); err != nil {
		fmt.Println("Error de lectura durante el escaneo del archivo:", err)
	}
}