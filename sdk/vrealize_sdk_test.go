package sdk

import (
	"fmt"
	"testing"

	"github.com/vmware/terraform-provider-vra7/utils"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var (
	client   APIClient
	user     = "admin@myvra.local"
	password = "pass!@#"
	tenant   = "vsphere.local"
	baseURL  = "http://localhost"
	insecure = true
)

func init() {
	fmt.Println("init")
	client = NewClient(user, password, tenant, baseURL, insecure)
}

func TestGetCatalogItemRequestTemplate(t *testing.T) {
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	requestTemplateResponse := `{"type":"com.vmware.vcac.catalog.domain.request.CatalogItemProvisioningRequest","catalogItemId":"feaedf73-560c-4612-a573-41667e017691","requestedFor":"fritz@coke.sqa-horizon.local","businessGroupId":"b2470b94-cbca-43db-be37-803cca7b0f1a","description":null,"reasons":null,"data":{"_leaseDays":null,"_number_of_instances":1,"machine2.vsphere":{"componentTypeId":"com.vmware.csp.component.cafe.composition","componentId":null,"classId":"Blueprint.Component.Declaration","typeFilter":"Prativa_CentOs*machine2.vsphere","data":{"_cluster":1,"_hasChildren":false,"cpu":1,"datacenter_location":null,"description":null,"disks":[{"componentTypeId":"com.vmware.csp.iaas.blueprint.service","componentId":null,"classId":"Infrastructure.Compute.Machine.MachineDisk","typeFilter":null,"data":{"capacity":8,"custom_properties":null,"id":1537301083687,"initial_location":"","is_clone":true,"label":"Hard disk 1","storage_reservation_policy":"","userCreated":false,"volumeId":0}}],"display_location":false,"guest_customization_specification":null,"location.loc":"","machine_prefix":null,"max_network_adapters":-1,"max_per_user":0,"max_volumes":60,"memory":1024,"nics":null,"os_arch":"x86_64","os_distribution":null,"os_type":"Linux","os_version":null,"ovfAuthNeeded":false,"ovf_proxy_endpoint":null,"ovf_url":null,"ovf_url_pwd":null,"ovf_url_username":null,"ovf_use_proxy":false,"property_groups":null,"reservation_policy":null,"security_groups":[],"security_tags":[],"snapshot_name":null,"source_machine_external_snapshot":null,"source_machine_vmsnapshot":null,"storage":8}}}}`

	catalogItemID := "feaedf73-560c-4612-a573-41667e017691"

	path := fmt.Sprintf(RequestTemplateAPI, catalogItemID)
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, requestTemplateResponse))

	catalogItemReqTemplate, err := client.GetCatalogItemRequestTemplate(catalogItemID)
	utils.AssertNilError(t, err)
	utils.AssertEqualsString(t, catalogItemID, catalogItemReqTemplate.CatalogItemID)

	catalogItemReqTemplate, err = client.GetCatalogItemRequestTemplate("635e5v-8e37efd60-hdgdh")
	utils.AssertNotNilError(t, err)

	requestTemplateErrorResponse := `{"errors":[{"code":20116,"source":null,"message":"Unable to find the specified catalog item in the service catalog: feaedf73-560c-4612-a573-0041667e0176.","systemMessage":"Unable to find the specified catalog item in the service catalog: feaedf73-560c-4612-a573-0041667e0176.","moreInfoUrl":null}]}`
	httpmock.Reset()
	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(20116, requestTemplateErrorResponse))
	invalidCatalogItemID := "feaedf73-560c-4612-a573-0041667e0176"
	catalogItemReqTemplate, err = client.GetCatalogItemRequestTemplate(invalidCatalogItemID)
	utils.AssertNotNilError(t, err)

}

