Mejores Prácticas de Go (Agente Antigravity)
Sintaxis y Estilo
Aplicar gofmt para imponer un estilo estándar y universal de alineación sin necesidad de debatir sobre el formato
.
Utilizar MayúsculasIntercaladas (CamelCase) en nombres de variables y funciones, evitando el uso de guiones bajos
.
Omitir prefijos como Obt (o Get) en funciones captadoras; por ejemplo, nombrar al método Propietario() en lugar de ObtPropietario()
.
Preferir la declaración corta := para declarar e inicializar variables locales dentro de las funciones
.
Manejo de Errores y Control de Flujo
Verificar siempre los errores retornados inmediatamente después de invocar la función
.
Jamás utilizar el identificador en blanco _ para descartar o ignorar mensajes de error
.
Eliminar bloques else redundantes si la instrucción if previa termina en return, break o continue
.
f, err := os.Open(nombre)
if err != nil {
    return err
}
d, err := f.Stat()
if err != nil {
    f.Close()
    return err
}
Concurrencia y Sincronización
Aplicar el principio de comunicación de procesos secuenciales (CSP): no comunicarse compartiendo memoria, sino compartir memoria comunicándose a través de canales
.
Utilizar siempre el detector de carreras (go run -race o go test -race) durante el desarrollo para identificar accesos concurrentes inseguros
.
Evitar fugas de goroutines garantizando el cierre de canales y utilizando el paquete context para manejar tiempos de espera o cancelaciones
.
Programar la limpieza de recursos (cerrar archivos, liberar mutexes) utilizando la instrucción defer inmediatamente después de su adquisición
.
func hazloConSeguridad(trabajo *Trabajo) {
    defer func() {
        if err := recover(); err != nil {
            log.Println(err)
        }
    }()
    haz(trabajo)
}
Optimización de Memoria
Priorizar la asignación de variables en el Stack (pila) evaluando el código con go build -gcflags="-m"
.
Minimizar las variables que escapan al Heap (montón) para reducir la carga de trabajo del recolector de basura (GC)
.
Utilizar sync.Pool para reciclar objetos de corta duración que se asignan y liberan con mucha frecuencia
.
Ciberseguridad
Validar y desinfectar rigurosamente todas las entradas de usuario para bloquear inyecciones SQL y ataques XSS
.
Usar siempre consultas parametrizadas con database/sql y emplear el paquete html/template para el escape automático de caracteres
.
Nunca exponer información de depuración, errores internos o stack traces en entornos de producción
.
Ejecutar la herramienta govulncheck de manera rutinaria para rastrear vulnerabilidades conocidas en las dependencias
.