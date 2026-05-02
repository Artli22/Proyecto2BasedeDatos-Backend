package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Respuesta struct {
	Status  int         `json:"status"`
	Mensaje string      `json:"mensaje"`
	Data    interface{} `json:"data,omitempty"`
}


const (
	MsgIDFaltante            = "Falta el parametro id"
	MsgJSONInvalido          = "El cuerpo del request no es un JSON valido"
	MsgNoEncontrado          = "no encontrado"
	MsgYaDesactivado         = "ya se encuentra desactivado"
	MsgDesactivado           = "se encuentra desactivado"
	MsgErrorConsulta         = "Error al consultar"
	MsgErrorLectura          = "Error al leer fila de"
	MsgErrorInsertar         = "Error al insertar"
	MsgErrorActualizar       = "Error al actualizar"
	MsgErrorDesactivar       = "Error al desactivar"
	MsgErrorTransaccion      = "Error al confirmar transaccion"
	MsgErrorRollback         = "se hizo rollback"

	MsgObtenidoCorrectamente = "obtenido correctamente"
	MsgObtenidosCorrectamente = "obtenidos correctamente"
	MsgCreadoCorrectamente   = "creado correctamente"
	MsgActualizadoCorrectamente = "actualizado correctamente"
	MsgDesactivadoCorrectamente = "desactivado correctamente"
)

// Escribe una respuesta JSON en el ResponseWriter.
func RespondJSON(w http.ResponseWriter, status int, mensaje string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Respuesta{
		Status:  status,
		Mensaje: mensaje,
		Data:    data,
	})
}

// Obtiene y valida el parametro "id" 
func ValidarIDParametro(r *http.Request, w http.ResponseWriter, recurso string) (string, bool) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			fmt.Sprintf("%s de %s", MsgIDFaltante, recurso), nil)
		return "", false
	}
	return idStr, true
}

// valida si el JSON es válido al decodificarlo
func ValidarJSONDecodificacion(err error, w http.ResponseWriter) bool {
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, MsgJSONInvalido, nil)
		return false
	}
	return true
}


// Manejador de errores de consultas SELECT
func ManejarErrorConsulta(err error, w http.ResponseWriter, recurso string) bool {
	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound,
			fmt.Sprintf("%s %s", recurso, MsgNoEncontrado), nil)
		return true
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			fmt.Sprintf("%s %s", MsgErrorConsulta, recurso), nil)
		return true
	}
	return false
}

// Manejador de errores de INSERT/UPDATE/DELETE
func ManejarErrorInsertActualizar(err error, w http.ResponseWriter, operacion, recurso string) bool {
	if err != nil {
		mensaje := fmt.Sprintf("%s %s %s", MsgErrorActualizar, recurso, "en la base de datos")
		if operacion == "insert" {
			mensaje = fmt.Sprintf("%s %s", MsgErrorInsertar, recurso)
		} else if operacion == "delete" {
			mensaje = fmt.Sprintf("%s %s", MsgErrorDesactivar, recurso)
		}
		RespondJSON(w, http.StatusInternalServerError, mensaje, nil)
		return true
	}
	return false
}

// Valida si la operación afectó registros en las filas de las entidades
func ValidarFilasAfectadas(result sql.Result, w http.ResponseWriter, recurso string) bool {
	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			fmt.Sprintf("%s %s o %s %s", recurso, MsgNoEncontrado, MsgDesactivado, recurso), nil)
		return false
	}
	return true
}

// Ejecuta una función dentro de una transacción
func EjecutarEnTransaccion(w http.ResponseWriter, fn func(*sql.Tx) (interface{}, error)) (interface{}, bool) {
	tx, err := DB.Begin()
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al iniciar transaccion", nil)
		return nil, false
	}

	resultado, err := fn(tx)
	if err != nil {
		tx.Rollback()
		RespondJSON(w, http.StatusInternalServerError,
			fmt.Sprintf("%v %s", err, MsgErrorRollback), nil)
		return nil, false
	}

	if err = tx.Commit(); err != nil {
		RespondJSON(w, http.StatusInternalServerError, MsgErrorTransaccion, nil)
		return nil, false
	}

	return resultado, true
}

