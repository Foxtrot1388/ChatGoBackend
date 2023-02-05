package app

import (
	"ChatGo/config"
	storage "ChatGo/internal/adapters/db/mongodb"
	controller "ChatGo/internal/controller/http/v1"
	entity "ChatGo/internal/domain/entity"
	message_usecase "ChatGo/internal/usecase/message"
	user_usecase "ChatGo/internal/usecase/user"
	auth "ChatGo/server/midleware"
	"context"
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type user struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type contact struct {
	Login string `json:"login"`
}

type message struct {
	Body      string `json:"body"`
	Recipient string `json:"recipient"`
}

type answer struct {
	Error string `json:"error"`
	Data  string `json:"data"`
}

func Run() error {

	validation.ErrorTag = "vall"

	r := mux.NewRouter()
	r.HandleFunc("/CreateUser", Create).Methods("POST", "PUT")
	r.HandleFunc("/LoginUser", Login).Methods("POST", "PUT")
	r.HandleFunc("/FindUser", FindUser).Methods("GET")
	r.HandleFunc("/AddContact", AddContact).Methods("POST", "PUT")
	r.HandleFunc("/CreateMessage", CreateMessage).Methods("POST", "PUT")
	r.HandleFunc("/ListMessages", ListMessages).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(PageNotFound)

	r.Use(auth.AuthMiddleware)

	err := http.ListenAndServe(":8090", r)
	if err != nil {
		return err
	}

	return nil
}

func PageNotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Page Not Found"))
}

func Create(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var NewUser user
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := user_usecase.New(repo)
	con := controller.NewUserUseCase(useCase)
	err = con.Create(req.Context(), &entity.User{
		Login: NewUser.Login,
		Pass:  NewUser.Pass,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&answer{})

}

func Login(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var NewUser user
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := user_usecase.New(repo)
	con := controller.NewUserUseCase(useCase)
	token, err := con.Login(req.Context(), &entity.User{
		Login: NewUser.Login,
		Pass:  NewUser.Pass,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{})
		return
	}

	json.NewEncoder(w).Encode(&answer{Data: token})

}

func FindUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	parUserUrl := req.URL.Query()["User"]
	if len(parUserUrl) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := user_usecase.New(repo)
	con := controller.NewUserUseCase(useCase)
	results, err := con.Find(req.Context(), parUserUrl[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&results)

}

func AddContact(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var NewUser contact
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := user_usecase.New(repo)
	con := controller.NewUserUseCase(useCase)
	err = con.AddContact(req.Context(), &entity.FindUser{Login: NewUser.Login})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&answer{})

}

func CreateMessage(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var NewMes message
	err := json.NewDecoder(req.Body).Decode(&NewMes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := message_usecase.New(repo)
	con := controller.NewMessageUseCase(useCase)
	err = con.CreateMessage(req.Context(), NewMes.Body, NewMes.Recipient)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&answer{})

}

func ListMessages(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	parOffsetUrl := req.URL.Query()["afterAt"]
	var offset interface{}
	if len(parOffsetUrl) != 0 {
		offset, err = time.Parse(time.RFC3339, parOffsetUrl[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&answer{Error: err.Error()})
			return
		}
	} else {
		offset = nil
	}

	parRecipientUrl := req.URL.Query()["Recipient"]
	if len(parRecipientUrl) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := message_usecase.New(repo)
	con := controller.NewMessageUseCase(useCase)
	results, err := con.ListMessages(req.Context(), parRecipientUrl[0], offset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&answer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&results)

}
