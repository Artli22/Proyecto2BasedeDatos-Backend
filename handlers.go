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

// Handler para traer todos los clientes 
func getClientes(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_cliente, nombre, telefono, correo, 
		 activo
		FROM cliente WHERE activo = TRUE
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar clientes en la base de datos", nil)
		return
	}
	defer rows.Close()

	clientes := []Cliente{}
	for rows.Next() {
		var c Cliente
		err := rows.Scan(
			&c.id_cliente, &c.nombre, &c.telefono, &c.correo, &c.activo,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de cliente", nil)
			return
		}
		clientes = append(clientes, c)
	}

	RespondJSON(w, http.StatusOK, "Clientes obtenidos correctamente", clientes)
}

// Handler para traer un cliente por ID 
func getClientePorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de cliente", nil)
		return
	}

	var c Cliente
	err := DB.QueryRow(`
		SELECT id_cliente, nombre, telefono, correo, activo
		FROM cliente WHERE id_cliente = $1 AND activo = TRUE
	`, idStr).Scan(
		&c.IDCliente, &c.Nombre, &c.Telefono,
		&c.Correo, &c.Activo,
	)

	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound,
			"Cliente no encontrado", nil)
		return
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar cliente", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Cliente obtenido correctamente", c)
}

// Handler para actualizar un cliente 
func actualizarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de cliente", nil)
		return
	}

	var c Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	result, err := DB.Exec(`
		UPDATE cliente 
		SET nombre=$1, telefono=$2, correo=$3, activo=$4
		WHERE id_cliente=$5 AND activo = TRUE
	`,
		c.Nombre, c.Telefono, c.Correo, c.Activo, idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al actualizar cliente en la base de datos", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Cliente no encontrado o se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Cliente actualizado correctamente", nil)
}

// Handler para eliminar un cliente (desactivarlo)
func eliminarCliente(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de cliente", nil)
		return
	}

	result, err := DB.Exec(
		"UPDATE cliente SET activo = FALSE WHERE id_cliente = $1 AND activo = TRUE",
		idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al desactivar cliente", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Cliente no encontrado o ya se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Cliente desactivado correctamente", nil)
}

// Handler para traer todos los empleados
func getEmpleados(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_empleado, nombre, telefono, correo, 
		 activo
		FROM empleado WHERE activo = TRUE
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar empleados en la base de datos", nil)
		return
	}
	defer rows.Close()

	empleados := []Empleado{}
	for rows.Next() {
		var e Empleado
		err := rows.Scan(
			&e.id_empleado, &e.nombre, &e.telefono, &e.correo, &e.activo,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de empleado", nil)
			return
		}
		empleados = append(empleados, e)
	}

	RespondJSON(w, http.StatusOK, "Empleados obtenidos correctamente", empleados)
}

// Handler para traer un empleado por ID 
func getEmpleadoPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de empleado", nil)
		return
	}

	var e Empleado
	err := DB.QueryRow(`
		SELECT id_empleado, nombre, telefono, correo, activo
		FROM empleado WHERE id_empleado = $1 AND activo = TRUE
	`, idStr).Scan(
		&e.IDEmpleado, &e.Nombre, &e.Telefono,
		&e.Correo, &e.Activo,
	)

	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound,
			"Empleado no encontrado", nil)
		return
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar empleado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Empleado obtenido correctamente", e)
}

// Handler para actualizar un empleado
func actualizarEmpleado(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de empleado", nil)
		return
	}

	var e Empleado
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	result, err := DB.Exec(`
		UPDATE empleado 
		SET nombre=$1, telefono=$2, correo=$3, activo=$4
		WHERE id_empleado=$5 AND activo = TRUE
	`,
		e.Nombre, e.Telefono, e.Correo, e.Activo, idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al actualizar empleado en la base de datos", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Empleado no encontrado o se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Empleado actualizado correctamente", nil)
}

// Handler para eliminar un empleado (desactivarlo)
func eliminarEmpleado(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de empleado", nil)
		return
	}

	result, err := DB.Exec(
		"UPDATE empleado SET activo = FALSE WHERE id_empleado = $1 AND activo = TRUE",
		idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al desactivar empleado", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Empleado no encontrado o ya se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Empleado desactivado correctamente", nil)
}

// Handler para traer todos los proveedores 
func getProveedores(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_proveedor, nombre, telefono, correo, 
		 activo
		FROM proveedor WHERE activo = TRUE
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar proveedores en la base de datos", nil)
		return
	}
	defer rows.Close()

	proveedores := []Proveedor{}
	for rows.Next() {
		var p Proveedor
		err := rows.Scan(
			&p.id_proveedor, &p.nombre, &p.telefono, &p.correo, &p.activo,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de proveedor", nil)
			return
		}
		proveedores = append(proveedores, p)
	}

	RespondJSON(w, http.StatusOK, "Proveedores obtenidos correctamente", proveedores)
}

// Handler para traer un proveedor por ID
func getProveedorPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de proveedor", nil)
		return
	}

	var p Proveedor
	err := DB.QueryRow(`
		SELECT id_proveedor, nombre, telefono, correo, activo
		FROM proveedor WHERE id_proveedor = $1 AND activo = TRUE
	`, idStr).Scan(
		&p.IDProveedor, &p.Nombre, &p.Telefono,
		&p.Correo, &p.Activo,
	)

	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound,
			"Proveedor no encontrado", nil)
		return
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar proveedor", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Proveedor obtenido correctamente", p)
}

// Handler para actualizar un proveedor
func actualizarProveedor(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de proveedor", nil)
		return
	}

	var p Proveedor
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	result, err := DB.Exec(`
		UPDATE proveedor 
		SET nombre=$1, telefono=$2, correo=$3, activo=$4
		WHERE id_proveedor=$5 AND activo = TRUE
	`,
		p.Nombre, p.Telefono, p.Correo, p.Activo, idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al actualizar proveedor en la base de datos", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Proveedor no encontrado o se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Proveedor actualizado correctamente", nil)
}

// Handler para eliminar un proveedor (desactivarlo)
func eliminarProveedor(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de proveedor", nil)
		return
	}

	result, err := DB.Exec(
		"UPDATE proveedor SET activo = FALSE WHERE id_proveedor = $1 AND activo = TRUE",
		idStr,
	)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al desactivar proveedor", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Proveedor no encontrado o ya se encuentra desactivado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Proveedor desactivado correctamente", nil)
}
