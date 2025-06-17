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

## Microservicio LDAP

Dentro del directorio `ldapservice` se encuentra un servicio en Python que se conecta a un servidor de Active Directory usando LDAP y sincroniza la informaci\u00f3n b\u00e1sica de los usuarios en una base de datos PostgreSQL.

Para instalar sus dependencias se puede usar `pip`:

```bash
pip install -r ldapservice/requirements.txt
```

El servicio utiliza variables de entorno para la conexi\u00f3n:

- `LDAP_SERVER`: URL del servidor LDAP.
- `LDAP_USER` y `LDAP_PASSWORD`: credenciales de conexi\u00f3n.
- `LDAP_BASE_DN`: DN base para la b\u00fasqueda.
- `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGPASSWORD`: datos de acceso a PostgreSQL.

Ejemplo de ejecuci\u00f3n:

```bash
export LDAP_SERVER=ldap://ldap.example.com
export LDAP_USER=binduser@example.com
export LDAP_PASSWORD=secret
export LDAP_BASE_DN="dc=example,dc=com"
export PGHOST=localhost
export PGDATABASE=agendador
python ldapservice/ldap_service.py
```

El script crear\u00e1 la tabla `users` si no existe e insertar\u00e1 o actualizar\u00e1 los registros obtenidos del LDAP.
