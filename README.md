# Functions with OCI Compute

Fn functions to list and update OCI Compute instances. These are written in [Go](https://golang.org) and use the [OCI Go SDK](https://github.com/oracle/oci-go-sdk) to interact with Compute service

**Your OCI API private key is required for authentication by the SDK** and it'll be embedded in the function (Docker) container (in `/function` directory). This is done using a `Dockerfile` for building the function which does exactly the same as **Fn** would do out-of-the-box except that it copies over the user provided OCI API private key in the second stage of the function build process

## Pre-requisites

Clone this repo

### Switch to correct context

- `fn use context <your context name>`
- Check using `fn ls apps`

### Create app

`fn create app --annotation oracle.com/oci/subnetIds=<SUBNETS> --config TENANT_OCID=<TENANT_OCID> --config USER_OCID=<USER_OCID> --config REGION=<REGION> --config FINGERPRINT=<FINGERPRINT> --config PRIVATE_KEY_NAME=<NAME_OF_YOUR_OCI_API_PRIVATE_KEY> --config PASSPHRASE=<PASSPHRASE> --syslog-url=<SYSLOG_ENDPOINT> fn-compute-app`

> `--syslog-url` is optional

e.g.

`fn create app --annotation oracle.com/oci/subnetIds='["ocid1.subnet.oc1.phx.aaaaaaaaghmsma7mpqhqdhbgnby25u2zo4wqlrrcskvu7jg56dryxt3hgvka"]' --config TENANT_OCID=ocid1.tenancy.oc1..aaaaaaaaydrjm77otncda2xn7qtv7l3hqnd3zxn2u6siwdhniibwfv4wwhta --config USER_OCID=ocid1.user.oc1..aaaaaaaa4seqx6jeyma46ldy4cbuv35q4l26scz5p4rkz3rauuoioo26qwmq --config REGION=us-phoenix-1 --config FINGERPRINT=41:82:5f:44:ca:a1:2e:58:d2:63:6a:af:52:d5:3d:04 --config PRIVATE_KEY_NAME=oci_private_key.pem --config PASSPHRASE=1987 --syslog-url=tcp://s3cr3t.papertrailapp.com:4242 fn-compute-app`

**Check**

`fn inspect app fn-compute-app`

## LIST function

**Before you proceed...**

`cd fn-oci-compute/list` and copy your private key file e.g. `cp /home/me/oci_private_key.pem .` 

**Deploy**

 `fn -v deploy --app fn-compute-app`

> Go to OCIR and make sure your repo is converted to PUBLIC

**Run**

List all instances in a specific OCI compartment. `CompartmentIDFilter` - OCID of the compartment in which your compute instance exists

> Replace value for `CompartmentIDFilter` as per your environment

`echo -n '{"CompartmentIDFilter":"ocid1.compartment.oc1..aaaaaaaaxmrampww6livwdw3usqxlrmn5fiwi3dbkwtl3waigzbwl5olu5pa"}' | DEBUG=1 fn invoke fn-compute-app list`

If successful, you'll get a JSON response similar to below

	[
	    {
	        "OCID": "ocid1.instance.oc1.iad.abuwcljrljpitmsa5qqqamotfmnrte6jazfkx2ovbikadgpx7bdgei4t27ia",
	        "DisplayName": "my-instance-1"
	    },
	    {
	        "OCID": "ocid1.instance.oc1.iad.abuwcljrdzmdqypsumg2iw7l5ebxpdjt5fbcf2kdu6msffmdnss2mfxmi6qq",
	        "DisplayName": "my-instance-2"
	    },
	    {
	        "OCID": "ocid1.instance.oc1.iad.abuwcljrcshge4wbbwisjt5nbrotpuegsz6wzsn7mg2xac5ci2gl7g3esffa",
	        "DisplayName": "my-instance-3"
	    }
	]

In case you specify an incorrect compartment ID, you might see

	{
	    "Message": "Problem listing instances",
	    "Error": "Service error:InvalidParameter. CompartmentId must be specified as a valid OCID. http status code: 400. Opc request id: 85a816e1db955c0a4869cf7b6a4d2231/24B49E46B2350CE2B04657A1AD9791BE/F8FED5D9F3F964E10B1E9C6E65FC6EEC"
	}

## UPDATE function

Now that you have the list of instances (along with their OCID), try updating the **display name** of one such instance

**Before you proceed...**

`cd fn-oci-compute/update` and copy your private key file e.g. `cp /home/me/oci_private_key.pem .` 

**Deploy**

 `fn -v deploy --app fn-compute-app`

> Go to OCIR and make sure your repo is converted to PUBLIC

**Run**

> Replace value for `OCID` as per your compute instance and use an appropriate value for `NewDisplayName`

`echo -n '{"OCID":"ocid1.instance.oc1.iad.abuwcljrljpitmsa5qqqamotfmnrte6jazfkx2ovbikadgpx7bdgei4t27ia", "NewDisplayName":"test-name"}' | DEBUG=1 fn invoke fn-compute-app update`

If successful, you'll get a `Updated Compute Instance information successfully` response. Now,

- you can log into your OCI console to check the same `Menu > Compute > Instances` (choose the correct compartment), or
- run the `LIST` function (which you previously deployed) to verify the same i.e. `echo -n '{"CompartmentIDFilter":"ocid1.compartment.oc1..aaaaaaaaxmrampww6livwdw3usqxlrmn5fiwi3dbkwtl3waigzbwl5olu5pa"}' | DEBUG=1 fn invoke fn-compute-app list`

In case you specify an incorrect OCID, you might see

	{
	    "Message": "Problem updating instance",
	    "Error": "Service error:InvalidParameter. AvailabilityDomain could not be inferred from the Request. http status code: 400. Opc request id: 598fc3eb768794557da1db50f8e14f69/8EF45B8946CE485E8776B02565F115E7/44928564719A83DFDD8C06DB14BCA86A"
	}

.. or an invalid (e.g. empty) display name will result in 

	{
	    "Message": "Problem updating instance",
	    "Error": "Service error:InvalidParameter. displayName size must be between 1 and 255. http status code: 400. Opc request id: c9aedd02d36f63226aee7e29540aec5c/C61A547B97A02B26E58143586DBC5169/9CEF05C1640A93BE640CE9EF4C5A781E"
	}
