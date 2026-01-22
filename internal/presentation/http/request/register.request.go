package request

type RegisterRequest struct {
	Username              string `json:"username" binding:"required,gte=3,lte=30"`
	Email                 string `json:"email" binding:"required,email,gte=5,lte=255"`
	Password              string `json:"password" binding:"required,gte=8,lte=50"`
	Password_confirmation string `json:"password_confirmation" binding:"required,gte=8,lte=50,eqfield=Password"`
}
