package dialogflowmap

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	"strconv"
)

// DialogFlowProcessor action information for connecting with DialogFlow
type DialogFlowProcessor struct {
	projectID        string
	authJSONFilePath string
	lang             string
	timeZone         string
	sessionClient    *dialogflow.SessionsClient
	ctx              context.Context
}

// NPLResponse action struct for response Diagnostic info
type NPLResponse struct {
	FulfillmentText string            `json:"FulfillmentText"`
	Intent          string            `json:"intent"`
	Confidence      float32           `json:"confidence"`
	Entities        map[string]string `json:"entities"`
}

var dp DialogFlowProcessor

func (dp *DialogFlowProcessor) Init(arr ...string) (*DialogFlowProcessor, error) {
	dp.projectID = arr[0]
	dp.authJSONFilePath = arr[1]
	dp.lang = arr[2]
	dp.timeZone = arr[3]

	dp.ctx = context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if err != nil {
		log.Fatal("Error in auth with DialogFlow")
	}
	dp.sessionClient = sessionClient

	return dp, nil
}

func (dp *DialogFlowProcessor) ProcessNPL(rawMessage, username string) (r NPLResponse) {
	sessionID := username
	req := dialogflowpb.DetectIntentRequest{
		Session: fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, sessionID),
		QueryParams: &dialogflowpb.QueryParameters{
			TimeZone: dp.timeZone,
		},
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         rawMessage,
					LanguageCode: dp.lang,
				},
			},
		},
	}

	log.Println(dp.ctx)
	response, err := dp.sessionClient.DetectIntent(dp.ctx, &req)
	if err != nil {
		log.Fatalf("Error comunication with DialogFlow %s", err.Error())
		return
	}

	queryResult := response.GetQueryResult()
	if queryResult.Intent != nil {
		r.Intent = queryResult.Intent.DisplayName
		r.Confidence = float32(queryResult.IntentDetectionConfidence)
	}

	r.FulfillmentText = queryResult.FulfillmentText
	r.Entities = make(map[string]string)
	// *structpb.Value
	if r.Intent == "GetDateWithTime" {
		outputContexts := queryResult.OutputContexts

		if len(outputContexts) > 0 {
			for _, value := range outputContexts {
				params := value.Parameters.GetFields()
				for name, p := range params {
					extractValue := extractDialogFlowEntities(p)
					r.Entities[name] = extractValue
				}
			}
		}
	}

	params := queryResult.Parameters.GetFields()
	if len(params) > 0 {
		for name, p := range params {
			extractValue := extractDialogFlowEntities(p)
			r.Entities[name] = extractValue
		}
	}
	return
}

func extractDialogFlowEntities(p *structpb.Value) (extractedEntity string) {
	kind := p.GetKind()
	switch kind.(type) {
	case *structpb.Value_StringValue:
		return p.GetStringValue()
	case *structpb.Value_NumberValue:
		return strconv.FormatFloat(p.GetNumberValue(), 'f', 6, 64)
	case *structpb.Value_BoolValue:
		return strconv.FormatBool(p.GetBoolValue())
	case *structpb.Value_StructValue:
		sValue := p.GetStructValue()
		fields := sValue.GetFields()
		extractedEntity := ""
		for k, v := range fields {
			if k == "amount" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, strconv.FormatFloat(v.GetNumberValue(), 'f', 6, 64))
			}
			if k == "unit" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, v.GetStringValue())
			}
			if k == "date_time" {
				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, v.GetStringValue())
			}
		}
		return extractedEntity

	default:
		return ""
	}
}
