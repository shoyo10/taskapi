package repository

type repository struct {
}

// NewRepository 依賴注入
func New() Repositorier {
	return &repository{}
}
