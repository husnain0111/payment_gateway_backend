package shared

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"testing"

	"strings"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/ftloc/exception"
	"github.com/goware/emailx"
	"github.com/pkg/errors"
	"github.com/smallnest/rpcx/server"
	gomail "gopkg.in/gomail.v2"
)

func ValidateEmail(email string) bool {
	err := emailx.Validate(email)
	if err != nil {
		return false
	}
	return true
}

func ValidateRequired(value interface{}) bool {

	if value == nil {
		return false
	}
	switch t := value.(type) {
	case int:
		return true
	case bool:
		return true
	case float64:
		return true
	case string:
		if value.(string) == "" {
			return false
		}
		return true
	default:
		_ = t
		return false
	}
}
func ValidateRequiredWithMessage(value interface{}, message string, code string) Response {
	valid := ValidateRequired(value)
	if !valid {
		return Response{Success: false, Message: message, Code: code}
	}
	return Response{Success: true}
}

func ReturnMessage(success bool, message interface{}, code string) Response {
	return Response{Success: success, Message: message, Code: code}
}
func Exe_cmd(cmd string) string {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	result := fmt.Sprintf("%s", out)
	result = strings.Replace(result, "\n", "", 1)
	return result

}

func In_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func ValidateToken(token string) (email string) {
	fmt.Println(token)
	sharedKey := []byte{54, 104, 215, 160, 189, 218, 171, 9, 214, 26, 111, 221, 170, 199, 125, 71}

	var pay Payload
	payload, _, _ := jose.Decode(token, sharedKey)
	err := json.NewDecoder(strings.NewReader(string(payload))).Decode(&pay)
	if err != nil {
		err2 := errors.WithStack(err)
		exception.Throw(fmt.Errorf("%+v", err2))

	}
	timeIAT := time.Now()
	fmt.Println(timeIAT.Unix())
	timeExp := time.Unix(pay.Exp, 0)
	timeex := timeIAT.Before(timeExp)
	if timeex == true {
		return pay.Email
	} else {
		return ""
	}
}
func ThrowErrorIFNotNil(e error) {
	if e != nil {
		err2 := errors.WithStack(e)
		exception.Throw(fmt.Errorf("%+v", err2))
	}
}

func ThrowErrorIFNotNilWithMessage(e error, message string) {
	if e != nil {
		//err2 := errors.WithStack(e)
		exception.Throw(fmt.Errorf(message))
	}
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}
func RpcxListener(serviceName string, handler interface{}, addr string) {
	cert, err := tls.LoadX509KeyPair("certificate.pem", "key.pem")
	if err != nil {
		log.Print(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	s := server.NewServer(server.WithTLSConfig(config))
	s.RegisterName(serviceName, handler, "")
	err = s.Serve("tcp", addr)
	if err != nil {

	}
}
func SendEmail(email string, body string, subject string) (response Response) {
	// from := mail.NewEmail("PaymentGateway", "no-reply@gig9.io")
	// to := mail.NewEmail(email, email)
	// message := mail.NewSingleEmail(from, subject, to, body, body)
	// client := sendgrid.NewSendClient("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")
	// _, err := client.Send(message)
	// ThrowErrorIFNotNilWithMessage(err, "Email not send")

	m := gomail.NewMessage()
	m.SetHeader("From", "testproject628@gmail.com")
	m.SetHeader("To", email)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "testproject628@gmail.com", "hello1234@")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
		ThrowErrorIFNotNilWithMessage(err, "Email not send")
	}

	response.Message = "success"
	response.Success = true
	return response
}

func GetIpGeoLocation(ip string) (GeoIP, error) {
	// get public ip
	// url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others

	// resp, err := http.Get(url)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// ip1, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	//get user locagion from ip
	geolocationapi := "http://api.ipstack.com/" + ip + "?access_key=c7f30ed035f07ddf9e121c7c363e5180"
	geo, err1 := http.Get(geolocationapi)
	if err1 != nil {
		panic(err1)
	}
	defer geo.Body.Close()
	var geolocation GeoIP
	geoip, err11 := ioutil.ReadAll(geo.Body)
	if err11 != nil {
		panic(err11)
	}
	err := json.Unmarshal(geoip, &geolocation)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("My IP is:%s\n", ip1)
	return geolocation, nil
}
