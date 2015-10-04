package main

import j "encoding/json"

import (
"github.com/gorilla/rpc/json"
"github.com/gorilla/rpc"
"io/ioutil"
"strings"
"strconv"
"net/http"
"net/url"
"fmt"   
)

type Stock struct{
    Symbol string
    Price  float32
    Share  int
    Percentage float32
}

type Reply struct{
    Message string
}

type BuyingResult struct{
    TRADE int
    Stocks []Stock
    BUD float32
    UnvestedAmount float32
}

type Args struct{
    SSP string 
    BUD float32                  
    TRADE int 
}

var buy []BuyingResult

type StockInfo struct {
Query struct {
Count   int                          `json:"count"`
Results struct {
Quote struct {
LastTradePriceOnly    string `json:"lasttradepriceonly"`
Symbol                string `json:"symbol"`
}                            `json:"quote"`
}                            `json:"results"`
}                            `json:"query"`
}

type Service struct{}

func main() {
    rpcHandler := rpc.NewServer()
    codec := json.NewCodec()
    rpcHandler.RegisterCodec(codec, "application/json")
    rpcHandler.RegisterCodec(codec, "application/json; charset=UTF-8")
    rpcHandler.RegisterService(new(Service), "")
    http.Handle("/rpc", rpcHandler)
    http.ListenAndServe("127.0.0.1:8080", nil)
}

func getPrice(Symbol string) string{
    queryStr := "select symbol, LastTradePriceOnly from yahoo.finance.quote where symbol in ('"+Symbol+"')"
    urlPath :=  "http://query.yahooapis.com/v1/public/yql?q="
    urlPath += url.QueryEscape(queryStr)
    urlPath += "&format=json&env=store://datatables.org/alltableswithkeys"
    res, err := http.Get(urlPath)
    if err!=nil {
        fmt.Println("getPrice: http.Get",err)
        panic(err)
    }
    defer res.Body.Close()
    body,err := ioutil.ReadAll(res.Body)
    if err!=nil {
        fmt.Println("getPrice: ioutil.ReadAll",err)
        panic(err)
    }
    var s StockInfo
    fmt.Println(string(body[:]))
    err = j.Unmarshal(body, &s)
    if err!=nil {
        fmt.Println("getPrice: json.Unmarshal",err)
        panic(err)
    }
    return s.Query.Results.Quote.LastTradePriceOnly
}

func (s *Service) Buying(r *http.Request, args *Args, reply *Reply) error {
    str := args.SSP
    budget := args.BUD
    message := ""
    var unvested float32
    var stock Stock
    var thisbuy BuyingResult
    if strings.Contains(str, ","){
        splitstr := strings.Split(str, ",")
        for i:=0;i<len(splitstr);i++{
            onestock := strings.Split(splitstr[i], ":")
            stock.Symbol = onestock[0]
            stockprice := getPrice(onestock[0])
            price,_ := strconv.ParseFloat(stockprice, 64)
            pricef32 := float32(price)
            stock.Price = pricef32
            percent := strings.Split(onestock[1], "%")
            percentage,_ := strconv.ParseFloat(percent[0], 64)
            percentagef32 := float32(percentage)
            stock.Percentage = percentagef32
            tempBudget := budget*percentagef32/100
            share := int(tempBudget/pricef32)
            stock.Share = share
            thisbuy.Stocks = append(thisbuy.Stocks, stock)
            if i==0{
                unvested = budget - pricef32*float32(share)
            }else{
                unvested = unvested - pricef32*float32(share)
            }
            strShare := strconv.Itoa(share)
            if i==0{
                message = onestock[0]+":"+strShare+":$"+stockprice
            }else{
                message = message+","+onestock[0]+":"+strShare+":$"+stockprice
            }
        }
    }else{
        onestock := strings.Split(str, ":")
        stock.Symbol = onestock[0]
        stockprice := getPrice(onestock[0])
        price,_ := strconv.ParseFloat(stockprice, 64)
        pricef32 := float32(price)
        stock.Price = pricef32
        percent := strings.Split(onestock[1], "%")
        percentage,_ := strconv.ParseFloat(percent[0], 64)
        percentagef32 := float32(percentage)
        stock.Percentage = percentagef32
        tempBudget := budget*percentagef32/100
        share := int(tempBudget/pricef32)
        stock.Share = share
        thisbuy.Stocks = append(thisbuy.Stocks, stock)
        unvested = budget - pricef32*float32(share)
        strShare := strconv.Itoa(share)
        message = onestock[0]+":"+strShare+":$"+stockprice
    }
    result := "tradeId : "+strconv.Itoa(args.TRADE)
    unvestedf64 := float64(unvested)
    result += "\nstocks : " +message+"\nunvestedAmount : " + strconv.FormatFloat(unvestedf64, 'f', 3, 64)
    reply.Message = result
    thisbuy.TRADE = args.TRADE
    thisbuy.BUD = args.BUD
    thisbuy.UnvestedAmount = unvested
    buy = append(buy,thisbuy)
    return nil
}

func (s *Service) Checking(r *http.Request, args *Args, reply *Reply) error{
    str := "Stocks : "
    isTrue := false
    var thisbuy BuyingResult
    for i:=0;i<len(buy);i++{
        if buy[i].TRADE == args.TRADE{
            isTrue = true
            thisbuy = buy[i]
            fmt.Println(thisbuy)
        }
    }
    if isTrue==false{
        str = "Trade : "+ strconv.Itoa(args.TRADE) + " not exists"
        reply.Message = str
        return nil
    }
    var currentMarketValue float32
    for _,v := range thisbuy.Stocks{
        stockprice := getPrice(v.Symbol)
        price,_ := strconv.ParseFloat(stockprice, 64)
        pricef32 := float32(price)
        str += v.Symbol+":"+strconv.Itoa(v.Share)+":"
        if pricef32 > v.Price{
            str += "+"
        }
        if pricef32 < v.Price{
            str += "-"
        }
        str += "$"+stockprice+","
        currentMarketValue += float32(v.Share)*pricef32
    } 
    str += "\ncurrentmarketValue : " + strconv.FormatFloat(float64(currentMarketValue), 'f', 3, 64)
    str += "\nunvestedamount : " + strconv.FormatFloat(float64(thisbuy.UnvestedAmount), 'f', 3, 64)
    reply.Message = str
    return nil
}




