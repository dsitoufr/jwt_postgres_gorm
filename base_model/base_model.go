package model


import (
	"fmt"
	"log"
	"errors"
	"github.com/jinzhu/gorm"
	//_"github.com/jinzhu/gorm/dialects/postgres"  //postgres driver
	"github/jwt_api_auth_2/auth"
	"github.com/badoux/checkmail"
	"github.com/twinj/uuid"
)

//
//objects 
//

type Server struct {
	modelInterface
	DB *gorm.DB
}

type User struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Email string `gorm:"size:255; not null" json:"email"`
}

type Todo struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	UserID uint64 `gorm:"not null" json:"user_id"`
	Title string `gorm:"size:255;not null" json:"title"`
}

type Auth struct {
	ID uint64 `gorm:"primary_key;auto_increment" json:"id"`
	UserID uint64 `gorm:"not null" json:"user_id"`
	AuthUUID string `gorm:"size:255;not null" json:"auth_uuid"`
}

type modelInterface interface {
   //db init
   Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*gorm.DB, error)

   //user methods
   ValidateEmail(string) error
   CreateUser(*User)(*User, error)
   GetUserByEmail(string) (*User, error)

   //todo methods
   CreateTodo(*Todo)(*Todo, error)

   //auth methods
   FetchAuth(*auth.AuthDetails) (*Auth, error)
   DeleteAuth(*auth.AuthDetails) error
   CreateAuth(uint64)(*Auth, error)
}

var (
	//Model modelInterface = &Server{}
 )
 
 
 func(s *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*gorm.DB, error) {
	 var err error
	 DBURL := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",DbHost, DbPort, DbUser, DbPassword)
	 
	 s.DB, err = gorm.Open(Dbdriver, DBURL)
     if err != nil {
		 return nil, err
	 }
	 
	 log.Println("connected to postgres database")

	 s.DB.Debug().AutoMigrate(
		 &User{},
		 &Auth{},
		 &Todo{},
	 )

	 return s.DB, nil

 }


 func(s *Server) ValidateEmail(email string) error {
	if email == "" {
		return errors.New("required email")
	}
 
	if email != "" {
		 
		if err := checkmail.ValidateFormat(email); err != nil {
			 return errors.New("invalid email")
		 }
	}
 
	return nil
 }
 
 //create user in database
 func(s *Server) CreateUser(user *User) (*User, error) {
	 
	 emailErr := s.ValidateEmail(user.Email)
	 if emailErr != nil {
		 return nil, emailErr
	 }
 
	 err := s.DB.Debug().Create(&user).Error
	 if err != nil {
		  return nil, err
	 }
 
	 return user, nil
 }
 
 //query user from database
 func(s *Server) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := s.DB.Debug().Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
 
	return user, nil
 }

 
 //create Todo in database
func (s *Server) CreateTodo(todo *Todo)(*Todo, error) {
   if todo.Title == "" {
	   return nil, errors.New("please provide a valid title!")
   }

   if todo.UserID == 0 {
	   return nil, errors.New("a valid user id is required !")
   }

   err := s.DB.Debug().Create(&todo).Error
   if err != nil {
	   return nil, err
   }

   return todo, nil
}

func (s *Server) FetchAuth(authD *auth.AuthDetails) (*Auth, error) {
	
	au := &Auth{}

	err := s.DB.Debug().Where("user_id = ? and auth_uuid = ?", authD.UserId, authD.AuthUuid).Take(&au).Error
	if err != nil {
		return nil, err
	}
	return au, nil

}

//Once a user row in the auth table
func (s *Server) DeleteAuth(authD *auth.AuthDetails) error {
	au := &Auth{}
	db := s.DB.Debug().Where("user_id = ? AND auth_uuid = ?", authD.UserId, authD.AuthUuid).Take(&au).Delete(&au)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

//Once the user signup/login, create a row in the auth table, with a new uuid
func (s *Server) CreateAuth(userId uint64) (*Auth, error) {
	au := &Auth{}
	au.AuthUUID = uuid.NewV4().String() //generate a new UUID each time
	au.UserID = userId
	err := s.DB.Debug().Create(&au).Error
	if err != nil {
		return nil, err
	}
	return au, nil
}