package responses

import "encoding/json"

var OK = New("OK", "200 OK", "Avito_Segment_Service-000200")
var Created = New("Content created", "201 Created", "Avito_Segment_Service-000201")
var Deleted = New("Content deleted", "204 No Content", "Avito_Segment_Service-000204")
var Updated = New("Content updated", "200 OK", "Avito_Segment_Service-000200")

type ApiResponse struct {
	Description  string `json:"description"`
	DeveloperMsg string `json:"developer_msg"`
	Code         string `json:"code"`
}

func (r *ApiResponse) Marshal() []byte {
	bytes, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	return bytes
}

func New(description, developerMsg, code string) *ApiResponse {
	return &ApiResponse{
		Description:  description,
		DeveloperMsg: developerMsg,
		Code:         code,
	}
}
