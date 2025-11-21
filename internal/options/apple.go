package options

import "github.com/spf13/pflag"

type AppleAuthOptions struct {
	ClientID string `json:"client-id" mapstructure:"client-id"`
	// KeyID       string `json:"key_id"`
	// Secret      string `json:"secret"`
	// TeamID      string `json:"team_id"`
	// RedirectURI string `json:"redirect_uri"`
}

func NewAppleAuthOptions() *AppleAuthOptions {
	return &AppleAuthOptions{}
}

func (o *AppleAuthOptions) Validate() []error {
	errs := []error{}

	return errs
}

func (o *AppleAuthOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ClientID, "client-id", o.ClientID, "Apple Auth Client ID (Bundle ID)")
	// fs.StringVar(&o.KeyID, "key-id", o.KeyID, "Apple Auth Key ID")
	// fs.StringVar(&o.Secret, "secret", o.Secret, "Apple Auth Private Key")
	// fs.StringVar(&o.TeamID, "team-id", o.TeamID, "Apple Auth Team ID")
	// fs.StringVar(&o.RedirectURI, "redirect-uri", o.RedirectURI, "Apple Auth Redirect URI")
}
