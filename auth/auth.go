package auth 

import (
    "fmt"
// "log"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
    "strconv"
    "strings"
)

type AuthDetails struct {
	AuthUuid string
	UserId uint64
}

func CreateToken(authD *AuthDetails) (string, error) {
	
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["auth_uuid"] = authD.AuthUuid
	claims["user_id"] = authD.UserId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
} 

func TokenValid(request *http.Request) error {
	 token, err := VerifyToken(request)
	 if err != nil {
		 return err
	 }

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {  
		  return err
	}

	 return nil
}


func VerifyToken(request *http.Request) (*jwt.Token, error) {
	 tokenString := ExtracToken(request)
	 token, err := jwt.Parse(tokenString, 
		                 func(token *jwt.Token,)(interface{}, error){
				             //does the token match with HMAC ?
			                 if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					           return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				             }
                            return []byte(os.Getenv("API_SECRET")), nil
						})
	  if err != nil {
		  return nil, err
	  }
    return token, nil
}

//get token from request body
func ExtracToken(request * http.Request) (string) {
   keys := request.URL.Query()
   token := keys.Get("token")
   
   if token != "" {
	   return token
   }

   bearToken := request.Header.Get("Authorization")
   //Authorization the token

   strArr :=  strings.Split(bearToken," ")
   if len(strArr) == 2 {
	   return strArr[1]
   }
   return ""
}

func ExtracTokenAuth(request *http.Request) (*AuthDetails, error) {
	token, err := VerifyToken(request)
	if err != nil {
		return nil, err
	}

	//token should be a MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	
	if ok && token.Valid {
		authUuid, ok := claims["auth_uuid"].(string) //convert interface to string
		if !ok {
			return nil, err
		}

		userId, err := strconv.ParseUint(fmt.Sprintf("%f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AuthDetails{
			AuthUuid: authUuid,
			UserId: userId,
		}, nil
	} 

	return nil, err
}