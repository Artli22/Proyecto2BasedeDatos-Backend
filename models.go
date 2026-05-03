package main 
type Cliente struct {
	IdCliente int `json:"id_cliente"`
	Nombre 	   string `json:"nombre"`
	Telefono   *string `json:"telefono"`
	Correo     *string `json:"correo"`
	Activo     bool `json:"activo"`
}

type Empleado struct {
	IdEmpleado int `json:"id_empleado"`
	Nombre      string `json:"nombre"`
	Telefono    *string `json:"telefono"`
	Correo      *string `json:"correo"`
	Activo      bool `json:"activo"`
}

type Categoria struct {
	IdCategoria int `json:"id_categoria"`
	Nombre       string `json:"nombre"`
}

type Proveedor struct {
	IDProveedor int `json:"id_proveedor"`
	Nombre        string `json:"nombre"`
	Telefono      *string `json:"telefono"`
	Correo        *string `json:"correo"`
	Activo        bool `json:"activo"`
}

type Compra struct {
	IDCompra   int `json:"id_compra"`
	Fecha       string `json:"fecha"`
	Total       float64 `json:"total"`
	MetodoPago *string `json:"metodo_pago"`
	Estado      *string `json:"estado"`
	NumFactura string `json:"num_factura"`
	IDCliente 	int `json:"id_cliente"`
	IDEmpleado int `json:"id_empleado"`
}

type Producto struct {
    IDProducto       int      `json:"id_producto"`
    Nombre           string   `json:"nombre"`
    Descripcion      *string  `json:"descripcion"`       
    PrecioActual     float64  `json:"precio_actual"`
    FechaVencimiento *string  `json:"fecha_vencimiento"` 
    Imagen           *string  `json:"imagen"`            
    Stock            int      `json:"stock"`
	Activo           bool     `json:"activo"`
    IDCategoria      int      `json:"id_categoria"`
    IDProveedor      int      `json:"id_proveedor"`
}

type ItemCompra struct {
    IDProducto int `json:"id_producto"`
    Cantidad   int `json:"cantidad"`
}

type CompraRequest struct {
    Fecha      string       `json:"fecha"`
    MetodoPago string       `json:"metodo_pago"`
    IDCliente  int          `json:"id_cliente"`
    IDEmpleado int          `json:"id_empleado"`
    Productos  []ItemCompra `json:"productos"`
}

type ItemDetalle struct {
    Producto       string  `json:"producto"`
    Cantidad       int     `json:"cantidad"`
    PrecioUnitario float64 `json:"precio_unitario"`
    SubTotal       float64 `json:"sub_total"`
}

type CompraDetalle struct {
    IDCompra   int           `json:"id_compra"`
    Fecha      string        `json:"fecha"`
    Total      float64       `json:"total"`
    MetodoPago *string       `json:"metodo_pago"`
    Estado     *string       `json:"estado"`
    NumFactura string        `json:"num_factura"`
    Cliente    string        `json:"cliente"`    
    Empleado   string        `json:"empleado"`   
    Productos  []ItemDetalle `json:"productos"`  
}

type AuditoriaVenta struct {
    IDCompra      int     `json:"id_compra"`
    NumFactura    string  `json:"num_factura"`
    Fecha         string  `json:"fecha"`
    MetodoPago    *string `json:"metodo_pago"`
    Estado        *string `json:"estado"`
    Total         float64 `json:"total"`
    Cliente       string  `json:"cliente"`
    CorreoCliente *string `json:"correo_cliente"`
    EmpleadoCajero string `json:"empleado_cajero"`
}

type DetalleTemp struct {
    IDProducto     int
    Cantidad       int
    PrecioUnitario float64
    SubTotal       float64
}

type RentabilidadProducto struct {
    IDProducto          int     `json:"id_producto"`
    Producto            string  `json:"producto"`
    Categoria           string  `json:"categoria"`
    UnidadesVendidas    int     `json:"unidades_vendidas"`
    IngresosTotales     float64 `json:"ingresos_totales"`
    PrecioPromedioVenta float64 `json:"precio_promedio_venta"`
}

type StockCritico struct {
    IDProducto        int     `json:"id_producto"`
    Producto          string  `json:"producto"`
    Categoria         string  `json:"categoria"`
    Proveedor         string  `json:"proveedor"`
    TelefonoProveedor *string `json:"telefono_proveedor"`
    StockActual       int     `json:"stock_actual"`
    FechaVencimiento  *string `json:"fecha_vencimiento"`
}

type DesempenoEmpleado struct {
    IDEmpleado          int     `json:"id_empleado"`
    Empleado            string  `json:"empleado"`
    TotalTransacciones  int     `json:"total_transacciones"`
    MontoTotalVendido   float64 `json:"monto_total_vendido"`
    TicketPromedio      float64 `json:"ticket_promedio"`
    UltimaVenta         *string `json:"ultima_venta"`
}

type DetalleCompra struct {
    IDCompra        int     `json:"id_compra"`
    IDProducto      int     `json:"id_producto"`
    Cantidad        int     `json:"cantidad"`
    PrecioUnitario  float64 `json:"precio_unitario"`
    SubTotal        float64 `json:"sub_total"`
}