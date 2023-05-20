package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserId       string             `json:"userid"`
	FirstName    *string            `json:"firstname"`
	LastName     *string            `json:"lastname"`
	Password     *string            `json:"password"`
	Email        *string            `json:"email"`
	Phone        *string            `json:"phone"`
	IsActive     bool               `json:"isactive"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshtoken"`
	UserType     *string            `json:"usertype"`
	CreatedBy    *string            `json:"createdby"`
	CreatedDate  time.Time          `json:"createddate"`
	ModifiedBy   *string            `json:"modifiedby"`
	ModifiedDate time.Time          `json:"modifieddate"`
}

type UpdateUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserId    string             `json:"userid"`
	FirstName *string            `json:"firstname"`
	LastName  *string            `json:"lastname"`
	//Password     *string            `json:"password"`
	Email        *string   `json:"email"`
	Phone        *string   `json:"phone"`
	IsActive     bool      `json:"isactive"`
	Token        *string   `json:"token"`
	RefreshToken *string   `json:"refreshtoken"`
	UserType     *string   `json:"usertype"`
	CreatedBy    *string   `json:"createdby"`
	CreatedDate  time.Time `json:"createddate"`
	ModifiedBy   *string   `json:"modifiedby"`
	ModifiedDate time.Time `json:"modifieddate"`
}

type LoginDto struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
