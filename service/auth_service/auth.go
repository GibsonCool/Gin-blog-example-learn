package auth_service

import "Gin-blog-example/models"

type Auth struct {
	Username string
	Password string
}

func (au *Auth) Check() (bool, error) {
	return models.CheckAuth(au.Username, au.Password)
}
