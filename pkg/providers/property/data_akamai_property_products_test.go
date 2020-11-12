package property

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
)

func TestVerifyProductsDataSourceSchema(t *testing.T) {
	t.Run("akamai_property_products - test data source required contract", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			Providers:  testAccProviders,
			IsUnitTest: true,
			Steps: []resource.TestStep{{
				Config:      testConfig(""),
				ExpectError: regexp.MustCompile("The argument \"contract_id\" is required, but no definition was found"),
			}},
		})
	})
}

func TestOutputProductsDataSource(t *testing.T) {

	t.Run("akamai_property_products - input OK - output OK", func(t *testing.T) {
		client := &mockpapi{}
		client.On("GetProducts", AnyCTX, mock.Anything).Return(&papi.GetProductsResponse{
			AccountID:  "act_anyAccount",
			ContractID: "ctr_AnyContract",
			Products: papi.ProductsItems{
				Items: []papi.ProductItem{{ProductName: "anyProduct", ProductID: "prd_anyProduct"}},
			},
		}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers:  testAccProviders,
				IsUnitTest: true,
				Steps: []resource.TestStep{{
					Config: testConfig("contract_id = \"ctr_test\""),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckOutput("product_name0", "anyProduct"),
						resource.TestCheckOutput("product_id0", "prd_anyProduct"),
					),
				}},
			})
		})
	})
}

func testConfig(contractIdConfig string) string {
	return fmt.Sprintf(`
	provider "akamai" {
		edgerc = "~/.edgerc"
	}

	data "akamai_property_products" "example" { %s }

    output "product_name0" {
		value = "${data.akamai_property_products.example.products[0].product_name}"
	}

    output "product_id0" {
		value = "${data.akamai_property_products.example.products[0].product_id}"
	}
`, contractIdConfig)
}