func TestReadCatalogItemNameByID(t *testing.T) {

	catalogItemResp := `{"catalogItem":{"callbacks":null,"catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint",
"label":"Composite Blueprint"},"dateCreated":"2015-12-22T03:16:19.289Z","description":"CentOS 6.3 IaaS Blueprint",
"forms":{"itemDetails":{"type":"external","formId":"composition.catalog.item.details"},"catalogRequestInfoHidden":true,
"requestFormScale":"BIG","requestSubmission":{"type":"extension","extensionId":"com.vmware.vcac.core.design.blueprints.requestForm",
"extensionPointId":null},"requestDetails":{"type":"extension","extensionId":"com.vmware.vcac.core.design.blueprints.requestDetailsForm",
"extensionPointId":null},"requestPreApproval":null,"requestPostApproval":null},"iconId":"e5dd4fba-45ed-4943-b1fc-7f96239286be",
"id":"e5dd4fba-45ed-4943-b1fc-7f96239286be","isNoteworthy":false,"lastUpdatedDate":"2017-01-06T05:12:56.690Z",
"name":"CentOS 6.3","organization":{"tenantRef":"vsphere.local","tenantLabel":"vsphere.local","subtenantRef":null,
"subtenantLabel":null},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"},
"providerBinding":{"bindingId":"vsphere.local!::!CentOS63","providerRef":{"id":"2fbaabc5-3a48-488a-9f2a-a42616345445",
"label":"Blueprint Service"}},"serviceRef":{"id":"baad0ad2-8b96-4347-b188-f534dad53a0d","label":"Infrastructure"},
"status":"PUBLISHED","statusName":"Published","quota":0,"version":4,"requestable":true},"entitledOrganizations":[{"tenantRef":"vsphere.local",
"tenantLabel":"vsphere.local","subtenantRef":"53619006-56bb-4788-9723-9eab79752cc1","subtenantLabel":"Content"}]}`

	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	catalogItemID := "e5dd4fba-45ed-4943-b1fc-7f96239286be"
	path := fmt.Sprintf(EntitledCatalogItems+"/"+"%s", catalogItemID)
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, catalogItemResp))

	catalogItemName, err := client.ReadCatalogItemNameByID(catalogItemID)
	utils.AssertNilError(t, err)
	utils.AssertEqualsString(t, "CentOS 6.3", catalogItemName)

	catalogItemName, err = client.ReadCatalogItemNameByID("84rg=73dv-dd8dhy-hg")
	utils.AssertNotNilError(t, err)
	utils.AssertEqualsString(t, "", catalogItemName)
}

func TestReadCatalogItemByName(t *testing.T) {
	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	entitledCatalogItemViewsResponse := `{"links":[],"content":[{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"qe","tenantLabel":"QETenant","subtenantRef":"b2470b94-cbca-43db-be37-803cca7b0f1a","subtenantLabel":"Development"}],"catalogItemId":"2e13fd45-a85e-4985-b89e-4ebb19ab272c","name":"NoSoftwareMachine","description":null,"isNoteworthy":false,"dateCreated":"2018-08-31T19:30:51.036Z","lastUpdatedDate":"2018-09-18T19:37:50.903Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/2e13fd45-a85e-4985-b89e-4ebb19ab272c/requests/template"},{"@type":"link","rel":"POST: Submit Request","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/2e13fd45-a85e-4985-b89e-4ebb19ab272c/requests"}],"iconId":"composition.blueprint.png","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"d33afbc2-954d-416d-8d4a-bc4ad1b66058","label":"test-service"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}},{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"qe","tenantLabel":"QETenant","subtenantRef":"b2470b94-cbca-43db-be37-803cca7b0f1a","subtenantLabel":"Development"}],"catalogItemId":"b7472041-d49b-44f2-a9b6-8ea714f55c55","name":"Dukes Bank Application","description":"Three-tier Dukes Bank application on Linux nodes; Apache HTTP Server for load balancer, JBoss Server for application server, and MySQL Server for database server. This is the basic version of Dukes Bank on Linux nodes.","isNoteworthy":false,"dateCreated":"2018-08-31T18:12:49.680Z","lastUpdatedDate":"2018-09-18T19:24:09.673Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/b7472041-d49b-44f2-a9b6-8ea714f55c55/requests/template"},{"@type":"link","rel":"POST: Submit Request","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/b7472041-d49b-44f2-a9b6-8ea714f55c55/requests"}],"iconId":"composition.blueprint.png","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"d33afbc2-954d-416d-8d4a-bc4ad1b66058","label":"test-service"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}},{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"qe","tenantLabel":"QETenant","subtenantRef":"b2470b94-cbca-43db-be37-803cca7b0f1a","subtenantLabel":"Development"}],"catalogItemId":"1eb8e1d4-152e-4a93-a3b6-265df6870555","name":"Azure Machine (2)","description":"Creates new Azure virtual machine","isNoteworthy":false,"dateCreated":"2018-09-18T19:40:19.172Z","lastUpdatedDate":"2018-09-18T19:41:39.432Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/1eb8e1d4-152e-4a93-a3b6-265df6870555/requests/template"},{"@type":"link","rel":"POST: Submit Request","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/1eb8e1d4-152e-4a93-a3b6-265df6870555/requests"}],"iconId":"_internal!::!dda3cc3f-f2da-45f5-bc64-a292c1ace3ba_icon","catalogItemTypeRef":{"id":"com.vmware.csp.core.designer.service.serviceblueprint","label":"XaaS Blueprint"},"serviceRef":{"id":"d33afbc2-954d-416d-8d4a-bc4ad1b66058","label":"test-service"},"outputResourceTypeRef":{"id":"_internal!::!dda3cc3f-f2da-45f5-bc64-a292c1ace3ba","label":"Azure Virtual Machine"}},{"@type":"ConsumerEntitledCatalogItemView","entitledOrganizations":[{"tenantRef":"qe","tenantLabel":"QETenant","subtenantRef":"b2470b94-cbca-43db-be37-803cca7b0f1a","subtenantLabel":"Development"}],"catalogItemId":"feaedf73-560c-4612-a573-41667e017691","name":"CentOs","description":"","isNoteworthy":false,"dateCreated":"2018-09-18T20:04:55.805Z","lastUpdatedDate":"2018-09-27T22:52:48.390Z","links":[{"@type":"link","rel":"GET: Request Template","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/feaedf73-560c-4612-a573-41667e017691/requests/template"},{"@type":"link","rel":"POST: Submit Request","href":"https://cava-n-80-072.eng.vmware.com/catalog-service/api/consumer/entitledCatalogItems/feaedf73-560c-4612-a573-41667e017691/requests"}],"iconId":"composition.blueprint.png","catalogItemTypeRef":{"id":"com.vmware.csp.component.cafe.composition.blueprint","label":"Composite Blueprint"},"serviceRef":{"id":"d33afbc2-954d-416d-8d4a-bc4ad1b66058","label":"test-service"},"outputResourceTypeRef":{"id":"composition.resource.type.deployment","label":"Deployment"}}],"metadata":{"size":20,"totalElements":4,"totalPages":1,"number":1,"offset":0}}`
	path := fmt.Sprintf(EntitledCatalogItemViewsAPI)
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url,
		httpmock.NewStringResponder(200, entitledCatalogItemViewsResponse))

	catalogItemID, err := client.ReadCatalogItemByName("CentOs")
	utils.AssertEqualsString(t, "feaedf73-560c-4612-a573-41667e017691", catalogItemID)
	utils.AssertNilError(t, err)

	catalogItemID, err = client.ReadCatalogItemByName("Invalid Catalog Item name")
	utils.AssertEqualsString(t, "", catalogItemID)
}

