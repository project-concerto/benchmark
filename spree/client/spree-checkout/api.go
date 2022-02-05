// The official document for API V1 is gone. For history version, see:
// https://web.archive.org/web/20200529131302/https://guides.spreecommerce.org/api/
package checkout

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"gap-lock/utils"

	"go.uber.org/zap"
)

const BASE_URL = "localhost:4000"

// API token for the admin(spree@example.com:spree123)
// It may change if you repopulate the database.
const APIToken = "9b5d25eb6d0397c436fbdb4c93d22e9e1732cb1a49c06708"
const MAX_VARIANT_ID = 232

const ADDRESS_PAYLOAD = `{
   "order":{
	  "email": "john@snow.org",
      "bill_address_attributes":{
		  "firstname": "John",
		  "lastname": "Snow",
		  "address1": "7735 Old Georgetown Road",
		  "city": "Bethesda",
		  "phone": "3014445002",
		  "zipcode": "20814",
		  "state_name": "MD",
		  "country_iso": "US"
      },
      "ship_address_attributes":{
		  "firstname": "John",
		  "lastname": "Snow",
		  "address1": "7735 Old Georgetown Road",
		  "city": "Bethesda",
		  "phone": "3014445002",
		  "zipcode": "20814",
		  "state_name": "MD",
		  "country_iso": "US"
      }
   }
}`
const SHIPPING_METHOD_PAYLOAD = `{
   "order":{
	  "shipments_attributes": [ {
		  "id": "%v",
		  "selected_shipping_rate_id": "%v"
	  } ]
   }
}`
const V1_ADDRESS_PAYLOAD = `{
	"order": {
		"bill_address_attributes": {
			"firstname": "John",
			"lastname": "Doe",
			"address1": "233 36th Ave Ne",
			"city": "St Petersburg",
			"phone": "3014445002",
			"zipcode": "33704-1535",
			"state_id": 516,
			"country_id": 224
		},
		"ship_address_attributes": {
			"firstname": "John",
			"lastname": "Doe",
			"address1": "233 36th Ave Ne",
			"city": "St Petersburg",
			"phone": "3014445002",
			"zipcode": "33704-1535",
			"state_id": 516,
			"country_id": 224
		}
	}
}`
const V1_PAYMENT_PAYLOAD = `{
	"order": {
		"payments_attributes": [{
			"payment_method_id": "2"
		}]
	},
	"payment_source": {
		"2": {
			"number": "4111111111111111",
			"month": "1",
			"year": "2017",
			"verification_value": "123",
			"name": "John Smith"
		}
	}
}
`

// payment_method_id: 1 is the bogus payment method.
const PAYMENT_METHOD_PAYLOAD = `{
"order": {
  "payments_attributes": [{
    "payment_method_id": "1",
    "source_attributes": {
      "gateway_payment_profile_id": "BGS-1JqvNB2eZvKYlo2C5OlqLV7S",
      "gateway_customer_profile_id": "BGS-1JqvNB2eZvKYlo2C5OlqLV7S",
	  "number": "1",
      "cc_type": "visa",
      "last_digits": "1111",
      "month": "10",
      "year": "2026",
      "name": "John Snow"
    }
  }]
}}`

type Checkout struct {
	Index              int32
	APIToken           string
	OrderToken         string
	OrderNumber        string
	access_token       string
	shippment_id       string
	shipping_method_id string
	logger             *zap.SugaredLogger
	httpClient         *http.Client
}

func NewCheckout(api_token string) *Checkout {
	return &Checkout{
		APIToken:   api_token,
		logger:     zap.L().Sugar(),
		httpClient: &http.Client{},
	}
}

func NewCheckoutBoth(index int32, api_token, order_number string) *Checkout {
	return &Checkout{
		Index:       index,
		APIToken:    api_token,
		OrderNumber: order_number,
		logger:      zap.L().Sugar(),
		httpClient:  &http.Client{},
	}
}

// Checkout once for users form utils.APITokens, one by one, from user000 to user999.
// The purpose is to maximize possibility of gap lock conflicts.
type CheckoutAPITokenFactory struct {
	counter int32
	maximum int32 // index used should be less than maximum
}

