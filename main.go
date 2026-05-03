package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	conectarDB()

	// Productos 
	http.HandleFunc("/productos", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductos(w, r)
		case http.MethodPost:
			crearProducto(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/productos/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductoPorID(w, r)
		case http.MethodPut:
			actualizarProducto(w, r)
		case http.MethodDelete:
			eliminarProducto(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Clientes 
	http.HandleFunc("/clientes", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getClientes(w, r)
		case http.MethodPost:
			crearCliente(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/clientes/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getClientePorID(w, r)
		case http.MethodPut:
			actualizarCliente(w, r)
		case http.MethodDelete:
			eliminarCliente(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Empleados 
	http.HandleFunc("/empleados", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getEmpleados(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/empleados/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getEmpleadoPorID(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Proveedores 
	http.HandleFunc("/proveedores", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProveedores(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/proveedores/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProveedorPorID(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Categorias 
	http.HandleFunc("/categorias", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCategorias(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/categorias/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCategoriaPorID(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Compras 
	http.HandleFunc("/compras", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCompras(w, r)
		case http.MethodPost:
			crearCompra(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/compras/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCompraPorID(w, r)
		case http.MethodPut:
			cancelarCompra(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Detalles de Compra 
	http.HandleFunc("/detalle-compra", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDetalleCompras(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/detalle-compra/detalle", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDetalleCompraPorID(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	// Reportes 
	http.HandleFunc("/reportes/auditoria-ventas", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAuditoriaVentas(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/reportes/rentabilidad-productos", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getRentabilidadProductos(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/reportes/stock-critico", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getControlStock(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	http.HandleFunc("/reportes/desempeno-empleados", habilitarCORS(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDesempenoEmpleados(w, r)
		default:
			RespondJSON(w, http.StatusMethodNotAllowed, "Metodo no permitido", nil)
		}
	}))

	fmt.Println("Servidor corriendo en puerto 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar servidor:", err)
	}
}

// habilitarCORS permite peticiones desde el frontend de Vite
func habilitarCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}