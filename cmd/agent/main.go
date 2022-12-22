package main

import (
	// "bytes"
	// "bytes"
	"encoding/json"
	"fmt"
	"log"

	// "runtime/metrics"

	"math/rand"
	// "net/http"
	"runtime"
	"sync"
	"time"
	// "io"
	"github.com/go-resty/resty/v2"
)

const (
	url            = "http://127.0.0.1:8080"
	pollInterval   = 2
	reportInterval = 4
)

// var mutex *sync.RWMutex
var wg sync.WaitGroup 
type Gauge float64
type Counter int64
type MonitorMap map[string]Gauge
var PollCount Counter
var RandomValue Gauge
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
// type MyApiError struct {
//     Code      int       `json:"code"`
//     Message   string    `json:"message"`
//     Timestamp time.Time `json:"timestamp"`
// }

func NewMonitor(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
	// fmt.Println(rtm)
	// mutex.Lock()
	(*m)["Alloc"] = Gauge(rtm.Alloc)
	(*m)["BuckHashSys"] = Gauge(rtm.BuckHashSys)
	(*m)["TotalAlloc"] = Gauge(rtm.TotalAlloc)
	(*m)["Sys"] = Gauge(rtm.Sys)
	(*m)["Mallocs"] = Gauge(rtm.Mallocs)
	(*m)["Frees"] = Gauge(rtm.Frees)
	(*m)["PauseTotalNs"] = Gauge(rtm.PauseTotalNs)
	(*m)["NumGC"] = Gauge(rtm.NumGC)
	(*m)["GCCPUFraction"] = Gauge(rtm.GCCPUFraction)
	(*m)["GCSys"] = Gauge(rtm.GCSys)
	(*m)["HeapInuse"] = Gauge(rtm.HeapInuse)
	(*m)["HeapObjects"] = Gauge(rtm.HeapObjects)
	(*m)["HeapReleased"] = Gauge(rtm.HeapReleased)
	(*m)["HeapSys"] = Gauge(rtm.HeapSys)
	(*m)["LastGC"] = Gauge(rtm.LastGC)
	(*m)["MSpanInuse"] = Gauge(rtm.MSpanInuse)
	(*m)["MCacheSys"] = Gauge(rtm.MCacheSys)
	(*m)["MCacheInuse"] = Gauge(rtm.MCacheInuse)
	(*m)["MSpanSys"] = Gauge(rtm.MSpanSys)
	(*m)["NextGC"] = Gauge(rtm.NextGC)
	(*m)["NumForcedGC"] = Gauge(rtm.NumForcedGC)
	(*m)["OtherSys"] = Gauge(rtm.OtherSys)
	(*m)["StackSys"] = Gauge(rtm.StackSys)
	(*m)["StackInuse"] = Gauge(rtm.StackInuse)
	(*m)["TotalAlloc"] = Gauge(rtm.TotalAlloc)
	// mutex.Unlock()
}
func main() {
	wg.Add(2)
	var rtm runtime.MemStats
	var m = make(MonitorMap)
	go func(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
		defer wg.Done()

		var interval = time.Duration(pollInterval) * time.Second
		for {
			<-time.After(interval)
			runtime.ReadMemStats(&rtm)
			NewMonitor(m, rtm)//, mutex)
			PollCount++
			RandomValue = Gauge(rand.Float64())
			// fmt.Println(m)
		}
	}(&m, rtm)
	
	
	go func() {
		defer wg.Done()
		var interval = time.Duration(reportInterval) * time.Second
		// for {

			<-time.After(interval) 
			//PollCount----------------------------------------------------------
			strURL := fmt.Sprintf("%s/update/", url)
			client := resty.New()
			// client.JSONMarshal = json.Marshal
			// client.JSONUnmarshal = json.Unmarshal
			pollCount := int64(PollCount)
			var varMetrics Metrics = Metrics{
				ID: "PollCount",
				MType: "Counter",
				Delta: &pollCount,
			}
			bodyBytes, err := json.Marshal(varMetrics)
			// fmt.Println(bodyBytes, "---------bodyBytes------------") 
			// fmt.Println(string(bodyBytes), "---------string(bodyBytes)------------") 
			
			if err != nil {
				log.Fatalf("Failed marshal json: %s", err)
			}
			var varMetrics1 Metrics
			// client := http.Client{}
			// resp, err := client.Post(strURL, "application/json; charset=UTF-8", bytes.NewBuffer(bodyBytes)) 
			// if err != nil {
			// 	log.Fatalf("Failed get response: %s", err)
			// }
			// if resp.StatusCode == http.StatusOK {
			// 	bodyBytes, _ := io.ReadAll(resp.Body)
			// 	fmt.Println(bodyBytes)
			// 	bodyString := string(bodyBytes)
			// 	fmt.Println(bodyBytes)
			// 	byteString := []byte(bodyString)
			// 	fmt.Println(byteString)
			// 	//fmt.Println(bodyString)
			// 	err = json.Unmarshal(byteString,&varMetrics1)
			// 	if err != nil {
			// 		log.Fatalf("Failed unmarshal response: %s", err)
			// 	}
			// }
			// if err != nil {
			// 	// fmt.Println(responseErr)
        	// 	// panic(err)
			// 	log.Fatalf("Failed sent request: %s", err)
			// 	}
			// var responseErr MyApiError
			// resp, err := client.R().
			// SetError(&responseErr).
			// SetResult(&varMetrics1).
			// SetBody(bodyBytes).
			// Post(strURL)
			// responseBytes, err := io.ReadAll(response.Body)
			// fmt.Println(string(responseBytes))
			// if err != nil {
			// 	log.Fatalf("Failed read body of response: %s, body: %v", err, string(responseBytes))
			// }
			// err = json.Unmarshal(responseBytes , &varMetrics)
			// if err != nil {
			// 	log.Fatalf("Failed unmarshal response: %s", err)
			// }
			// fmt.Println(varMetrics1) 

			// metricsStruct := new(Metrics) 
			var metricsStruct Metrics
			// fmt.Println(response, "------------response-------------") 
			// fmt.Println(response.Body(), "------------response.Body()-------------") 
			// fmt.Println(string(response.Body()), "---------string(response.Body()------------") 
			if err != nil {
				log.Fatalf("Failed unmarshal response: %s", err)
			}
			fmt.Println(varMetrics) 
			fmt.Println(string(bodyBytes), "========string(bodyBytes)============")

			//RandomValue------------------------------------------------------------
			// strURL = fmt.Sprintf("%s/update/", url)
			client = resty.New()
			randomValue := float64(RandomValue)
			varMetrics = Metrics{
				ID: "RandomValue",
				MType: "Gauge",
				Value: &randomValue,
			}
			bodyBytes, err := json.Marshal(varMetrics)
			fmt.Println(string(bodyBytes), "--------------------------------------------") 

			if err != nil {
				log.Fatalf("Failed marshal json: %s", err)
			}
			var varMetrics1 Metrics
			_, err = client.R().
			SetResult(&varMetrics1).
			SetBody(bodyBytes).
			Post(strURL)
			if err != nil {
				log.Fatalf("Failed unmarshall response: %s", err)
			}
			// var metricsStruct Metrics
			// err = json.Unmarshal(response.Body(), &metricsStruct)
			// if err != nil {
			// 	log.Fatalf("Failed unmarshal response: %s", err)
			// }
			fmt.Println(varMetrics1) 
			//---------------------------------------------------------------
			// fmt.Println(response) 
			// response, err = client.R().Post("http://localhost:8080/update/counter/testSetGet33/187")
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println(response)

			// response, err = client.R().Get("http://localhost:8080/value/counter/testSetGet33")
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println(response)

			// n := 0
			// for key, val := range m {
			// 	// n++
			// 	// if n > 1 {
			// 	// 	break
			// 	// }
			// 	strURL := fmt.Sprintf("%s/update/", url)
			// 	client := resty.New()
			// 	val1 := float64(val)
			// 	varMetrics := Metrics{
			// 		ID: key,
			// 		MType: "Gauge",
			// 		Value: &val1,
			// 	}
			// 	bodyBytes, err := json.Marshal(varMetrics)
			// 	if err != nil {
			// 		log.Fatalf("Failed marshal json: %s", err)
			// 	}
			// 	response, err := client.R().
			// 	SetBody(bodyBytes).
			// 	Post(strURL)
			// 	if err != nil {
			// 		log.Fatalf("Failed sent request: %s", err)
			// 	}
			// 	// var metricsStruct Metrics
			// 	err = json.Unmarshal(response.Body(), &metricsStruct)
			// 	if err != nil {
			// 		log.Fatalf("Failed unmarshal response: %s", err)
			// 	}
			// 	fmt.Println(metricsStruct) 
			// 	fmt.Println(response)
			// }
		// }
	}()
	wg.Wait()
}