func NewCheckoutSequentialFactory() *CheckoutAPITokenFactory {
	return &CheckoutAPITokenFactory{
		counter: -1,
		maximum: int32(len(utils.APITokens)),
	}
}

func (f *CheckoutAPITokenFactory) Make() utils.API {
	idx := atomic.AddInt32(&f.counter, 1)
	if idx >= f.maximum {
		log.Fatal("Used up all users!")
	}
	if idx%10 == 0 {
		zap.L().Sugar().Infof("Generate Checkout for %v", idx)
	}
	return NewCheckout(utils.APITokens[idx])
}

type CheckoutBothTokenFactory struct {
	counter int32
	maximum int32
}

func NewCheckoutBothTokenFactory() *CheckoutBothTokenFactory {
	return &CheckoutBothTokenFactory{
		counter: -1,
		maximum: int32(len(utils.APITokens)),
	}
}

func (f *CheckoutBothTokenFactory) Make(threadId int) utils.API {
	idx := atomic.AddInt32(&f.counter, 1)
	if idx >= f.maximum {
		zap.L().Fatal("Used up all users!")
	}
	if idx%10 == 0 {
		zap.L().Sugar().Infof("Generate Checkout for %v", idx)
	}
	return NewCheckoutBoth(idx, utils.APITokens[idx], utils.OrderNumbers[idx])
}

func (f *CheckoutBothTokenFactory) Prepare() {
	refreshDockerVolume("")
}

func (f *CheckoutBothTokenFactory) Stop() {}

// Compared to CheckoutBothTokenFactory, CheckoutAPITokenFactory do orders randomly.
type CheckoutNoContentionFactory struct {
	taskChan chan int32
	stop     context.CancelFunc
}

func NewCheckoutNoContentionFactory() *CheckoutNoContentionFactory {
	tasks := make([]int32, 0)
	for i := int32(0); i < int32(len(utils.APITokens)); i++ {
		tasks = append(tasks, i)
	}
	rand.Shuffle(len(tasks), func(i, j int) {
		tasks[i], tasks[j] = tasks[j], tasks[i]
	})
	ctx, cancel := context.WithCancel(context.Background())

	taskChan := make(chan int32, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case taskChan <- tasks[0]:
				tasks = tasks[1:]
			}
		}
	}()
	return &CheckoutNoContentionFactory{
		taskChan: taskChan,
		stop:     cancel,
	}
}

func (f *CheckoutNoContentionFactory) Make(threadId int) utils.API {
	id := <-f.taskChan
	return NewCheckoutBoth(id, utils.APITokens[id], utils.OrderNumbers[id])
}

func (f *CheckoutNoContentionFactory) Prepare() {
	refreshDockerVolume("spree_postgres_payment")
}

func (f *CheckoutNoContentionFactory) Stop() {
	f.stop()
}

// Experimental
func (c *Checkout) CreateToken() bool {
	type CreateTokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		CreatedAt    int    `json:"created_at"`
	}

	var data = strings.NewReader(`{
  "grant_type": "password",
  "username": "002@example.com",
  "password": "jkl;jkl;"
}`)

	url := fmt.Sprintf("http://%s/spree_oauth/token", BASE_URL)

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res CreateTokenResponse
	json.Unmarshal(bodyText, &res)

	c.access_token = res.AccessToken
	c.logger.Debugw("CreateToken", "access_token", c.access_token)

	return true
}

func (c *Checkout) CreateCart() bool {
	type CreateCartResponse struct {
		Data struct {
			Attributes struct {
				Token string `json:"token"`
			} `json:"attributes"`
		} `json:"data"`
	}

	url := fmt.Sprintf("http://%s/api/v2/storefront/cart", BASE_URL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	res := CreateCartResponse{}
	json.Unmarshal(body, &res)

	c.OrderToken = res.Data.Attributes.Token
	c.logger.Debugw("CreateCart", "order token", c.OrderToken)

	return true
}

// For debugging only, not used in benchmarking.
func (c *Checkout) GetCart() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/cart", BASE_URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

func (c *Checkout) AddItem() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/cart/add_item", BASE_URL)

	variant_id := rand.Intn(MAX_VARIANT_ID) + 1
	var data = strings.NewReader(
		fmt.Sprintf(`{ "variant_id": %d, "quantity": 1 }`, variant_id))

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return false
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	return true
}

