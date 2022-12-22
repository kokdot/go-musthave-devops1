package main

import (
	// "bytes"
	"fmt"
	"log"
	"encoding/json"
	"github.com/caarlos0/env/v6"
	// "net/http"
	"runtime"
	"flag"
	"sync"
	"time"
	"math/rand"
	// "io"
	"github.com/go-resty/resty/v2"
)

const (
	Url            = "127.0.0.1:8080"
	PollInterval   = 2
	ReportInterval = 10
)
type Config struct {
    Address  string 		`env:"ADDRESS" envDefault:"127.0.0.1:8080"`
    ReportInterval int	 `env:"REPORT_INTERVAL" envDefault:"10"`
    PollInterval int	 `env:"POLL_INTERVAL" envDefault:"2"`
}

// var mutex *sync.RWMutex
var wg sync.WaitGroup 

type Gauge float64
type Counter int64
type MonitorMap map[string]Gauge
var PollCount Counter
var RandomValue Gauge

var( 
	pollIntervalReal = PollInterval
	reportIntervalReal = ReportInterval
	urlReal = Url
	// urlReal = "http://" + Url
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
func NewMonitor(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
	// runtime.ReadMemStats(&rtm)
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
func onboarding() {
	var cfg Config

    err := env.Parse(&cfg)
    if err != nil {
        log.Fatal(err)
    }
	urlReal	= cfg.Address
	reportIntervalReal	= cfg.ReportInterval
	pollIntervalReal	= cfg.PollInterval

	urlRealPtr := flag.String("a", "127.0.0.1:8080", "ip adddress of server")
    reportIntervalRealPtr := flag.Int("r", 10, "interval of perort")
    pollIntervalRealPtr := flag.Int("p", 2, "interval of poll")

    flag.Parse()
	if urlReal == Url {
        urlReal = *urlRealPtr
	}
	urlReal = "http://" + urlReal
	if reportIntervalReal == ReportInterval {
		reportIntervalReal = *reportIntervalRealPtr
	}
	if pollIntervalReal == PollInterval {
		pollIntervalReal = *pollIntervalRealPtr
	}

	// fmt.Printf("Current address is %s\n", cfg.Address)
    // fmt.Printf("Current report_interval is %d\n", cfg.ReportInterval)
    // fmt.Printf("Current poll_interval is %d\n", cfg.PollInterval)

}
func main() {
	wg.Add(2)
	onboarding()

	var rtm runtime.MemStats
	var m = make(MonitorMap)
	go func(m *MonitorMap, rtm runtime.MemStats) {//}, mutex *sync.RWMutex) {
		defer wg.Done()

		var interval = time.Duration(pollIntervalReal) * time.Second
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
		var interval = time.Duration(reportIntervalReal) * time.Second
		for {

			<-time.After(interval) 
			// //text----------------------------------------------------------

			// //PollCount Update-------------------------------------------------
			// client := resty.New()
			// strURL1 := fmt.Sprintf("%s/update/counter/%s/%v", urlReal, "PollCount", PollCount)
			// response, err := client.R().Post(strURL1)
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println("PollCount response: ", response) 

			
			// // //PollCount get --------------------------------------------------------
			// response, err = client.R().Get("http://localhost:8080/value/counter/PollCount")
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println("PollCount Get response: ", response)

			// //RandomValue Update--------------------------------------------------
			// strURL1 = fmt.Sprintf("%s/update/gauge/%s/%v", urlReal, "RandomValue", RandomValue)
			// client = resty.New()
			// response, err = client.R().Post(strURL1)
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println("RandomValue Update response: ", response) 

			
			// //RandomValue get --------------------------------------------------------
			// response, err = client.R().Get("http://localhost:8080/value/gauge/RandomValue")
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println("RandomValue Get response: ", response)

			// //testSetGet33 Update ---------------------------------------------------
			// response, err = client.R().Post("http://localhost:8080/update/counter/testSetGet33/187")
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println("testSetGet33 Update response: ", response)

			// //testSetGet33 get --------------------------------------------------------
			// response, err = client.R().Get("http://localhost:8080/value/counter/testSetGet33")
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println("testSetGet33 Get response: ", response)
			// //text----------------------------------------------------------
			
			
			//PollCount----------------------------------------------------------
			strURL := fmt.Sprintf("%s/update/", urlReal)
			// strURL := fmt.Sprintf("%s/update/counter/%s/%v", url, "PollCount", PollCount)
			var varMetrics Metrics = Metrics{
				ID: "PollCount",
				MType: "Counter",
				Delta: &PollCount,
			}
			bodyBytes, err := json.Marshal(varMetrics)
			if err != nil {
				log.Fatalf("Failed marshal json: %s", err)
			}
			// var metricsStruct Metrics
			client := resty.New()
			_, err = client.R().
			SetResult(&varMetrics).
			SetBody(bodyBytes).
			Post(strURL)
			if err != nil {
				log.Fatalf("Failed unmarshall response PollCount: %s", err)
			}
			fmt.Println("PollCount: ", *varMetrics.Delta) 

			// //--------------------
			// client = resty.New()
			// _, _ = client.R().
			// SetResult(&varMetrics).
			// SetBody(bodyBytes).
			// Post(strURL)
			// //-------------------------
			// client = resty.New()
			// _, _ = client.R().
			// SetResult(&varMetrics).
			// SetBody(bodyBytes).
			// Post(strURL)
			//  if err != nil {
			// 	log.Fatalf("Failed unmarshal response: %s", err)
			// }
			// fmt.Println(int(*varMetrics.Delta)) 
			//RandomValue----------------------------------------------------------
			client = resty.New()
			varMetrics = Metrics{
				ID: "RandomValue",
				MType: "Gauge",
				Value: &RandomValue,
			}
			bodyBytes, err = json.Marshal(varMetrics)
			if err != nil {
				log.Fatalf("Failed marshal json: %s", err)
			}
			var metricsStruct Metrics
			_, err = client.R().
			SetResult(&metricsStruct).
			SetBody(bodyBytes).
			Post(strURL)
			if err != nil {
				log.Fatalf("Failed unmarshall response RandomValue: %s", err)
			}
			fmt.Println("RandomValue: ", *varMetrics.Value) 
			//RandomValueGet---------------------------------------------------
			// strURLGet := fmt.Sprintf("%s/value/", urlReal)
			// var metricsStructGet Metrics
			// client = resty.New()
			// // randomValue := float64(RandomValue)
			// varMetrics = Metrics{
			// 	ID: "RandomValue",
			// 	MType: "Gauge",
			// }
			// bodyBytes, err = json.Marshal(varMetrics)
			// if err != nil {
			// 	log.Fatalf("Failed marshal json: %s", err)
			// }
			// // var varMetrics1 Metrics
			// _, err = client.R().
			// SetResult(&metricsStructGet).
			// SetBody(bodyBytes).
			// Post(strURLGet)
			// if err != nil {
			// 	log.Fatalf("Failed unmarshall response: %s", err)
			// }
			// fmt.Println("RandomValueGet:  ", float64(*metricsStructGet.Value)) 
			// Gauge ----------------------------------------------------------
			// // n := 0
			for key, val := range m {
				// n++
				// if n > 1 {
				// 	break
				// }
				client = resty.New()
				varMetrics = Metrics{
					ID: key,
					MType: "Gauge",
					Value: &val,
				}
				bodyBytes, err = json.Marshal(varMetrics)
				if err != nil {
					log.Fatalf("Failed marshal json: %s", err)
				}
				_, err = client.R().
				// SetResult(&metricsStruct).
				// ForceContentType("application/json").
				SetBody(bodyBytes).
				Post(strURL)
				if err != nil {
					log.Fatalf("Failed unmarshall response: %s", err)
				}
				fmt.Println(varMetrics) 
			}

			// response, err := client.R().Post(strURL)
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
			// fmt.Println(response) 
			// strURL = fmt.Sprintf("%s/update/gauge/%s/%v", url, "RandomValue", RandomValue)
			// client = resty.New()
			// response, err = client.R().Post(strURL)
			// if err != nil {
			// 	log.Fatalf("Failed sent request: %s", err)
			// }
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

			// // n := 0
			// for key, val := range m {
			// 	// n++
			// 	// if n > 1 {
			// 	// 	break
			// 	// }
			// 	strURL := fmt.Sprintf("%s/update/gauge/%s/%v", url, key, val)
			// 	client := resty.New()
			// 	response, err := client.R().Post(strURL)
			// 	if err != nil {
			// 		log.Fatalf("Failed sent request: %s", err)
			// 	}
			// 	fmt.Println(response)
			// }
		}
	}()
	wg.Wait()
}
