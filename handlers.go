// handlers.go
package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// Handler para traer todos los productos
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

// Handler para traer un producto por ID
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

// Handler para eliminar un producto (desactivarlo)
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
		SELECT id_cliente, nombre, telefono, correo, activo
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
			&c.Id_cliente, &c.Nombre, &c.Telefono, &c.Correo, &c.Activo,
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
		&c.Id_cliente, &c.Nombre, &c.Telefono,
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

// Handler para crear un cliente
func crearCliente(w http.ResponseWriter, r *http.Request) {
	var c Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	err := DB.QueryRow(`
		INSERT INTO cliente (nombre, telefono, correo)
		VALUES ($1, $2, $3)
		RETURNING id_cliente
	`,
		c.Nombre, c.Telefono, c.Correo,
	).Scan(&c.Id_cliente)

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al insertar cliente", nil)
		return
	}

	RespondJSON(w, http.StatusCreated, "Cliente creado correctamente", c)
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
		SET nombre=$1, telefono=$2, correo=$3
		WHERE id_cliente=$4 AND activo = TRUE
	`,
		c.Nombre, c.Telefono, c.Correo, idStr,
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
		SELECT id_empleado, nombre, telefono, correo, activo
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
			&e.Id_empleado, &e.Nombre, &e.Telefono, &e.Correo, &e.Activo,
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
		&e.Id_empleado, &e.Nombre, &e.Telefono,
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

// Handler para crear un empleado
func crearEmpleado(w http.ResponseWriter, r *http.Request) {
	var e Empleado
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	err := DB.QueryRow(`
		INSERT INTO empleado (nombre, telefono, correo)
		VALUES ($1, $2, $3)
		RETURNING id_empleado
	`,
		e.Nombre, e.Telefono, e.Correo,
	).Scan(&e.Id_empleado)

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al insertar empleado", nil)
		return
	}

	RespondJSON(w, http.StatusCreated, "Empleado creado correctamente", e)
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
		SET nombre=$1, telefono=$2, correo=$3
		WHERE id_empleado=$4 AND activo = TRUE
	`,
		e.Nombre, e.Telefono, e.Correo, idStr,
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
		SELECT id_proveedor, nombre, telefono, correo, activo
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
		var prov Proveedor
		err := rows.Scan(
			&prov.IDProveedor, &prov.Nombre, &prov.Telefono, &prov.Correo, &prov.Activo,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de proveedor", nil)
			return
		}
		proveedores = append(proveedores, prov)
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

	var prov Proveedor
	err := DB.QueryRow(`
		SELECT id_proveedor, nombre, telefono, correo, activo
		FROM proveedor WHERE id_proveedor = $1 AND activo = TRUE
	`, idStr).Scan(
		&prov.IDProveedor, &prov.Nombre, &prov.Telefono,
		&prov.Correo, &prov.Activo,
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

	RespondJSON(w, http.StatusOK, "Proveedor obtenido correctamente", prov)
}

// Handler para crear un proveedor
func crearProveedor(w http.ResponseWriter, r *http.Request) {
	var prov Proveedor
	if err := json.NewDecoder(r.Body).Decode(&prov); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	err := DB.QueryRow(`
		INSERT INTO proveedor (nombre, telefono, correo)
		VALUES ($1, $2, $3)
		RETURNING id_proveedor
	`,
		prov.Nombre, prov.Telefono, prov.Correo,
	).Scan(&prov.IDProveedor)

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al insertar proveedor", nil)
		return
	}

	RespondJSON(w, http.StatusCreated, "Proveedor creado correctamente", prov)
}

// Handler para actualizar un proveedor
func actualizarProveedor(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de proveedor", nil)
		return
	}

	var prov Proveedor
	if err := json.NewDecoder(r.Body).Decode(&prov); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	result, err := DB.Exec(`
		UPDATE proveedor 
		SET nombre=$1, telefono=$2, correo=$3
		WHERE id_proveedor=$4 AND activo = TRUE
	`,
		prov.Nombre, prov.Telefono, prov.Correo, idStr,
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

// Handler para traer todas las categorias
func getCategorias(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_categoria, nombre
		FROM categoria
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar categorias en la base de datos", nil)
		return
	}
	defer rows.Close()

	categorias := []Categoria{}
	for rows.Next() {
		var c Categoria
		err := rows.Scan(
			&c.Id_categoria, &c.Nombre,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de categoria", nil)
			return
		}
		categorias = append(categorias, c)
	}

	RespondJSON(w, http.StatusOK, "Categorias obtenidas correctamente", categorias)
}

// Handler para traer una categoria por ID
func getCategoriaPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id de categoria", nil)
		return
	}

	var c Categoria
	err := DB.QueryRow(`
		SELECT id_categoria, nombre
		FROM categoria WHERE id_categoria = $1
	`, idStr).Scan(
		&c.Id_categoria, &c.Nombre,
	)

	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound,
			"Categoria no encontrada", nil)
		return
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar categoria", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Categoria obtenida correctamente", c)
}

