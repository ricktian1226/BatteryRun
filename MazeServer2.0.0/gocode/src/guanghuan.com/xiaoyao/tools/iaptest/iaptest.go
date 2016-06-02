package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	IAP_SUCCESS = 0

	IAP_STEP_VELIDATE_DATA         = 1
	IAP_INVALID_USER               = IAP_STEP_VELIDATE_DATA + 1
	IAP_INVALID_ORDER              = IAP_STEP_VELIDATE_DATA + 2
	IAP_INVALID_ORDER_SERVER_STATE = IAP_STEP_VELIDATE_DATA + 3
	IAP_INVALID_GOODS              = IAP_STEP_VELIDATE_DATA + 4
	IAP_INVALID_RECEIPT            = IAP_STEP_VELIDATE_DATA + 5

	IAP_STEP_UPDATE_SERVER_STATE = 10
	IAP_UPDATE_SERVER_STATE_FAIL = IAP_STEP_UPDATE_SERVER_STATE + 1

	IAP_STEP_VERIFY_RECEIPT_PREPARE       = 20
	IAP_FAILED_BEFORE_INVOKE_APPLE_VARIFY = IAP_STEP_VERIFY_RECEIPT_PREPARE + 1
	IAP_FAILED_TO_INVOKE_APPLE_VERIFY     = IAP_STEP_VERIFY_RECEIPT_PREPARE + 2
	IAP_APPLE_VARIFY_INVALID_RESPONSE     = IAP_STEP_VERIFY_RECEIPT_PREPARE + 3

	IAP_STEP_VERIFY_RECEIPT_POST = 30
	IAP_APPLE_VARIFY_FAIL        = IAP_STEP_VERIFY_RECEIPT_POST + 1

	IAP_STEP_ADUITION         = 40
	IAP_UPDATE_IAP_STATE_FAIL = IAP_STEP_ADUITION + 1
	IAP_UPDATE_ACCOUNT_FAIL   = IAP_STEP_ADUITION + 2
)

// 向app store 请求验证
func VerifyReceipt(receipt_data string, isProduction bool) (fail_reason int32, err error) {
	type Receipt struct {
		Receipt_data string `json:"receipt-data,omitempty"`
	}
	var (
		verifyUrl string
		receipt   Receipt
		rcpt      []byte
	)
	fail_reason = IAP_SUCCESS
	log.Printf("Verify receipt via apple, receipt(%d): \n%s", len(receipt_data), receipt_data)

	if isProduction {
		verifyUrl = "https://buy.itunes.apple.com/verifyReceipt"
	} else {
		verifyUrl = "https://sandbox.itunes.apple.com/verifyReceipt"
	}
	//var verifyUrl = "https://buy.itunes.apple.com/verifyReceipt"

	receipt.Receipt_data = receipt_data
	rcpt, err = json.Marshal(receipt)
	if err != nil {
		fail_reason = IAP_FAILED_BEFORE_INVOKE_APPLE_VARIFY
		log.Printf("[VerifyReceipt] Error marshal json msg: %s", err.Error())
		return
	}
	if len(rcpt) > 20 {
		log.Printf("[VerifyReceipt] app store verify receipt:%s...%s (%d)", string(rcpt[:10]), string(rcpt[len(rcpt)-10:]), len(rcpt))
	} else {
		log.Printf("[VerifyReceipt] app store verify receipt:%s (%d)", string(rcpt), len(rcpt))
	}

	var verifyResp *http.Response
	verifyResp, err = http.Post(verifyUrl, "application/json", bytes.NewReader(rcpt))
	defer verifyResp.Body.Close()

	log.Printf("[VerifyReceipt] get response from apple, err: %v", err)

	if err != nil {
		fail_reason = IAP_FAILED_TO_INVOKE_APPLE_VERIFY
		log.Printf("[VerifyReceipt] app store verify rst err: %v", err)
	} else {
		// IAP_STEP_VERIFY_RECEIPT_POST
		var (
			r           interface{}
			res         map[string]interface{}
			value       interface{}
			isType      bool
			verify_data []byte
		)

		verify_data, err = ioutil.ReadAll(verifyResp.Body)
		log.Printf("Apple returns [byte]: \n%v", verify_data)
		decoder := json.NewDecoder(bytes.NewReader(verify_data))
		decoder.Decode(&r)
		log.Printf("Apple returns [json]: \n%v", r)
		res, isType = r.(map[string]interface{})
		if !isType {
			fail_reason = IAP_APPLE_VARIFY_INVALID_RESPONSE
			log.Printf("respone is not a valid json msg: %s", string(verify_data))
			err = errors.New("respone is not a valid json msg")
		} else {
			value = res["status"]
			log.Printf("[VerifyReceipt] status: %v", value)
			if value == float64(0) {
				fail_reason = IAP_SUCCESS
			} else {
				fail_reason = IAP_APPLE_VARIFY_FAIL
				s := fmt.Sprintf("apple returns: %0.0f", value.(float64))
				err = errors.New(s)
				log.Printf("[VerifyReceipt] failed with status: %v", value)
			}
		}
	}
	return
}

