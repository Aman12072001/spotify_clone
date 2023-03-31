package cont

import (
	"main/db"
	"main/models"
)

func Membership_Distribution(){

	//default memberships with their prices
	var membership models.Memberships

	membership.Membership_name="Individual"
	membership.Price=119

	db.DB.Create(&membership)

	membership.Membership_name="Duo"
	membership.Price=149

	db.DB.Create(&membership)

}



