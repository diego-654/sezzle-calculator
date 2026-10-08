# Calculadora Sezzle

[![CI](https://github.com/diego-654/sezzle-calculator/actions/workflows/ci.yml/badge.svg)](https://github.com/diego-654/sezzle-calculator/actions/workflows/ci.yml)

Calculadora full-stack: un frontend en React + TypeScript que consume un microservicio REST en Go.
Soporta suma, resta, multiplicación, división, potencia, raíz cuadrada y porcentaje.

## Estructura del proyecto

```
.
├── backend/                     # Microservicio en Go (solo librería estándar)
│   ├── cmd/server/              # punto de entrada: configuración, servidor HTTP, apagado ordenado
│   └── internal/
│       ├── calculator/          # lógica de dominio pura, devuelve (resultado, error)
│       └── transport/http/      # handlers, JSON, mapeo de errores, middleware
├── frontend/                    # Vite + React + TypeScript, servido por nginx en Docker
├── docker-compose.yml           # levanta los dos servicios juntos
├── Makefile                     # atajos para tests, cobertura y Docker
└── .github/workflows/ci.yml     # lint, tests y build en cada push y PR
```

## Requisitos

- Go 1.27+
- Node.js 22+
- Docker (opcional, para levantar todo con un solo comando)
- make (opcional; cada target es un comando simple que también se puede correr a mano)

## Cómo ejecutarlo

### Opción 1: Docker Compose (recomendada)

```bash
docker compose up --build
```

Abre http://localhost:3000. nginx sirve el frontend y redirige `/api` al backend,
que no se expone fuera de Docker.

### Opción 2: cada parte por separado

Backend (escucha en el puerto 8080; se cambia con la variable de entorno `PORT`):

```bash
cd backend
go run ./cmd/server
```

Frontend, en otra terminal:

```bash
cd frontend
npm install
npm run dev
```

Abre http://localhost:5173. El servidor de desarrollo de Vite redirige `/api` a `http://localhost:8080`.

## Tests y cobertura

```bash
make test        # tests unitarios de backend y frontend
make coverage    # lo mismo, con informe de cobertura
```

Sin make:

```bash
cd backend && go test -cover ./...
cd frontend && npm run coverage
```

Cobertura actual: 94% en la capa HTTP y 65% en el paquete `calculator` (backend), y cerca
del 96% de las sentencias en el frontend. El informe HTML del frontend queda en
`frontend/coverage/`. La CI corre los mismos chequeos, además de `gofmt`, `go vet`, ESLint
y el build de producción.

## API

Todas las operaciones usan `POST /api/v1/{operación}` con un cuerpo JSON.

| Operación  | Cuerpo                | Resultado      |
|------------|-----------------------|----------------|
| `add`      | `{"a": 2, "b": 3}`    | `a + b`        |
| `subtract` | `{"a": 2, "b": 3}`    | `a - b`        |
| `multiply` | `{"a": 2, "b": 3}`    | `a * b`        |
| `divide`   | `{"a": 6, "b": 3}`    | `a / b`        |
| `power`    | `{"a": 2, "b": 3}`    | `a ^ b`        |
| `percent`  | `{"a": 200, "b": 15}` | `b% de a`      |
| `sqrt`     | `{"a": 16}`           | `√a`           |

`GET /health` devuelve `{"status":"ok"}` para los health checks.

### Ejemplos

Con curl (Linux, macOS o Git Bash):

```bash
curl -X POST http://localhost:8080/api/v1/add \
  -H "Content-Type: application/json" -d '{"a": 2, "b": 3}'
# 200 {"result":5}

curl -X POST http://localhost:8080/api/v1/sqrt \
  -H "Content-Type: application/json" -d '{"a": 16}'
# 200 {"result":4}

curl -X POST http://localhost:8080/api/v1/divide \
  -H "Content-Type: application/json" -d '{"a": 1, "b": 0}'
# 422 {"error":{"code":"DIVISION_BY_ZERO","message":"division by zero"}}
```

Con PowerShell (Windows):

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/v1/add `
  -ContentType "application/json" -Body '{"a": 2, "b": 3}'
# result : 5
```

Con Docker Compose, usa `http://localhost:3000/api/v1/...` en lugar del puerto 8080.

### Errores

Todos los errores tienen la misma forma: `{"error": {"code": "...", "message": "..."}}`.

| Estado | Cuándo | Códigos |
|--------|--------|---------|
| 400 | Petición mal formada: JSON inválido, falta un operando, campo desconocido, NaN/Infinity | `INVALID_INPUT`, `INVALID_OPERAND` |
| 404 | Operación desconocida | `INVALID_OPERATION` |
| 405 | Método HTTP incorrecto | (respuesta estándar del router de Go) |
| 422 | Petición válida que la matemática no puede resolver | `DIVISION_BY_ZERO`, `NEGATIVE_SQUARE_ROOT`, `INVALID_EXPONENT`, `OVERFLOW`, `UNDEFINED_RESULT` |
| 500 | Error inesperado (el detalle va al log, no a la respuesta) | `INTERNAL_ERROR` |

## Decisiones de diseño

- **Un solo microservicio, no uno por operación.** El backend es un único servicio sin estado,
  con su propio `go.mod` y Dockerfile, health check, configuración por variables de entorno,
  logs estructurados (`log/slog`) y apagado ordenado al recibir SIGTERM. Separar cada operación
  en un servicio solo agregaría llamadas de red y trabajo de despliegue sin ningún beneficio.
- **Solo la librería estándar.** El router de Go 1.22+ (`POST /api/v1/{op}`) cubre todo lo que
  necesita esta API, así que no hay un framework que aprender ni mantener.
- **Dominio separado de HTTP.** `internal/calculator` no sabe nada de HTTP y devuelve errores en
  lugar de hacer panic. La capa HTTP traduce cada error de dominio a un código de estado en una
  sola tabla, así los handlers quedan pequeños y las dos capas son fáciles de testear.
- **Validación en ambas capas.** El frontend valida para dar una buena experiencia de uso; el
  backend valida de nuevo porque es la fuente de verdad y se puede llamar directamente.
- **Mismo origen, sin CORS.** nginx (en Docker) y el proxy de Vite (en desarrollo) mandan `/api`
  al backend, así que el navegador solo habla con un origen.
- **Casos borde.** Se rechazan entradas NaN y ±Infinity, un resultado que desborda devuelve
  `OVERFLOW`, `0` elevado a una potencia negativa es una división por cero, una base negativa
  con exponente fraccionario se rechaza y `-0` se normaliza a `0`.

### Precisión numérica

El servicio usa `float64`, así que los resultados siguen IEEE 754: `0.1 + 0.2` devuelve
`0.30000000000000004`. Es una concesión consciente para una calculadora simple. Para manejar
dinero usaría un tipo decimal (por ejemplo `shopspring/decimal`) o centavos enteros, y
definiría una regla de redondeo.

## Qué haría con más tiempo

- Aritmética decimal y una política de redondeo explícita.
- Más tests unitarios para los casos borde del paquete `calculator`.
- Una especificación OpenAPI como contrato de la API.
- Rate limiting, métricas con Prometheus y tracing.
- Despliegue en Kubernetes aprovechando el health check existente.

## Uso de IA

En [AI_USAGE.md](AI_USAGE.md) explico cómo usé herramientas de IA en este proyecto.
