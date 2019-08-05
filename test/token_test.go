package mas_test

import (
	"fmt"
	"testing"

	"github.com/nminhquan/mas/credentials"
)

const testData = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIifQ.X0-Uu2UO7r_l4HiOS06V4Xs_VTTDabZxRAYRVRHLEJA"
const testKey = "quannm4@vng.com.vng"

func TestValidateToken(t *testing.T) {
	token, error := credentials.ValidateToken(testData, "quannm4@vng.com.vng")

	if error != nil {
		panic(error)
	}

	fmt.Println("Token claims: ", token.Claims)
	fmt.Println("Token method: ", token.Method)
	fmt.Println("Token : ", token.Valid)
}
