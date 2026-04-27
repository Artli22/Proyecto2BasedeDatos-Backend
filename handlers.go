// handlers.go
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Handler para traer todos los producctos
func getProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_producto, nombre, descripcion, precio_actual, 
		fecha_vencimiento, imagen, stock, activo, id_categoria, id_proveedor 
		FROM producto WHERE activo = TRUE
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar productos en la base de datos", nil)
		return
	}
	defer rows.Close()

	productos := []Producto{}
	for rows.Next() {
		var p Producto
		err := rows.Scan(
			&p.IDProducto, &p.Nombre, &p.Descripcion,
			&p.PrecioActual, &p.FechaVencimiento, &p.Imagen,
			&p.Stock, &p.Activo, &p.IDCategoria, &p.IDProveedor,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de producto", nil)
			return
		}
		productos = append(productos, p)
	}

	RespondJSON(w, http.StatusOK, "Productos obtenidos correctamente", productos)
}

// Handler para traer un producto po ID 
func getProductoPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de producto", nil)
		return
	}

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

	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound,
			"Producto no encontrado", nil)
		return
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar producto", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Producto obtenido correctamente", p)
}

// Handler para crear un producto 
func crearProducto(w http.ResponseWriter, r *http.Request) {
	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	err := DB.QueryRow(`
		INSERT INTO producto (nombre, descripcion, precio_actual, fecha_vencimiento, imagen, stock, id_categoria, id_proveedor)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_producto
	`,
		p.Nombre, p.Descripcion, p.PrecioActual,
		p.FechaVencimiento, p.Imagen, p.Stock,
		p.IDCategoria, p.IDProveedor,
	).Scan(&p.IDProducto)

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al insertar producto, verifique que id_categoria e id_proveedor existan", nil)
		return
	}

	RespondJSON(w, http.StatusCreated, "Producto creado correctamente", p)
}

// Handler para actualizar un producto
func actualizarProducto(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de producto", nil)
		return
	}

	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	result, err := DB.Exec(`
		UPDATE producto 
		SET nombre=$1, descripcion=$2, precio_actual=$3, fecha_vencimiento=$4,
			imagen=$5, stock=$6, id_categoria=$7, id_proveedor=$8
		WHERE id_producto=$9 AND activo = TRUE
	`,
		p.Nombre, p.Descripcion, p.PrecioActual,
		p.FechaVencimiento, p.Imagen, p.Stock,
		p.IDCategoria, p.IDProveedor, idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al actualizar producto en la base de datos", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Producto no encontrado o se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Producto actualizado correctamente", nil)
}

// Handler para elimianr un producto (desactivarlo)
func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de producto", nil)
		return
	}

	result, err := DB.Exec(
		"UPDATE producto SET activo = FALSE WHERE id_producto = $1 AND activo = TRUE",
		idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al desactivar producto", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Producto no encontrado o ya se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Producto desactivado correctamente", nil)
}
