package prismic

// Webhook represents the webhook object received by Prismic
type Webhook struct {
	Type      string `json:"type"`
	Secret    string `json:"secret"`
	MasterRef string `json:"masterRef"`
	Domain    string `json:"domain"`
	ApiUrl    string `json:"apiUrl"`
}

// VerifySecret checks that the webhook secret matches the one in the configuration
func (w *Webhook) VerifySecret(secret string) bool {
	return w.Secret == secret
}