func (c *Checkout) NextCheckoutStep() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout/next", BASE_URL)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

func (c *Checkout) AddAddress() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout", BASE_URL)

	payload := strings.NewReader(ADDRESS_PAYLOAD)

	req, err := http.NewRequest("PATCH", url, payload)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

func (c *Checkout) AdvanceCheckout() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout/advance", BASE_URL)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

// See https://api.spreecommerce.org/docs/api-v2/b3A6MzE0Mjc2MQ-list-shipping-rates
func (c *Checkout) ListShippingRates() bool {
	type ListShippingRatesResponse struct {
		Data []struct {
			ID            string `json:"id"`
			Relationships struct {
				ShippingRates struct {
					Data []struct {
						ID string `json:"id"`
					} `json:"data"`
				} `json:"shipping_rates"`
			} `json:"relationships"`
		} `json:"data"`
	}

	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout/shipping_rates", BASE_URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	var j ListShippingRatesResponse
	json.Unmarshal(body, &j)

	c.shippment_id = j.Data[0].ID
	c.shipping_method_id = j.Data[0].Relationships.ShippingRates.Data[0].ID

	c.logger.Debug(string(body))

	return true
}

func (c *Checkout) SelectsShippingMethod() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout", BASE_URL)

	payload := strings.NewReader(fmt.Sprintf(SHIPPING_METHOD_PAYLOAD, c.shippment_id, c.shipping_method_id))

	req, err := http.NewRequest("PATCH", url, payload)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

func (c *Checkout) SelectPayment() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout", BASE_URL)

	payload := strings.NewReader(PAYMENT_METHOD_PAYLOAD)

	req, err := http.NewRequest("PATCH", url, payload)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

func (c *Checkout) CompleteCheckout() bool {
	url := fmt.Sprintf("http://%s/api/v2/storefront/checkout/complete", BASE_URL)

	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Order-Token", c.OrderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debug(string(body))

	return true
}

