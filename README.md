# agendador

Una app para conectar con nuestro LDAP y generar una agenda telefónica.

## Frontend

Este repositorio incluye una implementación sencilla del frontend usando React
(a través de CDN). Se encuentra dentro del directorio `frontend` y puede
abrirse directamente en un navegador sin pasos de compilación.

Para probarlo, abra `frontend/index.html` y podrá añadir, editar y eliminar
contactos en una tabla de manera local.

## API Gateway

Se incluye un gateway sencillo escrito en Go ubicado en `gateway`. Este componente expone un punto de inicio en `:8080` con autenticación JWT básica.

Para ejecutarlo:

```bash
cd gateway
go run .
```

El endpoint `/login` acepta credenciales (admin/password) y devuelve un token JWT. Las rutas bajo `/api/` se proxéan al servicio configurado en el código (por defecto `http://localhost:8000`) y requieren el token en la cabecera `Authorization`.
