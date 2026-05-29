
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Endpoint para autenticación de usuarios
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondJSON(w, http.StatusMethodNotAllowed, "Método no permitido", nil)
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, "El cuerpo del request no es un JSON válido", nil)
		return
	}

	if req.Usuario == "" || req.Contraseña == "" {
		RespondJSON(w, http.StatusBadRequest, "Usuario y contraseña son requeridos", nil)
		return
	}

	rol, err := ValidarCredenciales(req.Usuario, req.Contraseña)
	if err != nil {
		RespondJSON(w, http.StatusUnauthorized, "Credenciales inválidas", nil)
		return
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "tu-secret-key-default" 
	}

	token, err := GenerarToken(req.Usuario, rol, secretKey)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Error generando token", nil)
		return
	}

	response := LoginResponse{
		Token:   token,
		Usuario: req.Usuario,
		Rol:     rol,
		Message: "Autenticación exitosa",
	}

	RespondJSON(w, http.StatusOK, "Login exitoso", response)
}

// Enpoint para traer todos los productos
// Enpoint para traer todos los productos con GORM
func getProductos(w http.ResponseWriter, r *http.Request) {
	var productos []Producto
	// GORM: Leer todos los productos activos de la base de datos
	result := DB.Where("activo = ?", true).Find(&productos)
	if result.Error != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar productos en la base de datos", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Productos obtenidos correctamente", productos)
}

// Enpoint para traer un producto por ID
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

// Enpoint para crear un producto con GORM
func crearProducto(w http.ResponseWriter, r *http.Request) {
	var p Producto
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&p), w) {
		return
	}

	// GORM: Establecer activo como true por defecto
	p.Activo = true
	// GORM: Crear el producto en la base de datos
	result := DB.Create(&p)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "insert", "producto")
		return
	}

	RespondJSON(w, http.StatusCreated, fmt.Sprintf("Producto %s", MsgCreadoCorrectamente), p)
}

// Enpoint para actualizar un producto con GORM
func actualizarProducto(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "producto")
	if !ok {
		return
	}

	var p Producto
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&p), w) {
		return
	}

	// Convertir string a int
	id := convertStringToInt(idStr)
	
	// GORM: Actualizar solo los campos no cero en el struct
	result := DB.Model(&Producto{}).Where("id_producto = ?", id).Updates(p)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "update", "producto")
		return
	}

	if result.RowsAffected == 0 {
		RespondJSON(w, http.StatusNotFound, "Producto no encontrado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Producto %s", MsgActualizadoCorrectamente), nil)
}

// Enpoint para eliminar un producto (desactivarlo)
func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "producto")
	if !ok {
		return
	}

	id := convertStringToInt(idStr)
	result := DB.Model(&Producto{}).Where("id_producto = ? AND activo = ?", id, true).Update("activo", false)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "delete", "producto")
		return
	}

	if result.RowsAffected == 0 {
		RespondJSON(w, http.StatusNotFound, "Producto no encontrado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Producto %s", MsgDesactivadoCorrectamente), nil)
}

// Enpoint para traer todos los clientes
func getClientes(w http.ResponseWriter, r *http.Request) {
	var clientes []Cliente
	result := DB.Find(&clientes)
	if result.Error != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar clientes en la base de datos", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Clientes obtenidos correctamente", clientes)
}

// Enpoint para traer un cliente por ID
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

// Enpoint para crear un cliente
func crearCliente(w http.ResponseWriter, r *http.Request) {
	var c Cliente
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&c), w) {
		return
	}

	result := DB.Create(&c)
	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "insert", "cliente")
		return
	}

	RespondJSON(w, http.StatusCreated, fmt.Sprintf("Cliente %s", MsgCreadoCorrectamente), c)
}

// Enpoint para actualizar un cliente
func actualizarCliente(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "cliente")
	if !ok {
		return
	}

	var c Cliente
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&c), w) {
		return
	}

	id := convertStringToInt(idStr)
	result := DB.Model(&Cliente{}).Where("id_cliente = ?", id).Updates(c)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "update", "cliente")
		return
	}

	if result.RowsAffected == 0 {
		RespondJSON(w, http.StatusNotFound, "Cliente no encontrado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Cliente %s", MsgActualizadoCorrectamente), nil)
}

// Enpoint para eliminar un cliente (desactivarlo)
func eliminarCliente(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "cliente")
	if !ok {
		return
	}

	id := convertStringToInt(idStr)
	result := DB.Model(&Cliente{}).Where("id_cliente = ? AND activo = ?", id, true).Update("activo", false)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "delete", "cliente")
		return
	}

	if result.RowsAffected == 0 {
		RespondJSON(w, http.StatusNotFound, "Cliente no encontrado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Cliente %s", MsgDesactivadoCorrectamente), nil)
}

