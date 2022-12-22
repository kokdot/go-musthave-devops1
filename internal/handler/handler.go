package handler

import (
	// "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "strconv"

	// "github.com/go-chi/chi/v5"
	"github.com/kokdot/go-musthave-devops1/internal/store"
	// "golang.org/x/text/cases"
)

type key int
 
const (
    nameDataKey key = iota
    valueDataKey
)
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
// func (t Metrics) MarshalJSON() ([]byte, error) {
//     // чтобы избежать рекурсии при json.Marshal, объявляем новый тип
//     type MetricsAlias Metrics
// 	delta1 := store.Counter(*t.Delta)
// 	value1 := store.Gauge(*t.Value)
//     aliasValue := struct {
//         MetricsAlias
//         // переопределяем поле внутри анонимной структуры
// 			Delta *store.Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
// 			Value *store.Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
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
var m store.Repo
var ms = new(store.MemStorage)

func init() {
	ms.GaugeMap = make(store.GaugeMap)
	ms.CounterMap = make(store.CounterMap)
	m = ms
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
        // fmt.Println(string(r, "-----------------------------------------------------------")

    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    var metrics Metrics
    err = json.Unmarshal(bodyBytes, &metrics)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    switch metrics.MType  {
    case "Gauge":
        m.SaveGaugeValue(metrics.ID, store.Gauge(*metrics.Value))
        w.Header().Set("content-type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", bodyBytes) 
    case "Counter":
        m.SaveCounterValue(metrics.ID, store.Counter(*metrics.Delta))
        delta, _ := m.GetCounterValue(metrics.ID)
        *metrics.Delta = int64(delta)
        bodyBytes, err := json.Marshal(metrics)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        fmt.Println(string(bodyBytes), "-----------------------------------------------------------")
        fmt.Println(bodyBytes, "-----------------------------------------------------------")
        w.Header().Set("content-type", "application/json; charset=UTF-8")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", bodyBytes) 
    default:
        http.Error(w, err.Error(), 500)
        return
    }
}

