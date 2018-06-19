package handlers

import (
	"log"
	"net/http"

	"../models"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type ProfileController struct {
	DB *gorm.DB
}

func GetUserFromContext(db *gorm.DB, c echo.Context) (*models.User, error) {
	userFc := c.Get("user")
	if userFc != nil {
		return userFc.(*models.User), nil
	}

	ses, _ := session.Get("session", c)
	id, ok := ses.Values["user_id"].(uint)
	if !ok {
		return nil, nil
	}

	user, err := models.GetUserByID(db, id)
	if err != nil {
		return nil, err
	}

	c.Set("user", user)
	return user, nil
}

func (pc *ProfileController) GetProfilePage(c echo.Context) error {
	user, err := GetUserFromContext(pc.DB, c)
	if err != nil {
		log.Fatal(err)
	}
	return c.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"user": user,
	})
}

func (pc *ProfileController) UpdateProfile(c echo.Context) error {
	user, err := GetUserFromContext(pc.DB, c)
	if err != nil {
		log.Fatal(err)
	}
	fpassword := c.FormValue("fpassword")
	if fpassword == "" {
		return c.Render(http.StatusBadRequest, "profile.html", map[string]interface{}{
			"error": "Пароль не может быть пустым",
		})
	}

	spassword := c.FormValue("spassword")
	if fpassword != spassword {
		return c.Render(http.StatusBadRequest, "profile.html", map[string]interface{}{
			"error": "Введеные пароли не совпадают",
		})
	}
	user.PasswordHash = models.HashPwd(fpassword)
	user.Save(pc.DB)

	return c.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"message": "Обновление профиля выполнено успешно",
	})
}
