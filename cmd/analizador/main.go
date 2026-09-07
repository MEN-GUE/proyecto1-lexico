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

	// 1. Conversión Infix a Postfix (Shunting Yard)[cite: 6]
	tokens := shuntingyard.Tokenizar(regex)
	tokensSimples := shuntingyard.SimplificarExtensiones(tokens)
	tokensConConcat := shuntingyard.AgregarConcatenacionExplicita(tokensSimples)
	postfix := shuntingyard.Convertir(tokensConConcat)

	// 2. Construcción del Árbol Sintáctico Abstracto
	raizAST := ast.ConstruirAST(postfix)

	// 3. Construcción y Simulación del AFN (Algoritmo de Thompson)[cite: 6]
	afn.ReiniciarContador()
	automataNoDeterminista := afn.ConstruirThompson(raizAST)
	graficarAFN(automataNoDeterminista, index)
	
	resultadoAFN := afn.Simular(automataNoDeterminista, cadena)
	imprimirResultado("AFN", cadena, resultadoAFN)

	// ====================================================================
	// FASE 2: CONSTRUCCIÓN DE SUBCONJUNTOS, GRAFICACIÓN Y SIMULACIÓN
	// ====================================================================
	
	alfabeto := afd.ObtenerAlfabeto(regex)
	
	// Generación de AFD con Subconjuntos[cite: 6]
	automataDeterminista := afd.ConstruirSubconjuntos(automataNoDeterminista, alfabeto)
	
	// Graficación del AFD generado[cite: 6]
	graficarAFD(automataDeterminista, index, "afd")
	
	// Simulación determinista de la cadena w[cite: 6]
	resultadoAFD := afd.Simular(automataDeterminista, cadena)
	imprimirResultado("AFD", cadena, resultadoAFD)

	// ====================================================================
	// FASE 3: MINIMIZACIÓN DE AFD
	// ====================================================================
	// Minimización de AFD[cite: 6]
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
	exec.Command("xdg-open", pngFilename).Start()
}

// graficarAFD centraliza la renderización de grafos para el AFD normal y el minimizado[cite: 6]
func graficarAFD(automata *afd.AFD, index int, prefijo string) {
	dotSource := afd.GenerarDOT(automata)
	dotFilename := fmt.Sprintf("%s_%d.dot", prefijo, index)
	pngFilename := fmt.Sprintf("%s_%d.png", prefijo, index)

	os.WriteFile(dotFilename, []byte(dotSource), 0644)
	err := exec.Command("dot", "-Tpng", dotFilename, "-o", pngFilename).Run()
	if err != nil {
		fmt.Printf("Error al generar la imagen %s.\n", pngFilename)
		return
	}
	
	fmt.Printf("✅ %s generado exitosamente en: %s\n", strings.ToUpper(prefijo), pngFilename)
	exec.Command("xdg-open", pngFilename).Start()
}

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
		fmt.Println("Error: No se encontró 'expresiones.txt'.")
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