# Proyecto Microservicios: Auth, Event y Signup

Este proyecto implementa un sistema distribuido basado en **microservicios** para autenticación, gestión de eventos e inscripciones.  
El objetivo principal es demostrar cómo un **logout invalida el token** de sesión mediante introspección centralizada.

---

## 📌 Arquitectura

El sistema está compuesto por tres microservicios independientes, cada uno con su propia base de datos:

- **Auth (go-auth):** manejo de usuarios, login/logout y validación de sesiones con JWT.  
- **Event:** creación y gestión de eventos.  
- **Signup:** inscripciones de usuarios en eventos.

Cada microservicio expone una API REST y se comunica con los demás vía **HTTP (Bearer Tokens)**.  
No existe base de datos compartida: cada servicio mantiene su propia persistencia.

---

## ⚙️ Tecnologías

- **Lenguaje:** Go  
- **Base de datos:** MongoDB (una por microservicio)  
- **Autenticación:** JWT con introspección  
- **Comunicación entre servicios:** HTTP + Bearer Tokens  
- **Contenedores:** Docker (opcional)

---

## 🔑 Gestión de Sesiones y JWT

- En el **login**, el servicio Auth:
  - Genera un `sid` único.  
  - Inserta la sesión en la base (`{ _id: sid, userId, expires }`).  
  - Devuelve un **JWT** que incluye `{ sid, uid, usr, exp }`.

- En el **logout**, Auth borra la sesión de la base.  
  - Como resultado, el token queda inválido: los demás servicios ya no lo aceptan.  

- En cada request protegida, Event y Signup consultan a Auth vía **`/token/introspect`** para verificar:
  1. Firma y expiración del JWT.  
  2. Que la sesión (`sid`) siga existiendo.  
  - Si no existe, Auth devuelve **401**, invalidando el token.

---

## 🔄 Comunicación entre microservicios

- **Cliente → Auth:** login, logout, user/current.  
- **Cliente → Event/Signup:** operaciones protegidas (Bearer JWT).  
- **Event/Signup → Auth:** introspección de tokens en `/token/introspect`.  
- **Signup → Event:** validación de evento en `/event/:id`.

---

## ▶️ Recorridos (ejemplo)

### A) Login → Crear evento
1. Cliente hace login en Auth → obtiene JWT con `sid`.  
2. Cliente llama a Event con el JWT.  
3. Event valida el token vía introspección en Auth.  
4. Si es válido, crea el evento asociado al usuario.

### B) Login → Inscribirse a un evento
1. Cliente hace login en Auth → obtiene JWT con `sid`.  
2. Cliente llama a Signup con el JWT y `eventId`.  
3. Signup valida el token vía introspección.  
4. Consulta a Event que el evento exista y no esté cancelado.  
5. Si todo es correcto, guarda la inscripción.

### C) Logout → Intentar crear evento/inscribirse
1. Cliente hace logout en Auth → se borra la sesión (`sid`).  
2. Cliente intenta usar el mismo JWT en Event o Signup.  
3. Introspección falla (Auth devuelve 401) → operación rechazada.

---

## ❓ ¿Por qué el logout invalida el token?

- Antes: cada microservicio validaba **localmente** la firma del JWT ⇒ no sabían que se había hecho logout.  
- Ahora: Event/Signup **preguntan a Auth** en cada request si el `sid` sigue activo.  
- Logout elimina la sesión ⇒ Auth responde **401** en `/token/introspect` ⇒ el token deja de ser válido automáticamente.

---

## 🚀 Puesta en marcha

1. Clonar el repositorio.  
2. Configurar variables de entorno:  

```env
AUTH_BASE_URL=http://localhost:3001
HTTP_TIMEOUT_MS=2000
```
3. Levantar los servicios (ejemplo con Docker Compose).
4. Probar el flujo completo: login → crear evento/inscribirse → logout → token inválido.

---

Autor: Matías Hansen
