package main 

import (
    "encoding/json"
    "net/http"
)

// Estructura estandar para todas las respuestas del API.
type Respuesta struct {
    Status  int         `json:"status"`        
    Mensaje string      `json:"mensaje"`        
    Data    interface{} `json:"data,omitempty"` 
}

// RespondJSON escribe una respuesta JSON estructurada en el ResponseWriter.
func RespondJSON(w http.ResponseWriter, status int, mensaje string, data interface{}) {

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
	
    json.NewEncoder(w).Encode(Respuesta{
        Status:  status,
        Mensaje: mensaje,
        Data:    data,
    })
}