# eventgo

Estructura inspirada en `authgo`:

```
internal/
  di/          # inyector (Mongo, servicios)
  env/         # variables de entorno
  domain/
    event/     # schema.go, repository.go, service.go
  rest/
    router.go
    server/
      engine.go
      auth_middleware.go
      post_event.go
      put_event.go
      delete_event.go
      get_event.go
      get_event_id.go
main.go
```

## Correr

```bash
cp .env.example .env
go mod tidy
go run ./main.go
```

Mongo debe estar en `mongodb://localhost:27017`. DB por defecto: `eventsdb`.

## API
- POST /event (auth)
- PUT /event/:id (auth)
- DELETE /event/:id (auth)
- GET /event
- GET /event/:id
