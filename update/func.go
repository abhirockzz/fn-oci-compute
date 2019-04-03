package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	fdk "github.com/fnproject/fdk-go"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(ociComputeEventHandler))
}

const privateKeyFolder string = "/function"
const successMsg string = "Updated Compute Instance information successfully"

func ociComputeEventHandler(ctx context.Context, in io.Reader, out io.Writer) {

	tenancy := os.Getenv("TENANT_OCID")
	user := os.Getenv("USER_OCID")
	region := os.Getenv("REGION")
	fingerprint := os.Getenv("FINGERPRINT")
	privateKeyName := os.Getenv("OCI_PRIVATE_KEY_FILE_NAME")
	privateKeyLocation := privateKeyFolder + "/" + privateKeyName
	passphrase := os.Getenv("PASSPHRASE")

	log.Println("TENANT_OCID ", tenancy)
	log.Println("USER_OCID ", user)
	log.Println("REGION ", region)
	log.Println("FINGERPRINT ", fingerprint)
	log.Println("OCI_PRIVATE_KEY_FILE_NAME ", privateKeyName)
	log.Println("PRIVATE_KEY_LOCATION ", privateKeyLocation)

	privateKey, err := ioutil.ReadFile(privateKeyLocation)
	if err == nil {
		log.Println("read private key from ", privateKeyLocation)
	} else {
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

	var updateInfo UpdateInfo
	json.NewDecoder(in).Decode(&updateInfo)
	log.Println("UpdateInfo ", updateInfo)

	_, updateErr := cc.UpdateInstance(context.Background(), core.UpdateInstanceRequest{InstanceId: common.String(updateInfo.OCID), UpdateInstanceDetails: core.UpdateInstanceDetails{DisplayName: common.String(updateInfo.NewDisplayName)}})
	if updateErr != nil {
		resp := FailedResponse{Message: "Problem updating instance", Error: updateErr.Error()}
		log.Println(resp.toString())
		json.NewEncoder(out).Encode(resp)
		return
	}

	log.Println(successMsg)

	out.Write([]byte(successMsg))
}

//UpdateInfo ...
type UpdateInfo struct {
	OCID           string
	NewDisplayName string
}

//FailedResponse ...
type FailedResponse struct {
	Message string
	Error   string
}

func (response FailedResponse) toString() string {
	return response.Message + " due to " + response.Error
}
