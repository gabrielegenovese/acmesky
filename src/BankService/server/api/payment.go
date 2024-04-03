package api

import (
	"bank/util"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Payment struct {
	ID          uuid.UUID    `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" gorm:"index"`
	User        string       `json:"user"`
	Description string       `json:"description"`
	Amount      uint32       `json:"amount"`
	Paid        bool         `json:"paid"`
}

type PaymentReq struct {
	User        string `json:"user"`
	Description string `json:"description"`
	Amount      uint32 `json:"amount"`
}

type Res struct {
	Res string `json:"res"`
}

// POST /payment/new
func NewPayment(c *gin.Context) {
	db := util.GetDb()

	var data PaymentReq
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couldn't parse data: " + err.Error()})
		return
	}

	// Non viene controllato se l'id già esiste ma non dovrebbe servire (finchè abbiamo pochi dati nel db)
	pay := Payment{
		ID:          uuid.New(),
		User:        data.User,
		Description: data.Description,
		Amount:      data.Amount,
		Paid:        false,
	}
	if err := db.Save(&pay).Error; err != nil {
		c.JSON(http.StatusBadRequest, Res{Res: "Couln't save payment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, pay)
}

// POST /payment/pay/:id
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

// DELETE /payment/:id
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

// GET /payment/:id
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