// Handler para obtener todas las compras
func getCompras(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_compra, fecha, total, metodo_pago, estado, num_factura, id_cliente, id_empleado
		FROM compra
		ORDER BY fecha DESC
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar compras", nil)
		return
	}
	defer rows.Close()

	compras := []Compra{}
	for rows.Next() {
		var c Compra
		err := rows.Scan(
			&c.IDCompra, &c.Fecha, &c.Total,
			&c.MetodoPago, &c.Estado, &c.NumFactura,
			&c.IDCliente, &c.IDEmpleado,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de compra", nil)
			return
		}
		compras = append(compras, c)
	}

	RespondJSON(w, http.StatusOK, "Compras obtenidas correctamente", compras)
}

// Handler para obtener compra por ID
func getCompraPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id en la URL, ejemplo: /compras/detalle?id=1", nil)
		return
	}

	type CompraDetalle struct {
		IDCompra   int    `json:"id_compra"`
		Fecha      string `json:"fecha"`
		Total      float64 `json:"total"`
		MetodoPago *string `json:"metodo_pago"`
		Estado     *string `json:"estado"`
		NumFactura string `json:"num_factura"`
		Cliente    string `json:"cliente"`
		Empleado   string `json:"empleado"`
		Productos  []ItemDetalle `json:"productos"`
	}

	var c CompraDetalle
	err := DB.QueryRow(`
		SELECT c.id_compra, c.fecha, c.total, c.metodo_pago, c.estado, c.num_factura,
		       cl.nombre AS cliente, e.nombre AS empleado
		FROM compra c
		JOIN cliente cl ON c.id_cliente = cl.id_cliente
		JOIN empleado e ON c.id_empleado = e.id_empleado
		WHERE c.id_compra = $1
	`, idStr).Scan(
		&c.IDCompra, &c.Fecha, &c.Total,
		&c.MetodoPago, &c.Estado, &c.NumFactura,
		&c.Cliente, &c.Empleado,
	)

	if err == sql.ErrNoRows {
		RespondJSON(w, http.StatusNotFound, "Compra no encontrada", nil)
		return
	}
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar compra", nil)
		return
	}

	rows, err := DB.Query(`
		SELECT p.nombre, dc.cantidad, dc.precio_unitario, dc.sub_total
		FROM detalle_compra dc
		JOIN producto p ON dc.id_producto = p.id_producto
		WHERE dc.id_compra = $1
	`, idStr)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar detalle de compra", nil)
		return
	}
	defer rows.Close()

	c.Productos = []ItemDetalle{}
	for rows.Next() {
		var item ItemDetalle
		rows.Scan(
			&item.Producto, &item.Cantidad,
			&item.PrecioUnitario, &item.SubTotal,
		)
		c.Productos = append(c.Productos, item)
	}

	RespondJSON(w, http.StatusOK, "Compra obtenida correctamente", c)
}

// Handler para crear una compra
func crearCompra(w http.ResponseWriter, r *http.Request) {
	var req CompraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondJSON(w, http.StatusBadRequest,
			"El cuerpo del request no es un JSON valido", nil)
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al iniciar transaccion", nil)
		return
	}

	type detalleTemp struct {
		idProducto     int
		cantidad       int
		precioUnitario float64
		subTotal       float64
	}

	var total float64
	var detalles []detalleTemp

	for _, item := range req.Productos {
		var precio float64
		err := tx.QueryRow(
			"SELECT precio_actual FROM producto WHERE id_producto = $1 AND activo = TRUE",
			item.IDProducto,
		).Scan(&precio)

		if err == sql.ErrNoRows {
			tx.Rollback()
			RespondJSON(w, http.StatusNotFound,
				"Producto no encontrado o inactivo", nil)
			return
		}
		if err != nil {
			tx.Rollback()
			RespondJSON(w, http.StatusInternalServerError,
				"Error al consultar precio de producto", nil)
			return
		}

		subTotal := precio * float64(item.Cantidad)
		total += subTotal
		detalles = append(detalles, detalleTemp{
			idProducto:     item.IDProducto,
			cantidad:       item.Cantidad,
			precioUnitario: precio,
			subTotal:       subTotal,
		})
	}

	var idCompra int
	err = tx.QueryRow(`
		INSERT INTO compra (fecha, total, metodo_pago, estado, num_factura, id_cliente, id_empleado)
		VALUES ($1, $2, $3, 'completado', $4, $5, $6)
		RETURNING id_compra
	`,
		req.Fecha, total, req.MetodoPago,
		generarNumFactura(),
		req.IDCliente, req.IDEmpleado,
	).Scan(&idCompra)

	if err != nil {
		tx.Rollback()
		RespondJSON(w, http.StatusInternalServerError,
			"Error al registrar compra, se hizo rollback", nil)
		return
	}

	for _, d := range detalles {
		_, err = tx.Exec(`
			INSERT INTO detalle_compra (id_compra, id_producto, cantidad, precio_unitario, sub_total)
			VALUES ($1, $2, $3, $4, $5)
		`, idCompra, d.idProducto, d.cantidad, d.precioUnitario, d.subTotal)

		if err != nil {
			tx.Rollback()
			RespondJSON(w, http.StatusInternalServerError,
				"Error al registrar detalle de compra, se hizo rollback", nil)
			return
		}

		_, err = tx.Exec(`
			UPDATE producto SET stock = stock - $1
			WHERE id_producto = $2
		`, d.cantidad, d.idProducto)

		if err != nil {
			tx.Rollback()
			RespondJSON(w, http.StatusInternalServerError,
				"Error al actualizar stock, se hizo rollback", nil)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al confirmar transaccion", nil)
		return
	}

	RespondJSON(w, http.StatusCreated, "Compra registrada correctamente",
		map[string]interface{}{
			"id_compra": idCompra,
			"total":     total,
		})
}

