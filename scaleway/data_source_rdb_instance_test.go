package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScalewayDataSourceRdbInstance_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	latestEngineVersion := testAccCheckScalewayRdbEngineGetLatestVersion(tt, postgreSQLEngineName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayRdbInstanceDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "scaleway_rdb_instance" "test" {
						name = "data-rdb-test-terraform"
						engine = %q
						node_type = "db-dev-s"
					}

					data "scaleway_rdb_instance" "test" {
						name = scaleway_rdb_instance.test.name
					}

					data "scaleway_rdb_instance" "test2" {
						instance_id = scaleway_rdb_instance.test.id
					}
				`, latestEngineVersion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayRdbExists(tt, "scaleway_rdb_instance.test"),

					resource.TestCheckResourceAttr("scaleway_rdb_instance.test", "name", "data-rdb-test-terraform"),
					resource.TestCheckResourceAttrSet("data.scaleway_rdb_instance.test", "id"),

					resource.TestCheckResourceAttr("data.scaleway_rdb_instance.test2", "name", "data-rdb-test-terraform"),
					resource.TestCheckResourceAttrSet("data.scaleway_rdb_instance.test2", "id"),
				),
			},
		},
	})
}
