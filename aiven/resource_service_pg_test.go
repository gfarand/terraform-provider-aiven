package aiven

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"strings"
	"testing"
)

// PG service tests
func TestAccAivenService_pg(t *testing.T) {
	t.Parallel()
	resourceName := "aiven_service.bar-pg"
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAivenServiceResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPGServiceResource(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAivenServiceCommonAttributes("data.aiven_service.service-pg"),
					testAccCheckAivenServicePGAttributes("data.aiven_service.service-pg"),
					resource.TestCheckResourceAttr(resourceName, "service_name", fmt.Sprintf("test-acc-sr-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "project", fmt.Sprintf("test-acc-pr-pg-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "service_type", "pg"),
					resource.TestCheckResourceAttr(resourceName, "cloud_name", "google-europe-west1"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window_dow", "monday"),
					resource.TestCheckResourceAttr(resourceName, "maintenance_window_time", "10:00:00"),
					resource.TestCheckResourceAttr(resourceName, "state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "termination_protection", "false"),
				),
			},
		},
	})
}

func testAccPGServiceResource(name string) string {
	return fmt.Sprintf(`
		resource "aiven_project" "foo-pg" {
			project = "test-acc-pr-pg-%s"
			card_id="%s"	
		}
		
		resource "aiven_service" "bar-pg" {
			project = aiven_project.foo-pg.project
			cloud_name = "google-europe-west1"
			plan = "startup-4"
			service_name = "test-acc-sr-%s"
			service_type = "pg"
			maintenance_window_dow = "monday"
			maintenance_window_time = "10:00:00"
			
			pg_user_config {
				pg_version = 11

				public_access {
					pg = true
					prometheus = false
				}

				pg {
					idle_in_transaction_session_timeout = 900
				}
			}
		}
		
		data "aiven_service" "service-pg" {
			service_name = aiven_service.bar-pg.service_name
			project = aiven_project.foo-pg.project
		}
		`, name, os.Getenv("AIVEN_CARD_ID"), name)
}

func testAccCheckAivenServicePGAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if !strings.Contains(a["service_type"], "pg") {
			return fmt.Errorf("expected to get a correct service type from Aiven, got :%s", a["service_type"])
		}

		if a["pg_user_config.0.pg_version"] != "11" {
			return fmt.Errorf("expected to get a correct PG version from Aiven, got :%s", a["pg_user_config.0.pg_version"])
		}

		if a["pg_user_config.0.ip_filter.0"] != "0.0.0.0/0" {
			return fmt.Errorf("expected to get a correct PG ip_filter from Aiven, got :%s", a["pg_user_config.0.ip_filter.0"])
		}

		if a["pg.0.dbname"] != "defaultdb" {
			return fmt.Errorf("expected to get a correct PG dbname from Aiven, got :%s", a["pg.0.dbname"])
		}

		if a["pg_user_config.0.pg.0.idle_in_transaction_session_timeout"] != "900" {
			return fmt.Errorf("expected to get a correct PG idle_in_transaction_session_timeout from Aiven")
		}

		if a["pg_user_config.0.public_access.0.pg"] != "true" {
			return fmt.Errorf("expected to get a correct PG public_access from Aiven")
		}

		if a["pg_user_config.0.public_access.0.pgbouncer"] != "<<value not set>>" {
			return fmt.Errorf("expected to get a correct PG public_access from Aiven")
		}

		if a["pg_user_config.0.public_access.0.prometheus"] != "false" {
			return fmt.Errorf("expected to get a correct PG public_access from Aiven")
		}

		if a["pg.0.host"] == "" {
			return fmt.Errorf("expected to get a PG host from Aiven")
		}

		if a["pg.0.user"] == "" {
			return fmt.Errorf("expected to get a PG user from Aiven")
		}

		if a["pg.0.password"] == "" {
			return fmt.Errorf("expected to get a PG password from Aiven")
		}

		if a["pg.0.port"] == "" {
			return fmt.Errorf("expected to get a PG port from Aiven")
		}

		if a["pg.0.uri"] == "" {
			return fmt.Errorf("expected to get a PG URI from Aiven")
		}

		if a["pg_user_config.0.service_to_fork_from"] != "<<value not set>>" {
			return fmt.Errorf("expected to get a PG service_to_fork_from not set to any value")
		}

		return nil
	}
}