// Handler para cancelar una compra (valido unicamente para estado completado)
func cancelarCompra(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		RespondJSON(w, http.StatusBadRequest,
			"Falta el parametro id en la URL, ejemplo: /compras/detalle?id=1", nil)
		return
	}

	result, err := DB.Exec(`
		UPDATE compra SET estado = 'cancelado'
		WHERE id_compra = $1 AND estado = 'completado'
	`, idStr)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al cancelar compra", nil)
		return
	}

	rowsAfectadas, _ := result.RowsAffected()
	if rowsAfectadas == 0 {
		RespondJSON(w, http.StatusNotFound,
			"Compra no encontrada o ya estaba cancelada", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Compra cancelada correctamente", nil)
}

// Handler para vista de auditoria de ventas
func getAuditoriaVentas(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_compra, num_factura, fecha, metodo_pago, estado, total, cliente, correo_cliente, empleado_cajero 
		FROM vista_auditoria_completa_ventas
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar auditoria de ventas", nil)
		return
	}
	defer rows.Close()

	ventas := []AuditoriaVenta{}
	for rows.Next() {
		var v AuditoriaVenta
		err := rows.Scan(
			&v.IDCompra, &v.NumFactura, &v.Fecha,
			&v.MetodoPago, &v.Estado, &v.Total,
			&v.Cliente, &v.CorreoCliente, &v.EmpleadoCajero,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de auditoria", nil)
			return
		}
		ventas = append(ventas, v)
	}

	RespondJSON(w, http.StatusOK, "Auditoria de ventas obtenida correctamente", ventas)
}

// Handler para vista de rentabilidad de productos
func getRentabilidadProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_producto, producto, categoria, unidades_vendidas, ingresos_totales, precio_promedio_venta 
		FROM vista_rentabilidad_productos
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar rentabilidad de productos", nil)
		return
	}
	defer rows.Close()

	productos := []RentabilidadProducto{}
	for rows.Next() {
		var p RentabilidadProducto
		err := rows.Scan(
			&p.IDProducto, &p.Producto, &p.Categoria,
			&p.UnidadesVendidas, &p.IngresosTotales, &p.PrecioPromedioVenta,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de rentabilidad", nil)
			return
		}
		productos = append(productos, p)
	}

	RespondJSON(w, http.StatusOK, "Rentabilidad de productos obtenida correctamente", productos)
}

// Handler para vista de control de stock
func getControlStock(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_producto, producto, categoria, proveedor, telefono_proveedor, stock_actual, fecha_vencimiento 
		FROM vista_stock_critico
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar stock critico", nil)
		return
	}
	defer rows.Close()

	productos := []StockCritico{}
	for rows.Next() {
		var p StockCritico
		err := rows.Scan(
			&p.IDProducto, &p.Producto, &p.Categoria,
			&p.Proveedor, &p.TelefonoProveedor,
			&p.StockActual, &p.FechaVencimiento,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de stock critico", nil)
			return
		}
		productos = append(productos, p)
	}

	RespondJSON(w, http.StatusOK, "Stock critico obtenido correctamente", productos)
}

// Handler para vista de desempenio laboral
func getDesempenoEmpleados(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`
		SELECT id_empleado, empleado, total_transacciones, monto_total_vendido, ticket_promedio, ultima_venta 
		FROM vista_desempeno_empleados
	`)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar desempeno de empleados", nil)
		return
	}
	defer rows.Close()

	empleados := []DesempenoEmpleado{}
	for rows.Next() {
		var e DesempenoEmpleado
		err := rows.Scan(
			&e.IDEmpleado, &e.Empleado, &e.TotalTransacciones,
			&e.MontoTotalVendido, &e.TicketPromedio, &e.UltimaVenta,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de desempeno", nil)
			return
		}
		empleados = append(empleados, e)
	}

	RespondJSON(w, http.StatusOK, "Desempeno de empleados obtenido correctamente", empleados)
}
