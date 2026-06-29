# Sistema de Gestión de E-Books - Proyecto Evaluación con el Docente

Este proyecto es una aplicación web para la administración, compra y descarga de libros electrónicos o Ebooks, desarrollado como parte del proyecto que hemos estado realizando durante el curso.

- **Fecha:** 28 de junio de 2026 (última actualización)
- **Estudiante:** Miguel David Martillo González

---

## ¿Por qué es un "Monolito Híbrido"?

Este sistema está diseñado bajo el concepto de **Monolito Híbrido** por dos razones principales:

1. **Es un Monolito:** Todo el sistema tanto las pantallas visuales HTML como la lógica de compra de libros, el manejo de la base de datos y la seguridad, viven y se ejecutan dentro de un único proyecto de software unificado. Esto facilita el desarrollo, las pruebas y la puesta en marcha.
2. **Es Híbrido:** El servidor realiza dos funciones diferentes de forma simultánea:
   - **Flujo Web Tradicional:** Genera y envía directamente las pantallas visuales HTML dinámico que el usuario ve en su navegador utilizando plantillas de Go, estilizadas de forma atractiva con TailwindCSS y DaisyUI.
   - **Servicios Web de Datos (API JSON):** Expone 8 endpoints independientes (bajo la ruta `/api/v1/...`) que envían y reciben datos en formato JSON. Esto permite que el mismo sistema sea compatible con navegadores web tradicionales y, al mismo tiempo, con aplicaciones móviles, programas externos o herramientas de prueba como Postman.

---

## Funcionalidades Principales

El proyecto se divide en cuatro módulos principales:

- **Autenticación y Seguridad:** Permite registrar usuarios, iniciar sesión y cerrar sesión de manera segura. El sistema es flexible: reconoce el inicio de sesión a través de las cookies del navegador web o mediante tokens de autorización Bearer enviados en cabeceras HTTP por clientes de API externos.
- **Catálogo de Libros:** Permite registrar nuevos libros electrónicos con validaciones de datos para evitar información incompleta o precios negativos y listar el catálogo disponible.
- **Perfil de Usuario y Finanzas:** Muestra la información del usuario, sus libros comprados y le permite recargar su saldo directamente.
- **Transacciones de Compra:** Permite comprar libros electrónicos, verificando que el usuario tenga saldo suficiente, evitando que compre el mismo libro dos veces de manera simultánea y permitiendo la descarga segura del contenido del libro.

---

## Documentación de las Funciones Principales (Capa de Servicios)

La lógica de negocio del sistema está organizada en tres servicios principales dentro de la capa `internal/services`. A continuación se detallan sus funciones y responsabilidades:

### 1. Servicio de Usuarios ([user_service.go](internal/services/user_service.go))

Este servicio gestiona el registro, autenticación y finanzas de los usuarios:

*   **`Register(ctx, username, email, password)`**: Registra a un nuevo usuario en el sistema. Valida que todos los campos estén llenos, verifica que el correo electrónico no esté duplicado en la base de datos y encripta de forma segura la contraseña usando **Bcrypt** antes de almacenarla.
*   **`Login(ctx, email, password)`**: Autentica al usuario en el sistema. Busca al usuario por su correo electrónico y compara la contraseña ingresada con el hash de contraseña almacenado usando la verificación de Bcrypt.
*   **`AddBalance(ctx, userID, amount)`**: Permite recargar el monedero digital del usuario. Verifica que el monto a recargar sea mayor a cero e incrementa el saldo total.
*   **`GetUserByID(ctx, id)`**: Obtiene toda la información de perfil de un usuario específico a partir de su identificador único.

### 2. Servicio de Libros ([book_service.go](internal/services/book_service.go))

Este servicio controla el catálogo de libros electrónicos:

