package app

import (
	storage "ChatGo/internal/adapters/db/mongodb"
	controller "ChatGo/internal/controller/http/v1"
	"ChatGo/internal/domain/entity"
	contact_usecase "ChatGo/internal/usecase/contact"
	message_usecase "ChatGo/internal/usecase/message"
	user_usecase "ChatGo/internal/usecase/user"
	"ChatGo/pkg/logging"
	"ChatGo/server/midleware"
	"context"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
	"github.com/kataras/versioning"
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

type DeleteRequest struct {
	Id string `json:"id"`
}

type message struct {
	Body      string `json:"body"`
	Recipient string `json:"recipient"`
}

type Answer struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func Run() error {

	logger := logging.GetLogger()
	logger.Info("Start app")

	validation.ErrorTag = "vall"

	r := createRouter()

	err := http.ListenAndServe(":8090", r)
	if err != nil {
		return err
	}

	return nil
}

func createRouter() *mux.Router {

	logger := logging.GetLogger()
	logger.Trace("createRouter")

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(PageNotFound)

	r.Handle("/CreateUser", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(Create),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("POST", "PUT")

	r.Handle("/LoginUser", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(Login),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("POST", "PUT")

	r.Handle("/FindUser", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(FindUser),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("GET")

	r.Handle("/AddContact", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(AddContact),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("POST", "PUT")

	r.Handle("/DeleteContact", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(DeleteContact),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("POST", "PUT")

	r.Handle("/ListContact", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(ListContact),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("GET")

	r.Handle("/CreateMessage", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(CreateMessage),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("POST", "PUT")

	r.Handle("/ListMessages", midleware.VersionMiddleware(versioning.NewMatcher(versioning.Map{
		"1":                 http.HandlerFunc(ListMessages),
		versioning.NotFound: http.HandlerFunc(PageNotFound),
	}))).Methods("GET")

	r.Use(midleware.AuthMiddleware)

	return r
}

func requestHandling(w http.ResponseWriter, result interface{}, code int) {
	logger := logging.GetLogger()
	w.Header().Set("Content-Type", "application/json")
	switch result.(type) {
	case error:
		errjson, _ := json.Marshal(&Answer{
			Error: result.(error).Error(),
			Data:  ""})
		http.Error(w, string(errjson), code)
		logger.Errorf("code %d result %s", code, result.(error).Error())
	default:
		w.WriteHeader(code)
		resultjson, _ := json.Marshal(&Answer{
			Error: "",
			Data:  result})
		w.Write(resultjson)
		logger.Debug(string(resultjson))
	}
}

func PageNotFound(w http.ResponseWriter, req *http.Request) {
	logger := logging.GetLogger()
	logger.Trace("PageNotFound")
	requestHandling(w, errors.New("404 Page Not Found"), http.StatusNotFound)
}

func Create(w http.ResponseWriter, req *http.Request) {

	logger := logging.GetLogger()
	logger.Trace("Create")

	var NewUser user
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := user_usecase.NewUserUseCase(repo)
	con := controller.NewUserUseCase(useCase)
	err = con.CreateUser(req.Context(), &entity.User{
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

	logger := logging.GetLogger()
	logger.Trace("Login")

	var NewUser user
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := user_usecase.NewUserUseCase(repo)
	con := controller.NewUserUseCase(useCase)
	token, err := con.LoginUser(req.Context(), &entity.User{
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

	logger := logging.GetLogger()
	logger.Trace("FindUser")

	parUserUrl := req.URL.Query()["User"]
	if len(parUserUrl) == 0 {
		requestHandling(w, errors.New("User parametrs is missing"), http.StatusBadRequest)
		return
	}

	logger.Debugf(" FindUser User %s", parUserUrl[0])

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := user_usecase.NewUserUseCase(repo)
	con := controller.NewUserUseCase(useCase)
	results, err := con.FindUser(req.Context(), parUserUrl[0])
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, results, http.StatusOK)

}

func AddContact(w http.ResponseWriter, req *http.Request) {

	logger := logging.GetLogger()
	logger.Trace("AddContact")

	var NewUser contact
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := contact_usecase.New(repo)
	con := controller.NewContactUseCase(useCase)
	result, err := con.AddContact(req.Context(), &entity.FindUser{Login: NewUser.Login})
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, result, http.StatusOK)

}

func DeleteContact(w http.ResponseWriter, req *http.Request) {

	logger := logging.GetLogger()
	logger.Trace("DeleteContact")

	var NewUser DeleteRequest
	err := json.NewDecoder(req.Body).Decode(&NewUser)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := contact_usecase.New(repo)
	con := controller.NewContactUseCase(useCase)
	err = con.DeleteContact(req.Context(), NewUser.Id)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, "Ok", http.StatusOK)
}

func ListContact(w http.ResponseWriter, req *http.Request) {

	logger := logging.GetLogger()
	logger.Trace("ListContact")

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := contact_usecase.New(repo)
	con := controller.NewContactUseCase(useCase)
	results, err := con.ListContact(req.Context())
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, results, http.StatusOK)

}

func CreateMessage(w http.ResponseWriter, req *http.Request) {

	logger := logging.GetLogger()
	logger.Trace("CreateMessage")

	var NewMes message
	err := json.NewDecoder(req.Body).Decode(&NewMes)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := message_usecase.New(repo)
	con := controller.NewMessageUseCase(useCase)
	resultID, err := con.CreateMessage(req.Context(), NewMes.Body, NewMes.Recipient)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, resultID, http.StatusOK)

}

func ListMessages(w http.ResponseWriter, req *http.Request) {

	logger := logging.GetLogger()
	logger.Trace("ListMessages")

	parOffsetUrl := req.URL.Query()["afterAt"]
	var offset interface{}
	var err error
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

	logger.Debugf("ListMessages Recipient %s offset %v", parRecipientUrl[0], offset)

	repo, err := storage.New(context.TODO())
	if err != nil {
		requestHandling(w, err, http.StatusInternalServerError)
		return
	}

	useCase := message_usecase.New(repo)
	con := controller.NewMessageUseCase(useCase)
	results, err := con.ListMessages(req.Context(), parRecipientUrl[0], offset)
	if err != nil {
		requestHandling(w, err, http.StatusBadRequest)
		return
	}

	requestHandling(w, results, http.StatusOK)

}
