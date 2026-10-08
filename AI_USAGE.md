# Uso de IA

El enunciado permite usar cualquier herramienta de IA y pide compartir las indicaciones usadas.
Aquí explico cómo usé la IA en este proyecto, qué decidí yo y dónde cambié o rechacé sus sugerencias.

## Contexto

Antes de este proyecto no había usado Go. Lo elegí igual porque es el lenguaje preferido para
el puesto, y usé un asistente de IA (Claude) como tutor para aprenderlo mientras construía.
[Con tus palabras: por qué elegiste Go y cómo te sentías al empezar.]

## Cómo la usé

- **Planificación:** compartí el enunciado y discutí la arquitectura antes de escribir código
  (un solo microservicio sin estado, solo librería estándar, dominio separado de HTTP).
- **Guía paso a paso:** en cada parte (dominio, API HTTP, Docker, frontend, Makefile, CI) la IA
  explicaba el concepto y por qué es idiomático, y proponía el código. Yo lo escribía, lo corría
  y hacía el commit, con una rama y un pull request por pieza.
- **Revisión de código:** le pedí a la IA que revisara mi código. Por ejemplo, la revisión de
  `calculator.go` encontró que un NaN podía pasar como resultado válido, que `-0` se devolvía
  como `-0` y que los archivos no estaban formateados con `gofmt`. Yo hice las correcciones.
- **Tests:** la IA abrió el pull request #2 con los tests de los handlers HTTP. Lo revisé y lo
  fusioné. [Cuenta qué revisaste antes de fusionarlo.]
- **Documentación:** la IA redactó un borrador del README a partir del código, y yo lo revisé y edité.

## Decisiones que tomé (y sugerencias que cambié o rechacé)

- **Librería estándar en lugar de un framework como Gin:** menos dependencias, y el router de
  Go 1.22+ cubre todo lo que necesita esta API.
- **Errores en lugar de panics:** las funciones de dominio devuelven `(resultado, error)`, así la
  capa HTTP puede traducir cada error a un código de estado claro.
- **Un solo microservicio, no uno por operación:** dividirlo sería sobreingeniería.
- **`float64` en lugar de una librería decimal:** lo mantuve simple a propósito y está
  documentado en el README.
- [Agrega alguna sugerencia que rechazaste o cambiaste, y por qué.]

## Qué verifiqué yo

Corrí cada cambio en mi máquina antes de hacer commit (`go test`, `npm test`), y la CI revisa
formato, vet, lint, tests y build en cada push. [Agrega lo que hayas comprobado aparte.]

## Prompts principales

Estos son los prompts más importantes que usé, en orden. Los redacté con más claridad para este
documento, manteniendo la intención original de cada uno.

1. **Arquitectura:** "¿Cómo estructurarías el backend en Go para esta calculadora? Quiero una
   arquitectura sencilla, idiomática y fácil de testear, con separación entre la lógica de negocio
   y la capa HTTP. Prioriza la librería estándar y evita dependencias innecesarias."
   *Resultado:* `cmd/server` para arrancar el servidor, `internal/calculator` para la lógica pura
   e `internal/transport/http` para HTTP.

2. **Revisión de código:** "Revisa `calculator.go` y dime qué está mal o qué se puede mejorar
   antes de seguir. Busca errores de lógica, casos límite, problemas con `float64` y código que
   no siga las convenciones de Go. Explica cada hallazgo y propón correcciones mínimas."
   *Resultado:* se identificaron aspectos relacionados con `NaN`, la representación de `-0` y
   archivos sin formato `gofmt`. Revisé los hallazgos y realicé las correcciones correspondientes.

3. **Aprender testing:** "Ayúdame a escribir `calculator_test.go` poco a poco, explicándome cada
   parte para entenderla y no solo copiarla. Utiliza pruebas basadas en tablas (*table-driven
   tests*) e incluye casos normales, valores límite y errores esperados."
   *Resultado:* aprendí a utilizar tests basados en tablas, organizar los casos de prueba y
   verificar tanto los resultados como los errores esperados.

4. **API HTTP:** "Explícame en detalle qué hacer en la rama `feat/http-api`, paso a paso.
   Quiero entender cómo implementar los handlers, validar las peticiones, convertir los errores
   del dominio en respuestas HTTP, utilizar middleware y realizar un apagado ordenado del servidor.
   Indícame cómo comprobar que cada parte funciona antes de continuar."
   *Resultado:* trabajé en los handlers, el mapeo de errores, los middleware y el apagado ordenado.
   Los tests de los handlers llegaron como PR #2, que revisé antes de fusionar.

5. **Docker, frontend y cierre:** "Ya terminé esta rama, ¿cómo continuamos? Indícame el siguiente
   paso del proyecto y guíame poco a poco, explicando qué debo implementar, qué comandos ejecutar
   y cómo verificar que funciona antes de seguir."
   *Resultado:* utilicé este enfoque en cada etapa (Docker, frontend con React, Makefile, CI y
   README), avanzando paso a paso y verificando los cambios antes de continuar.

## Cómo escribiría estos prompts hoy

Mis prompts fueron cortos y directos. Funcionaron porque la conversación ya tenía contexto, pero
aprendí que un buen prompt dice el contexto, las restricciones y qué formato de respuesta espero.

- **Antes:** "Agreguemos un hilo revisando calculator.go."
  **Hoy:** "Revisa `calculator.go`, un paquete de dominio en Go sin dependencias. Busca errores de
  corrección (NaN, infinito, `-0`, división por cero) y código no idiomático. Ordena los hallazgos
  por gravedad y propón la corrección mínima de cada uno, sin reescribir el archivo."
- **Antes:** "Listo, ya terminamos con esa rama, ahora continuemos."
  **Hoy:** "Terminé `feat/http-api` y está fusionada. Según el plan, sigue Docker. Dame los pasos
  de esa rama, un commit por paso con su mensaje en Conventional Commits, y cómo verificar cada uno."

## Qué aprendí

[2 o 3 líneas con tus palabras: qué aprendiste de Go, de testing o del flujo de trabajo.]
