package main 
type Cliente struct {
	id_cliente int `json:"id_cliente"`
	nombre 	   string `json:"nombre"`
	telefono   *string `json:"telefono"`
	correo     *string `json:"correo"`
	activo     bool `json:"activo"`
}

type Empleado struct {
	id_empleado int `json:"id_empleado"`
	nombre      string `json:"nombre"`
	telefono    *string `json:"telefono"`
	correo      *string `json:"correo"`
	activo      bool `json:"activo"`
}

type Categoria struct {
	id_categoria int `json:"id_categoria"`
	nombre       string `json:"nombre"`
}

type Proveedor struct {
	id_proveedor int `json:"id_proveedor"`
	nombre        string `json:"nombre"`
	telefono      *string `json:"telefono"`
	correo        *string `json:"correo"`
	activo        bool `json:"activo"`
}

type Compra struct {
	id_compra   int `json:"id_compra"`
	fecha       string `json:"fecha"`
	total       float64 `json:"total"`
	metodo_pago *string `json:"metodo_pago"`
	estado 		*string `json:"estado"`
	num_factura string `json:"num_factura"`
	id_cliente 	int `json:"id_cliente"`
	id_empleado int `json:"id_empleado"`
}

type DetalleCompra struct {
	id_compra 		int `json:"id_compra"`
	id_producto 	int `json:"id_producto"`
	cantidad 		int `json:"cantidad"`
	precio_unitario float64 `json:"precio_unitario"`
	sub_total       float64 `json:"sub_total"`	

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