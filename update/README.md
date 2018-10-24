# Functions interacting with OCI Compute 

## Pre-requisites

Clone this repo

### Switch to correct context

- `fn use context <your context name>`
- Check using `fn ls apps`

### Create app

`fn create app --annotation oracle.com/oci/subnetIds=<SUBNETS> --config TENANT_OCID=<TENANT_OCID> --config USER_OCID=<USER_OCID> --config REGION=<REGION> --config FINGERPRINT=<FINGERPRINT> --config PRIVATE_KEY_NAME=oci_private_key.pem --config PASSPHRASE=<PASSPHRASE> --syslog-url=<SYSLOG_ENDPOINT> fn-compute-app`

> `--syslog-url` is optional. Use your own!

e.g.

`fn create app --annotation oracle.com/oci/subnetIds='["ocid1.subnet.oc1.phx.aaaaaaaaghmsma7mpqhqdhbgnby25u2zo4wqlrrcskvu7jg56dryxt3hgvka"]' --config TENANT_OCID=ocid1.tenancy.oc1..aaaaaaaaydrjm77otncda2xn7qtv7l3hqnd3zxn2u6siwdhniibwfv4wwhta --config USER_OCID=ocid1.user.oc1..aaaaaaaa4seqx6jeyma46ldy4cbuv35q4l26scz5p4rkz3rauuoioo26qwmq --config REGION=us-phoenix-1 --config FINGERPRINT=41:82:5f:44:ca:a1:2e:58:d2:63:6a:af:52:d5:3d:04 --config PRIVATE_KEY_NAME=oci_private_key.pem --config PASSPHRASE=1987 --syslog-url=tcp://s3cr3t.papertrailapp.com:4242 fn-compute-app`

**Check**

`fn inspect app fn-compute-app`

## LIST function

**Before you proceed...**

`cd list` and copy your private key file e.g. `cp /home/me/oci_private_key.pem .` 

**Deploy**

 `fn -v deploy --app fn-compute-app`

> Go to OCIR and make sure your repo is converted to PUBLIC

**Run**

List all instances in a specific OCI compartment. `CompartmentIDFilter` - OCID of the compartment in which your compute instance exists

> Replace value for `CompartmentIDFilter` as per your environment

`echo -n '{"CompartmentIDFilter":"ocid1.compartment.oc1..aaaaaaaaokbzj2jn3hf5kwdwqoxl2dq7u54p3tsmxrjd7s3uu7x23tkegiua"}' | DEBUG=1 fn invoke fn-compute-app list`

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
