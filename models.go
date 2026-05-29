package main 
type Cliente struct {
	IdCliente int `gorm:"primaryKey;column:id_cliente" json:"id_cliente"`
	Nombre 	   string `gorm:"column:nombre" json:"nombre"`
	Telefono   *string `gorm:"column:telefono" json:"telefono"`
	Correo     *string `gorm:"column:correo" json:"correo"`
	Activo     bool `gorm:"column:activo" json:"activo"`
}

func (Cliente) TableName() string {
	return "cliente"
}

type Empleado struct {
	IdEmpleado int `gorm:"primaryKey;column:id_empleado" json:"id_empleado"`
	Nombre      string `gorm:"column:nombre" json:"nombre"`
	Telefono    *string `gorm:"column:telefono" json:"telefono"`
	Correo      *string `gorm:"column:correo" json:"correo"`
	Activo      bool `gorm:"column:activo" json:"activo"`
}

func (Empleado) TableName() string {
	return "empleado"
}

type Categoria struct {
	IdCategoria int `gorm:"primaryKey;column:id_categoria" json:"id_categoria"`
	Nombre       string `gorm:"column:nombre" json:"nombre"`
}

func (Categoria) TableName() string {
	return "categoria"
}

type Proveedor struct {
	IDProveedor int `gorm:"primaryKey;column:id_proveedor" json:"id_proveedor"`
	Nombre        string `gorm:"column:nombre" json:"nombre"`
	Telefono      *string `gorm:"column:telefono" json:"telefono"`
	Correo        *string `gorm:"column:correo" json:"correo"`
	Activo        bool `gorm:"column:activo" json:"activo"`
}

func (Proveedor) TableName() string {
	return "proveedor"
}

type Compra struct {
	IDCompra   int `gorm:"primaryKey;column:id_compra" json:"id_compra"`
	Fecha      string `gorm:"column:fecha" json:"fecha"`
	Total      float64 `gorm:"column:total" json:"total"`
	MetodoPago *string `gorm:"column:metodo_pago" json:"metodo_pago"`
	Estado     *string `gorm:"column:estado" json:"estado"`
	NumFactura *string `gorm:"column:num_factura" json:"num_factura"`
	IDCliente  int `gorm:"column:id_cliente" json:"id_cliente"`
	IDEmpleado int `gorm:"column:id_empleado" json:"id_empleado"`
}

func (Compra) TableName() string {
	return "compra"
}

type Producto struct {
    IDProducto       int `gorm:"primaryKey;column:id_producto" json:"id_producto"`
    Nombre           string `gorm:"column:nombre" json:"nombre"`
    Descripcion      *string `gorm:"column:descripcion" json:"descripcion"`       
    PrecioActual     float64 `gorm:"column:precio_actual" json:"precio_actual"`
    FechaVencimiento *string `gorm:"column:fecha_vencimiento" json:"fecha_vencimiento"` 
    Imagen           *string `gorm:"column:imagen" json:"imagen"`            
    Stock            int `gorm:"column:stock_actual" json:"stock"`
	Activo           bool `gorm:"column:activo" json:"activo"`
    IDCategoria      int `gorm:"column:id_categoria" json:"id_categoria"`
    IDProveedor      int `gorm:"column:id_proveedor" json:"id_proveedor"`
}

func (Producto) TableName() string {
	return "producto"
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
    IDCompra       int     `json:"id_compra"`
    NumFactura     *string `json:"num_factura"`
    Fecha          string  `json:"fecha"`
    MetodoPago     *string `json:"metodo_pago"`
    Estado         *string `json:"estado"`
    Total          float64 `json:"total"`
    Cliente        string  `json:"cliente"`
    CorreoCliente  *string `json:"correo_cliente"`
    EmpleadoCajero string  `json:"empleado_cajero"`
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

type LoginRequest struct {
    Usuario    string `json:"usuario"`
    Contraseña string `json:"contraseña"`
}

type LoginResponse struct {
    Token   string `json:"token"`
    Usuario string `json:"usuario"`
    Rol     string `json:"rol"`
    Message string `json:"message"`
}