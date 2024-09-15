package servidor

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"tarefas/banco"

	"github.com/gorilla/mux"
)

//criar tarefas
//buscar tarefas OK
//buscar tarefa por ID
//editar tarefa
//apagar tarefa

type tarefa struct {
	ID              uint32 `json:"id"`
	Titulo          string `json:"titulo"`
	Descricao       string `json:"descricao"`
	Data_vencimento string `json:"data_vencimento"`
	Status          string `json:"status"`
	Criado_em       string `json:"criado_em"`
	Atualizado_em   string `json:"atualizado_em"`
}

func CriarTarefas(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Não foi possivel ler o corpo da requisição", http.StatusBadRequest)
		return
	}

	var tarefa tarefa
	if err := json.Unmarshal(body, &tarefa); err != nil {
		http.Error(w, "Não foi possivel converter para JSON", http.StatusBadRequest)
		return
	}

	db, err := banco.Connection()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	statement, err := db.Prepare("INSERT INTO tarefas(titulo, descricao, data_vencimento, status, criado_em, atualizado_em) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Printf("Erro ao preparar a query SQL %v", err)
		http.Error(w, "Erro ao inserir a tarefa", http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	inserir, err := statement.Exec(tarefa.Titulo, tarefa.Descricao, tarefa.Data_vencimento, tarefa.Status, tarefa.Criado_em, tarefa.Atualizado_em)
	if err != nil {

		http.Error(w, "Erro ao inserir", http.StatusInternalServerError)
		return
	}

	idInserido, err := inserir.LastInsertId()
	if err != nil {
		log.Printf("Erro ao buscar o id %v", err)
		return
	}

	w.Write([]byte(fmt.Sprintf("Tarefa inserida com sucesso! ID: %d", idInserido)))
	w.WriteHeader(http.StatusOK)
}

func BuscarTarefas(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Connection()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	rows, err := db.Query("select * from tarefas")
	if err != nil {
		log.Fatal("Erro ao realizar o select", err)
	}

	defer db.Close()

	var tarefas []tarefa
	for rows.Next() {
		var tarefa tarefa

		if err := rows.Scan(&tarefa.ID, &tarefa.Titulo, &tarefa.Descricao, &tarefa.Data_vencimento, &tarefa.Status, &tarefa.Criado_em, &tarefa.Atualizado_em); err != nil {
			http.Error(w, "Erro ao buscar as tarefas", http.StatusBadRequest)
			return
		}

		tarefas = append(tarefas, tarefa)
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tarefas); err != nil {
		http.Error(w, "Erro ao onverter as tarefas", http.StatusBadRequest)
		return
	}
}

func BuscarTarefa(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, err := strconv.ParseUint(parametros["id"], 10, 32)
	if err != nil {
		http.Error(w, "Erro ao buscar ID", http.StatusBadRequest)
	}

	db, err := banco.Connection()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco", http.StatusInternalServerError)
	}

	linha, err := db.Query("select * from tarefas where id = ?", ID)
	if err != nil {
		http.Error(w, "Erro ao realizar o select", http.StatusBadRequest)
	}

	var tarefa tarefa

	if linha.Next() {
		if err := linha.Scan(&tarefa.ID, &tarefa.Titulo, &tarefa.Descricao, &tarefa.Data_vencimento, &tarefa.Status, &tarefa.Criado_em, &tarefa.Atualizado_em); err != nil {
			http.Error(w, "Erro ao scanear usuário", http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tarefa); err != nil {
		http.Error(w, "Erro ao converter para JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AlteraTarefa(w http.ResponseWriter, r *http.Request) {
	paramentros := mux.Vars(r)

	ID, err := strconv.ParseUint(paramentros["id"], 10, 32)
	if err != nil {
		http.Error(w, "Erro ao buscar o ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao obter o body", http.StatusBadRequest)
		return
	}

	var tarefa tarefa
	if err := json.Unmarshal(body, &tarefa); err != nil {
		http.Error(w, "Erro ao converter para Json!", http.StatusBadRequest)
		return
	}

	db, err := banco.Connection()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	statement, err := db.Prepare("update tarefas set titulo = ?, descricao = ?, data_vencimento = ?, status = ?, criado_em = ?, atualizado_em = ? where id = ?")
	if err != nil {
		http.Error(w, "Erro ao preparar o update", http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	if _, err := statement.Exec(tarefa.Titulo, tarefa.Descricao, tarefa.Data_vencimento, tarefa.Status, tarefa.Criado_em, tarefa.Atualizado_em, ID); err != nil {
		log.Printf("Erro ao executar o insert %v", err)
		http.Error(w, "Erro ao realizar o update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletarTarefa(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, err := strconv.ParseUint(parametros["id"], 10, 32)
	if err != nil {
		http.Error(w, "Erro ao converter o ID para inteiro", http.StatusBadRequest)
		return
	}

	db, err := banco.Connection()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco de dados!", http.StatusInternalServerError)
	}

	defer db.Close()

	statement, err := db.Prepare("delete * from tarefas where id = ? ")
	if err != nil {
		http.Error(w, "Erro ao preparar o delete", http.StatusInternalServerError)
		return
	}

	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		http.Error(w, "Erro executar o delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
