package tonicpow

// APIEnvironment is used to differentiate the environment when making requests
type APIEnvironment string

const (

	// Field key names for various model requests
	fieldAdvertiserProfileID = "advertiser_profile_id"
	fieldAmount              = "amount"
	fieldApiKey              = "api_key"
	fieldCampaignID          = "campaign_id"
	fieldCurrency            = "currency"
	fieldCustomDimensions    = "custom_dimensions"
	fieldDelayInMinutes      = "delay_in_minutes"
	fieldEmail               = "email"
	fieldGoalID              = "goal_id"
	fieldID                  = "id"
	fieldLastBalance         = "last_balance"
	fieldLinkID              = "link_id"
	fieldName                = "name"
	fieldPassword            = "password"
	fieldPasswordConfirm     = "password_confirm"
	fieldPhone               = "phone"
	fieldPhoneCode           = "phone_code"
	fieldReason              = "reason"
	fieldShortCode           = "short_code"
	fieldTargetURL           = "target_url"
	fieldToken               = "token"
	fieldUserID              = "user_id"
	fieldVisitorSessionGUID  = "tncpw_session"

	// Model names (used for request endpoints)
	modelAdvertiser     = "advertisers"
	modelCampaign       = "campaigns"
	modelConversion     = "conversions"
	modelGoal           = "goals"
	modelLink           = "links"
	modelRates          = "rates"
	modelUser           = "users"
	modelVisitors       = "visitors"
	modelVisitorSession = "sessions"

	// apiVersion current version for all endpoints
	apiVersion = "v1"

	// defaultUserAgent is the default user agent for all requests
	defaultUserAgent string = "go-tonicpow: " + apiVersion

	// LiveEnvironment is the live production environment
	LiveEnvironment APIEnvironment = "https://api.tonicpow.com/" + apiVersion + "/"

	// LocalEnvironment is for testing locally using your own api instance
	LocalEnvironment APIEnvironment = "http://localhost:3000/" + apiVersion + "/"

	// StagingEnvironment is used for production-like testing
	StagingEnvironment APIEnvironment = "https://apistaging.tonicpow.com/" + apiVersion + "/"

	// sessionCookie is the cookie name for session tokens
	sessionCookie = "session_token"

	// TestEnvironment is a test-only environment
	//TestEnvironment APIEnvironment = "https://test.tonicpow.com/"+apiVersion+"/"
)

// Error is the universal error response from the API
//
// For more information: https://docs.tonicpow.com/#d7fe13a3-2b6d-4399-8d0f-1d6b8ad6ebd9
type Error struct {
	Code        int    `json:"code"`
	Data        string `json:"data"`
	IPAddress   string `json:"ip_address"`
	Method      string `json:"method"`
	Message     string `json:"message"`
	RequestGUID string `json:"request_guid"`
	URL         string `json:"url"`
}

// User is the user model
//
// For more information: https://docs.tonicpow.com/#50b3c130-7254-4a05-b312-b14647736e38
type User struct {
	Balance            uint64 `json:"balance,omitempty"`
	Earned             uint64 `json:"earned,omitempty"`
	Email              string `json:"email,omitempty"`
	EmailVerified      bool   `json:"email_verified,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	ID                 uint64 `json:"id,omitempty"`
	InternalAddress    string `json:"internal_address,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	MiddleName         string `json:"middle_name,omitempty"`
	NewPassword        string `json:"new_password,omitempty"`
	NewPasswordConfirm string `json:"new_password_confirm,omitempty"`
	Password           string `json:"password,omitempty"`
	PayoutAddress      string `json:"payout_address,omitempty"`
	Phone              string `json:"phone,omitempty"`
	PhoneVerified      bool   `json:"phone_verified,omitempty"`
	ReferredByUserID   uint64 `json:"referred_by_user_id,omitempty"`
	Status             string `json:"status,omitempty"`
}