func TestGetBusinessGroupID(t *testing.T) {

	httpmock.ActivateNonDefault(client.Client)
	defer httpmock.DeactivateAndReset()

	subtenants := `{"links":[],"content":[{"@type":"Subtenant","id":"b2470b94-cbca-43db-be37-803cca7b0f1a","name":"Development","description":"created by demo content","subtenantRoles":null,"tenant":"qe","extensionData":{"entries":[{"key":"iaas-manager-emails","value":{"type":"string","value":"astoyanov@vcac.sqa-horizon.local"}}]}},{"@type":"Subtenant","id":"ff371ec6-d4d8-4dee-aa73-f09e1bb3a4fd","name":"Quality Engineering","description":"created by demo content","subtenantRoles":null,"tenant":"qe","extensionData":{"entries":[{"key":"iaas-manager-emails","value":{"type":"string","value":"astoyanov@vcac.sqa-horizon.local"}}]}},{"@type":"Subtenant","id":"e11a782e-cff3-4707-b1e1-88940507dab3","name":"Finance","description":"created by demo content","subtenantRoles":null,"tenant":"qe","extensionData":{"entries":[{"key":"iaas-manager-emails","value":{"type":"string","value":"gloria@coke.sqa-horizon.local"}}]}},{"@type":"Subtenant","id":"6273f684-210c-4839-85ec-573346f53799","name":"cloudclient-ITs-bg","description":"","subtenantRoles":null,"tenant":"qe","extensionData":{"entries":[{"key":"iaas-machine-prefix","value":{"type":"string","value":"b272b729-1583-4e6e-b299-15c293a6d0fa"}},{"key":"iaas-ad-container","value":{"type":"string","value":""}},{"key":"iaas-manager-emails","value":{"type":"string","value":"fritz@sqa-local.com"}}]}}],"metadata":{"size":20,"totalElements":4,"totalPages":1,"number":1,"offset":0}}`

	path := Tenants + "/" + tenant + "/subtenants"
	url := client.BuildEncodedURL(path, nil)

	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, subtenants))

	id, err := client.GetBusinessGroupID("Development", tenant)

	if id == "b2470b94-cbca-43db-be37-803cca7b0f1a" {
		fmt.Println("Passed")
	}

	if err != nil {
		t.Errorf("Error fetching is %v ", err)
	}
}
