# signupgo (Servicio Inscripciones)

Estructura similar a `eventgo`:

- `internal/rest/server/` → un archivo por endpoint (`post_signup.go`, `delete_signup.go`, `get_signup.go`) + `auth_middleware.go`
- `internal/domain/signup/` → `schema.go`, `repository.go`, `service.go`
- `internal/domain/event/` → `reader.go` (para validar evento activo)
- `internal/di/`, `internal/env/`, `main.go` en raíz

## API

- `POST /signup` (auth) — crea inscripción en `eventId`. Rechaza duplicadas y eventos cancelados/no existentes.
- `DELETE /signup/:id` (auth) — cancela la inscripción (solo el usuario dueño).
- `GET /signup?eventId=...` — lista inscriptos activos `{ userName, userId, id, signupDate }`.

## Correr

```bash
cp .env.example .env
go mod tidy
go run ./main.go
```

Requiere Mongo `mongodb://localhost:27017` (DB `eventsdb` por defecto) para leer la colección `events`. `JWT_SECRET` debe coincidir con Auth/Event.
