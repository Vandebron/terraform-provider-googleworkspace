// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

//TODO: Fix tests
//func TestAccGroupResource(t *testing.T) {
//	//	resource.Test(t, resource.TestCase{
//	//		PreCheck:                 func() { testAccPreCheck(t) },
//	//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//	//		Steps: []resource.TestStep{
//	//			// Create and Read testing
//	//			{
//	//				Config: testAccGroupResourceConfig("test", "test@example.com"),
//	//				ConfigStateChecks: []statecheck.StateCheck{
//	//					statecheck.ExpectKnownValue(
//	//						"googleworkspace_group.test",
//	//						tfjsonpath.New("id"),
//	//						knownvalue.StringExact("example-id"),
//	//					),
//	//					statecheck.ExpectKnownValue(
//	//						"googleworkspace_group.test",
//	//						tfjsonpath.New("defaulted"),
//	//						knownvalue.StringExact("example value when not configured"),
//	//					),
//	//					statecheck.ExpectKnownValue(
//	//						"googleworkspace_group.test",
//	//						tfjsonpath.New("configurable_attribute"),
//	//						knownvalue.StringExact("one"),
//	//					),
//	//				},
//	//			},
//	//			// ImportState testing
//	//			{
//	//				ResourceName:      "googleworkspace_group.test",
//	//				ImportState:       true,
//	//				ImportStateVerify: true,
//	//				// This is not normally necessary, but is here because this
//	//				// example code does not have an actual upstream service.
//	//				// Once the Read method is able to refresh information from
//	//				// the upstream service, this can be removed.
//	//				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
//	//			},
//	//			// Update and Read testing
//	//			{
//	//				Config: testAccGroupResourceConfig("test", "test@example.com"),
//	//				ConfigStateChecks: []statecheck.StateCheck{
//	//					statecheck.ExpectKnownValue(
//	//						"googleworkspace_group.test",
//	//						tfjsonpath.New("id"),
//	//						knownvalue.StringExact("example-id"),
//	//					),
//	//					statecheck.ExpectKnownValue(
//	//						"googleworkspace_group.test",
//	//						tfjsonpath.New("defaulted"),
//	//						knownvalue.StringExact("example value when not configured"),
//	//					),
//	//					statecheck.ExpectKnownValue(
//	//						"googleworkspace_group.test",
//	//						tfjsonpath.New("configurable_attribute"),
//	//						knownvalue.StringExact("two"),
//	//					),
//	//				},
//	//			},
//	//			// Delete testing automatically occurs in TestCase
//	//		},
//	//	})
//}
//
//func testAccGroupResourceConfig(name, email string) string {
//	return fmt.Sprintf(`
//resource "googleworkspace_group" "test" {
//  name = %[1]q
//	email = %[2]q
//}
//`, name, email)
//}
