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
	"errors"
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
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
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

func requestHandling(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	switch result.(type) {
	case error:
		errjson, _ := json.Marshal(&answer{
			Error: result.(error).Error(),
			Data:  ""})
		http.Error(w, string(errjson), code)
	default:
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(&answer{
			Error: "",
			Data:  result})
	}
}

func PageNotFound(w http.ResponseWriter, req *http.Request) {
	requestHandling(w, errors.New("404 Page Not Found"), http.StatusNotFound)
}

func Create(w http.ResponseWriter, req *http.Request) {

	var NewUser user
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
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
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, "Ok", http.StatusOK)

}

func Login(w http.ResponseWriter, req *http.Request) {

	var NewUser user
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
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
		requestHandling(w, nil, http.StatusBadRequest)
		return
	}

	requestHandling(w, token, http.StatusOK)

}

func FindUser(w http.ResponseWriter, req *http.Request) {

	parUserUrl := req.URL.Query()["User"]
	if len(parUserUrl) == 0 {
		requestHandling(w, errors.New("User parametrs is missing"), http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := user_usecase.New(repo)
	con := controller.NewUserUseCase(useCase)
	results, err := con.Find(req.Context(), parUserUrl[0])
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, results, http.StatusOK)

}

func AddContact(w http.ResponseWriter, req *http.Request) {

	var NewUser contact
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := user_usecase.New(repo)
	con := controller.NewUserUseCase(useCase)
	err = con.AddContact(req.Context(), &entity.FindUser{Login: NewUser.Login})
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, "Ok", http.StatusOK)

}

func CreateMessage(w http.ResponseWriter, req *http.Request) {

	var NewMes message
	err := json.NewDecoder(req.Body).Decode(&NewMes)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := message_usecase.New(repo)
	con := controller.NewMessageUseCase(useCase)
	err = con.CreateMessage(req.Context(), NewMes.Body, NewMes.Recipient)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, "Ok", http.StatusOK)

}

func ListMessages(w http.ResponseWriter, req *http.Request) {

	cfg := config.Get()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	parOffsetUrl := req.URL.Query()["afterAt"]
	var offset interface{}
	if len(parOffsetUrl) != 0 {
		offset, err = time.Parse(time.RFC3339, parOffsetUrl[0])
		if err != nil {
			requestHandling(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		offset = nil
	}

	parRecipientUrl := req.URL.Query()["Recipient"]
	if len(parRecipientUrl) == 0 {
		requestHandling(w, errors.New("Recipient parametrs is missing"), http.StatusBadRequest)
		return
	}

	repo := storage.New(client.Database("ChatGo"))
	useCase := message_usecase.New(repo)
	con := controller.NewMessageUseCase(useCase)
	results, err := con.ListMessages(req.Context(), parRecipientUrl[0], offset)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, results, http.StatusOK)

}
