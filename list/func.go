package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(ociComputeEventHandler))
}

const privateKeyFolder string = "/function"

func ociComputeEventHandler(ctx context.Context, in io.Reader, out io.Writer) {

	fnCtx := fdk.GetContext(ctx)

	tenancy := fnCtx.Config()["TENANT_OCID"]
	user := fnCtx.Config()["USER_OCID"]
	region := fnCtx.Config()["REGION"]
	fingerprint := fnCtx.Config()["FINGERPRINT"]
	privateKeyName := fnCtx.Config()["PRIVATE_KEY_NAME"]
	privateKeyLocation := privateKeyFolder + "/" + privateKeyName
	passphrase := fnCtx.Config()["PASSPHRASE"]

	log.Println("TENANT_OCID ", tenancy)
	log.Println("USER_OCID ", user)
	log.Println("REGION ", region)
	log.Println("FINGERPRINT ", fingerprint)
	log.Println("PRIVATE_KEY_NAME ", privateKeyName)
	log.Println("PRIVATE_KEY_LOCATION ", privateKeyLocation)

	privateKey, err := ioutil.ReadFile(privateKeyLocation)
	if err == nil {
		log.Println("read private key from ", privateKeyLocation)
	} else {
		//log.Println("unable read private key from " + privateKeyLocation + " due to " + err.Error())
		resp := FailedResponse{Message: "Unable to read private Key", Error: err.Error()}
		log.Println(resp.toString())
		json.NewEncoder(out).Encode(resp)
		return
	}

	rawConfigProvider := common.NewRawConfigurationProvider(tenancy, user, region, fingerprint, string(privateKey), common.String(passphrase))
	cc, err := core.NewComputeClientWithConfigurationProvider(rawConfigProvider)

	if err != nil {
		resp := FailedResponse{Message: "Problem getting Compute Client handle", Error: err.Error()}
		log.Println(resp.toString())
		json.NewEncoder(out).Encode(resp)
		return
	}

	var criteria SearchCriteria
	json.NewDecoder(in).Decode(&criteria)
	log.Println("Criteria ", criteria)
	cmpt := criteria.CompartmentIDFilter
	//name := criteria.DisplayNameFilter

	instances := []Instance{}

	//resp, err := cc.ListInstances(context.Background(), core.ListInstancesRequest{CompartmentId: common.String(cmpt), DisplayName: common.String(name)})
	resp, err := cc.ListInstances(context.Background(), core.ListInstancesRequest{CompartmentId: common.String(cmpt)})
	if err != nil {
		resp := FailedResponse{Message: "Problem listing instances", Error: err.Error()}
		log.Println(resp.toString())
		json.NewEncoder(out).Encode(resp)
		return
	}
	numInstances := len(resp.Items)
	if numInstances == 0 {
		log.Println("No instances for search criteria ", criteria)
		json.NewEncoder(out).Encode(instances)
		return
	}

	log.Println("Found instances ", numInstances)

	for _, instance := range resp.Items {
		ins := Instance{OCID: *instance.Id, DisplayName: *instance.DisplayName}
		log.Println("Added instance ", ins)
		instances = append(instances, ins)
	}
	json.NewEncoder(out).Encode(instances)
}

//Instance ...
type Instance struct {
	OCID        string
	DisplayName string
}

//SearchCriteria ...
type SearchCriteria struct {
	CompartmentIDFilter string
}

//FailedResponse ...
type FailedResponse struct {
	Message string
	Error   string
}

func (response FailedResponse) toString() string {
	return response.Message + " due to " + response.Error
}
