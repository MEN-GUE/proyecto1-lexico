# Proyecto No. 1 - Analizador Léxico (Fase Inicial)
**Autor:** Juan Fernando Menéndez Guerra

Este proyecto implementa la construcción de autómatas finitos (AFN y AFD) a partir de expresiones regulares, incluyendo la minimización del AFD y la simulación de cadenas para verificar su aceptación.

## Símbolo Épsilon
El símbolo designado para representar a $\epsilon$ en este analizador es `E`.

## Requisitos Previos
Sistema Linux con Graphviz instalado para la generación de autómatas:
`sudo dnf install graphviz`

## Ejecución
Para procesar el archivo de expresiones:
`go run cmd/analizador/main.go`