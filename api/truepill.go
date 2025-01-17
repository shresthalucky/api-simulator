package api

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sarweshmaharjan/api-simulator.git/common"
	"github.com/sarweshmaharjan/api-simulator.git/storage"
	"github.com/sarweshmaharjan/api-simulator.git/types"
)

const (
	RequestID                 = "5aca912a7g1491a446f5d554212346"
	PrescriptionToken         = "z3acq212jrg1s1312356"
	PatientToken              = "452ac6d1901g4a12356"
	InsuranceToken            = "skhsacy21q83g1111rk4d3uht691235"
	DirectTransferID          = "c8a73aca21554dgc61325"
	TransferPrescriptionToken = "z3q2jrac22131g4211365"
	DirectTransferToken       = "dff0a98ac12124g441112635"
	CopayPrescriptionToken    = "an7hj7gpac62123gyj411h6gpw1235"
	CopayRequestToken         = "ef6p8y6wzac7x123g4h11p6ng1235"
	OrderToken                = "b392d"
	MedicationToken           = "310ca0fa"
)

func GetDirectTransfer(ctx *gin.Context) {
	directTransfer := &types.DirectTransferResponse{
		RequestID: RequestID,
		Status:    "success",
		TimeStamp: float64(time.Now().Unix()),
		Details: &types.DirectTransferDetails{
			DirectTransferID: DirectTransferID,
		},
	}
	sql := storage.PrimayConnection()
	storage.Init(sql)
	// Insert a JSON value into the "simulator" table
	insertSQL := `
		INSERT INTO api_simulator_v1.direct_transfer_response (requestId,status,detailsTransferId) VALUES ($1,$2,$3);
	`
	_, err := sql.Exec(insertSQL, directTransfer.RequestID, directTransfer.Status, directTransfer.Details.DirectTransferID)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(http.StatusCreated, directTransfer)
	go func() {
		SendNotifyRxWebhook(ctx)
		time.Sleep(5 * time.Second)
		SendDirectTransferWebhook(ctx)
	}()
}

