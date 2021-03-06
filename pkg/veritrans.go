package pkg

import (
	"os"

	"github.com/david1992121/veritrans-microservice/internal/veritrans"
)

// ServiceConfig struct
type ServiceConfig struct {
	MDKConfig        veritrans.MDKConfig
	ConnectionConfig veritrans.ConnectionConfig
}

// GetServiceConfig initializes the service configuration
func GetServiceConfig() *ServiceConfig {
	mdkConfig := veritrans.MDKConfig{
		APIURL:   os.Getenv("MDK_API_URL"),
		APIToken: os.Getenv("MDK_API_TOKEN"),
	}
	connectionConfig := veritrans.ConnectionConfig{
		MerchantCCID:     os.Getenv("MERCHANT_CCID"),
		MerchantPassword: os.Getenv("MERCHANT_PASSWORD"),
		AccountAPIURL:    os.Getenv("ACCOUNT_API_URL"),
		PaymentAPIURL:    os.Getenv("PAYMENT_API_URL"),
		SearchAPIURL:     os.Getenv("SEARCH_API_URL"),
		TxnVersion:       os.Getenv("TXN_VERSION"),
		DummyRequest:     os.Getenv("DUMMY_REQUEST"),
	}

	serviceConfig := &ServiceConfig{
		MDKConfig:        mdkConfig,
		ConnectionConfig: connectionConfig,
	}

	return serviceConfig
}

type veritransService struct {
	MDKService     *veritrans.MDKService
	AccountService *veritrans.AccountService
	PaymentService *veritrans.PaymentService
}

// NewService initializes the veritrans service
func NewService(config *ServiceConfig) Service {
	mdkService := veritrans.NewMDKService(config.MDKConfig)

	paymentService, _ := veritrans.NewPaymentService(config.ConnectionConfig)
	accountService := veritrans.NewAccountService(config.ConnectionConfig)
	return &veritransService{
		MDKService:     mdkService,
		AccountService: accountService,
		PaymentService: paymentService,
	}
}

func (v *veritransService) GetMDKToken(cardInfo *veritrans.ClientCardInfo) (string, error) {
	return v.MDKService.GetCardToken(cardInfo)
}

func (v *veritransService) CreateAccount(accountParam *veritrans.AccountParam) (*veritrans.Account, error) {
	return v.AccountService.CreateAccount(accountParam)
}

func (v *veritransService) UpdateAccount(accountParam *veritrans.AccountParam) (*veritrans.Account, error) {
	return v.AccountService.UpdateAccount(accountParam)
}

func (v *veritransService) CreateCard(accountParam *veritrans.AccountParam) (*veritrans.Account, error) {
	return v.AccountService.CreateCard(accountParam)
}

func (v *veritransService) UpdateCard(accountParam *veritrans.AccountParam) (*veritrans.Account, error) {
	return v.AccountService.UpdateCard(accountParam)
}

func (v *veritransService) DeleteCard(accountParam *veritrans.AccountParam) (*veritrans.Account, error) {
	return v.AccountService.DeleteCard(accountParam)
}

func (v *veritransService) GetCard(accountParam *veritrans.AccountParam) (*veritrans.Account, error) {
	return v.AccountService.GetCard(accountParam)
}

func (v *veritransService) Authorize(param *veritrans.Params) error {
	_, err := v.PaymentService.Authorize(param, veritrans.PaymentServiceType(veritrans.PayCard))
	return err
}

func (v *veritransService) Capture(param *veritrans.Params) error {
	_, err := v.PaymentService.Capture(param, veritrans.PaymentServiceType(veritrans.PayCard))
	return err
}

func (v *veritransService) Cancel(param *veritrans.Params) error {
	_, err := v.PaymentService.Cancel(param, veritrans.PaymentServiceType(veritrans.PayCard))
	return err
}
