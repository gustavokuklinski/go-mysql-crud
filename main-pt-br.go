package main

import (
	"database/sql"  // Pacote Database SQL para realizar query
	"log"           // Mostra informações no console
	"net/http"      // Gerencia URLs e Servidor Web
	"text/template" // Gerencia templates

	_ "github.com/go-sql-driver/mysql" // Driver Mysql para Go
)

// Struct utilizada para exibir dados no template
// Essa struct deve ter os mesmos campos do banco de dados
type Names struct {
	Id    int
	Name  string
	Email string
}

// Função dbConn, abre conexão com banco de dados
func dbConn() (db *sql.DB) {

	dbDriver := "mysql" // Driver do banco de dados
	dbUser := ""        // Usuário
	dbPass := ""        // Senha
	dbName := ""        // Nome do banco

	// Realiza a conexão com banco de dados:
	// sql.Open("DRIVER", "Usuario:Senha/@BancoDeDados)
	// A variavel `db` é utilizada junto com pacote `database/sql`
	// para a montagem de Query.
	// A variavel `err` é utilizada no tratamento de erros
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

// Usa a variável tmpl para renderizar todos os templates da pasta `tmpl`
// independente da extenção
var tmpl = template.Must(template.ParseGlob("tmpl/*"))

// Função Index, usada para renderizar o arquivo Index
func Index(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com banco de dados utilizando a função dbConn()
	db := dbConn()

	// Realiza a consulta com banco de dados e trata erros
	selDB, err := db.Query("SELECT * FROM names ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Names{}

	// Monta um array para guardar os valores da struct
	res := []Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email

		// Junta a Struct com Array
		res = append(res, n)

	}

	// Abre a página Index e exibe todos os registrados na tela
	tmpl.ExecuteTemplate(w, "Index", res)

	// Fecha conexão
	defer db.Close()
}

// Função Show exibe apenas um resultado
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	// Usa o ID para fazer a consulta e trata erros
	selDB, err := db.Query("SELECT * FROM names WHERE id=" + nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email

	}

	// Mostra o template
	tmpl.ExecuteTemplate(w, "Show", n)

	// Fecha a conexão
	defer db.Close()

}

// Função New apenas exibe o formulário para inserir dados
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Função Edit, edita uma linha
func Edit(w http.ResponseWriter, r *http.Request) {
	// Abre conexão com banco de dados
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	// Realiza consulta usando o ID e trata erros
	selDB, err := db.Query("SELECT * FROM names WHERE id=" + nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		// Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email

	}

	// Mostra o template com formulário preenchido para edição
	tmpl.ExecuteTemplate(w, "Edit", n)

	// Fecha conexão com banco de dados
	defer db.Close()
}

// Função Insert, insere valores no banco de dados
func Insert(w http.ResponseWriter, r *http.Request) {

	// Abre conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do fomulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulário com a SQL tratada e verifica erros
		insForm.Exec(name, email)

		// Exibe um log com o valores digitados no formulário
		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// Função Insert, insere valores no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {

	// Abre conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do fomulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulário com a SQL tratada e verifica erros
		insForm.Exec(name, email, id)

		// Exibe um log com o valores digitados no formulário
		log.Println("UPDATE: Name: " + name + " | E-mail: " + email)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// Função Insert, insere valores no banco de dados
func Delete(w http.ResponseWriter, r *http.Request) {

	// Abre conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	nId := r.URL.Query().Get("id")

	// Prepara a SQL e verifica errors
	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	// Insere valores do formulário com a SQL tratada e verifica erros
	delForm.Exec(nId)

	// Exibe um log com o valores digitados no formulário
	log.Println("DELETE")

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func main() {

	// Exibe mensagem que o servidor iniciou
	log.Println("Server started on: http://localhost:9000")

	// Gerencia as URLs
	http.HandleFunc("/", Index)    // Mostra todos os registros
	http.HandleFunc("/show", Show) // Template: Um registro
	http.HandleFunc("/new", New)   // Template: Novo registro
	http.HandleFunc("/edit", Edit) // Template: Edita registro

	// Ações
	http.HandleFunc("/insert", Insert) // Ação: Novo Registro
	http.HandleFunc("/update", Update) // Ação: Edita registro
	http.HandleFunc("/delete", Delete) // Ação: Deletar registro
	// Inicia o servidor na porta 9000
	http.ListenAndServe(":9000", nil)

}
