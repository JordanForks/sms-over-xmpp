package aaisp

import (
	"net/http"
	"net/url"
	"io/ioutil"

	"src.agwa.name/sms-over-xmpp"
	"src.agwa.name/sms-over-xmpp/httputil"
)

type Provider struct {
	service *smsxmpp.Service
	username string
	password string
	httpPassword string
}

func (provider *Provider) Type() string {
	return "aaisp"
}

func (provider *Provider) Send(message *smsxmpp.Message) error {
	// https://support.aa.net.uk/SMS_API
	request := make(url.Values)
	request.Set("username", provider.username)
	request.Set("password", provider.password)
	request.Set("da", message.To)
	request.Set("ud", message.Body)
	request.Set("oa", message.From)

	if err := provider.sendSMS(request); err != nil {
		return err
	}

	return nil
}

func (provider *Provider) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/inbound-sms", provider.handleInboundSMS)
	//mux.HandleFunc("/delivery-receipt", provider.handleDeliveryReceipt) TODO: handle delivery receipts
	return httputil.RequireHTTPAuthHandler(provider.httpPassword, mux)
}

func (provider *Provider) handleInboundSMS(w http.ResponseWriter, req *http.Request) {
	requestBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "400 Bad Request: unable to read request body", 400)
		return
	}

	request, err := url.ParseQuery(string(requestBytes))
	if err != nil {
		http.Error(w, "400 Bad Request: unable to parse request body", 400)
		return
	}

	message := smsxmpp.Message{
		From: request.Get("oa"),
		To: request.Get("da"),
		Body: request.Get("ud"),
	}
	if err := provider.service.Receive(&message); err != nil {
		http.Error(w, "500 Internal Server Error: failed to receive message", 500)
		return
	}

	w.WriteHeader(204)
}

func MakeProvider(service *smsxmpp.Service, config smsxmpp.ProviderConfig) (smsxmpp.Provider, error) {
	return &Provider {
		service: service,
		username: config["username"],
		password: config["password"],
		httpPassword: config["http_password"],
	}, nil
}

func init() {
	smsxmpp.RegisterProviderType("aaisp", MakeProvider)
}

