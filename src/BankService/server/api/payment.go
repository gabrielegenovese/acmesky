package api

import (
	"bank/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		New payment
// @Description	Create a new unpaid payment
// @Tags			payment
// @Accept			json
// @Produce		json
// @Param			request	body		PaymentReq	true	"payment data"
// @Success		200		{object}	Payment
// @Failure		400		{object}	Res
// @Router			/payment/new [put]
func NewPayment(c *gin.Context) {
	db := util.GetDb()

	var data PaymentReq
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couldn't parse data: " + err.Error()})
		return
	}

	// Non viene controllato se l'id già esiste ma non dovrebbe servire (finchè abbiamo pochi dati nel db)
	var id = uuid.New()
	pay := Payment{
		ID:          id,
		User:        data.User,
		Description: data.Description,
		Amount:      data.Amount,
		Link:        "http://localhost:8083/pay/" + id.String() + "?redirecturi=",
		Paid:        false,
	}
	if err := db.Save(&pay).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't save payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, pay)
}

// @Summary		Pay a payment
// @Description	Pay an unpaid payment
// @Tags			payment
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"payment ID"
// @Success		200	{object}	Res
// @Failure		400	{object}	Res
// @Router			/payment/pay/{id} [post]
func PayPaymentById(c *gin.Context) {
	db := util.GetDb()
	userid := c.Param("id")

	uid, err := uuid.Parse(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't parse id: " + err.Error()})
		return
	}

	var res Payment
	if err := db.Where(Payment{ID: uid}).First(&res).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't get payment: " + err.Error()})
		return
	}

	if err := db.Model(Payment{}).Where(&res).Update("paid", true).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't update payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Res{Res: "Done"})
}

// @Summary		Delete a payment
// @Description	Given a payment ID, find the corresponding payment and delete it.
// @Tags			payment
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"payment ID"
// @Success		200	{object}	Res
// @Failure		400	{object}	Res
// @Router			/payment/{id} [delete]
func DelPaymentById(c *gin.Context) {
	db := util.GetDb()
	userid := c.Param("id")

	uid, err := uuid.Parse(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't parse id: " + err.Error()})
		return
	}

	var res Payment
	if err := db.Where(Payment{ID: uid}).Delete(&res).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't delete payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Res{Res: "Done"})
}

// @Summary		Get a payment
// @Description	Given a payment ID, find the corresponding payment and return it.
// @Tags			payment
// @Accept			json
// @Produce		json
// @Param			id	path		string	true	"payment ID"
// @Success		200	{object}	Res
// @Failure		400	{object}	Res
// @Router			/payment/{id} [get]
func GetPaymentById(c *gin.Context) {
	db := util.GetDb()
	userid := c.Param("id")

	uid, err := uuid.Parse(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't parse id: " + err.Error()})
		return
	}

	var res Payment
	if err := db.Where(Payment{ID: uid}).First(&res).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't get payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