// Obtiencion de un producto por ID
func ObtenerProductoPorID(idStr string) (*Producto, error) {
	var p Producto
	err := DB.QueryRow(`
		SELECT id_producto, nombre, descripcion, precio_actual,
		fecha_vencimiento, imagen, stock, activo, id_categoria, id_proveedor
		FROM producto WHERE id_producto = $1 AND activo = TRUE
	`, idStr).Scan(
		&p.IDProducto, &p.Nombre, &p.Descripcion,
		&p.PrecioActual, &p.FechaVencimiento, &p.Imagen,
		&p.Stock, &p.Activo, &p.IDCategoria, &p.IDProveedor,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Obtiencion del precio actual de un producto
func ObtenerPrecioProducto(tx *sql.Tx, idProducto int) (float64, int, error) {
	var precio float64
	var stock int 

	err := tx.QueryRow(
		"SELECT precio_actual, stock FROM producto WHERE id_producto = $1 AND activo = TRUE",
		idProducto,
	).Scan(&precio, &stock)
	return precio, stock, err
}

// Obtiencion de un cliente por ID
func ObtenerClientePorID(idStr string) (*Cliente, error) {
	var c Cliente
	err := DB.QueryRow(`
		SELECT id_cliente, nombre, telefono, correo, activo
		FROM cliente WHERE id_cliente = $1 AND activo = TRUE
	`, idStr).Scan(
		&c.IdCliente, &c.Nombre, &c.Telefono, &c.Correo, &c.Activo,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Obtiencion de un empleado por ID
func ObtenerEmpleadoPorID(idStr string) (*Empleado, error) {
	var e Empleado
	err := DB.QueryRow(`
		SELECT id_empleado, nombre, telefono, correo, activo
		FROM empleado WHERE id_empleado = $1 AND activo = TRUE
	`, idStr).Scan(
		&e.IdEmpleado, &e.Nombre, &e.Telefono, &e.Correo, &e.Activo,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Obtiencion de un proveedor por ID
func ObtenerProveedorPorID(idStr string) (*Proveedor, error) {
	var prov Proveedor
	err := DB.QueryRow(`
		SELECT id_proveedor, nombre, telefono, correo, activo
		FROM proveedor WHERE id_proveedor = $1 AND activo = TRUE
	`, idStr).Scan(
		&prov.IDProveedor, &prov.Nombre, &prov.Telefono, &prov.Correo, &prov.Activo,
	)
	if err != nil {
		return nil, err
	}
	return &prov, nil
}

// Obtiencion de una categoría por ID
func ObtenerCategoriaPorID(idStr string) (*Categoria, error) {
	var c Categoria
	err := DB.QueryRow(`
		SELECT id_categoria, nombre
		FROM categoria WHERE id_categoria = $1
	`, idStr).Scan(&c.IdCategoria, &c.Nombre)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Genera un número de factura único basado en la fecha/hora
func generarNumFactura() string {
	return fmt.Sprintf("FAC-%s", time.Now().Format("20060102-150405"))
}

// Helper para procesar compra dentro de transacción
func procesarNuevaCompra(tx *sql.Tx, req CompraRequest) (interface{}, error) {
	var total float64
	var detalles []DetalleTemp

	if len(req.Productos) == 0 {
	 return nil, fmt.Errorf("La compra debe incluir al menos un producto")
	}

	for _, item := range req.Productos {
		if item.Cantidad <= 0 {
			return nil, fmt.Errorf(
				"La cantidad del producto %d debe ser mayor a 0",
				item.IDProducto,
			)
		}

		precio, stock, err := ObtenerPrecioProducto(tx, item.IDProducto)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(
				"Producto %d no encontrado o inactivo",
				item.IDProducto,
			)
		}

		if err != nil {
			return nil, fmt.Errorf(
				"Error al consultar producto %d",
				item.IDProducto,
			)
		}

		if stock < item.Cantidad {
			return nil, fmt.Errorf(
				"Stock insuficiente para producto %d. Disponible: %d, solicitado: %d",
				item.IDProducto,
				stock,
				item.Cantidad,
			)
		}

		detalles = append(detalles, DetalleTemp{
			IDProducto:     item.IDProducto,
			Cantidad:       item.Cantidad,
			PrecioUnitario: precio,
			SubTotal:       precio * float64(item.Cantidad),
		})
		total += precio * float64(item.Cantidad)
	}
	
	var idCompra int
	err := tx.QueryRow(`
		INSERT INTO compra (fecha, total, metodo_pago, estado, num_factura, id_cliente, id_empleado)
		VALUES ($1, $2, $3, 'completado', $4, $5, $6)
		RETURNING id_compra
	`,
		req.Fecha, total, req.MetodoPago,
		generarNumFactura(),
		req.IDCliente, req.IDEmpleado,
	).Scan(&idCompra)

	if err != nil {
		return nil, fmt.Errorf("Error al registrar compra")
	}

	for _, d := range detalles {
		_, err = tx.Exec(`
			INSERT INTO detalle_compra (id_compra, id_producto, cantidad, precio_unitario, sub_total)
			VALUES ($1, $2, $3, $4, $5)
		`, idCompra, d.IDProducto, d.Cantidad, d.PrecioUnitario, d.SubTotal)

		if err != nil {
			return nil, fmt.Errorf("Error al registrar detalle de compra")
		}

		result, err := tx.Exec(`
			UPDATE producto SET stock = stock - $1
			WHERE id_producto = $2
			 AND activo = TRUE
	  		 AND stock >= $1
		`, d.Cantidad, d.IDProducto)

		if err != nil {
			return nil, fmt.Errorf("Error al actualizar stock")
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return nil, fmt.Errorf(
				"Stock insuficiente al actualizar producto %d",
				d.IDProducto,
			)
		}
		
	}

	return map[string]interface{}{
		"id_compra": idCompra,
		"total":     total,
	}, nil
}

// Helper para procesar cancelación de compra dentro de transacción
func procesarCancelacionCompra(tx *sql.Tx, idStr string) (interface{}, error) {
	result, err := tx.Exec(`
		UPDATE compra 
		SET estado = 'cancelado'
		WHERE id_compra = $1 
		  AND estado = 'completado'
	`, idStr)

	if err != nil {
		return nil, fmt.Errorf("Error al cancelar compra")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("Compra no encontrada o ya cancelada")
	}

	rows, err := tx.Query(`
		SELECT id_producto, cantidad
		FROM detalle_compra
		WHERE id_compra = $1
	`, idStr)

	if err != nil {
		return nil, fmt.Errorf("Error al consultar detalle de compra")
	}
	defer rows.Close()

	productosRestaurados := 0

	for rows.Next() {
		var idProducto int
		var cantidad int

		if err := rows.Scan(&idProducto, &cantidad); err != nil {
			return nil, fmt.Errorf("Error al leer detalle de compra")
		}

		_, err = tx.Exec(`
			UPDATE producto
			SET stock = stock + $1
			WHERE id_producto = $2
		`, cantidad, idProducto)

		if err != nil {
			return nil, fmt.Errorf("Error al restaurar stock del producto %d", idProducto)
		}

		productosRestaurados++
	}

	return map[string]interface{}{
		"id_compra":              idStr,
		"productos_restaurados":  productosRestaurados,
		"estado":                 "cancelado",
	}, nil
}