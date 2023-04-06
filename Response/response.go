package Resp

import (
	"encoding/json"
	"main/models"
	"net/http"
)


func Response(status string,code int64,message string,data interface{},w http.ResponseWriter){

	var res models.Response
	res.Status = status
	res.Code=code
	res.Message=message
	res.Data=data	

	json.NewEncoder(w).Encode(&res)


}

// Res.Response("Method Not Allowed ",405,"use correct http method","",w)
// Res.Response("Bad gateway",502,er.Error(),"",w)
// Res.Response("OK",200,"token provided successfully","",w)
//Res.Response("Unauthorized",401,"token not valid","",w)

