package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// 1 - Entity - Representa o formato exato que vai ser salvo no banco de dados como se fosse um molde
type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

// 2 - DTOs Usadas exclusivamente para entrada e saida de na camada web
type CreateUserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// 3 - Repository - A Interface que define O QUE o repositorio faz, não como ele faz.

type UserRepository interface {
	Save(user User) error
	FindByEmail(email string) (*User, error)
}

type memoryUserRepository struct {
	users map[string]User
}

func NewMemoryUserRepository() *memoryUserRepository {
	return &memoryUserRepository{users: make(map[string]User)}
}

func (r *memoryUserRepository) Save(u User) error {
	r.users[u.ID] = u
	return nil
}

func (r *memoryUserRepository) FindByEmail(email string) (*User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, nil
}

//  4 - Service - Aqui é onde a "mágica" e as regras de negócio acontecem. Se um usuário não pode se cadastrar com um e-mail que já existe, ou se um cálculo de juros precisa ser feito antes de salvar algo, é o Service que faz. Ele orquestra a lógica real da aplicação.

type UserService struct {
	repo UserRepository // Depende da Interface, e não de um banco específico!
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(dto CreateUserDTO) (*User, error) {
	// Regra de Negócio 1: Validação básica
	if dto.Name == "" || dto.Email == "" {
		return nil, errors.New("nome e email são obrigatórios")
	}

	// Regra de Negócio 2: Não permitir e-mails duplicados
	existingUser, _ := s.repo.FindByEmail(dto.Email)
	if existingUser != nil {
		return nil, errors.New("este email já está em uso")
	}

	// Transforma o DTO em uma Entity para ser salva
	newUser := User{
		ID:       "uuid-gerado-1234", // Num cenário real, gerariamos um UUID aqui
		Name:     dto.Name,
		Email:    dto.Email,
		Password: "hashed-" + dto.Password, // Simulando um hash de senha
	}

	// Manda o Repository salvar a Entity no banco
	err := s.repo.Save(newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// 5 - Controller - É a porta de entrada da sua aplicação. Ele recebe a requisição HTTP (o pedido do cliente), valida se os dados básicos estão ali e repassa o trabalho pesado para a próxima camada. Depois, ele pega a resposta e entrega de volta ao cliente com o Status Code correto

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Recebe e decodifica o JSON para o DTO
	var dto CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// 2. Repassa o DTO para o Service aplicar as regras e salvar
	createdUser, err := h.service.CreateUser(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	// 3. Converte a Entity retornada em um DTO de Resposta (escondendo a senha)
	response := UserResponseDTO{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}

	// 4. Entrega a resposta para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// 6 - Main

func main() {
	// Aqui nós "montamos" o quebra-cabeça de dentro para fora
	repo := NewMemoryUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)

	// Registra a rota HTTP apontando para o nosso Controller
	http.HandleFunc("/users", handler.HandleCreateUser)

	fmt.Println("Servidor rodando na porta 8080... Mande um POST para http://localhost:8080/users")
	http.ListenAndServe(":8080", nil)
}
