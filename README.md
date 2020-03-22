# Abuse
This package provides some simple functions to help detect bots or users that are using a temporary or disposable email address.

Bot detection is performed using Google's [reCaptcha V3](https://developers.google.com/recaptcha/intro).  You will need to visit the link, register your website and generate server and client keys.

The list of domains associated with temporary email addresses is a static list that we update from time-to-time.  Create issues if you discover or contact us on @packetriot on Twitter to let us know about more.  Note, there is some logic you can use to automatically maintain your list that we'll share below.

This package was built for our purposes but if it can be improved further to serve a more general audience I'd be glad to maintain this package.

## Captcha
The captcha functions provided in this package require the global variable `Captcha`, defined in `captcha.go`, to be correctly populated with a Server and Client key values.  

```
var (
	Captcha CaptchaKey
)

type CaptchaKey struct {
	Server string `json:"server"`
	Client string `json:"client`
}
```

The function `VerifyCaptcha()` will be useless without valid values.  This function uses a token that is generated on the client-side (browser) and is normally sent using a POST request to your backend.  This function requires this token value and the IP address of the user to verify the captcha with Google.  You'll be provided a score in the response.  A score of `0.0` indicates a bot, although, you can configure your tolerance with the global variable `MinCaptchaThreshold` in captcha.go.

## Abusive Domains
These functions are extremely easily to use.  Use the function `IsTempEmail()` and `IsAbusiveDomain()` to check if an email address or domain is associated to a known temporary email address provider.  

These functions don't require any initialization to work, but you may want to consider initialize the package.  

Many of these providers buy and setup new domains for their email addresses.  However, the service is still being hosted and served by `temp-mail.org`.  If you're onboarding process includes an email confirmation link, many times the `Referrer` in those HTTP GET requests will be set to the primary domain of the temporary email provider.

Here is an example to demonstrate this idea:

```
import (
	"abuse"
)

func confirmEmail(s *Session, w http.ResponseWriter, r *http.Request) {
	user := getUserFromSession(s)
	if referrer :=  r.Header.Get("Referrer"); len(referrer) > 0 {
		u, _ := url.Parse(referrer)
		if abuse.IsAbusiveDomain(u.Host) {
			// Yes, must be a new domain for emails but coming from a provider we know
			abuse.Add(u.Host)

			serveNoTempEmail(s, w, r)
			return
		}
	}

	// ...
}

func getUserFromSession(s *Session) *User {
	// ...
	return u
}

type User struct {
	Email string
}

```

You will want to initialize the package with a path where to write new domains that are collection during runtime.

```
abuse.Init("/path/to/app-data/abusive-domains.txt"
defer abuse.Close()
```

## Updates

As we discover more domains associated to temporary email provider, we add them to the file `domains.txt`.  We use the small utility in `cmd/gen-static.go` to read this input and it will write to standard out a Go source file that can be used to initialize the internal map `tempDomain` then the package is initialized. 

Example usage:
```
go build cmd/gen-static.go

./cmd/gen-static -input domains.txt | gofmw > domains-static.go
```