func SendNotifyRxWebhook(ctx *gin.Context) {
	notifyRx := &types.NotifyRx{
		PrescriptionToken:         PrescriptionToken,
		PatientToken:              PatientToken,
		TransferPrescriptionToken: TransferPrescriptionToken,
		MedicationName:            "BONJESTA 20-20 MG",
		Prescriber:                "Dr. Strange",
		Location:                  "Hayward, CA",
	}

	notifyRxMap := map[string]interface{}{
		"prescription_token":          notifyRx.PrescriptionToken,
		"patient_token":               notifyRx.PatientToken,
		"transfer_prescription_token": notifyRx.TransferPrescriptionToken,
		"medication_name":             notifyRx.MedicationName,
		"prescriber":                  notifyRx.Prescriber,
		"location":                    notifyRx.Location,
	}
	callbackResponse := &types.CallbackRequest{
		Timestamp:    float64(time.Now().Unix()),
		CallbackType: "NOTIFY_RX",
		Details:      notifyRxMap,
	}

	webHookURL := "http://localhost:8888/api/v1/truepill/callback"

	callbackJson := common.ToJSON(callbackResponse)
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(callbackJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()

}

func SendDirectTransferWebhook(ctx *gin.Context) {
	notifyRx := &types.DirectTransferWebhook{
		PrescriptionToken:   PrescriptionToken,
		PatientToken:        PatientToken,
		Metadata:            TransferPrescriptionToken,
		DirectTransferToken: DirectTransferToken,
		Message:             "Direct Transfer Accepted.",
	}

	notifyRxMap := map[string]interface{}{
		"prescription_token":    notifyRx.PrescriptionToken,
		"patient_token":         notifyRx.PatientToken,
		"metadata":              notifyRx.Metadata,
		"message":               notifyRx.Message,
		"direct_transfer_token": notifyRx.DirectTransferToken,
	}
	callbackResponse := &types.CallbackRequest{
		RequestID:    RequestID,
		Status:       "success",
		Timestamp:    float64(time.Now().Unix()),
		CallbackType: "DIRECT_TRANSFER",
		Details:      notifyRxMap,
	}

	webHookURL := "http://localhost:8888/api/v1/truepill/callback"

	callbackJson := common.ToJSON(callbackResponse)
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(callbackJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func GetPrescriptionDetails(ctx *gin.Context) {
	prescription := &types.Prescription{
		DateWritten:            "2023-09-17T17:35:40.083Z",
		ExpirationDate:         "2023-12-30T17:35:40.083Z",
		DaysSupply:             "84",
		IsRefill:               0,
		LastFilledDate:         "2023-08-30T17:35:40.083Z",
		MedicationSIG:          "Take one Tablet by mouth at the same time daily",
		NumberOfRefillsAllowed: 0,
		Origin:                 "Electronic",
		PrescribedBrandName:    "BONJESTA 20-20 MG",
		PrescribedDrugStrength: "20-20 MG",
		PrescribedGenericName:  "TAKE 1 TABLET DAILY|TOME 1 TABLETA AL DIA",
		PrescribedNDC:          "55494012060",
		PrescribedQuantity:     84,
		PrescribedWrittenName:  "Lutera-28 Tablet",
		Prescriber:             "Dr. Strange",
		PrescriberAddress: &types.PrescriberAddress{
			Name:    "Dr. Strange",
			Company: "",
			Street1: "12345 Avengers Rd",
			Street2: "",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94402",
			Country: "US",
			Phone:   "(800) 888-8888",
			Email:   "dr.strange@avengers.com",
		},
		PrescriberNPI:         "1639349256",
		PrescriberOrderNumber: "11111111",
		QuantityRemaining:     0,
		RefillsRemaining:      0,
		RxNumber:              "1966785",
		PrescriptionToken:     PrescriptionToken,
		Notes:                 "prescription notes",
		ICDCodes: &types.ICDCodes{
			ICD10: []string{"Z30.09"},
			ICD9:  []string{"N94.6"},
		},
		IsDAW:                  false,
		Fillable:               true,
		DEASchedule:            0,
		OriginalPrescribedNDC:  "43547041110",
		DateFilledUTC:          "2020-07-30T20:08:55.000Z",
		PrescribedQuantityUnit: "EA",
	}

	resp := &types.PrescriptionResponse{
		Prescription: prescription,
	}

	ctx.JSON(http.StatusCreated, resp)
}

func GetInsuranceDetails(ctx *gin.Context) {
	insuranceResponse := &types.InsuranceResponse{
		RequestID: RequestID,
		Status:    "success",
		TimeStamp: float64(time.Now().Unix()),
		Details: &types.InsuranceDetails{
			InsuranceToken: InsuranceToken,
		},
	}

	ctx.JSON(http.StatusCreated, insuranceResponse)
}

func GetCopayRequest(ctx *gin.Context) {
	copayRxToken := make([]*types.CopayRequestPrescriptionToken, 0)
	copayPrescriptionToken := &types.CopayRequestPrescriptionToken{
		PrescriptionToken:      PrescriptionToken,
		CopayPrescriptionToken: CopayPrescriptionToken,
		Status:                 "pending",
	}
	copayRxToken = append(copayRxToken, copayPrescriptionToken)
	insuranceResponse := &types.CopayResponse{
		RequestID: RequestID,
		Status:    "success",
		TimeStamp: float64(time.Now().Unix()),
		Details: &types.CopayResponseDetails{
			CopayRequestToken:              CopayRequestToken,
			CopayRequestPrescriptionTokens: copayRxToken,
		},
	}

	ctx.JSON(http.StatusCreated, insuranceResponse)
	go func() {
		SendCopayRejectedWebhook(ctx)
	}()
}

func SendCopayWebhook(ctx *gin.Context) {
	copayCallback := &types.CopayWebhook{
		PatientToken: PatientToken,
		Metadata:     TransferPrescriptionToken,
		Message:      "Direct Transfer Accepted.",
	}

	copayRxMap := map[string]interface{}{
		"prescriptions": []map[string]interface{}{
			{
				"prescription_token":               PrescriptionToken,
				"insurance_token":                  InsuranceToken,
				"status":                           "completed",
				"copay_request_prescription_token": CopayPrescriptionToken,
				"next_fill_date":                   nil,
				"copay_amount":                     "10",
				"fill_number":                      1,
				"days_supply":                      "30",
				"quantity":                         30,
				"claims": []map[string]interface{}{
					{
						"insurance_token": InsuranceToken,
						"billing_order":   1,
						"status":          "Paid",
						"copay_amount":    "50",
					},
					{
						"insurance_token": InsuranceToken,
						"billing_order":   2,
						"status":          "Paid",
						"copay_amount":    "10",
					},
				},
			},
		},
		"patient_token": copayCallback.PatientToken,
		"metadata":      copayCallback.Metadata,
		"message":       copayCallback.Message,
	}

	callbackResponse := &types.CallbackRequest{
		RequestID:    RequestID,
		Status:       "success",
		Timestamp:    float64(time.Now().Unix()),
		CallbackType: "COPAY",
		Details:      copayRxMap,
	}

	webHookURL := "http://localhost:8888/api/v1/truepill/callback"

	callbackJson := common.ToJSON(callbackResponse)
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(callbackJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func SendCopayRejectedWebhook(ctx *gin.Context) {
	copayCallback := &types.CopayWebhook{
		PatientToken: PatientToken,
		Metadata:     TransferPrescriptionToken,
		Message:      "Direct Transfer Accepted.",
	}

	copayRxMap := map[string]interface{}{
		"prescriptions": []map[string]interface{}{
			{
				"prescription_token":               PrescriptionToken,
				"insurance_token":                  InsuranceToken,
				"status":                           "completed",
				"copay_request_prescription_token": CopayPrescriptionToken,
				"next_fill_date":                   nil,
				"copay_amount":                     "10",
				"fill_number":                      1,
				"days_supply":                      "30",
				"quantity":                         30,
				"claims": []map[string]interface{}{
					{
						"insurance_token": InsuranceToken,
						"billing_order":   1,
						"status":          "Rejected",
						"claim_reject_codes": []map[string]interface{}{
							{
								"code":    "70",
								"message": "Product/Service Not Covered - Plan/Benefit Exclusion",
							}, {
								"code":    "MR",
								"message": "Product Not On Formulary",
							},
						},
					},
					{
						"insurance_token": InsuranceToken,
						"billing_order":   2,
						"status":          "Paid",
						"copay_amount":    "10",
					},
				},
				"claim_reject_codes": []map[string]interface{}{
					{
						"code":    "70",
						"message": "Product/Service Not Covered - Plan/Benefit Exclusion",
					}, {
						"code":    "MR",
						"message": "Product Not On Formulary",
					},
				},
			},
		},
		"patient_token": copayCallback.PatientToken,
		"metadata":      copayCallback.Metadata,
		"message":       copayCallback.Message,
	}

	callbackResponse := &types.CallbackRequest{
		RequestID:    RequestID,
		Status:       "success",
		Timestamp:    float64(time.Now().Unix()),
		CallbackType: "COPAY",
		Details:      copayRxMap,
	}

	webHookURL := "http://localhost:8888/api/v1/truepill/callback"

	callbackJson := common.ToJSON(callbackResponse)
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(callbackJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func GetClaimsDetails(ctx *gin.Context) {
	copayRxMap := []map[string]interface{}{
		{
			"patient_copay_amount": 2.42,
			"transaction_code":     "B1",
			"transaction_status":   "PAID",
			"claim_type":           "BILLING",
			"transmission_date":    "2023-07-16T12:41:30.000Z",
			"insurance_type":       "PRIMARY",
			"insurance": map[string]interface{}{
				"group_number":  "RX881A",
				"bin":           "17142",
				"pcn":           nil,
				"cardholder_id": "ABC123456",
			},
			"payor":                     "Health Partners",
			"payor_adjudication_amount": 2.5,
		},
		{
			"patient_copay_amount": 2.42,
			"transaction_code":     "B2",
			"transaction_status":   "APPROVED",
			"claim_type":           "REVERSAL",
			"transmission_date":    "2023-07-16T12:41:30.000Z",
			"insurance_type":       "PRIMARY",
			"insurance": map[string]interface{}{
				"group":         "RX881A",
				"bin":           "17142",
				"pcn":           nil,
				"cardholder_id": "ABC123456",
			},
			"payor":                     "Health Partners",
			"payor_adjudication_amount": 2.5,
		},
	}

	ctx.JSON(http.StatusCreated, copayRxMap)
}

func GetFillRequest(ctx *gin.Context) {
	insuranceResponse := &types.FillResponse{
		RequestID: RequestID,
		Status:    "success",
		TimeStamp: float64(time.Now().Unix()),
		Details: &types.FillDetails{
			Message: "Your fill request has been processed successfully.",
		},
	}

	ctx.JSON(http.StatusCreated, insuranceResponse)
	go func() {
		SendFillWebhook(ctx)
		time.Sleep(5 * time.Second)
		SendShippedWebhook(ctx)
	}()
}

func SendFillWebhook(ctx *gin.Context) {
	fillResp := &types.Order{
		RequestID:       RequestID,
		Status:          "success",
		Timestamp:       1571251314,
		RxID:            "",
		PhilOrderNumber: "",
		PhilOrderID:     "",
		Details: &types.FillDetail{
			Metadata:   TransferPrescriptionToken,
			Message:    "Your fill request was processed and is pending shipment.",
			DateFilled: "Mon, 01 Jun 2020 17:45:06 GMT",
			OrderToken: OrderToken,
			Medications: []*types.FillDetailMedication{
				{
					MedicationName:          "Metronidazole 500 mg oral tablet",
					DispensedMedicationName: "Metronidazole 500 mg oral tablet",
					RequestedMedicationName: "Metronidazole 500 mg oral tablet",
					DaysSupply:              func() *float64 { v := float64(7); return &v }(),
					Quantity:                14.0,
					FillNumber:              "0",
					RxNumber:                "1484947",
					TotalRefillsAllowed:     3.0,
					PrescriptionToken:       PrescriptionToken,
					MedicationToken:         MedicationToken,
					RemainingRefills: &types.RemainingRefills{
						TotalRemainingRefills:  0.0,
						TotalQuantityRemaining: 0.0,
					},
					PatientCopayAmount: 0.0,
					PatientCashAmount:  0.0,
				},
			},
			TrackingURL:  "https://tools.usps.com/go/TrackConfirmAction_input?origTrackNum=92001902453595000012688153",
			ErrorCode:    "",
			Description:  "",
			PatientToken: "4526d90a",
		},
	}

	webHookURL := "http://localhost:8888/api/v1/truepill/callback"

	callbackJson := common.ToJSON(fillResp)
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(callbackJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func SendShippedWebhook(ctx *gin.Context) {

	copayRxMap := map[string]interface{}{
		"metadata":        TransferPrescriptionToken,
		"status":          "DELIVERED",
		"message":         "Your shipment has been delivered at the destination mailbox.",
		"eta":             "2020-06-02T00:54:31.838Z",
		"tracking_number": "43904456187100000000000000",
		"tracking_url":    "https://tools.usps.com/go/TrackConfirmAction_input?origTrackNum=43904456187100000000000000",
		"carrier":         "usps",
	}

	callbackResponse := &types.CallbackRequest{
		RequestID:    RequestID,
		Status:       "success",
		Timestamp:    float64(time.Now().Unix()),
		CallbackType: "SHIPMENT",
		Details:      copayRxMap,
	}

	webHookURL := "http://localhost:8888/api/v1/truepill/callback"

	callbackJson := common.ToJSON(callbackResponse)
	resp, err := http.Post(webHookURL, "application/json", bytes.NewBuffer(callbackJson))
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
