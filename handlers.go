// handlers.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	idStr, ok := ValidarIDParametro(r, w, "producto")
	if !ok {
		return
	}

	p, err := ObtenerProductoPorID(idStr)
	if ManejarErrorConsulta(err, w, "Producto") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Producto %s", MsgObtenidoCorrectamente), p)
}

// Handler para crear un producto
func crearProducto(w http.ResponseWriter, r *http.Request) {
	var p Producto
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&p), w) {
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

	if ManejarErrorInsertActualizar(err, w, "insert", "producto") {
		return
	}

	RespondJSON(w, http.StatusCreated, fmt.Sprintf("Producto %s", MsgCreadoCorrectamente), p)
}

// Handler para actualizar un producto
func actualizarProducto(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "producto")
	if !ok {
		return
	}

	var p Producto
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&p), w) {
		return
	}

	result, err := DB.Exec(`
		UPDATE producto 
		SET nombre=$1, descripcion=$2, precio_actual=$3, fecha_vencimiento=$4,
			imagen=$5, stock=$6, id_categoria=$7, id_proveedor=$8, activo=$9
		WHERE id_producto=$10
	`,
		p.Nombre, p.Descripcion, p.PrecioActual,
		p.FechaVencimiento, p.Imagen, p.Stock,
		p.IDCategoria, p.IDProveedor, p.Activo, idStr,
	)

	if ManejarErrorInsertActualizar(err, w, "update", "producto") {
		return
	}

	if !ValidarFilasAfectadas(result, w, "Producto") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Producto %s", MsgActualizadoCorrectamente), nil)
}

// Handler para eliminar un producto (desactivarlo)
func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "producto")
	if !ok {
		return
	}

	result, err := DB.Exec(
		"UPDATE producto SET activo = FALSE WHERE id_producto = $1 AND activo = TRUE",
		idStr,
	)

	if ManejarErrorInsertActualizar(err, w, "delete", "producto") {
		return
	}

	if !ValidarFilasAfectadas(result, w, "Producto") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Producto %s", MsgDesactivadoCorrectamente), nil)
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
			&c.IdCliente, &c.Nombre, &c.Telefono, &c.Correo, &c.Activo,
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
	idStr, ok := ValidarIDParametro(r, w, "cliente")
	if !ok {
		return
	}

	c, err := ObtenerClientePorID(idStr)
	if ManejarErrorConsulta(err, w, "Cliente") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Cliente %s", MsgObtenidoCorrectamente), c)
}

// Handler para crear un cliente
func crearCliente(w http.ResponseWriter, r *http.Request) {
	var c Cliente
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&c), w) {
		return
	}

	err := DB.QueryRow(`
		INSERT INTO cliente (nombre, telefono, correo)
		VALUES ($1, $2, $3)
		RETURNING id_cliente
	`,
		c.Nombre, c.Telefono, c.Correo,
	).Scan(&c.IdCliente)

	if ManejarErrorInsertActualizar(err, w, "insert", "cliente") {
		return
	}

	RespondJSON(w, http.StatusCreated, fmt.Sprintf("Cliente %s", MsgCreadoCorrectamente), c)
}

// Handler para actualizar un cliente
func actualizarCliente(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "cliente")
	if !ok {
		return
	}

	var c Cliente
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&c), w) {
		return
	}

	result, err := DB.Exec(`
		UPDATE cliente 
		SET nombre=$1, telefono=$2, correo=$3, activo=$4
		WHERE id_cliente=$5
	`,
		c.Nombre, c.Telefono, c.Correo, c.Activo, idStr,
	)

	if ManejarErrorInsertActualizar(err, w, "update", "cliente") {
		return
	}

	if !ValidarFilasAfectadas(result, w, "Cliente") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Cliente %s", MsgActualizadoCorrectamente), nil)
}

// Handler para eliminar un cliente (desactivarlo)
func eliminarCliente(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "cliente")
	if !ok {
		return
	}

	result, err := DB.Exec(
		"UPDATE cliente SET activo = FALSE WHERE id_cliente = $1 AND activo = TRUE",
		idStr,
	)

	if ManejarErrorInsertActualizar(err, w, "delete", "cliente") {
		return
	}

	if !ValidarFilasAfectadas(result, w, "Cliente") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Cliente %s", MsgDesactivadoCorrectamente), nil)
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
			&e.IdEmpleado, &e.Nombre, &e.Telefono, &e.Correo, &e.Activo,
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
	idStr, ok := ValidarIDParametro(r, w, "empleado")
	if !ok {
		return
	}

	e, err := ObtenerEmpleadoPorID(idStr)
	if ManejarErrorConsulta(err, w, "Empleado") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Empleado %s", MsgObtenidoCorrectamente), e)
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
	idStr, ok := ValidarIDParametro(r, w, "proveedor")
	if !ok {
		return
	}

	prov, err := ObtenerProveedorPorID(idStr)
	if ManejarErrorConsulta(err, w, "Proveedor") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Proveedor %s", MsgObtenidoCorrectamente), prov)
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
			&c.IdCategoria, &c.Nombre,
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
	idStr, ok := ValidarIDParametro(r, w, "categoria")
	if !ok {
		return
	}

	c, err := ObtenerCategoriaPorID(idStr)
	if ManejarErrorConsulta(err, w, "Categoria") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Categoria %s", MsgObtenidoCorrectamente), c)
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
	idStr, ok := ValidarIDParametro(r, w, "compra")
	if !ok {
		return
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

	if ManejarErrorConsulta(err, w, "Compra") {
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

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Compra %s", MsgObtenidoCorrectamente), c)
}

// Handler para crear una compra
func crearCompra(w http.ResponseWriter, r *http.Request) {
	var req CompraRequest
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&req), w) {
		return
	}

	resultado, ok := EjecutarEnTransaccion(w, func(tx *sql.Tx) (interface{}, error) {
		return procesarNuevaCompra(tx, req)
	})

	if !ok {
		return
	}

	RespondJSON(w, http.StatusCreated, fmt.Sprintf("Compra %s", MsgCreadoCorrectamente), resultado)
}

// Handler para cancelar una compra (valido unicamente para estado completado)
func cancelarCompra(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "compra")
	if !ok {
		return
	}

	resultado, ok := EjecutarEnTransaccion(w, func(tx *sql.Tx) (interface{}, error) {
		return procesarCancelacionCompra(tx, idStr)
	})

	if !ok {
		return
	}

	RespondJSON(w, http.StatusOK, "Compra cancelada correctamente y stock restaurado", resultado)
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

// Handler para vista de desempeno laboral
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