func GetValue(w http.ResponseWriter, r *http.Request) {
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    var metrics Metrics
    err = json.Unmarshal(bodyBytes, &metrics)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    switch metrics.MType  {
    case "Gauge":
        gaugeValue, err := m.GetGaugeValue(metrics.ID)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        *metrics.Value = float64(gaugeValue)
        bodyBytes, err := json.Marshal(metrics)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", bodyBytes) 
    case "Counter":
        delta, _ := m.GetCounterValue(metrics.ID)
        *metrics.Delta = int64(delta)
        bodyBytes, err := json.Marshal(metrics)
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        w.Header().Set("content-type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "%v", bodyBytes) 
    default:
        http.Error(w, err.Error(), 500)
        return
    }
}
func GetAll(w http.ResponseWriter, r *http.Request) {
    metricSlise := make([]Metrics, 0)
    gaugeMap, counterMap := m.GetAllValuesJson()
    for key, val := range gaugeMap {
        val1 := float64(val)
        metricSlise = append(metricSlise, Metrics{
            ID: key,
            MType: "Gauge",
            Value: &val1,
        })
    }
     for key, val := range counterMap {
        delta1 := int64(val)
        metricSlise = append(metricSlise, Metrics{
            ID: key,
            MType: "Counter",
            Delta: &delta1,
        })
    }
     bodyBytes, err := json.Marshal(metricSlise)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.Header().Set("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "%v", bodyBytes) 
}
// func PostCounterCtx(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         var nameData string
//         var valueData int

// 		nameDataStr := chi.URLParam(r, "nameData")
// 		valueDataStr := chi.URLParam(r, "valueData")

//         if nameDataStr == "" || valueDataStr == "" {
//             w.Header().Set("content-type", "text/plain; charset=utf-8")
// 		    w.WriteHeader(http.StatusNotFound)
//             fmt.Fprint(w, "http.StatusNotFound")
//             return
//         }
//         nameData = nameDataStr
//         valueData, err := strconv.Atoi(valueDataStr)
//         if err != nil {
//             w.Header().Set("content-type", "text/plain; charset=utf-8")
// 			w.WriteHeader(http.StatusBadRequest)
//             fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }

// 		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
// 		ctx = context.WithValue(ctx, valueDataKey, valueData)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
// func GetCtx(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         var nameData string
// 		nameDataStr := chi.URLParam(r, "nameData")
//         if nameDataStr == "" {
//             w.Header().Set("content-type", "text/plain; charset=utf-8")
// 		    w.WriteHeader(http.StatusNotFound)
//             fmt.Fprint(w, "line: 115; http.StatusNotFound")
//             return
//         }
//         nameData = nameDataStr

// 		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func PostGaugeCtx(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         var nameData string
//         var valueData float64

// 		nameDataStr := chi.URLParam(r, "nameData")
// 		valueDataStr := chi.URLParam(r, "valueData")

//         if nameDataStr == "" || valueDataStr == "" {
//             w.Header().Set("content-type", "text/plain; charset=utf-8")
// 		    w.WriteHeader(http.StatusNotFound)
//             fmt.Fprint(w, "http.StatusNotFound")
//             return
//         }
//         nameData = nameDataStr
//         valueData, err := strconv.ParseFloat(valueDataStr, 64)
//         if err != nil {
//             w.Header().Set("content-type", "text/plain; charset=utf-8")
// 			w.WriteHeader(http.StatusBadRequest)
//             fmt.Fprint(w, "http.StatusBadRequest")
//             return
//         }

// 		ctx := context.WithValue(r.Context(), nameDataKey, nameData)
// 		ctx = context.WithValue(ctx, valueDataKey, valueData)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
// func PostUpdateCounter(w http.ResponseWriter, r *http.Request) {
// 	valueData := r.Context().Value(valueDataKey).(int)
// 	nameData := r.Context().Value(nameDataKey).(string)
// 	// fmt.Println("__________________________________", m)
//     m.SaveCounterValue(nameData, store.Counter(valueData))
// 	w.Header().Set("content-type", "text/plain; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
//     fmt.Fprint(w, "http.StatusOK")
// }
// func PostUpdateGauge(w http.ResponseWriter, r *http.Request) {
// 	// valueData := r.Context().Value(valueDataKey).(float64)
// 	// nameData := r.Context().Value(nameDataKey).(string)
//     m.SaveGaugeValue(nameData, store.Gauge(valueData))
// 	w.Header().Set("content-type", "text/plain; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
//     fmt.Fprint(w, "http.StatusOK")
// }
// func GetCounter(w http.ResponseWriter, r *http.Request) {
//     nameData := r.Context().Value(nameDataKey).(string)
//     n, err := m.GetCounterValue(nameData)
//     if err != nil {
//         w.Header().Set("content-type", "text/plain; charset=utf-8")
//         w.WriteHeader(http.StatusNotFound)
//         fmt.Fprint(w, "line: 175; http.StatusNotFound")
//     } else {
//         w.Header().Set("content-type", "text/plain; charset=utf-8")
//         w.WriteHeader(http.StatusOK)
//         fmt.Fprintf(w, "%v", n)
//     }
// }
// func GetGauge(w http.ResponseWriter, r *http.Request) {
//     nameData := r.Context().Value(nameDataKey).(string)
//     n, err := m.GetGaugeValue(nameData)
//     if err != nil {
//         w.Header().Set("content-type", "text/plain; charset=utf-8")
//         w.WriteHeader(http.StatusNotFound)
//         fmt.Fprint(w, "line: 188; http.StatusNotFound")
//     } else {
//         w.Header().Set("content-type", "text/plain; charset=utf-8")
//         w.WriteHeader(http.StatusOK)
//         fmt.Fprintf(w, "%v", n)
//     }    
// }