// func (t Metrics) MarshalJSON() ([]byte, error) {
//     // чтобы избежать рекурсии при json.Marshal, объявляем новый тип
//     type MetricsAlias Metrics
// 	delta1 := Counter(*t.Delta)
// 	value1 := Gauge(*t.Value)
//     aliasValue := struct {
//         MetricsAlias
//         // переопределяем поле внутри анонимной структуры
// 			Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
// 			Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
//     }{
//         // встраиваем значение всех полей изначального объекта (embedding)
//         MetricsAlias: MetricsAlias(t),
//         // задаём значение для переопределённого поля

//         Delta: &delta1,
//         Value: &value1,
//     }

//     return json.Marshal(aliasValue) // вызываем стандартный Marshal
// }

// // UnmarshalJSON реализует интерфейс json.Unmarshaler.
// func (t *Metrics) UnmarshalJSON(data []byte) error {
//     type MetricsAlias Metrics

//     aliasValue := &struct {
//         *MetricsAlias
//         // переопределяем поле внутри анонимной структуры
//      Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
// 		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
//     }{
//         // задаём указатель на целевой объект
//         MetricsAlias: (*MetricsAlias)(t),
//     }

//     // вызываем стандартный Unmarshal
//     if err := json.Unmarshal(data, aliasValue); err != nil {
//         return err
//     }
//     // вручную задаём значение полей Delta и Value
// 	delta1 := int64(*aliasValue.Delta)
// 	value1 := float64(*aliasValue.Value)
//     t.Delta = &delta1
//     t.Value = &value1

//     return nil
// } 
// var varMetrics Metrics
// runtime.ReadMemStats(&rtm)
