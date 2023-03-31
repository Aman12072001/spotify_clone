package cont

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/models"
	cons "main/utils"
	"net/http"

	// "github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)



func Create_Payment(w http.ResponseWriter, r *http.Request){

	stripe.Key=cons.Stripe_Key

	//get the membership name from body
	var membership models.Memberships
	json.NewDecoder(r.Body).Decode(&membership)
	db.DB.Where("membership_name=?",membership.Membership_name).First(&membership)
	



	params := &stripe.PaymentIntentParams{
	Amount: stripe.Int64(int64(membership.Price*100)),
	//   AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
	//     Enabled: stripe.Bool(true),
	//   },
	Currency: stripe.String(string(stripe.CurrencyINR)),
	}
	pi, _ := paymentintent.New(params)

	fmt.Println("",pi)
	json.NewEncoder(w).Encode(&pi)

}

func Verify_Payment(w http.ResponseWriter, r *http.Request){


	stripe.Key=cons.Stripe_Key
	

	// To create a PaymentIntent for confirmation, see our guide at: https://stripe.com/docs/payments/payment-intents/creating-payment-intents#creating-for-automatic
	params := &stripe.PaymentIntentConfirmParams{
	PaymentMethod: stripe.String("pm_card_visa"),
	}
	pi, _ := paymentintent.Confirm(
	"pi_3MrISTSGSCYafJiX1durch3U",
	params,
	)

	fmt.Println("payment pi :",pi)
	json.NewEncoder(w).Encode(&pi)




}
