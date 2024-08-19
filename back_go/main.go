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

type GroupStatus struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	ProcessingEnd bool   `json:"processing_end"`
}

var (
	dataReady = false
	dataCond  = sync.NewCond(&sync.Mutex{})
	groupData = []GroupStatus{
		{ID: uuid.New().String(), Status: "completed", Name: "First"},
		{ID: uuid.New().String(), Status: "completed", Name: "New Task"},
		{ID: uuid.New().String(), Status: "completed", Name: "SSE"},
	}
	newGroupData []GroupStatus // хранение новых данных
)

func handleObjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newData []GroupStatus
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Сохраняем новые данные отдельно
	newGroupData = newData
	groupData = append(groupData, newData...)
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

	clientClosed := make(chan bool)

	// Ожидание закрытия соединения клиентом
	go func() {
		if c, ok := w.(http.CloseNotifier); ok {
			<-c.CloseNotify()
			log.Println("Клинет закрыл соединение - 81")
			clientClosed <- true
		}
	}()

	dataCond.L.Lock()
	for !dataReady {
		dataCond.Wait() // Ждать пока данные не будут готовы
	}
	dataCond.L.Unlock()

	for i := 0; i < len(newGroupData); i++ {
		select {
		case <-clientClosed:
			log.Println("Клмент закрыл соединени, прекращаем обработку")
			return
		default:
			// изменяем статус объекта
			newGroupData[i].Status = getRandomStatus()

			// Если это последний элемент, устанавливаем флаг ProcessingEnd
			if i == len(newGroupData)-1 {
				newGroupData[i].ProcessingEnd = true
			}

			// Отправить обновленные данные на клиент
			data, err := json.Marshal(newGroupData[i])
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

	// Очистить `newGroupData` после отправки, чтобы быть готовым к новым данным
	newGroupData = nil
	dataReady = false
}

// func handleSSE(w http.ResponseWriter, r *http.Request) {
// 	flusher, ok := w.(http.Flusher)
// 	if !ok {
// 		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	clientClosed := make(chan bool)

// 	// Ожидание закрытия соединения клиентом
// 	go func() {
// 		if c, ok := w.(http.CloseNotifier); ok {
// 			<-c.CloseNotify()
// 			log.Println("Клиент закрыл соединение")
// 			clientClosed <- true
// 		}
// 	}()

// 	dataCond.L.Lock()
// 	for !dataReady {
// 		dataCond.Wait() // Ждать пока данные не будут готовы
// 	}
// 	dataCond.L.Unlock()

// 	for i := 0; i < len(newGroupData); i++ {
// 		log.Println("Старт обработки сущностей")
// 		newGroupData[i].Status = getRandomStatus()

// 		if i == len(newGroupData)-1 {
// 			newGroupData[i].ProcessingEnd = true
// 		}

// 		groupData[len(groupData)-len(newGroupData)+1] = newGroupData[i]

// 		// Отправить обновленные данные на клиент
// 		data, err := json.Marshal(newGroupData[i])
// 		log.Println("Отправили сущность", newGroupData[i])
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		fmt.Fprintf(w, "data: %s\n\n", data)
// 		flusher.Flush()
// 		// искусственная задержка перед обработкой следующего элемента
// 		select {
// 		case <-time.After(1 * time.Second):
// 			// продолжаем цикл
// 		case <-clientClosed:
// 			// если клиент закрыл соединение, завершаем обработку
// 			return
// 		}
// 	}

// 	// // Изменить статусы только новых объектов
// 	// for i := range newGroupData {
// 	// 	newGroupData[i].Status = getRandomStatus()
// 	// }

// 	// // Обновляем groupData с измененными статусами
// 	// groupData = append(groupData[:len(groupData)-len(newGroupData)], newGroupData...)

// 	// // Отправить обновленные данные на клиент
// 	// data, err := json.Marshal(newGroupData)
// 	// log.Println("Данные для клиента", newGroupData)
// 	// if err != nil {
// 	// 	http.Error(w, fmt.Sprintf("Failed to encode JSON: %v", err), http.StatusInternalServerError)
// 	// 	return
// 	// }
// 	// fmt.Fprintf(w, "data: %s\n\n", data)
// 	// flusher.Flush()

// 	// // Оставляем соединение открытым, пока клиент сам его не закроет
// 	// if c, ok := w.(http.CloseNotifier); ok {
// 	// 	<-c.CloseNotify()
// 	// 	log.Println("Клиент закрыл соединение")
// 	// }

// 	// Очистить newGroupData после отправки, чтобы быть готовым к новым данным
// 	newGroupData = nil
// 	dataReady = false
// }

func getRandomStatus() string {
	statuses := []string{"in_progress", "completed"}
	log.Println("Обработка статусов")
	return statuses[rand.Intn(len(statuses))]
}

func getConnect(w http.ResponseWriter, r *http.Request) {
	responseData, err := json.Marshal(groupData)
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
