package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time" // Para la lógica de reintentos

    _ "github.com/lib/pq" // Driver PostgreSQL
    //"github.com/gin-gonic/gin" // Descomentar si usas Gin
)

type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var db *sql.DB

func main() {
    // --- Conexión a la Base de Datos con Reintentos ---
    var err error
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    dbHost := os.Getenv("DB_HOST") // Nombre del servicio de la DB en Docker Compose
    dbPort := "5432"

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    maxRetries := 10
    for i := 0; i < maxRetries; i++ {
        db, err = sql.Open("postgres", connStr)
        if err != nil {
            log.Printf("Error abriendo conexión a la DB: %v. Reintentando en 5s...", err)
            time.Sleep(5 * time.Second)
            continue
        }
        err = db.Ping()
        if err == nil {
            log.Println("Conexión a la base de datos establecida exitosamente!")
            break // Conexión exitosa
        }
        log.Printf("Error haciendo ping a la DB: %v. Reintentando en 5s... (Intento %d/%d)", err, i+1, maxRetries)
        db.Close() // Cierra la conexión fallida antes de reintentar
        time.Sleep(5 * time.Second)
    }
    if err != nil { // Si después de los reintentos aún hay error
         log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
    }
    defer db.Close()

    // Crear tabla si no existe (¡Solo para demostración!)
    createTable()

    // --- Configuración del Servidor HTTP (usando net/http estándar) ---
    http.HandleFunc("/api/items", getItemsHandler)
    // Si necesitas más endpoints, añádelos aquí
    // http.HandleFunc("/api/items", createItemHandler) // POST

    log.Println("Servidor API Go escuchando en el puerto 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))

    /* --- Alternativa usando Gin ---
    r := gin.Default()
    // Middleware CORS (¡importante para desarrollo con React!)
    r.Use(func(c *gin.Context) {
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // ¡Sé más específico en producción!
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
            c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

            if c.Request.Method == "OPTIONS" {
                    c.AbortWithStatus(204)
                    return
            }
            c.Next()
    })

    r.GET("/api/items", getItemsGin)
    log.Println("Servidor API Gin escuchando en el puerto 8080")
    r.Run(":8080")
    */
}

func createTable() {
    // Crear tabla (simplificado, sin manejo de errores robusto)
    // En una app real, usa migraciones.
    query := `
    CREATE TABLE IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL
    );
    -- Insertar datos de ejemplo si la tabla está vacía
    INSERT INTO items (name) SELECT 'Item Ejemplo 1' WHERE NOT EXISTS (SELECT 1 FROM items);
    INSERT INTO items (name) SELECT 'Item Ejemplo 2' WHERE NOT EXISTS (SELECT 1 FROM items);
    `
    _, err := db.Exec(query)
    if err != nil {
        log.Printf("Error creando/verificando tabla 'items': %v", err)
        // Decide si quieres detener la app aquí o continuar
    } else {
         log.Println("Tabla 'items' verificada/creada y datos de ejemplo insertados si era necesario.")
    }
}


// --- Handler con net/http ---
func getItemsHandler(w http.ResponseWriter, r *http.Request) {
    // Configurar CORS (necesario si React y API se sirven desde diferentes orígenes/puertos en desarrollo)
    w.Header().Set("Access-Control-Allow-Origin", "*") // ¡Sé más específico en producción!
    w.Header().Set("Content-Type", "application/json")

     if r.Method != http.MethodGet {
            http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
            return
     }

    rows, err := db.Query("SELECT id, name FROM items ORDER BY id ASC")
    if err != nil {
        log.Printf("Error al consultar items: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    items := []Item{}
    for rows.Next() {
        var item Item
        if err := rows.Scan(&item.ID, &item.Name); err != nil {
            log.Printf("Error al escanear fila: %v", err)
            http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
            return
        }
        items = append(items, item)
    }

    if err := rows.Err(); err != nil { // Chequear errores después del bucle
        log.Printf("Error después de iterar filas: %v", err)
        http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(items)
}


/* --- Handler con Gin (alternativa) ---
func getItemsGin(c *gin.Context) {
    rows, err := db.Query("SELECT id, name FROM items ORDER BY id ASC")
    if err != nil {
        log.Printf("Error al consultar items: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
        return
    }
    defer rows.Close()

    items := []Item{}
    for rows.Next() {
        var item Item
        if err := rows.Scan(&item.ID, &item.Name); err != nil {
            log.Printf("Error al escanear fila: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
            return
        }
        items = append(items, item)
    }
     if err := rows.Err(); err != nil { // Chequear errores después del bucle
         log.Printf("Error después de iterar filas: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
         return
    }

    c.JSON(http.StatusOK, items)
}
*/
