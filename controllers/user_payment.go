package cont

import (
	"encoding/json"
	"fmt"
	"html/template"
	"main/models"
	cons "main/utils"
	"net/http"

	// "github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

type  html_updates struct{


	Email string 
	Membership string 
	Client_Secret string 
} 

func Create_Payment(w http.ResponseWriter, r *http.Request){

	stripe.Key=cons.Stripe_Key

	//get the membership name from body
	
	var membership models.Memberships
	membershipname := r.URL.Query().Get("membership")
	membership.Membership_name=membershipname
	
	// json.NewDecoder(r.Body).Decode(&membership)
	// db.DB.Where("membership_name=?",membership.Membership_name).First(&membership)
	



	params := &stripe.PaymentIntentParams{
	Amount: stripe.Int64(1000),  //int64(membership.Price*100)
	//   AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
	//     Enabled: stripe.Bool(true),
	//   },
	Currency: stripe.String(string(stripe.CurrencyINR)),
	}
	pi, _ := paymentintent.New(params)

	fmt.Println("pi client secret",pi.ClientSecret)
	var html_elements html_updates

	html_elements.Client_Secret=pi.ClientSecret
	html_elements.Membership=membership.Membership_name
	//json.NewEncoder(w).Encode(&pi)
	t, err := template.ParseFiles("controllers/checkout.html")

	if err != nil {

		fmt.Println("template parsing err", err)
	}

	err = t.Execute(w, html_elements)
	if err != nil {

		fmt.Println("template executing error", err)
	}



}

func Verify_Payment(w http.ResponseWriter, r *http.Request){


	stripe.Key=cons.Stripe_Key
	

	// To create a PaymentIntent for confirmation, see our guide at: https://stripe.com/docs/payments/payment-intents/creating-payment-intents#creating-for-automatic
	params := &stripe.PaymentIntentConfirmParams{
	PaymentMethod: stripe.String("pm_card_visa"),
	}
	pi, _ := paymentintent.Confirm(
	"pi_3MsjSsSGSCYafJiX0cLT4eEJ",
	params,
	)

	fmt.Println("payment pi :",pi)
	json.NewEncoder(w).Encode(&pi)




}
