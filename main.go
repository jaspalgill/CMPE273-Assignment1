package main

import j "encoding/json"

import (
    "net/url"
"github.com/gorilla/rpc/json"
"github.com/gorilla/rpc"
"net/http"
"io/ioutil"
"strings"
"fmt"  
"strconv"
)


type Reply struct{
    Message string
}

type Buy struct{
    tid int
    Stocks []Stock
    BUD float32
    UnvestedAmount float32
}
type Stock struct{
    sbl string
    Share  int
    perc float32
    cost  float32
}
var buy []Buy
type Arguments struct{
    SSP string 
    BUD float32                  
    tid int 
}

type Response struct{
   Query struct{
        Count int `json:"count"`
        Resource struct{
            Field struct {
                 Quote struct {
                                        Name string `json:"name"`
                                        cost string `json:"cost"`
                                        
                                } `json:"quote"`
                        } `json:"field"`
                } `json:"resources"`
        }`json:"query"`
}


type Service struct{}
func GetInput(Arguments *Arguments) ([]string, []float64,[]float64) {
	ssp:=strings.Split(Arguments.Ssp,",")
	sbl :=make([]string,len(ssp))
	perc:=make([]float64,len(ssp))
	amt:=make([]float64,len(ssp))
	
	for i:=range ssp{
		divide:=strings.Split(ssp[i],":")
		sbl[i] = divide[0]
		perc[i] , _=strconv.ParseFloat(strings.TrimSuffix(divide[1],"%"),64)
		amt[i] = Arguments.BUD*perc[i]/100
	}
	return sbl,perc,cost
}

func main() {
    rpcHandler := rpc.NewServer()
    codec := json.NewCodec()
    rpcHandler.RegisterCodec(codec, "application/json")
    rpcHandler.RegisterCodec(codec, "application/json; charset=UTF-8")
    rpcHandler.RegisterService(new(Service), "")
    http.Handle("/rpc", rpcHandler)
    http.ListenAndServe("127.0.0.1:8080", nil)
}

func getcost(sbl string) string{
    queryStr := "select sbl, LasttidcostOnly from yahoo.finance.quote where sbl in ('"+sbl+"')"
    urlPath :=  "http://query.yahooapis.com/v1/public/yql?q="
    urlPath += url.QueryEscape(queryStr)
    urlPath += "&format=json&env=store://datatables.org/alltableswithkeys"
    res, err := http.Get(urlPath)
    if err!=nil {
        fmt.Println("getcost: http.Get",err)
        panic(err)
    }
    defer res.Body.Close()
    body,err := ioutil.ReadAll(res.Body)
    if err!=nil {
        fmt.Println("getcost: ioutil.ReadAll",err)
        panic(err)
    }
    var s Stock
    fmt.Println(string(body[:]))
    err = j.Unmarshal(body, &s)
    if err!=nil {
        fmt.Println("getcost: json.Unmarshal",err)
        panic(err)
    }
    return s.Query.Resouce.Field.Quote.cost
}

func (s *Service) Buying(r *http.Request, Arguments *Arguments, reply *Reply) error {
    BUD := Arguments.BUD
    str := Arguments.SSP
    
    message := ""
    var stock Stock
    var thisbuy 
    var unvested float32
    ssp,perc,amt := GetInputs(Arguments);
    
    stockcost := getcost(ssp)
            cost,_ := strconv.ParseFloat(stockcost, 64)
            costf32 := float32(cost)
            stock.cost = costf32
            perc := strings.Split(ssp, "%")
            percage,_ := strconv.ParseFloat(perc, 64)
            percagef32 := float32(percage)
            stock.percage = percagef32
            tempBUD := BUD*percagef32/100
            share := int(tempBUD/costf32)
            stock.Share = share
            thisbuy.Stocks = append(thisbuy.Stocks, stock)
            if i==0{
                unvested = BUD - costf32*float32(share)
            }else{
                unvested = unvested - costf32*float32(share)
            }
            strShare := strconv.Itoa(share)
            if i==0{
                message = onestock[0]+":"+strShare+":$"+stockcost
            }else{
                message = message+","+onestock[0]+":"+strShare+":$"+stockcost
            }
            result := "tid : "+strconv.Itoa(argument.tid)
    unvestedf64 := float64(unvested)
    result += "\nstocks : " +message+"\nunvestedAmount : " + strconv.FormatFloat(unvestedf64, 'f', 3, 64)
    
    reply.Message = result
    thisbuy.tid = args.tid
    thisbuy.BUD = args.BUD
    thisbuy.UnvestedAmount = unvested
    buy = append(buy,thisbuy)
    return nil
}
       


func (s *Service) Checking(r *http.Request, Arguments *Arguments, reply *Reply) error{
    str := "Stocks : "
     var thisbuy Buy
    flag := false
   
    for i:=0;i<len(buy);i++{
        if buy[i].tid == Arguments.tid{
            flag = true
            thisbuy = buy[i]
            fmt.Println(thisbuy)
        }
    }
    if flag==false{
        str = "tid : "+ strconv.Itoa(Arguments.tid) + " not exists"
        reply.Message = str
        return nil
    }
    var currentMarketValue float32
    for _,v := range thisbuy.Stocks{
        stockcost := getcost(v.sbl)
        cost,_ := strconv.ParseFloat(stockcost, 64)
        costf32 := float32(cost)
        str += v.sbl+":"+strconv.Itoa(v.Share)+":"
        if costf32 > v.cost{
            str += "+"
        }
        if costf32 < v.cost{
            str += "-"
        }
        str += "$"+stockcost+","
        currentMarketValue += float32(v.Share)*costf32
    } 
    str += "\ncurrentmarketValue : " + strconv.FormatFloat(float64(currentMarketValue), 'f', 3, 64)
    str += "\nunvestedamount : " + strconv.FormatFloat(float64(thisbuy.UnvestedAmount), 'f', 3, 64)
    reply.Message = str
    return nil
}




