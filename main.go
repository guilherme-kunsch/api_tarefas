package main

import (
	"fmt"
	"net/http"
	"tarefas/banco"
	"tarefas/servidor"

	"github.com/gorilla/mux"
)

func main() {
	db, err := banco.Connection()
	if err != nil {
		fmt.Println("Erro ao conectar no banco", err)
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/v1/tarefas", servidor.CriarTarefas).Methods(http.MethodPost)
	r.HandleFunc("/v1/tarefas", servidor.BuscarTarefas).Methods(http.MethodGet)
	r.HandleFunc("/v1/tarefas/{id}", servidor.BuscarTarefa).Methods(http.MethodGet)
	r.HandleFunc("/v1/tarefas/{id}", servidor.AlteraTarefa).Methods(http.MethodPut)
	r.HandleFunc("/v1/tarefas/{id}", servidor.DeletarTarefa).Methods(http.MethodDelete)

	fmt.Print("Servidor Ligado\n")
	http.ListenAndServe(":5001", r)

}