*   **`AddBook(ctx, book)`**: Registra un nuevo libro electrónico en el catálogo general. Implementa validaciones estrictas: el título y autor no pueden estar vacíos, el año de publicación debe ser válido (mayor a cero) y el precio de venta no puede ser negativo.
*   **`ListBooks(ctx)`**: Recupera y lista todos los libros disponibles en el catálogo de la tienda para su visualización.
*   **`GetBookByID(ctx, id)`**: Busca y recupera la información detallada de un libro en particular utilizando su ID.

### 3. Servicio de Compras ([purchase_service.go](internal/services/purchase_service.go))

Este servicio administra las transacciones financieras y la entrega del material digital:

*   **`BuyBook(ctx, userID, bookID)`**: Procesa de forma segura la compra de un libro electrónico por parte de un usuario.
    1.  Verifica que el usuario no posea el libro previamente para evitar compras duplicadas.
    2.  Comprueba la existencia del libro y del usuario en el sistema.
    3.  Valida si el saldo del usuario es suficiente para cubrir el precio del libro.
    4.  Descuenta el saldo del usuario y registra la transacción de compra de forma persistente. *Si ocurre un error en el registro de la compra, se revierte el descuento del saldo.*
*   **`GetUserPurchases(ctx, userID)`**: Lista detalladamente todos los libros electrónicos que han sido adquiridos exitosamente por un usuario en particular.
*   **`DownloadBook(ctx, userID, bookID)`**: Permite al usuario descargar el contenido del libro. Verifica rigurosamente que el usuario haya comprado previamente el libro e incrementa el contador de descargas asociado a esa adquisición antes de entregar los datos.

---

## Tecnologías y Paquetes Utilizados

El sistema está construido sobre las siguientes tecnologías y librerías externas:

- **Go o Golang:** Es el lenguaje de programación principal del proyecto, elegido por su velocidad, eficiencia y simplicidad.
- **PostgreSQL:** Es el motor de base de datos relacional para almacenar de forma segura la información de los usuarios, libros y transacciones de compra. Se encuentra dockerizado para facilitar su inicialización rápida.
- **TailwindCSS y DaisyUI:** Son frameworks de diseño utilizados para crear una interfaz de usuario atractiva, moderna y adaptada a dispositivos móviles en la parte web tradicional.

### Paquetes Externos de Go (Librerías):

- **`github.com/go-chi/chi/v5`:** Es un enrutador de peticiones muy ligero. Lo usamos para organizar de manera ordenada las rutas web y las rutas de la API del servidor.
- **`github.com/jackc/pgx/v5`:** Es el controlador oficial y eficiente para conectar Go con la base de datos PostgreSQL. Se encarga de enviar las consultas SQL y gestionar la conexión.
- **`github.com/joho/godotenv`:** Una librería para cargar variables de configuración desde un archivo `.env`. Nos sirve para no exponer contraseñas ni puertos de base de datos directamente en el código.
- **`golang.org/x/crypto/bcrypt`:** Paquete utilizado para encriptar y proteger las contraseñas de los usuarios antes de guardarlas en la base de datos, garantizando la seguridad de la información.

---

## Cómo Iniciar el Proyecto

### 1. Requisitos Previos

- Tener instalado **Go** versión 1.26 o superior.
- Tener instalado **Docker** para la base de datos.

### 2. Configuración

1. Copia el archivo de plantilla `.env-template` a un nuevo archivo llamado `.env` en la raíz del proyecto.
2. Configura los accesos a tu base de datos PostgreSQL si es necesario.

### 3. Levantar la Base de Datos

Ejecuta el siguiente comando para levantar PostgreSQL de forma automática mediante Docker:

```bash
docker-compose up -d
```

_(Esto creará la base de datos e importará automáticamente las tablas del archivo `init.sql`)_.

### 4. Iniciar el Servidor

Ejecuta el servidor con el siguiente comando en la raíz del proyecto:

```bash
go run ./cmd/api/main.go
```

El servidor estará escuchando en `http://localhost:8080`.

- Puedes acceder desde tu navegador para usar la interfaz visual clásica.
- Puedes usar Postman para interactuar con los servicios JSON bajo `http://localhost:8080/api/v1`.
