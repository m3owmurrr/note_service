package handlers

import (
	"cloud_technologies/internal/models"
	"cloud_technologies/internal/storage"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type NotesHandler struct {
	storage storage.NoteStorage
}

func NewNotesHandler(s storage.NoteStorage) *NotesHandler {
	return &NotesHandler{
		storage: s,
	}
}

func (nh *NotesHandler) PostNoteHandler(w http.ResponseWriter, r *http.Request) {

	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !CheckCaptcha(note.Token, r.RemoteAddr) {
		http.Error(w, "Ошибка капчи, повторите попытку", http.StatusForbidden)
		return
	}

	id := uuid.NewString()

	note.Id = id
	if err := nh.storage.UploadNote(&note); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idResp := models.Note{Id: id}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(idResp)
}

func (nh *NotesHandler) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	note, err := nh.storage.GetNote(idString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*note)
}

// структура для ответа от Яндекс.СмартКапчи
type CaptchaResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func CheckCaptcha(token, addr string) bool {
	fmt.Println(token, addr)
	// Получаем секретный ключ из переменных окружения или жестко прописываем
	secretKey := "<SECRET_KEY>" // секретный ключ Яндекс капчи
	if secretKey == "" {
		log.Fatal("YANDEX_CAPTCHA_SECRET_KEY is not set")
	}

	// Формируем запрос для проверки капчи
	url := "https://smartcaptcha.yandexcloud.net/validate"

	payloadBytes := fmt.Sprintf("secret=%v&token=%v&ip=%v", secretKey, token, addr)

	// Отправляем запрос на сервер Яндекса
	resp, err := http.Post(url, "x-www-form-urlencoded", strings.NewReader(payloadBytes))
	if err != nil {
		log.Println("Error sending captcha validation request:", err)
		return false
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return false
	}

	var captchaResp CaptchaResponse
	err = json.Unmarshal(body, &captchaResp)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return false
	}

	fmt.Println("----", captchaResp)
	// Проверяем статус капчи
	return captchaResp.Status == "ok"
}