// Enpoint para traer todos los empleados
func getEmpleados(w http.ResponseWriter, r *http.Request) {
	var empleados []Empleado
	result := DB.Find(&empleados)
	if result.Error != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar empleados en la base de datos", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Empleados obtenidos correctamente", empleados)
}

// Enpoint para traer un empleado por ID
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

// Enpoint para crear un empleado
func crearEmpleado(w http.ResponseWriter, r *http.Request) {
	var e Empleado
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&e), w) {
		return
	}

	result := DB.Create(&e)
	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "insert", "empleado")
		return
	}

	RespondJSON(w, http.StatusCreated, fmt.Sprintf("Empleado %s", MsgCreadoCorrectamente), e)
}

// Enpoint para actualizar un empleado
func actualizarEmpleado(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "empleado")
	if !ok {
		return
	}

	var e Empleado
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&e), w) {
		return
	}

	id := convertStringToInt(idStr)
	result := DB.Model(&Empleado{}).Where("id_empleado = ?", id).Updates(e)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "update", "empleado")
		return
	}

	if result.RowsAffected == 0 {
		RespondJSON(w, http.StatusNotFound, "Empleado no encontrado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Empleado %s", MsgActualizadoCorrectamente), nil)
}

// Enpoint para eliminar un empleado (desactivarlo)
func eliminarEmpleado(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "empleado")
	if !ok {
		return
	}

	id := convertStringToInt(idStr)
	result := DB.Model(&Empleado{}).Where("id_empleado = ? AND activo = ?", id, true).Update("activo", false)

	if result.Error != nil {
		ManejarErrorInsertActualizar(result.Error, w, "delete", "empleado")
		return
	}

	if result.RowsAffected == 0 {
		RespondJSON(w, http.StatusNotFound, "Empleado no encontrado", nil)
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Empleado %s", MsgDesactivadoCorrectamente), nil)
}

// Enpoint para traer todos los proveedores
func getProveedores(w http.ResponseWriter, r *http.Request) {
	var proveedores []Proveedor
	result := DB.Find(&proveedores)
	if result.Error != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar proveedores en la base de datos", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Proveedores obtenidos correctamente", proveedores)
}

// Enpoint para traer un proveedor por ID
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

// Enpoint para traer todas las categorias
func getCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []Categoria
	result := DB.Find(&categorias)
	if result.Error != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar categorias en la base de datos", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Categorias obtenidas correctamente", categorias)
}

// Enpoint para traer una categoria por ID
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

// Enpoint para obtener todas las compras
func getCompras(w http.ResponseWriter, r *http.Request) {
	var compras []Compra
	result := DB.Order("fecha DESC").Find(&compras)
	if result.Error != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar compras", nil)
		return
	}

	RespondJSON(w, http.StatusOK, "Compras obtenidas correctamente", compras)
}

// Enpoint para obtener compra por ID
func getCompraPorID(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "compra")
	if !ok {
		return
	}

	var c CompraDetalle
	err := DB.Raw(`
		SELECT c.id_compra, c.fecha, c.total, c.metodo_pago, c.estado, c.num_factura,
		       cl.nombre AS cliente, e.nombre AS empleado
		FROM compra c
		JOIN cliente cl ON c.id_cliente = cl.id_cliente
		JOIN empleado e ON c.id_empleado = e.id_empleado
		WHERE c.id_compra = ?
	`, idStr).Scan(&c).Error

	if err != nil {
		ManejarErrorConsulta(err, w, "Compra")
		return
	}

	rows, err := DB.Raw(`
		SELECT p.nombre, dc.cantidad, dc.precio_unitario, dc.sub_total
		FROM detalle_compra dc
		JOIN producto p ON dc.id_producto = p.id_producto
		WHERE dc.id_compra = ?
	`, idStr).Rows()
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

