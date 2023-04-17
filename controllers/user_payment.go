package cont

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	Res "main/Response"
	"main/db"
	"main/models"
	"main/utils"
	"net/http"
	"os"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	razorpay "github.com/razorpay/razorpay-go"
	
)

type Pagevar struct {
	Orderid string
}
var pagevar Pagevar


type PaymentStatusUpdate struct {
	Entity    string   `json:"entity"`
	AccountID string   `json:"account_id"`
	Event     string   `json:"event"`
	Contains  []string `json:"contains"`
	Payload   struct {
		Payment struct {
			Entity struct {
				ID             string `json:"id"`
				Entity         string `json:"entity"`
				Amount         int    `json:"amount"`
				Currency       string `json:"currency"`
				Status         string `json:"status"`
				OrderID        string `json:"order_id"`
				InvoiceID      string `json:"invoice_id"`
				International  bool   `json:"international"`
				Method         string `json:"method"`
				AmountRefunded int    `json:"amount_refunded"`
				RefundStatus   string `json:"refund_status"`
				Captured       bool   `json:"captured"`
				Description    string `json:"description"`
				CardID         string `json:"card_id"`
				Bank           string `json:"bank"`
				Wallet         string `json:"wallet"`
				Vpa            string `json:"vpa"`
				Email          string `json:"email"`
				Contact        string `json:"contact"`
				Notes          struct {
					Address string `json:"address"`
				} `json:"notes"`
				Fee              int    `json:"fee"`
				Tax              int    `json:"tax"`
				ErrorCode        string `json:"error_code"`
				ErrorDescription string `json:"error_description"`
				ErrorSource      string `json:"error_source"`
				ErrorStep        string `json:"error_step"`
				ErrorReason      string `json:"error_reason"`
				AcquirerData     struct {
					BankTransactionID string `json:"bank_transaction_id"`
				} `json:"acquirer_data"`
				CreatedAt  int64 `json:"created_at"`
				BaseAmount int   `json:"base_amount"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
	CreatedAt int64 `json:"created_at"` 
}

//payment response struct

type paymentresponse struct {

	paymentID string	
	Amount int	
	Status string	
	orderId string

}
var paymentRes paymentresponse




// @Description Initiate Payment
// @Accept json
// @Produce json
// @Param  details body string true "enter membership_name SchemaExample({ "membership_name":"Individual"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /make-payment [post]
func MakepaymentHandler(w http.ResponseWriter, r *http.Request) {

	utils.SetHeader(w)


	if r.Method != http.MethodPost {
		
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}
	var user models.User
	var membership models.Memberships

	json.NewDecoder(r.Body).Decode(&membership)

	membership_name:=membership.Membership_name

	fmt.Println("memship_name: ",membership_name)
	
	

	err1 := validation.Validate(membership_name,validation.In("Duo","Individual"))
	
	if err1!=nil{

		Res.Response("Bad Request",400,err1.Error(),"",w)
		return
	}

	membership.Membership_name=membership_name
	


	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	user.User_id=userDetails.User_id
	

	
	var payment models.Payments
	payment.Membership_name=membership.Membership_name
	payment.User_id=user.User_id

	er:=db.DB.Create(&payment).Error
	if er!=nil {
		Res.Response("server error",500,er.Error(),"",w)
	}



	order_creation(user.User_id,membership,w)

	

}
func order_creation(user_id string,membership models.Memberships ,writer http.ResponseWriter){

	//ORDER CREATION------------------------------------------------------>

	var memship models.Memberships
	memship.Membership_name=membership.Membership_name
	er:=db.DB.Where("membership_name=?",membership.Membership_name).First(&memship).Error
	if er!=nil {
		Res.Response("server error",500,er.Error(),"",writer)
	}
	client := razorpay.NewClient("rzp_test_MLjFMJxEVuaLjd", os.Getenv("Razorpay_Key"))

	data := map[string]interface{}{
		"amount":   memship.Price,        
		"currency": "INR",
		"notes": map[string]interface{}{

        "subscription":membership.Membership_name,
		},
	}
	Body, err := client.Order.Create(data, nil)

	if err != nil {
		fmt.Println("error in order create request")
		Res.Response("access denied",402,err.Error(),"",writer)
		return
	}

	order_id := Body["id"].(string)

	pagevar.Orderid = order_id

	
// Template
	t, err := template.ParseFiles("controllers/payment.html")
	if err!=nil{
		fmt.Println("template parsing error",err)
	}

	err = t.Execute(writer, pagevar)
	if err != nil {

		fmt.Println("template executing error", err)
	}



	fmt.Println("body response", Body)
	fmt.Println("")

		

	//update during order creation
		var payment models.Payments
		
		payment.Order_id=order_id
		
		Er:=db.DB.Where("user_id=?",user_id).Updates(&payment).Error
		if Er!=nil{
			Res.Response("server error",500,er.Error(),"",writer)
		}


}



func Razorpay_Response(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Response function called./....")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	


	

	
	var response PaymentStatusUpdate
	json.Unmarshal(body, &response)
	fmt.Println("")

	fmt.Println("id", response.Payload.Payment.Entity.ID)
	fmt.Println("order_id",response.Payload.Payment.Entity.OrderID)
	fmt.Println("amount", (response.Payload.Payment.Entity.Amount)/100)
	fmt.Println("status", response.Payload.Payment.Entity.Status)

	

	//put all the response data to paymentresponse struct
	

	var payment models.Payments
	err1:=db.DB.Where("order_id=?",response.Payload.Payment.Entity.OrderID).First(&payment).Error
	if err1!=nil{
		fmt.Println("error is ",err1)
		Res.Response("server error",500,err1.Error(),"",w)
	}
	//updates after response
	payment.Payment_id=response.Payload.Payment.Entity.ID
	payment.Status=response.Payload.Payment.Entity.Status
	payment.Time=time.Now()

	if payment.Status=="captured"{

		var user models.User

		db.DB.Where("user_id=?",payment.User_id).First(&user)

		user.Membership=payment.Membership_name

		db.DB.Where("user_id=?",payment.User_id).Updates(&user)


		
	}

	fmt.Println("Payments is;",payment)
	dbErr:=db.DB.Where("order_id=?",response.Payload.Payment.Entity.OrderID).Updates(&payment).Error
	if dbErr!=nil{
		fmt.Println("db error",dbErr)
		Res.Response("Bad gateway",500,dbErr.Error(),"",w)
	}

	//Signature verification
	signature := r.Header.Get("X-Razorpay-Signature")
	fmt.Println("signature", signature)
	if !VerifyWebhookSignature(body, signature, os.Getenv("API_SecretKey")) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		Res.Response("Unauthorized",401,"Invalid signature","",w)
		return
	} else {

		fmt.Println("signature verified")
		Res.Response("OK",200,"Success","",w)
	}


	



	

}

func VerifyWebhookSignature(body []byte, signature string, secret string) bool {


	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	hash := hex.EncodeToString(h.Sum(nil))

	return hash == signature
}
