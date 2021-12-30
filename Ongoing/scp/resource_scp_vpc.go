package scp

import (
	"fmt"
	"log"
	"time"

	"github.com/jhp0204/Provider_test1/scp-clilib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func init() {
	RegisterResource("scp_vpc", resourceScpVpc())
}

func resourceScpVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceScpVpcCreate,
		Read:   resourceScpVpcRead,
		Delete: resourceScpVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// API에 필요6한 내용들 확인, 필요 시 수정 필요 
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validateInstanceName),
				Description:      "Subnet name to create. default: Assigned by SAMSUNG CLOUD PLATFORM.",
			},
			"ipv4_cidr_block": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: ToDiagFunc(validation.IsCIDRNetwork(16, 28)),
				Description:      "The CIDR block for the vpc.",
			},
			"vpc_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_network_acl_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_access_control_group_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_public_route_table_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_private_route_table_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceScpVpcCreate(d *schema.ResourceData, meta interface{}) error {
	fmt.println("resourceScpVpcCreate function is begin")
	config := meta.(*ProviderConfig)

	if !config.SupportVPC {
		return NotSupportClassic("resource `scp_vpc`")
	}

	reqParams := &vpc.CreateVpcRequest{
		RegionCode:    &config.RegionCode,
		Ipv4CidrBlock: scp-clilib.String(d.Get("ipv4_cidr_block").(string)),
	}

	if v, ok := d.GetOk("name"); ok {
		reqParams.VpcName = scp-clilib.String(v.(string))
	}

//	logCommonRequest("CreateVpc", reqParams)
	resp, err := config.Client.vpc.V2Api.CreateVpc(reqParams)
//	if err != nil {
//		logErrorResponse("Create Vpc Instance", err, reqParams)
//		return err
//	}

//	logCommonResponse("CreateVpc", GetCommonResponse(resp))

	vpcInstance := resp.VpcList[0]
	d.SetId(*vpcInstance.VpcNo)
	log.Printf("[INFO] VPC ID: %s", d.Id())

//	if err := waitForScpVpcCreation(config, d.Id()); err != nil {
//		return err
//	}

	return resourceScpVpcRead(d, meta)
}

func resourceScpVpcRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)
	instance, err := (config, d.Id())
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(*instance.VpcNo)
	d.Set("vpc_no", instance.VpcNo)
	d.Set("name", instance.VpcName)
	d.Set("ipv4_cidr_block", instance.Ipv4CidrBlock)
/*
	if *instance.VpcStatus.Code != "TERMTING" {
		defaultNetworkACLNo, err := getDefaultNetworkACL(config, d.Id())
		if err != nil {
			return fmt.Errorf("error get default network acl for VPC (%s): %s", d.Id(), err)
		}

		d.Set("default_network_acl_no", defaultNetworkACLNo)

		defaultAcgNo, err := getDefaultAccessControlGroup(config, d.Id())
		if err != nil {
			return fmt.Errorf("error get default Access Control Group for VPC (%s): %s", d.Id(), err)
		}
		d.Set("default_access_control_group_no", defaultAcgNo)

		publicRouteTableNo, privateRouteTableNo, err := getDefaultRouteTable(config, d.Id())
		if err != nil {
			return fmt.Errorf("error get default Route Table for VPC (%s): %s", d.Id(), err)
		}
		d.Set("default_public_route_table_no", publicRouteTableNo)
		d.Set("default_private_route_table_no", privateRouteTableNo)
	}
*/
	return nil
}
/*
func getDefaultNetworkACL(config *ProviderConfig, id string) (string, error) {
	reqParams := &vpc.GetNetworkAclListRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      scp.String(id),
	}

	logCommonRequest("GetNetworkAclList", reqParams)
	resp, err := config.Client.vpc.V2Api.GetNetworkAclList(reqParams)

	if err != nil {
		logErrorResponse("GetNetworkAclList", err, reqParams)
		return "", err
	}

	logResponse("GetNetworkAclList", resp)

	if resp == nil || len(resp.NetworkAclList) == 0 {
		return "", fmt.Errorf("no matching Network ACL found")
	}

	for _, i := range resp.NetworkAclList {
		if *i.IsDefault {
			return *i.NetworkAclNo, nil
		}
	}

	return "", fmt.Errorf("No matching default network ACL found")
}

//SCP VPC 생성 시, RouteTable 필요여부 확인 
func getDefaultRouteTable(config *ProviderConfig, id string) (publicRouteTableNo string, privateRouteTableNo string, error error) {
	reqParams := &vpc.GetRouteTableListRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      scp.String(id),
	}

	logCommonRequest("getDefaultRouteTable", reqParams)
	resp, err := config.Client.vpc.V2Api.GetRouteTableList(reqParams)

	if err != nil {
		logErrorResponse("getDefaultRouteTable", err, reqParams)
		return "", "", err
	}

	logResponse("getDefaultRouteTable", resp)

	for _, i := range resp.RouteTableList {
		if *i.IsDefault && *i.SupportedSubnetType.Code == "PRIVATE" {
			privateRouteTableNo = *i.RouteTableNo
		} else if *i.IsDefault && *i.SupportedSubnetType.Code == "PUBLIC" {
			publicRouteTableNo = *i.RouteTableNo
		}
	}

	return publicRouteTableNo, privateRouteTableNo, nil
}

func resourceScpVpcDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*ProviderConfig)

	reqParams := &vpc.DeleteVpcRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      scp.String(d.Get("vpc_no").(string)),
	}

	logCommonRequest("DeleteVpc", reqParams)
	resp, err := config.Client.vpc.V2Api.DeleteVpc(reqParams)
	if err != nil {
		logErrorResponse("DeleteVpc Vpc Instance", err, reqParams)
		return err
	}
	logResponse("DeleteVpc", resp)

	if err := waitForNcloudVpcDeletion(config, d.Id()); err != nil {
		return err
	}

	return nil
}

// 해당 func 사용상황 및 필요여부 확인 (Creation func 만으로 안되는 상황이 언제인지?)
func waitForScpVpcCreation(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INIT", "CREATING"},
		Target:  []string{"RUN"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getVpcInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcStatus")
		},
		Timeout:    DefaultCreateTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC (%s) to become available: %s", id, err)
	}

	return nil
}

func waitForScpVpcDeletion(config *ProviderConfig, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"RUN", "TERMTING"},
		Target:  []string{"TERMINATED"},
		Refresh: func() (interface{}, string, error) {
			instance, err := getVpcInstance(config, id)
			return VpcCommonStateRefreshFunc(instance, err, "VpcStatus")
		},
		Timeout:    DefaultTimeout,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for VPC (%s) to become termintaing: %s", id, err)
	}

	return nil
}

func getVpcInstance(config *ProviderConfig, id string) (*vpc.Vpc, error) {
	reqParams := &vpc.GetVpcDetailRequest{
		RegionCode: &config.RegionCode,
		VpcNo:      scp.String(id),
	}

	resp, err := config.Client.vpc.V2Api.GetVpcDetail(reqParams)
	if err != nil {
		logErrorResponse("Get Vpc Instance", err, reqParams)
		return nil, err
	}
	logResponse("GetVpcDetail", resp)

	if len(resp.VpcList) > 0 {
		vpc := resp.VpcList[0]
		return vpc, nil
	}

	return nil, nil
}
*/
