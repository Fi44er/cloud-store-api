package auth_utils

import (
	"encoding/json"
	"regexp"
	"strings"

	kratos "github.com/ory/kratos-client-go"
)

var reValueInQuotes = regexp.MustCompile(`^".*?"\s+`)

type kratosFlowWithUI struct {
	Ui struct {
		Messages []kratos.UiText `json:"messages"`
		Nodes    []kratos.UiNode `json:"nodes"`
	} `json:"ui"`
}

func ExtractKratosErrors(err error) map[string]string {
	errorsMap := make(map[string]string)
	if err == nil {
		return errorsMap
	}

	openApiErr, ok := err.(*kratos.GenericOpenAPIError)
	if !ok {
		errorsMap["generic"] = err.Error()
		return errorsMap
	}

	var flow kratosFlowWithUI
	if unmarshalErr := json.Unmarshal(openApiErr.Body(), &flow); unmarshalErr != nil {
		errorsMap["generic"] = err.Error()
		return errorsMap
	}

	if len(flow.Ui.Messages) > 0 {
		errorsMap["generic"] = joinKratosMessages(flow.Ui.Messages)
	}

	for _, node := range flow.Ui.Nodes {
		if len(node.Messages) == 0 {
			continue
		}

		if attr := node.Attributes.UiNodeInputAttributes; attr != nil {
			fieldName := strings.TrimPrefix(attr.Name, "traits.")
			errorsMap[fieldName] = joinKratosMessages(node.Messages)
		}
	}

	return errorsMap
}

func joinKratosMessages(messages []kratos.UiText) string {
	var msgs []string
	for _, m := range messages {
		msgs = append(msgs, cleanKratosMessage(m.Text))
	}
	return strings.Join(msgs, "; ")
}

func cleanKratosMessage(msg string) string {
	msg = reValueInQuotes.ReplaceAllString(msg, "")
	msg = strings.ReplaceAll(msg, "\"", "")

	replacements := map[string]string{
		"is not valid email": "invalid email format",
		"Property ":          "",
		"is missing":         "is required",
		"is shorter than":    "is too short",
	}

	for old, new := range replacements {
		msg = strings.ReplaceAll(msg, old, new)
	}

	msg = strings.TrimSpace(msg)
	if len(msg) > 0 {
		return strings.ToUpper(msg[:1]) + msg[1:]
	}
	return msg
}