// AdvertiserProfile is the advertiser_profile model (child of User)
//
// For more information: https://docs.tonicpow.com/#2f9ec542-0f88-4671-b47c-d0ee390af5ea
type AdvertiserProfile struct {
	HomepageURL string `json:"homepage_url,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	ID          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	UserID      uint64 `json:"user_id,omitempty"`
}

// Campaign is the campaign model (child of AdvertiserProfile)
//
// For more information: https://docs.tonicpow.com/#5aca2fc7-b3c8-445b-aa88-f62a681f8e0c
type Campaign struct {
	AdvertiserProfile   *AdvertiserProfile `json:"advertiser_profile,omitempty"`
	AdvertiserProfileID uint64             `json:"advertiser_profile_id,omitempty"`
	Balance             float64            `json:"balance,omitempty"`
	BalanceSatoshis     int64              `json:"balance_satoshis,omitempty"`
	BotProtection       bool               `json:"bot_protection,omitempty"`
	Clicks              uint64             `json:"clicks,omitempty"`
	Currency            string             `json:"currency,omitempty"`
	Description         string             `json:"description,omitempty"`
	ExpiresAt           string             `json:"expires_at,omitempty"` // time.RFC3339
	FundingAddress      string             `json:"funding_address,omitempty"`
	Goals               []*Goal            `json:"goals,omitempty"`
	ID                  uint64             `json:"id,omitempty"`
	ImageURL            string             `json:"image_url,omitempty"`
	LinksCreated        uint64             `json:"links_created,omitempty"`
	PayPerClickRate     float64            `json:"pay_per_click_rate,omitempty"`
	PublicGUID          string             `json:"public_guid,omitempty"`
	TargetURL           string             `json:"target_url,omitempty"`
	Title               string             `json:"title,omitempty"`
}

// Goal is the goal model (child of Campaign)
//
// For more information: https://docs.tonicpow.com/#316b77ab-4900-4f3d-96a7-e67c00af10ca
type Goal struct {
	CampaignID     uint64  `json:"campaign_id,omitempty"`
	Description    string  `json:"description,omitempty"`
	ID             uint64  `json:"id,omitempty"`
	MaxPerPromoter int16   `json:"max_per_promoter,omitempty"`
	MaxPerVisitor  int16   `json:"max_per_visitor,omitempty"`
	Name           string  `json:"name,omitempty"`
	PayoutRate     float64 `json:"payout_rate,omitempty"`
	Payouts        int     `json:"payouts,omitempty"`
	PayoutType     string  `json:"payout_type,omitempty"`
	Title          string  `json:"title,omitempty"`
}

// Conversion is the response of getting a conversion
//
// For more information: https://docs.tonicpow.com/#75c837d5-3336-4d87-a686-d80c6f8938b9
type Conversion struct {
	Amount           uint64 `json:"amount,omitempty"`
	CustomDimensions string `json:"custom_dimensions,omitempty"`
	GoalID           uint64 `json:"goal_id,omitempty"`
	GoalName         string `json:"goal_name,omitempty"`
	ID               uint64 `json:"ID,omitempty"`
	PayoutAfter      string `json:"payout_after,omitempty"`
	Status           string `json:"status,omitempty"`
	TxID             string `json:"tx_id,omitempty"`
	UserID           uint64 `json:"user_id,omitempty"`
}

// Link is the link model (child of User) (relates Campaign to User)
// Use the CustomShortCode on create for using your own short code
//
// For more information: https://docs.tonicpow.com/#ee74c3ce-b4df-4d57-abf2-ccf3a80e4e1e
type Link struct {
	CampaignID      uint64 `json:"campaign_id,omitempty"`
	CustomShortCode string `json:"custom_short_code,omitempty"`
	ID              uint64 `json:"id,omitempty"`
	ShortCode       string `json:"short_code,omitempty"`
	ShortCodeURL    string `json:"short_code_url,omitempty"`
	UserID          uint64 `json:"user_id,omitempty"`
}

// VisitorSession is the session for any visitor or user (related to link and campaign)
//
// For more information: https://docs.tonicpow.com/#d0d9055a-0c92-4f55-a370-762d44acf801
type VisitorSession struct {
	CampaignID       uint64 `json:"campaign_id,omitempty"`
	CustomDimensions string `json:"custom_dimensions,omitempty"`
	IPAddress        string `json:"ip_address,omitempty"`
	LinkID           uint64 `json:"link_id,omitempty"`
	LinkUserID       uint64 `json:"link_user_id,omitempty"`
	Provider         string `json:"provider,omitempty"`
	Referer          string `json:"referer,omitempty"`
	TncpwSession     string `json:"tncpw_session,omitempty"`
	UserAgent        string `json:"user_agent,omitempty"`
}

// Rate is the rate results
//
// For more information: https://docs.tonicpow.com/#fb00736e-61b9-4ec9-acaf-e3f9bb046c89
type Rate struct {
	Currency            string  `json:"currency,omitempty"`
	CurrencyAmount      float64 `json:"currency_amount,omitempty"`
	CurrencyLastUpdated string  `json:"currency_last_updated,omitempty"`
	CurrencyName        string  `json:"currency_name,omitempty"`
	Price               float64 `json:"price,omitempty"`
	PriceInSatoshis     int64   `json:"price_in_satoshis,omitempty"`
	RateLastUpdated     string  `json:"rate_last_updated,omitempty"`
}
