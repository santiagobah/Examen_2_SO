package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time" 

    _ "github.com/lib/pq"
)

type Item struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var db *sql.DB

func main() {
    var err error
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    dbHost := os.Getenv("DB_HOST") 
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
            break 
        }
        log.Printf("Error haciendo ping a la DB: %v. Reintentando en 5s... (Intento %d/%d)", err, i+1, maxRetries)
        db.Close() 
        time.Sleep(5 * time.Second)
    }
    if err != nil { 
         log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
    }
    defer db.Close()

    createTable()

    http.HandleFunc("/api/items", getItemsHandler)

    log.Println("Servidor API Go escuchando en el puerto 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
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
