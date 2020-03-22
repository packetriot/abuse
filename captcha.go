package abuse

/*
Visit the following site to generate client/server captcha keys
that you can use for your site:

https://developers.google.com/recaptcha/intro

Usage summary:
Include the snippet below in your website to track user-behavior and determine
if this is a bot or not.  You'll need to customize the callback to
handle the token value.  Typically, you can set can the token received
as the value to a hidden element in an HTML form.  It'll be submitted to
your server where it can be verified.

<script src="https://www.google.com/recaptcha/api.js?render=_reCAPTCHA_site_key"></script>
<script>
grecaptcha.ready(function() {
    grecaptcha.execute('_reCAPTCHA_site_key_', {action: 'homepage'}).then(function(token) {
       ...
    });
});
</script>
*/

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// Server-side token verification URL
	recaptchaURL = "https://www.google.com/recaptcha/api/siteverify"
)

var (
	Captcha CaptchaKey

	// Google indicates 0.0 is a bot
	MinCaptchaThreshold = 0.1

	// Token verification settings
	MaxRetryTokenVerification               = 5
	TokenVerificationInterval time.Duration = time.Second
)

type CaptchaKey struct {
	Server string `json:"server"`
	Client string `json:"client`
}

func (c *CaptchaKey) IsValid() bool {
	return len(c.Server) > 0 && len(c.Client) > 0
}

type CaptchaResponse struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"challenge_ts"`
	Score     float64   `json:"score"`
	Hostname  string    `json:"hostname"`
	Errors    []string  `json:"error-codes"`
}

func (c *CaptchaResponse) IsBot() bool {
	return (c.Score <= MinCaptchaThreshold)
}

func VerifyCaptcha(token, remoteIP string) (*CaptchaResponse, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("invalid (zero-length) response token passed in")
	} else if !Captcha.IsValid() {
		return nil, fmt.Errorf("invalid recaptcha configuration, expecting non-zero-length server/client keys")
	}

	// Check the captcha token, attempt 5 times in case the network fails
	// or something else happens on Google's side.
	var resp *http.Response
	var err error

	for i := 0; i < MaxRetryTokenVerification; i++ {
		// Wait between retrying the request.
		if i > 0 {
			time.Sleep(TokenVerificationInterval)
		}

		resp, err = http.PostForm(recaptchaURL, url.Values{
			"secret":   {Captcha.Server},
			"response": {token},
			"remoteip": {remoteIP},
		})

		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// Decode the response and return it the user.
	cr := &CaptchaResponse{}
	err = jsonDecode(resp.Body, cr)
	if err != nil {
		return nil, err
	}

	return cr, nil
}
