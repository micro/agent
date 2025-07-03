package handler

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/micro/agent/proto"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/errors"
	"go-micro.dev/v5/genai"
	"go-micro.dev/v5/registry"
)

type Agent struct{}

type Call struct {
	Service  string      `json:"service"`
	Endpoint string      `json:"endpoint"`
	Request  interface{} `json:"request"`
	Error    string      `json:"error"`
}

type Result struct {
	Answer interface{} `json:"answer"`
	Error  string      `json:"error"`
}

func New() *Agent {
	return new(Agent)
}

func (a *Agent) Query(ctx context.Context, req *pb.QueryRequest, rsp *pb.QueryResponse) error {
	fmt.Println("agent.query", req.Question)
	resp, err := genai.DefaultGenAI.Generate(req.Question)
	if err != nil {
		return err
	}
	rsp.Answer = resp.Text

	return nil
}

func (a *Agent) Command(ctx context.Context, req *pb.CommandRequest, rsp *pb.CommandResponse) error {
	// get list of services
	serviceList, err := registry.ListServices()
	if err != nil {
		return errors.InternalServerError("agent.registry", err.Error())
	}

	// get all the endpoints per service
	var services []*registry.Service
	for _, service := range serviceList {
		srv, err := registry.GetService(service.Name)
		if err != nil {
			return errors.InternalServerError("agent.registry", err.Error())
		}
		services = append(services, srv...)
	}

	b, _ := json.Marshal(services)

	prompt := `The user is requesting a certain action. You have a list of services
	which will enable you to perform that action. Based on the list of services,
	endpoints and request/response, determine which services or set of services 
	must be called to satisfy the user request. The format of your response 
	should be JSON with fields service, endpoint and request. Where service 
	is the name of the service to call. The endpoint is the name of the endpoint
	and the request is a JSON formatted request that can be passed in to satisfy 
	the service/endpoint request format and fields required, including the data 
	to get the expected response. 

	Note: the service request/response has field "values" for the name of fields. 
	When you provide your request, ensure your JSON for your request is only the 
	field names in the values of the request or response and then the type of value 
	you provide is based on the "type" of value it defines e.g if it's the field 
	"name" with type "string" then you provide {"name": "john"} as your request. 
	The same goes for the response format. The "service" and "endpoint" fields 
	you provide should be from the top level "name" field and the endpoint "name".

	Here is the service list:

	%s

	Here is the user request:

	%s

	Only respond with the JSON format for service/endpoint/request so that we
	can parse out the JSON. In the event of an error or problem define an 
	error field e.g if nothing satisfies the request`

	text := fmt.Sprintf(prompt, string(b), req.Request)

	resp, err := genai.DefaultGenAI.Generate(text)
	if err != nil {
		return errors.InternalServerError("agent.generate", err.Error())
	}

	var call Call
	json.Unmarshal([]byte(resp.Text), &call)

	fmt.Println("agent.call\n", resp.Text, call)

	if len(call.Error) > 0 {
		rsp.Error = call.Error
		return nil
	}

	// TODO: validate is a service in the list
	if len(call.Service) == 0 || len(call.Endpoint) == 0 {
		return errors.InternalServerError("agent.call", "unable to satisfy request")
	}

	// call a service
	request := client.NewRequest(call.Service, call.Endpoint, call.Request)
	var respB json.RawMessage
	if err := client.Call(ctx, request, &respB); err != nil {
		rsp.Error = err.Error()
		return nil
	}

	fmt.Println("agent.call-response\n", string(respB))

	// feed response into LLM
	prompt = `Here is the response from our service call 
	to service: %s endpoint %s. Format the response as 
	per the user's request, in the event no format is
	specified, parse the request JSON into a string. This 
	may require interpreting the response. Whatever you provide
	as output we will send to the user. In the event you need 
	to respond with an error specify an error field. Return 
	your response to me as JSON with fields answer for what 
	will be sent to the user and error in the case of errors.

	Request:

	%s

	Response:

	%s
	`
	text = fmt.Sprintf(prompt, call.Service, call.Endpoint, call.Request, string(respB))

	resp, err = genai.DefaultGenAI.Generate(text)
	if err != nil {
		return errors.InternalServerError("agent.generate", err.Error())
	}

	var res Result
	json.Unmarshal([]byte(resp.Text), &res)

	fmt.Println("agent.result", res, resp.Text)

	if len(res.Error) > 0 {
		rsp.Error = err.Error()
		return nil
	}

	// got the answer
	rsp.Response = fmt.Sprintf("%v", res.Answer)

	return nil
}
