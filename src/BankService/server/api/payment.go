package api

import (
	"bank/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	User        string `json:"user"`
	Description string `json:"description"`
	Amount      uint32 `json:"amount"`
	Paid        bool   `json:"paid"`
}

type PaymentReq struct {
	User        string `json:"user"`
	Description string `json:"description"`
	Amount      uint32 `json:"amount"`
}

type Res struct {
	Res string `json:"res"`
}

// POST /payment
func NewPayment(c *gin.Context) {
	db := util.GetDb()

	var data PaymentReq
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couldn't parse data: " + err.Error()})
		return
	}

	pay := Payment{
		User:        data.User,
		Description: data.Description,
		Amount:      data.Amount,
		Paid:        false,
	}
	if err := db.Save(&pay).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't save payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Res{Res: "Done"})
}

// POST /payment/:id/pay
func PayPaymentById(c *gin.Context) {
	db := util.GetDb()
	userid := c.Param("id")

	i, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't understand id: " + err.Error()})
		return
	}

	var res Payment
	if err := db.Model(Payment{}).Where("id = ?", i).First(&res).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't save payment: " + err.Error()})
		return
	}

	if err := db.Model(Payment{}).Where(&res).Update("paid", true).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't update payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Res{Res: "Done"})
}

// DELETE /payment/:id
func DelPaymentById(c *gin.Context) {
	db := util.GetDb()
	userid := c.Param("id")

	i, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't understand id: " + err.Error()})
		return
	}

	var res Payment
	if err := db.Model(Payment{}).Where("id = ?", i).Delete(&res).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't delete payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, Res{Res: "Done"})
}

// GET /payment/:id
func GetPaymentById(c *gin.Context) {
	db := util.GetDb()
	userid := c.Param("id")

	i, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't understand id: " + err.Error()})
		return
	}

	var res Payment
	if err := db.Model(Payment{}).Where("id = ?", i).First(&res).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't get payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