// Enpoint para crear una compra
func crearCompra(w http.ResponseWriter, r *http.Request) {
	var req CompraRequest
	if !ValidarJSONDecodificacion(json.NewDecoder(r.Body).Decode(&req), w) {
		return
	}

	// Convertir el slice de productos a JSON (string para que pq lo envíe como text→jsonb)
	productosJSON, err := json.Marshal(req.Productos)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, "Error al procesar productos", nil)
		return
	}

	// Llamar al stored procedure sp_crear_compra_con_validacion
	// El cast ::jsonb es necesario porque pq envía strings como text
	var idCompra int
	var total float64
	var mensaje string

	err = DB.Raw(
		"SELECT * FROM sp_crear_compra_con_validacion($1, $2, $3, $4, $5::jsonb)",
		req.Fecha, req.MetodoPago, req.IDCliente, req.IDEmpleado, string(productosJSON),
	).Scan(map[string]interface{}{"id_compra": &idCompra, "total": &total, "mensaje": &mensaje}).Error

	if err != nil {
		// err.Error() contiene el mensaje real de PostgreSQL cuando el SP lanza RAISE EXCEPTION
		RespondJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// p_id_compra = 0 indica error de negocio retornado por el SP
	if idCompra == 0 {
		RespondJSON(w, http.StatusBadRequest, mensaje, nil)
		return
	}

	resultado := map[string]interface{}{
		"id_compra": idCompra,
		"total":     total,
	}

	RespondJSON(w, http.StatusCreated, mensaje, resultado)
}

// Enpoint para cancelar una compra (valido unicamente para estado completado)
func cancelarCompra(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "compra")
	if !ok {
		return
	}

	// Llamar al stored procedure sp_cancelar_compra
	var success bool
	var mensaje string

	err := DB.Raw(
		"SELECT * FROM sp_cancelar_compra($1)",
		idStr,
	).Scan(map[string]interface{}{"success": &success, "mensaje": &mensaje}).Error

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Error al ejecutar SP", nil)
		return
	}

	if !success {
		RespondJSON(w, http.StatusBadRequest, mensaje, nil)
		return
	}

	RespondJSON(w, http.StatusOK, mensaje, map[string]interface{}{
		"id_compra": idStr,
		"estado":    "cancelado",
	})
}

// Enpoint para vista de auditoria de ventas
func getAuditoriaVentas(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Raw(`
		SELECT id_compra, num_factura, fecha, metodo_pago, estado, total, cliente, correo_cliente, empleado_cajero 
		FROM vista_auditoria_completa_ventas
	`).Rows()
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

// Enpoint para vista de rentabilidad de productos
func getRentabilidadProductos(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Raw(`
		SELECT id_producto, producto, categoria, unidades_vendidas, ingresos_totales, precio_promedio_venta 
		FROM vista_rentabilidad_productos
	`).Rows()
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

// Enpoint para vista de control de stock
func getControlStock(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Raw(`
		SELECT id_producto, producto, categoria, proveedor, telefono_proveedor, stock_actual, fecha_vencimiento 
		FROM vista_stock_critico
	`).Rows()
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

// Enpoint para vista de desempeno laboral
func getDesempenoEmpleados(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Raw(`
		SELECT id_empleado, empleado, total_transacciones, monto_total_vendido, ticket_promedio, ultima_venta 
		FROM vista_desempeno_empleados
	`).Rows()
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

// Enpoint para obtener todos los detalles de compra
func getDetalleCompras(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Raw(`
		SELECT id_compra, id_producto, cantidad, precio_unitario, sub_total
		FROM detalle_compra
		ORDER BY id_compra DESC, id_producto ASC
	`).Rows()
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError,
			"Error al consultar detalles de compra", nil)
		return
	}
	defer rows.Close()

	detalles := []DetalleCompra{}
	for rows.Next() {
		var dc DetalleCompra
		err := rows.Scan(
			&dc.IDCompra, &dc.IDProducto, &dc.Cantidad,
			&dc.PrecioUnitario, &dc.SubTotal,
		)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError,
				"Error al leer fila de detalle de compra", nil)
			return
		}
		detalles = append(detalles, dc)
	}

	RespondJSON(w, http.StatusOK, "Detalles de compra obtenidos correctamente", detalles)
}

// Enpoint para obtener detalle de compra por ID
func getDetalleCompraPorID(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "detalle de compra")
	if !ok {
		return
	}

	var dc DetalleCompra
	err := DB.Raw(`
		SELECT id_compra, id_producto, cantidad, precio_unitario, sub_total
		FROM detalle_compra WHERE id_compra = ?
		LIMIT 1
	`, idStr).Scan(&dc).Error

	if ManejarErrorConsulta(err, w, "Detalle de compra") {
		return
	}

	RespondJSON(w, http.StatusOK, fmt.Sprintf("Detalle de compra %s", MsgObtenidoCorrectamente), dc)
}

