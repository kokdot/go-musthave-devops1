package main

import (
	// "oleiade/reflections"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"github.com/oleiade/reflections"
	// "github.com/kokdot/go-musthave-devops/internal/store"
)

// type RunTimeMemStats struct {
// 	Alloc Guage
// 	TotalAlloc Guage
// 	Sys Guage
// 	Mallocs Guage
// }
// type FuncMem func(rtm *runtime.MemStats, value string) *runtime.MemStats {
// 	value: func(rtm runtime.MemStats) {

// }

// func ()

// var MonMap MonitorMap = MonitorMap{
// 	"Alloc": Gauge(1.3),
// 	"TotalAlloc": Gauge(1.4),
// 	"Sys": Gauge(1.5),
// 	"Mallocs": Gauge(1.6),
// }

func TestHandler(t *testing.T) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	var fields []string
	var waitMap MonitorMap = make(MonitorMap)

	// Fields will list every structure exportable fields.
	// Here, it's content would be equal to:
	// []string{"FirstField", "SecondField", "ThirdField"}
	fields, _ = reflections.Fields(rtm)
	fmt.Println(fields)
	for _, field := range fields {
		RandomValue := rand.Uint64()
		fmt.Println(RandomValue)
		err := reflections.SetField(&rtm, field, RandomValue)
		if err != nil {
			log.Fatalf("Failed set field: %s", err)
		}
		// waitMap[field] = Guage(float64(RandomValue))
	}
	// fmt.Println(waitMap)
	fmt.Println(rtm)
	_ = waitMap

}
// type want struct {
// 		monMap MonitorMap
// 		runTimeMemStats RunTimeMemStats
// 	}
// 	tests := []struct {
// 		name   string
// 		want   want
// 	}{
// 		{
// 			name: "norm gauge",
// 			want: {
// 				monMap: MonitorMap{
// 					"Alloc": Gauge(1.3),
// 					"TotalAlloc": Gauge(1.4),
// 					"Sys": Gauge(1.5),
// 					"Mallocs": Gauge(1.6),
// 				},
// 				runTimeMemStats: RunTimeMemStats{
// 					Alloc: Gauge(1.3),
// 					TotalAlloc: Gauge(1.4),
// 					Sys: Gauge(1.5),
// 					Mallocs: Gauge(1.6),
// 				},
// 			}
// 		}
// 	}

// type want struct {
// 		StatusCode  int
// 		contentType string
// 		result      string
// 	}
// 	tests := []struct {
// 		name   string
// 		want   want
// 		url    string
// 		method string
// 	}{
// 		{
// 			name: "counter norm",
// 			want: want{
// 				StatusCode:  http.StatusOK,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 			url:    "/update/counter/PollCount/5",
// 			method: http.MethodPost,
// 		},
// 		{
// 			name: "counter error",
// 			want: want{
// 				StatusCode:  http.StatusBadRequest,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 			url:    "/update/counter/PollCount/none",
// 			method: http.MethodPost,
// 		},
// 		{
// 			name: "gauge norm",
// 			want: want{
// 				StatusCode:  http.StatusOK,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 			url:    "/update/gauge/Alloc/3.6",
// 			method: http.MethodPost,
// 		},
// 		{
// 			name: "gauge error",
// 			want: want{
// 				StatusCode:  http.StatusBadRequest,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 			url:    "/update/gauge/Alloc/none",
// 			method: http.MethodPost,
// 		},
// 		{
// 			name: "no counter no gauge",
// 			want: want{
// 				StatusCode:  http.StatusNotImplemented,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 			url:    "/update/error/PollCount/100",
// 			method: http.MethodPost,
// 		},
// 		{
// 			name: "default",
// 			want: want{
// 				StatusCode:  http.StatusMethodNotAllowed,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 			url:    "/",
// 			method: http.MethodPost,
// 		},
// 		{
// 			name: "get counter",
// 			want: want{
// 				StatusCode:  http.StatusOK,
// 				contentType: "text/plain; charset=utf-8",
// 				result:      "5",
// 			},
// 			method: http.MethodGet,
// 			url:    "/value/counter/PollCount",
// 		},
// 		{
// 			name: "get gauge",
// 			want: want{
// 				StatusCode:  http.StatusOK,
// 				contentType: "text/plain; charset=utf-8",
// 				result:      "3.6",
// 			},
// 			url:    "/value/gauge/Alloc",
// 			method: http.MethodGet,
// 		},
// 	}