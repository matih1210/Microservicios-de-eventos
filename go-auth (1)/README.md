# go-auth

Mini servicio **Auth** en Go + Gin + MongoDB cumpliendo la API:

- `POST /user` (crear usuario)
- `POST /user/login` (Basic Auth → JWT con `sid`)
- `GET /user/current` (Bearer)
- `POST /user/logout` (Bearer)

## Config

```bash
cp .env.example .env
# Edita variables si hace falta
```

## Correr local

```bash
go mod tidy
go run ./cmd/server
```

## Docker

```bash
cp .env.example .env
docker compose up --build
```

## PowerShell ejemplos

```powershell
# Crear usuario
$body = @{ name="Matias"; username="mati"; password="secreto123" } | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:3001/user" -Method POST -ContentType "application/json" -Body $body

# Login (Basic)
$basic = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("mati:secreto123"))
Invoke-RestMethod -Uri "http://localhost:3001/user/login" -Method POST -Headers @{ "Authorization"="Basic $basic" }

# Current (Bearer)
$jwt = "<pega_token>"
Invoke-RestMethod -Uri "http://localhost:3001/user/current" -Method GET -Headers @{ "Authorization"="Bearer $jwt" }

# Logout (Bearer)
Invoke-RestMethod -Uri "http://localhost:3001/user/logout" -Method POST -Headers @{ "Authorization"="Bearer $jwt" }
```
