package models

type Tarefa struct {
	ID              uint32 `json:"id"`
	Titulo          string `json:"titulo"`
	Descricao       string `json:"descricao"`
	Data_vencimento string `json:"data_vencimento"`
	Status          string `json:"status"`
	Criado_em       string `json:"criado_em"`
	Atualizado_em   string `json:"atualizado_em"`
}
