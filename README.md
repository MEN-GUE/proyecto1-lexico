# Analizador Léxico Inicial - Verificador de Lenguajes Regulares

Este repositorio privado contiene el desarrollo arquitectónico del **Proyecto No. 1** del curso, construido en lenguaje Go. Su finalidad principal consiste en la implementación computacional de algoritmos básicos orientados a la síntesis y transformación de autómatas finitos a partir de expresiones regulares complejas.

El sistema ha sido estructurado para validar la aceptación matemática de cadenas finitas en lenguajes regulares mediante la instanciación de autómatas deterministas y no deterministas.

## Autor:
*   **Juan Fernando Menéndez Guerra**


## Definiciones y Convenciones del Entorno

### Operador de Transición Vacía
Se requiere, por parte del equipo desarrollador, la elección de una designación para las transiciones vacías. En cumplimiento de la directriz de evitar colisiones con letras comunes o números que estadísticamente formen parte de otros aspectos del proyecto, se ha designado de forma universal y estricta el carácter numérico de control (`#`) para representar el símbolo épsilon $\epsilon$.

### Extracción del Alfabeto
El algoritmo ha sido diseñado para mapear rigurosamente los símbolos, todos los caracteres recolectados del alfabeto poseen una longitud estricta de 1 y están conformados por entidades distintas encontradas dentro de la expresión regular, aislando todos los caracteres operativos.

## Estructura Algorítmica Funcional
El motor de evaluación encapsula y ejecuta secuencialmente los siguientes algoritmos base:
1.  **Algoritmo de Shunting Yard:** Ejecuta la transformación de la notación original Infix hacia la notación posfija (Postfix).
2.  **Construcción de Thompson:** Genera la primera iteración gráfica y lógica de un Autómata Finito No Determinista (AFN).
3.  **Algoritmo de Subconjuntos:** Elimina la ambigüedad, generando y mapeando el Autómata Finito Determinista (AFD).
4.  **Minimización por Particiones:** Implementa la teoría de Hopcroft para colapsar los estados matemáticamente equivalentes, reduciendo la memoria lógica requerida.

## Requisitos de Entorno y Compilación
La orquestación gráfica subyacente para desplegar los estados (inicial, aceptación y adicionales) de los AFN, AFDs y los autómatas reducidos requiere el motor de renderización *Graphviz*.

En infraestructuras Linux operadas por DNF, se recomienda aplicar el siguiente comando como superusuario para resolver la dependencia:

```bash
sudo dnf install graphviz
```
# Instrucciones Operativas:
El sistema procesa secuencialmente un archivo de texto en plano, donde recibe como entrada una expresión regular r y una cadena w en cada línea de evaluación. Reemplace el contenido del archivo interno expresiones.txt con las líneas de evaluación de formato [regex] [cadena], las cuales serán entregadas al momento de presentar el software. Posiciónese en el entorno en la raíz del proyecto.Ejecute la orquestación principal del programa:

```Bash
go run cmd/analizador/main.go
```
El sistema arrojará en el canal de salida estándar la validación binaria pertinente para indicar si la cadena pertenece (sí) o no pertenece (no) al universo validado $L(r)$.