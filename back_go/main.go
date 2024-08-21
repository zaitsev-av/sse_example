package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EntityType struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	ProcessingEnd bool   `json:"processing_end"`
}

var (
	dataReady  = false
	dataCond   = sync.NewCond(&sync.Mutex{})
	entityData = []EntityType{
		{ID: uuid.New().String(), Status: "completed", Name: "First"},
		{ID: uuid.New().String(), Status: "completed", Name: "New Task"},
		{ID: uuid.New().String(), Status: "completed", Name: "SSE"},
	}
	newEntityData []EntityType // хранение новых данных
)

func handleObjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newData []EntityType
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Сохраняем новые данные отдельно
	newEntityData = newData
	entityData = append(entityData, newData...)
	dataReady = true
	dataCond.Broadcast() // Сигнализировать, что данные готовы

	responseData, err := json.Marshal(newData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientClosedConnect := make(chan bool)

	// Ожидание закрытия соединения клиентом
	go func() {
		if c, ok := w.(http.CloseNotifier); ok {
			<-c.CloseNotify()
			log.Println("Клинет закрыл соединение - 81")
			clientClosedConnect <- true
		}
	}()

	dataCond.L.Lock()
	for !dataReady {
		dataCond.Wait() // Ждать пока данные не будут готовы
	}
	dataCond.L.Unlock()

	for i := 0; i < len(newEntityData); i++ {
		select {
		case <-clientClosedConnect:
			log.Println("Клмент закрыл соединени, прекращаем обработку")
			return
		default:
			// изменяем статус объекта
			itemId := newEntityData[i].ID

			newItemStatus := getRandomStatus()
			newEntityData[i].Status = newItemStatus
			FindById(entityData, itemId).Status = newItemStatus

			// Если это последний элемент, устанавливаем флаг ProcessingEnd
			if i == len(newEntityData)-1 {
				newEntityData[i].ProcessingEnd = true
			}

			// Отправить обновленные данные на клиент
			data, err := json.Marshal(newEntityData[i])
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
				return
			}

			_, err = fmt.Fprintf(w, "data: %s\n\n", data)
			if err != nil {
				log.Println("Ошибка при отправке данных клиенту:", err)
				return
			}
			flusher.Flush()

			// Пауза перед обработкой следующего элемента
			time.Sleep(3 * time.Second)
		}
	}

	// Очистить `newEntityData` после отправки, чтобы быть готовым к новым данным
	newEntityData = nil
	dataReady = false
}

func getRandomStatus() string {
	statuses := []string{"in_progress", "completed"}
	log.Println("Обработка статусов")
	return statuses[rand.Intn(len(statuses))]
}

func getConnect(w http.ResponseWriter, r *http.Request) {
	responseData, err := json.Marshal(entityData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func withCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/objects", handleObjects)
	mux.HandleFunc("/events", handleSSE)
	mux.HandleFunc("/connect", getConnect)

	server := withCORS(mux)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}

func FindById(list []EntityType, id string) *EntityType {
	for i := range list {
		if list[i].ID == id {
			return &list[i]
		}
	}
	return nil
}
