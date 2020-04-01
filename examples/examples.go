package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/tonicpow/go-tonicpow"
)

var (
	// TonicPowAPI is the client we will load on start-up
	TonicPowAPI *tonicpow.Client
)

// Load the TonicPow API Client once when the application loads
func init() {

	//
	// Get the API key (from env or your own config)
	//
	apiKey := os.Getenv("TONICPOW_API_KEY")
	if len(apiKey) == 0 {
		log.Fatalf("api key is required: %s", "TONICPOW_API_KEY")
	}

	// Set the environment
	environmentString := os.Getenv("TONICPOW_ENVIRONMENT")
	var environment tonicpow.APIEnvironment
	if environmentString == "staging" {
		environment = tonicpow.StagingEnvironment
	} else if environmentString == "live" {
		environment = tonicpow.LiveEnvironment
	} else {
		environment = tonicpow.LocalEnvironment
	}

	//
	// Load the api client (creates a new session)
	// You can also set the environment or client options
	//
	var err error
	TonicPowAPI, err = tonicpow.NewClient(apiKey, environment, nil)
	if err != nil {
		log.Fatalf("error in NewClient: %s", err.Error())
	}
}

func main() {
	//
	// Example for ending the api session for the application
	// This is not needed, sessions will automatically expire
	//
	defer func() {
		_ = TonicPowAPI.EndSession("")
	}()

	// Example vars
	var err error
	var userSessionToken string
	testPassword := "ExamplePassForNow0!"

	//
	// Example: Prolong a session
	//
	if err = TonicPowAPI.ProlongSession(""); err != nil {
		log.Fatalf("ProlongSession: %s", err.Error())
	} else {
		log.Println("session created and prolonged...")
	}

	//
	// Example: Create a user
	//
	user := &tonicpow.User{
		Email:    fmt.Sprintf("Tes_ti-ng+%d@TonicPow.com", rand.Intn(100000)),
		Password: testPassword,
	}
	if user, err = TonicPowAPI.CreateUser(user); err != nil {
		log.Fatalf("create user failed error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else {
		log.Printf("user %d created - referral url %s", user.ID, user.ReferralURL)
	}

	//
	// Example: Get a user (id)
	//
	if user, err = TonicPowAPI.GetUser(user.ID, ""); err != nil {
		log.Fatalf("get user failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("got user by id %d", user.ID)
	}

	//
	// Example: Get a user (email)
	//
	if user, err = TonicPowAPI.GetUser(0, user.Email); err != nil {
		log.Fatalf("get user failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("got user by email %s", user.Email)
	}

	//
	// Example: Update a user
	//
	user.FirstName = "Austin"
	user.PayoutAddress = "mrz@moneybutton.com"
	if user, err = TonicPowAPI.UpdateUser(user, ""); err != nil {
		log.Fatalf("update user failed error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else {
		log.Printf("user %d updated - first_name: %s", user.ID, user.FirstName)
	}

	//
	// Example: Create a user  (referred by another user) (alternate is a visitor tncpw_session)
	//
	refUser := &tonicpow.User{
		Email:            fmt.Sprintf("Tes_ti-ng+%d@TonicPow.com", rand.Intn(100000)),
		Password:         testPassword,
		ReferredByUserID: user.ID,
		// TncpwSession:     "1e74da42752127b60422bdcfe6741eeadc1a1c5dab22ca8f17ca9ee62dc21e26",
	}
	if refUser, err = TonicPowAPI.CreateUser(refUser); err != nil {
		log.Fatalf("create user failed error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else {
		log.Printf("user %d created and referred by user %d", refUser.ID, refUser.ReferredByUserID)
	}

	//
	// Example: Accept/Approve the User (if required by application config)
	//
	if err = TonicPowAPI.AcceptUser(user.ID, "", "awesome person"); err != nil {
		log.Fatalf("accept user failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("user accepted: %d", user.ID)
	}

	//
	// Example: Get new updated balance for user
	//
	if user, err = TonicPowAPI.GetUserBalance(user.ID, 0); err != nil {
		log.Fatalf("get user balance failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("user balance: %d", user.Balance)
	}

	//
	// Example: Release internal balance to user's payout_address
	//
	if err = TonicPowAPI.ReleaseUserBalance(user.ID, "they deserve it all"); err != nil {
		log.Fatalf("release balance failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Println("release balance was successful")
	}

	//
	// Example: Refund user's balance back to the corresponding campaigns
	//
	/*
		if err = TonicPowAPI.RefundUserBalance(user.ID, "account was deemed fraudulent"); err != nil {
			log.Fatalf("refund failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
		} else {
			log.Println("refund started...")
		}
	*/

	//
	// Example: Forgot password
	//
	if err = TonicPowAPI.ForgotPassword(user.Email); err != nil {
		log.Fatalf("forgot password failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("sent forgot password: %s", user.Email)
	}

	//
	// Example: Get all referrals from a user
	//
	var referrals []*tonicpow.UserReferral
	if referrals, err = TonicPowAPI.GetUserReferrals(user.ID, ""); err != nil {
		log.Fatalf("get user referrals failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("got %d referrals", len(referrals))
	}

	//
	// Example: List all referrals
	//
	var referralResults *tonicpow.ReferralResults
	if referralResults, err = TonicPowAPI.ListUserReferrals(1, 5, tonicpow.SortByFieldReferrals, tonicpow.SortOrderDesc); err != nil {
		log.Fatalf("list referrals failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("referrals found: %d page: %d", referralResults.Results, referralResults.CurrentPage)
	}

	//
	// Example: Create an advertiser
	//
	advertiser := &tonicpow.AdvertiserProfile{
		UserID:      user.ID,
		Name:        "Acme Advertising",
		HomepageURL: "https://tonicpow.com",
		IconURL:     "https://tonicpow.com/images/logos/apple-touch-icon.png",
	}
	if advertiser, err = createAdvertiserProfile(advertiser, ""); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get an advertiser profile
	//
	if advertiser, err = getAdvertiserProfile(advertiser.ID, ""); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Login for a user
	//
	user.Password = testPassword
	userSessionToken, err = TonicPowAPI.LoginUser(user)
	if err != nil {
		log.Fatalf("user login failed error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else {
		log.Printf("user login: %s token: %s", user.Email, userSessionToken)
	}

	//
	// Example: Logout (just for our example)
	//
	defer func(token string) {
		_ = TonicPowAPI.LogoutUser(token)
		log.Println("user logout complete")
	}(userSessionToken)

	//
	// Example: Current user details
	//
	user, err = TonicPowAPI.CurrentUser(userSessionToken)
	if err != nil {
		log.Fatalf("current user failed error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else {
		log.Printf("current user: %s", user.Email)
	}

	//
	// Example: Does User Exist?
	//
	var exists *tonicpow.UserExists
	exists, err = TonicPowAPI.UserExists(user.Email)
	if err != nil {
		log.Fatalf("check user exists error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else if exists.ID > 0 {
		log.Printf("user %d exists with status %s", exists.ID, exists.Status)
	}

	//
	// Example: Does User Exist?  (does not exist example)
	//
	exists, err = TonicPowAPI.UserExists("user@doesnotexist.com")
	if err != nil {
		log.Fatalf("check user exists error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else if exists == nil || exists.ID <= 0 {
		log.Printf("user does not exist!")
	}

	//
	// Example: Create an advertiser as a user
	//
	advertiser = &tonicpow.AdvertiserProfile{
		UserID:      user.ID,
		Name:        "Acme User Advertising",
		HomepageURL: "https://tonicpow.com",
		IconURL:     "https://tonicpow.com/images/logos/apple-touch-icon.png",
	}
	if advertiser, err = createAdvertiserProfile(advertiser, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get Advertiser Profile
	//
	if advertiser, err = getAdvertiserProfile(advertiser.ID, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Update advertising profile
	//
	advertiser.Name = "Acme New User Advertising"
	if advertiser, err = TonicPowAPI.UpdateAdvertiserProfile(advertiser, userSessionToken); err != nil {
		log.Fatalf("update advertiser failed error %s - api error: %s data: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data)
	} else {
		log.Printf("advertiser profile %s id %d updated", advertiser.Name, advertiser.ID)
	}

	//
	// Example: Get a rate
	//
	if _, err = getCurrentRate("usd"); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create a campaign
	//
	expiresAt := time.Now().UTC().Add(500 * time.Hour) // Optional expiration date
	campaign := &tonicpow.Campaign{
		AdvertiserProfileID: advertiser.ID,
		Currency:            "USD",
		Description:         "Earn BSV for sharing things you like.",
		ImageURL:            "https://i.imgur.com/TbRFiaR.png",
		PayPerClickRate:     0.01,
		TargetURL:           "https://offers.tonicpow.com",
		Title:               "TonicPow Offers",
		ExpiresAt:           expiresAt.Format(time.RFC3339),
	}
	if campaign, err = createCampaign(campaign, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get Campaign
	//
	if campaign, err = getCampaign(campaign.ID, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create a Goal (flat payout)
	//
	goal := &tonicpow.Goal{
		CampaignID:  campaign.ID,
		Description: "Bring leads and get paid!",
		Name:        "new-lead-landing-page",
		PayoutRate:  0.04,
		PayoutType:  "flat",
		Title:       "Landing Page Leads",
	}
	if goal, err = createGoal(goal, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get Goal
	//
	if goal, err = getGoal(goal.ID, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create a Goal (percent payout)
	//
	goalPercent := &tonicpow.Goal{
		// CampaignID: 1,
		CampaignID:  campaign.ID,
		Description: "Get 10% of all action!",
		// Name:        fmt.Sprintf("%s%d", "all-action", rand.Intn(100000)),
		Name:           "all-action",
		PayoutRate:     0.10, // 10% as a float
		PayoutType:     "percent",
		Title:          "10% Commissions",
		MaxPerPromoter: 2,
		MaxPerVisitor:  2,
	}
	if goalPercent, err = createGoal(goalPercent, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create a Link
	//
	link := &tonicpow.Link{
		CampaignID: campaign.ID,
		UserID:     user.ID,
	}
	if link, err = createLink(link, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get Link
	//
	if link, err = getLink(link.ID, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create a Link (making a custom short code using the user's name)
	//
	link = &tonicpow.Link{
		CampaignID:      campaign.ID,
		UserID:          user.ID,
		CustomShortCode: fmt.Sprintf("%s%d", user.FirstName, rand.Intn(100000)),
	}
	if link, err = createLink(link, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create another custom link
	//
	link = &tonicpow.Link{
		CampaignID:      campaign.ID,
		UserID:          user.ID,
		CustomShortCode: fmt.Sprintf("%s%d", user.FirstName+"zzz", rand.Intn(100000)),
	}
	if link, err = createLink(link, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create a link by url
	//
	linkByURL := &tonicpow.Link{
		TargetURL: "https://amazon.com",
		UserID:    user.ID,
	}
	if linkByURL, err = createLinkByURL(linkByURL, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get Link(s) by user
	//
	var linkResults *tonicpow.LinkResults
	if linkResults, err = getLinks(link.ID, userSessionToken); err != nil {
		log.Fatal(err.Error())
	}
	if len(linkResults.Links) > 0 {
		log.Printf("found links %v", linkResults.Results)
	}

	//
	// Example: List active campaigns
	//
	var campaignResults *tonicpow.CampaignResults
	if campaignResults, err = TonicPowAPI.ListCampaigns("", 1, 5, tonicpow.SortByFieldCreatedAt, tonicpow.SortOrderDesc); err != nil {
		log.Fatalf("list campaign failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("campaigns found: %d page: %d", campaignResults.Results, campaignResults.CurrentPage)
	}

	//
	// Example: List Campaigns  (by advertiser)
	//

	if campaignResults, err = TonicPowAPI.ListCampaignsByAdvertiserProfile(advertiser.ID, 1, 5, tonicpow.SortByFieldCreatedAt, tonicpow.SortOrderAsc); err != nil {
		log.Fatalf("list campaign failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("campaigns by advertiser found: %d page: %d", campaignResults.Results, campaignResults.CurrentPage)
	}

	//
	// Example: Get campaigns by url
	//
	if campaignResults, err = TonicPowAPI.ListCampaignsByURL(campaignResults.Campaigns[0].TargetURL, 1, 5, tonicpow.SortByFieldCreatedAt, tonicpow.SortOrderDesc); err != nil {
		log.Fatalf("get campaign by url failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("campaigns found: %d", campaignResults.Results)
	}

	//
	// Example: Campaigns Feed in XML
	//
	xmlResults, _ := TonicPowAPI.CampaignsFeed(tonicpow.FeedTypeRSS)
	if len(xmlResults) > 0 {
		log.Printf("campaigns feed found - length: %d", len(xmlResults))
	} else {
		log.Fatalf("failed to get campaigns feed: %s", TonicPowAPI.LastRequest.Error.Message)
	}

	//
	// Example: Campaign Statistics
	//
	var statistics *tonicpow.CampaignStatistics
	statistics, err = TonicPowAPI.CampaignStatistics()
	if err != nil {
		log.Fatalf("get campaign statistics failed error %s - api error: %s", err.Error(), TonicPowAPI.LastRequest.Error.Message)
	} else {
		log.Printf("statistics found: active %d", statistics.Active)
	}

	//
	// Example: Activate User
	//
	/*
		if err = TonicPowAPI.ActivateUser(user.ID, ""); err != nil {
			log.Fatalf("activate user failed - api error: %s", TonicPowAPI.LastRequest.Error.Message)
		}
	*/

	//
	// Example: Pause User
	//
	/*
		if err = TonicPowAPI.PauseUser(user.ID, "", "reviewing account"); err != nil {
			log.Fatalf("pause user failed - api error: %s", TonicPowAPI.LastRequest.Error.Message)
		}
	*/

	//
	// Example: Create a Visitor Session
	//
	visitorSession := &tonicpow.VisitorSession{
		CustomDimensions: "my custom dimensions",
		LinkID:           link.ID,
	}
	if visitorSession, err = createVisitorSession(visitorSession); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Get Visitor Session
	//
	if visitorSession, err = getVisitorSession(visitorSession.TncpwSession); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("visitor session found: %s", visitorSession.TncpwSession)

	//
	// Example: Create a Visitor Session (Full Anti-Bot Options)
	//
	visitorSession = &tonicpow.VisitorSession{
		CustomDimensions: "my custom dimensions",
		LinkID:           link.ID,
		IPAddress:        "123.123.123.123",                               // Visitor's IP Address
		Referer:          "https://somewebsite.com/page",                  // If there was a referer
		UserAgent:        "Mozilla/5.0 Chrome/51.0.2704.64 Safari/537.36", // Visitor's user agent
	}
	if visitorSession, err = createVisitorSession(visitorSession); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Fire a conversion on a goal (by user id)
	//
	var newConversion *tonicpow.Conversion
	if newConversion, err = TonicPowAPI.CreateConversionByUserID(goal.ID, user.ID, "", 0.00, 5); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("successful conversion event: %d", newConversion.ID)

	//
	// Example: Fire a conversion on a goal (by visitor)
	//
	if newConversion, err = TonicPowAPI.CreateConversionByGoalID(goal.ID, visitorSession.TncpwSession, "", 0.00, 10); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("successful conversion event: %d payout after: %s", newConversion.ID, newConversion.PayoutAfter)

	//
	// Example: Fire a conversion on a goal (by name)
	//
	/*
		if newConversion, err = TonicPowAPI.CreateConversionByGoalName(goal.Name, visitorSession.TncpwSession, "", 0.00, 10); err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("successful conversion event: %d payout after: %s", newConversion.ID, newConversion.PayoutAfter)
	*/

	//
	// Example: Get conversion
	//
	var conversion *tonicpow.Conversion
	if conversion, err = TonicPowAPI.GetConversion(newConversion.ID); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("got conversion: %d", conversion.ID)

	//
	// Example: Cancel a Delayed Conversion
	//
	if conversion, err = TonicPowAPI.CancelConversion(conversion.ID, "not needed anymore"); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("conversion status: %s", conversion.Status)

	//
	// Example: Create Conversion by User ID
	//
	if newConversion, err = TonicPowAPI.CreateConversionByUserID(1, 2, "", 0.00, 0); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create Conversion by User ID (another!)
	//
	if newConversion, err = TonicPowAPI.CreateConversionByUserID(1, 2, "", 0.00, 0); err != nil {
		log.Fatal(err.Error())
	}

	//
	// Example: Create Conversion by User ID (percent based conversion - using a purchase amount, or some tx of value)
	//
	/*
		if newConversion, err = TonicPowAPI.CreateConversionByUserID(goalPercent.ID, 2, "", 5.00, 0); err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("new conversion: %d", newConversion.ID)
	*/

	log.Println("examples complete!")
}

//
// Example Functions
//

func createAdvertiserProfile(profile *tonicpow.AdvertiserProfile, userSessionToken string) (createdProfile *tonicpow.AdvertiserProfile, err error) {
	if createdProfile, err = TonicPowAPI.CreateAdvertiserProfile(profile, userSessionToken); err != nil {
		log.Fatalf("create advertiser failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("advertiser profile %s id %d created", createdProfile.Name, createdProfile.ID)
	}
	return
}

func getAdvertiserProfile(advertiserProfileID uint64, userSessionToken string) (advertiserProfile *tonicpow.AdvertiserProfile, err error) {
	if advertiserProfile, err = TonicPowAPI.GetAdvertiserProfile(advertiserProfileID, userSessionToken); err != nil {
		log.Fatalf("get advertiser profile failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got advertiser profile by id %d", advertiserProfile.ID)
	}
	return
}

func createCampaign(campaign *tonicpow.Campaign, userSessionToken string) (createdCampaign *tonicpow.Campaign, err error) {
	if createdCampaign, err = TonicPowAPI.CreateCampaign(campaign, userSessionToken); err != nil {
		log.Fatalf("create campaign failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("campaign %s id %d created", createdCampaign.Title, createdCampaign.ID)
	}
	return
}

func getCampaign(campaignID uint64, userSessionToken string) (campaign *tonicpow.Campaign, err error) {
	if campaign, err = TonicPowAPI.GetCampaign(campaignID, userSessionToken); err != nil {
		log.Fatalf("get campaign failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got campaign by id %d", campaign.ID)
	}
	return
}

func createGoal(goal *tonicpow.Goal, userSessionToken string) (createdGoal *tonicpow.Goal, err error) {
	if createdGoal, err = TonicPowAPI.CreateGoal(goal, userSessionToken); err != nil {
		log.Fatalf("create goal failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("goal %s id %d created", createdGoal.Title, createdGoal.ID)
	}
	return
}

func getGoal(goalID uint64, userSessionToken string) (goal *tonicpow.Goal, err error) {
	if goal, err = TonicPowAPI.GetGoal(goalID, userSessionToken); err != nil {
		log.Fatalf("get goal failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got goal by id %d", goal.ID)
	}
	return
}

func createLink(link *tonicpow.Link, userSessionToken string) (createdLink *tonicpow.Link, err error) {
	if createdLink, err = TonicPowAPI.CreateLink(link, userSessionToken); err != nil {
		log.Fatalf("create link failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("link %s id %d created", createdLink.ShortCode, createdLink.ID)
	}
	return
}

func createLinkByURL(link *tonicpow.Link, userSessionToken string) (createdLink *tonicpow.Link, err error) {
	if createdLink, err = TonicPowAPI.CreateLinkByURL(link, userSessionToken); err != nil {
		log.Fatalf("create link failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("link %s id %d created", createdLink.ShortCode, createdLink.ID)
	}
	return
}

func getLink(linkID uint64, userSessionToken string) (link *tonicpow.Link, err error) {
	if link, err = TonicPowAPI.GetLink(linkID, userSessionToken); err != nil {
		log.Fatalf("get link failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got link by id %d", link.ID)
	}
	return
}

func getLinks(userID uint64, userSessionToken string) (results *tonicpow.LinkResults, err error) {
	if results, err = TonicPowAPI.ListLinksByUserID(userID, userSessionToken, 1, 20); err != nil {
		log.Fatalf("get link failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got link(s) %d - page: %d", results.Results, results.CurrentPage)
	}
	return
}

func createVisitorSession(visitorSession *tonicpow.VisitorSession) (createdSession *tonicpow.VisitorSession, err error) {
	if createdSession, err = TonicPowAPI.CreateVisitorSession(visitorSession); err != nil {
		log.Fatalf("create visitor session failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("visitor session %s created", createdSession.TncpwSession)
	}
	return
}

func getVisitorSession(tncpwSession string) (visitorSession *tonicpow.VisitorSession, err error) {
	if visitorSession, err = TonicPowAPI.GetVisitorSession(tncpwSession); err != nil {
		log.Fatalf("get visitor session failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got visitor session by %s", tncpwSession)
	}
	return
}

func getCurrentRate(currency string) (rate *tonicpow.Rate, err error) {
	if rate, err = TonicPowAPI.GetCurrentRate(currency, 0); err != nil {
		log.Fatalf("get rate failed - api error: %s data: %s - local error: %s", TonicPowAPI.LastRequest.Error.Message, TonicPowAPI.LastRequest.Error.Data, err.Error())
	} else {
		log.Printf("got rate by currency %s price in sats: %d", currency, rate.PriceInSatoshis)
	}
	return
}