func main() {
	var (
		receipt string
	)
	receipt = `ewoJInNpZ25hdHVyZSIgPSAiQXBkeEpkdE53UFUyckE1L2NuM2tJTzFPVGsyNWZlREthMGFhZ3l5UnZlV2xjRmxnbHY2UkY2em5raUJTM3VtOVVjN3BWb2IrUHFaUjJUOHd5VnJITnBsb2YzRFgzSXFET2xXcSs5MGE3WWwrcXJSN0E3ald3dml3NzA4UFMrNjdQeUhSbmhPL0c3YlZxZ1JwRXI2RXVGeWJpVTFGWEFpWEpjNmxzMVlBc3NReEFBQURWekNDQTFNd2dnSTdvQU1DQVFJQ0NHVVVrVTNaV0FTMU1BMEdDU3FHU0liM0RRRUJCUVVBTUg4eEN6QUpCZ05WQkFZVEFsVlRNUk13RVFZRFZRUUtEQXBCY0hCc1pTQkpibU11TVNZd0pBWURWUVFMREIxQmNIQnNaU0JEWlhKMGFXWnBZMkYwYVc5dUlFRjFkR2h2Y21sMGVURXpNREVHQTFVRUF3d3FRWEJ3YkdVZ2FWUjFibVZ6SUZOMGIzSmxJRU5sY25ScFptbGpZWFJwYjI0Z1FYVjBhRzl5YVhSNU1CNFhEVEE1TURZeE5USXlNRFUxTmxvWERURTBNRFl4TkRJeU1EVTFObG93WkRFak1DRUdBMVVFQXd3YVVIVnlZMmhoYzJWU1pXTmxhWEIwUTJWeWRHbG1hV05oZEdVeEd6QVpCZ05WQkFzTUVrRndjR3hsSUdsVWRXNWxjeUJUZEc5eVpURVRNQkVHQTFVRUNnd0tRWEJ3YkdVZ1NXNWpMakVMTUFrR0ExVUVCaE1DVlZNd2daOHdEUVlKS29aSWh2Y05BUUVCQlFBRGdZMEFNSUdKQW9HQkFNclJqRjJjdDRJclNkaVRDaGFJMGc4cHd2L2NtSHM4cC9Sd1YvcnQvOTFYS1ZoTmw0WElCaW1LalFRTmZnSHNEczZ5anUrK0RyS0pFN3VLc3BoTWRkS1lmRkU1ckdYc0FkQkVqQndSSXhleFRldngzSExFRkdBdDFtb0t4NTA5ZGh4dGlJZERnSnYyWWFWczQ5QjB1SnZOZHk2U01xTk5MSHNETHpEUzlvWkhBZ01CQUFHamNqQndNQXdHQTFVZEV3RUIvd1FDTUFBd0h3WURWUjBqQkJnd0ZvQVVOaDNvNHAyQzBnRVl0VEpyRHRkREM1RllRem93RGdZRFZSMFBBUUgvQkFRREFnZUFNQjBHQTFVZERnUVdCQlNwZzRQeUdVakZQaEpYQ0JUTXphTittVjhrOVRBUUJnb3Foa2lHOTJOa0JnVUJCQUlGQURBTkJna3Foa2lHOXcwQkFRVUZBQU9DQVFFQUVhU2JQanRtTjRDL0lCM1FFcEszMlJ4YWNDRFhkVlhBZVZSZVM1RmFaeGMrdDg4cFFQOTNCaUF4dmRXLzNlVFNNR1k1RmJlQVlMM2V0cVA1Z204d3JGb2pYMGlreVZSU3RRKy9BUTBLRWp0cUIwN2tMczlRVWU4Y3pSOFVHZmRNMUV1bVYvVWd2RGQ0TndOWXhMUU1nNFdUUWZna1FRVnk4R1had1ZIZ2JFL1VDNlk3MDUzcEdYQms1MU5QTTN3b3hoZDNnU1JMdlhqK2xvSHNTdGNURXFlOXBCRHBtRzUrc2s0dHcrR0szR01lRU41LytlMVFUOW5wL0tsMW5qK2FCdzdDMHhzeTBiRm5hQWQxY1NTNnhkb3J5L0NVdk02Z3RLc21uT09kcVRlc2JwMGJzOHNuNldxczBDOWRnY3hSSHVPTVoydG04bnBMVW03YXJnT1N6UT09IjsKCSJwdXJjaGFzZS1pbmZvIiA9ICJld29KSW05eWFXZHBibUZzTFhCMWNtTm9ZWE5sTFdSaGRHVXRjSE4wSWlBOUlDSXlNREV5TFRBM0xURXlJREExT2pVME9qTTFJRUZ0WlhKcFkyRXZURzl6WDBGdVoyVnNaWE1pT3dvSkluQjFjbU5vWVhObExXUmhkR1V0YlhNaUlEMGdJakV6TkRJd09UYzJOelU0T0RJaU93b0pJbTl5YVdkcGJtRnNMWFJ5WVc1ellXTjBhVzl1TFdsa0lpQTlJQ0l4TnpBd01EQXdNamswTkRrME1qQWlPd29KSW1KMmNuTWlJRDBnSWpFdU5DSTdDZ2tpWVhCd0xXbDBaVzB0YVdRaUlEMGdJalExTURVME1qSXpNeUk3Q2draWRISmhibk5oWTNScGIyNHRhV1FpSUQwZ0lqRTNNREF3TURBeU9UUTBPVFF5TUNJN0Nna2ljWFZoYm5ScGRIa2lJRDBnSWpFaU93b0pJbTl5YVdkcGJtRnNMWEIxY21Ob1lYTmxMV1JoZEdVdGJYTWlJRDBnSWpFek5ESXdPVGMyTnpVNE9ESWlPd29KSW1sMFpXMHRhV1FpSUQwZ0lqVXpOREU0TlRBME1pSTdDZ2tpZG1WeWMybHZiaTFsZUhSbGNtNWhiQzFwWkdWdWRHbG1hV1Z5SWlBOUlDSTVNRFV4TWpNMklqc0tDU0p3Y205a2RXTjBMV2xrSWlBOUlDSmpiMjB1ZW1Wd2RHOXNZV0l1WTNSeVltOXVkWE11YzNWd1pYSndiM2RsY2pFaU93b0pJbkIxY21Ob1lYTmxMV1JoZEdVaUlEMGdJakl3TVRJdE1EY3RNVElnTVRJNk5UUTZNelVnUlhSakwwZE5WQ0k3Q2draWIzSnBaMmx1WVd3dGNIVnlZMmhoYzJVdFpHRjBaU0lnUFNBaU1qQXhNaTB3TnkweE1pQXhNam8xTkRvek5TQkZkR012UjAxVUlqc0tDU0ppYVdRaUlEMGdJbU52YlM1NlpYQjBiMnhoWWk1amRISmxlSEJsY21sdFpXNTBjeUk3Q2draWNIVnlZMmhoYzJVdFpHRjBaUzF3YzNRaUlEMGdJakl3TVRJdE1EY3RNVElnTURVNk5UUTZNelVnUVcxbGNtbGpZUzlNYjNOZlFXNW5aV3hsY3lJN0NuMD0iOwoJInBvZCIgPSAiMTciOwoJInNpZ25pbmctc3RhdHVzIiA9ICIwIjsKfQ==`
	code, err := VerifyReceipt(receipt, true)
	log.Printf("code: %d, err: %v", code, err)
}