// Experimental
func (c *Checkout) CreateAddress() bool {
	var data = strings.NewReader(`{
  "address": {
    "firstname": "Mark",
    "lastname": "Winterburn",
    "company": "Paper Street Soap Co.",
    "address1": "775 Old Georgetown Road",
    "address2": "3rd Floor",
    "city": "Qethesda",
    "phone": "3488545445002",
    "zipcode": "90210",
    "state_name": "CA",
    "country_iso": "US",
    "label": "Work"
  }
}`)

	url := fmt.Sprintf("http://%s/api/v2/storefront/account/addresses", BASE_URL)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	var bearer = "Bearer " + c.access_token
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	c.logger.Debug(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	return true
}

func (c *Checkout) V1CreateOrder() bool {
	type CreateOrderResponse struct {
		ID     int    `json:"id"`
		Number string `json:"number"`
	}

	url := fmt.Sprintf("http://%s/api/v1/orders.json", BASE_URL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	var res CreateOrderResponse
	json.Unmarshal(body, &res)

	c.OrderNumber = res.Number

	c.logger.Debugw("CreateOrder", "order_number", c.OrderNumber)

	return true
}

func (c *Checkout) V1AddItem() bool {
	url := fmt.Sprintf("http://%s/api/v1/orders/%v/line_items.json", BASE_URL, c.OrderNumber)

	variant_id := rand.Intn(MAX_VARIANT_ID) + 1
	var data = strings.NewReader(
		fmt.Sprintf(`{"line_item": { "variant_id": %d, "quantity": 1 }}`, variant_id))

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		c.logger.Errorw("V1AddItem", "StatusCode", res.StatusCode)
		return false
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	return true
}

func (c *Checkout) V1Advance() bool {
	url := fmt.Sprintf("http://%s/api/v1/checkouts/%v/next.json", BASE_URL, c.OrderNumber)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return false
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	return true
}

func (c *Checkout) V1AddAddress() bool {
	url := fmt.Sprintf("http://%s/api/v1/checkouts/%v.json", BASE_URL, c.OrderNumber)

	payload := strings.NewReader(V1_ADDRESS_PAYLOAD)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return false
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	return true
}

func (c *Checkout) V1AddPayment() bool {
	url := fmt.Sprintf("http://%s/api/v1/checkouts/%v.json", BASE_URL, c.OrderNumber)

	payload := strings.NewReader(V1_PAYMENT_PAYLOAD)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	// c.logger.Infow("V1AddPayment", "StatusCode", res.StatusCode)
	if res.StatusCode/100 != 2 {
		body, _ := io.ReadAll(res.Body)
		c.logger.Warnw("Failed V1AddPayment", "reponse", string(body), "Index", c.Index)
		return false
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	c.logger.Debugw("V1AddPayment", "response", string(resBody))

	return true
}

func (c *Checkout) V1ConfirmOrder() bool {
	url := fmt.Sprintf("http://%s/api/v1/checkouts/%v.json", BASE_URL, c.OrderNumber)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		c.logger.Fatal(err)
	}
	req.Header.Set("X-Spree-Token", c.APIToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return false
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		c.logger.Fatal(err)
	}

	return true
}

func (c *Checkout) RunV2(ctx context.Context) utils.APIResult {
	steps := []utils.StepFn{
		c.CreateCart,
		c.AddItem,
		c.NextCheckoutStep,
		c.AddAddress,
		c.AdvanceCheckout,
		c.ListShippingRates,
		c.SelectsShippingMethod,
		c.SelectPayment,
		c.AdvanceCheckout,
		c.CompleteCheckout,
	}
	return c.RunSteps(ctx, steps)
}

func (c *Checkout) RunV1(ctx context.Context) utils.APIResult {
	steps := []utils.StepFn{
		c.V1CreateOrder,
		c.V1AddItem,
		c.V1Advance,
		c.V1AddAddress,
		c.V1AddPayment,
		// c.V1Advance,
	}
	return c.RunSteps(ctx, steps)
}

func (c *Checkout) RunPrepare(ctx context.Context) utils.APIResult {
	steps := []utils.StepFn{
		c.V1CreateOrder,
		c.V1AddItem,
		c.V1Advance,
		c.V1AddAddress,
		// c.V1Advance,
	}
	return c.RunSteps(ctx, steps)
}

func (c *Checkout) RunV1AddAddress(ctx context.Context) utils.APIResult {
	steps := []utils.StepFn{
		c.V1AddAddress,
	}
	return c.RunSteps(ctx, steps)
}

func (c *Checkout) RunV1AddPayment(ctx context.Context) utils.APIResult {
	// Retry
	// steps := []utils.StepFn{
	// 	c.V1AddPayment,
	// }
	// return c.RunSteps(ctx, steps)

	// No Retry
	result := utils.APIResult{
		SuccessLatency: nil,
		FailLatencies:  make([]int64, 0),
	}
	start := time.Now()
	ok := c.V1AddPayment()
	elapsed := time.Since(start).Microseconds()
	if ok {
		result.SuccessLatency = &elapsed
	} else {
		result.FailLatencies = append(result.FailLatencies, elapsed)
	}
	return result
}

func (c *Checkout) RunSteps(ctx context.Context, steps []utils.StepFn) utils.APIResult {
	result := utils.APIResult{
		SuccessLatency: nil,
		FailLatencies:  make([]int64, 0),
	}
	for _, step := range steps {
		for {
			select {
			case <-ctx.Done():
				goto end
			default:
				start := time.Now()
				ok := step()
				elapsed := time.Since(start).Microseconds()
				if ok {
					result.SuccessLatency = &elapsed
					goto finishon
				} else {
					c.logger.Warnw("Failed step", "step", utils.GetFunctionName(step))
					result.FailLatencies = append(result.FailLatencies, elapsed)
				}
			}
		}
	finishon:
	}
end:
	return result
}

func (c *Checkout) Run(ctx context.Context) utils.APIResult {
	// return c.RunV2(ctx)
	// return c.RunV1(ctx)
	// return c.RunV1AddAddress(ctx)
	return c.RunV1AddPayment(ctx)
}

func refreshDockerVolume(volume string) {
	logger := zap.L()
	logger.Info("Start refresh-docker-volume", zap.String("volume", volume))
	err := exec.Command("./refresh-docker-volume", volume).Run()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Finish refresh-docker-volume")
}
