package tv

import (
    "fmt"
    "encoding/json"
    "time"
    "strings"
    "github.com/Jeffail/gabs"
    "github.com/mikemintang/go-curl"
    "net/http"
)

var PEATIO_BASE_URL = "https://xxx.xxxx.xxx"

type Buffer struct {
        records    map[string][]interface{}
}

// ProcessRecord adds a message to the buffer.
func (buffer *Buffer) AddRecord(key string, record interface{}) {
    if buffer.records == nil {
        buffer.records = make(map[string][]interface{})
    }
        _, ok := buffer.records[key]
        if !ok {
                buffer.records[key] = make([]interface{}, 0)
        }

        buffer.records[key] = append(buffer.records[key], record)
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]interface{}{
        "supports_time":true,
        "supports_search":true,
        "supports_marks":false,
        "supports_group_request":false,
        "supports_timescale_marks": false,
        "supported_resolutions":[]string{"1","5","15","30","60","360","720","D"},
    })
}

func GetTime(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        json.NewEncoder(w).Encode(time.Now().Unix());
}

/* example url parameter : symbols?symbol=BTC-ETH */
func GetSymbol(w http.ResponseWriter, r *http.Request) {
    symbol := r.URL.Query().Get("symbol")
    symbol = strings.ToLower(symbol)
    currency := strings.Split(symbol, "-")
    //symbol = strings.Replace(symbol, "-", "", -1)


    url := PEATIO_BASE_URL + "/api/v2/peatio/public/currencies/" + currency[1]

    headers := map[string]string{
        "Accept":  "application/json",
    }

    req := curl.NewRequest()
    resp, err := req.
        SetHeaders(headers).
        SetUrl(url).
        Get()

    if err != nil {
        fmt.Println(err)
    } else {
        if resp.IsOk() {
            json.NewEncoder(w).Encode(map[string]interface{}{
                "symbol": r.URL.Query().Get("symbol"),
                "ticker": r.URL.Query().Get("symbol"),
                "minmovement": 1,
                "minmovement2": 0,
                "session":"24x7",
                "timezone":"Asia/Bangkok",
                "has_intraday":true,
                "description": r.URL.Query().Get("symbol"),
                //"supported_resolutions": []string{"1","5","15","30","60","480","D","W"},
                "type":"stock",
                "currency_code":"BTC",
                "exchange-listed":"",
                "volume_precision": 8,
                "pointvalue":1,
                "name": r.URL.Query().Get("symbol"),
                "exchange-traded":r.URL.Query().Get("symbol"),
                "pricescale":1,
                "data_status":"`delayed_streaming`",
                "has_no_volume": false,
            })
        } else {
            fmt.Println(resp.Raw)
        }
    }
}

/* example url parameter : history?symbol=BTC-ETH&period=30&time_from=1551883158&time_to=1554475218 */
func GetHistory(w http.ResponseWriter, r *http.Request) {
    var limit string = "440"

    symbol := r.URL.Query().Get("symbol")
    resolution := r.URL.Query().Get("resolution")
    //from := r.URL.Query().Get("from")
    //to := r.URL.Query().Get("to")

    symbol = strings.ToLower(symbol)
    symbol = strings.Replace(symbol, "-", "", -1)

    if resolution == "360"{
        limit = "800"
    }
    if resolution == "1D"{
        resolution = "1440"
        limit = "10"
    }
    url := PEATIO_BASE_URL + "/api/v2/peatio/public/markets/"+ symbol +"/k-line"

    headers := map[string]string{
        "Accept":  "application/json",
    }

    queries := map[string]string{
        //"market": symbol,
        "limit": limit,
        //"timestamp": to * 1000
        "period": resolution,
    }

    req := curl.NewRequest()
    resp, err := req.
        SetHeaders(headers).
        SetUrl(url).
        SetQueries(queries).
        Get()

    if err != nil {
        fmt.Println(err)
    } else {
        if resp.IsOk() {
            var t []float64
            var c []float64
            var o []float64
            var h []float64
            var l []float64
            var v []float64

            jsonParsed, _ := gabs.ParseJSON([]byte(resp.Body))

            children, _ := jsonParsed.Children()

              for _, child := range children {
                      t = append(t, child.Index(0).Data().(float64))
                      o = append(o, child.Index(1).Data().(float64))
                      h = append(h, child.Index(2).Data().(float64))
                      l = append(l, child.Index(3).Data().(float64))
                      c = append(c, child.Index(4).Data().(float64))
                      v = append(v, child.Index(5).Data().(float64))
              }

              json.NewEncoder(w).Encode(map[string]interface{}{
                  "s":"ok",
                  "t": t,
                  "o": o,
                  "h": h,
                  "c": c,
                  "l": l,
                  "v": v,
              })
        } else {
            fmt.Println(resp.Raw)
        }
    }
}

func GetQuotes(w http.ResponseWriter, r *http.Request) {
    var a []int
    json.NewEncoder(w).Encode(map[string]interface{}{
        "s":"ok",
        "d":a,
    })
}
