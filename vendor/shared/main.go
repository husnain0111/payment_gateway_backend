package shared

const DB_URL = "couchbase://206.189.227.175"
const DB_USERNAME = "dbuser"
const DB_PASSWORD = "R4a4*S.R8m2f&87"
const Raven = "https://a6c253e601444cf399bb1566ba2f3860:53c8a5098b1b4da88f243c28ea7c8133@sentry.io/1218607"

const BUCKET = "PaymentGateway"

type Request struct {
	Data   string `json:"data"`
	Action string `json:"action"`
}

type Response struct {
	Logs    []error     `json:"logs"`
	Message interface{} `json:"message"`
	Success bool        `json:"success"`
	Code    string      `json:"code"`
}

type UserSignup struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Username              string `json:"username"`
	Email                 string `json:"email"`
	Password              string `json:"password"`
	IsActive              bool   `json:"is_active"`
	Country               string `json:"country"`
	PhoneNumber           string `json:"phone_number"`
	Refferal              string `json:"refferal"`
	IP                    string `json:"ip"`
	Token                 string `json:"token"`
	GToken                string `json:"g_token"`
	Gender                string `json:"gender"`
	Position              string `json:"position"`
	CompanyName           string `json:"companyname"`
	CompanyWebsite        string `json:"companywebsite"`
	AboutCompany          string `json:"aboutcompany"`
	State                 string `json:"state"`
	City                  string `json:"city"`
	Skype                 string `json:"skype"`
	Nickname              string `json:"nickname"`
	NewsletterEmail       bool   `json:"newsletteremail"`
	MonetiserURLEmail     bool   `json:"monetiserurlemail"`
	AffiliateProgramEmail bool   `json:"affiliateprogramemail"`
}

type UserLogin struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	IP        string `json:"ip"`
	GToken    string `json:"g_token"`
	UserAgent string `json:"useragent"`
}

type ForgetPasswordRequest struct {
	Email    string `json:"email"`
	IP       string `json:"ip"`
	Password string `json:"password"`
	Token    string `json:"token"`
	GToken   string `json:"g_token"`
}

type Payload struct {
	Action string `json:"action"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	IAT    int64  `json:"iat"`
}
type AddAddress struct {
	Token    string `json:"token"`
	Address  string `json:"address"`
	Currency string `json:"currency"`
	Email    string `json:"email"`
	Generate bool   `json:"bool"`
}
type ListAddress struct {
	Currency string `json:"currency"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type Address struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
	Email    string `json:"email"`
	Inuse    bool   `json:"inuse"`
	Network  string `json:"network"`
	ID       string `json:"id"`
}
type ChangePassword struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type Wallet struct {
	Client   string
	DataDir  string
	Conf     string
	Currency string
	Args     string
	Network  string
	Price    float64
}

type GeoIP struct {
	Ip          string `json:"ip"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name""`
	City        string `json:"city"`
}
