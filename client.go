package main 
 
import ( 
"github.com/gorilla/rpc/json"	
  "flag"
  "bytes"
  "strconv"
  "net/http"
  "fmt"

) 

type Reply struct {
    Message string
}


type Args struct {
    SSP string 
    BUD float32                  
    TRADE int                       
}



func main() {

var checkarguments Args

var buyarguments Args
    
flag.Parse() 
method := "Service."+flag.Arg(0)

if flag.Arg(0) == "Checking" {
fmt.Println(method)
id,_ := strconv.Atoi(flag.Arg(1))
checkarguments.TRADE = id
reply,_ := CReq(method, checkarguments)
fmt.Println(reply.Message)
} 


if flag.Arg(0) == "Buying" {
fmt.Println(method)
buyarguments.SSP = flag.Arg(1)
budget,_ := strconv.ParseFloat(flag.Arg(2), 64)
id,_ := strconv.Atoi(flag.Arg(3))
buyarguments.TRADE = id;
buyarguments.BUD = float32(budget)
reply,_ := CReq(method, buyarguments)
fmt.Println(reply.Message)
}
  
} 

func CReq(method string, args Args) (reply Reply, err error) {
req, err := json.EncodeClientRequest(method, args)
res, err := http.Post("http://127.0.0.1:8080/rpc", "application/json", bytes.NewBuffer(req))
  
if err != nil{
fmt.Println("Error in POST request")
fmt.Println(err)
return reply, err
}


if err != nil {
fmt.Println("Error in Client Request")
fmt.Println(err)
return reply, err
}
    
defer res.Body.Close()
err = json.DecodeClientResponse(res.Body, &reply)
return reply, err
}