// Proceso Interno: Obtener resumen de compras con totales
func getResumenCompras(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	idClienteStr := query.Get("id_cliente")
	var idCliente interface{} = nil

	if idClienteStr != "" {
		id, err := strconv.Atoi(idClienteStr)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, "ID de cliente inválido", nil)
			return
		}
		idCliente = id
	}

	var totalCompras int
	var montoTotal *float64
	var montoPromedio *float64
	var mensaje string

	err := DB.Raw(
		"SELECT * FROM sp_obtener_resumen_compras(?)",
		idCliente,
	).Scan(map[string]interface{}{"total_compras": &totalCompras, "monto_total": &montoTotal, "monto_promedio": &montoPromedio, "mensaje": &mensaje}).Error

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Error al ejecutar SP", nil)
		return
	}

	mTotal := 0.0
	if montoTotal != nil {
		mTotal = *montoTotal
	}

	mPromedio := 0.0
	if montoPromedio != nil {
		mPromedio = *montoPromedio
	}

	resultado := map[string]interface{}{
		"total_compras":     totalCompras,
		"monto_total":       mTotal,
		"monto_promedio":    mPromedio,
		"mensaje":           mensaje,
	}

	RespondJSON(w, http.StatusOK, "Resumen de compras obtenido", resultado)
}

// Proceso Interno: Obtener reporte de inventario crítico
func getReporteInventarioCritico(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitStockStr := query.Get("limite_stock")
	limitStock := 20 // valor por defecto

	if limitStockStr != "" {
		limite, err := strconv.Atoi(limitStockStr)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, "Límite de stock inválido", nil)
			return
		}
		limitStock = limite
	}

	var totalProductos int
	var mensaje string

	err := DB.Raw(
		"SELECT * FROM sp_reporte_inventario_critico(?)",
		limitStock,
	).Scan(map[string]interface{}{"total_productos": &totalProductos, "mensaje": &mensaje}).Error

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Error al ejecutar SP", nil)
		return
	}

	// Obtener la lista de productos con stock crítico
	rows, err := DB.Raw(`
		SELECT id_producto, producto, categoria, proveedor, telefono_proveedor, stock_actual, fecha_vencimiento
		FROM vista_stock_critico
		WHERE stock_actual < ?
		ORDER BY stock_actual ASC
	`, limitStock).Rows()

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Error al consultar productos", nil)
		return
	}
	defer rows.Close()

	type ProductoCritico struct {
		IDProducto         int     `json:"id_producto"`
		Nombre             string  `json:"nombre"`
		Categoria          string  `json:"categoria"`
		Proveedor          string  `json:"proveedor"`
		TelefonoProveedor  string  `json:"telefono_proveedor"`
		StockActual        int     `json:"stock_actual"`
		FechaVencimiento   *string `json:"fecha_vencimiento"`
	}

	productos := []ProductoCritico{}
	for rows.Next() {
		var p ProductoCritico
		err := rows.Scan(&p.IDProducto, &p.Nombre, &p.Categoria, &p.Proveedor,
			&p.TelefonoProveedor, &p.StockActual, &p.FechaVencimiento)
		if err != nil {
			RespondJSON(w, http.StatusInternalServerError, "Error al leer producto", nil)
			return
		}
		productos = append(productos, p)
	}

	resultado := map[string]interface{}{
		"total_productos": totalProductos,
		"mensaje":         mensaje,
		"productos":       productos,
	}

	RespondJSON(w, http.StatusOK, "Reporte de inventario crítico", resultado)
}

// Proceso Interno: Obtener cliente con historial de compras
func getClienteConHistorial(w http.ResponseWriter, r *http.Request) {
	idStr, ok := ValidarIDParametro(r, w, "cliente")
	if !ok {
		return
	}

	var nombreCliente string
	var telefonoCliente string
	var correoCliente string
	var totalCompras int
	var montoTotalGastado float64
	var mensaje string

	err := DB.Raw(
		"SELECT * FROM sp_obtener_cliente_con_historial(?)",
		idStr,
	).Scan(map[string]interface{}{"nombre_cliente": &nombreCliente, "telefono_cliente": &telefonoCliente, "correo_cliente": &correoCliente, "total_compras": &totalCompras, "monto_total_gastado": &montoTotalGastado, "mensaje": &mensaje}).Error

	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, "Error al ejecutar SP", nil)
		return
	}

	resultado := map[string]interface{}{
		"nombre_cliente":           nombreCliente,
		"telefono_cliente":         telefonoCliente,
		"correo_cliente":           correoCliente,
		"total_compras":            totalCompras,
		"monto_total_gastado":      montoTotalGastado,
		"mensaje":                  mensaje,
	}

	RespondJSON(w, http.StatusOK, "Historial del cliente obtenido", resultado)
}
