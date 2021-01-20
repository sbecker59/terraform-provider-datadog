package datadog

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestAccDatadogMonitorsDatasource(t *testing.T) {
	accProviders, clock, cleanup := testAccProviders(t, initRecorder(t))
	uniq := strings.ReplaceAll(uniqueEntityName(clock, t), "-", "_")
	defer cleanup(t)
	accProvider := testAccProvider(t, accProviders)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    accProviders,
		CheckDestroy: testAccCheckDatadogMonitorDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMonitorsNameFilterConfig(uniq),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check:              checkDatasourceMonitorsAttrs(accProvider, uniq),
			},
			{
				Config: testAccDatasourceMonitorsTagsFilterConfig(uniq),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check:              checkDatasourceMonitorsAttrs(accProvider, uniq),
			},
			{
				Config: testAccDatasourceMonitorsMonitorTagsFilterConfig(uniq),
				// Because of the `depends_on` in the datasource, the plan cannot be empty.
				// See https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
				ExpectNonEmptyPlan: true,
				Check:              checkDatasourceMonitorsAttrs(accProvider, uniq),
			},
		},
	})
}

func checkDatasourceMonitorsAttrs(accProvider *schema.Provider, uniq string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testAccCheckDatadogMonitorExists(accProvider, "datadog_monitors.foo"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "name", uniq),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "message", "some message Notify: @hipchat-channel"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "type", "query alert"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "query", fmt.Sprintf("avg(last_4h):anomalies(avg:aws.ec2.cpu{environment:foo,host:foo,test_datasource_monitors_scope:%s} by {host}, 'basic', 2, direction='both', alert_window='last_15m', interval=60, count_default_zero='true') >= 1", uniq)),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "notify_no_data", "false"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "new_host_delay", "600"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "evaluation_delay", "700"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "renotify_interval", "60"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "thresholds.warning", "0.5"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "monitor_thresholds.0.warning", "0.5"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "thresholds.warning_recovery", "0.3"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "monitor_thresholds.0.warning_recovery", "0.3"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "thresholds.critical", "1"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "monitor_thresholds.0.critical", "1"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "thresholds.critical_recovery", "0.7"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "monitor_thresholds.0.critical_recovery", "0.7"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "require_full_window", "true"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "locked", "false"),
		// Tags are a TypeSet => use a weird way to access members by their hash
		// TF TypeSet is internally represented as a map that maps computed hashes
		// to actual values. Since the hashes are always the same for one value,
		// this is the way to get them.
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "tags.#", "2"),
		resource.TestCheckResourceAttr(
			"data.datadog_monitors.foo", "tags.2644851163", "baz"),
	)
}

func testAccMonitorsConfig(uniq string) string {
	return fmt.Sprintf(`
resource "datadog_monitors" "foo" {
  name = "%s"
  type = "metric alert"
  message = "some message Notify: @hipchat-channel"
  escalation_message = "the situation has escalated @pagerduty"

  query = "avg(last_4h):anomalies(avg:aws.ec2.cpu{environment:foo,host:foo,test_datasource_monitor_scope:%s} by {host}, 'basic', 2, direction='both', alert_window='last_15m', interval=60, count_default_zero='true') >= 1"

  thresholds = {
	warning = "0.5"
	critical = "1.0"
	warning_recovery = "0.3"
	critical_recovery = "0.7"
  }

	threshold_windows = {
		trigger_window = "last_15m"
		recovery_window = "last_15m"
	}

  renotify_interval = 60

  notify_audit = false
  timeout_h = 60
  new_host_delay = 600
  evaluation_delay = 700
  include_tags = true
  require_full_window = true
  locked = false
  tags = ["test_datasource_monitor:%s", "baz"]
}`, uniq, uniq, uniq)
}

func testAccDatasourceMonitorsNameFilterConfig(uniq string) string {
	return fmt.Sprintf(`
%s
data "datadog_monitors" "foo" {
  depends_on = [
    datadog_monitor.foo,
  ]
  name_filter = "%s"
}`, testAccMonitorsConfig(uniq), uniq)
}

func testAccDatasourceMonitorsTagsFilterConfig(uniq string) string {
	return fmt.Sprintf(`
%s
data "datadog_monitor" "foo" {
  depends_on = [
    datadog_monitor.foo,
  ]
  tags_filter = [
    "test_datasource_monitor_scope:%s",
  ]
}
`, testAccMonitorsConfig(uniq), uniq)
}

func testAccDatasourceMonitorsMonitorTagsFilterConfig(uniq string) string {
	return fmt.Sprintf(`
%s
data "datadog_monitors" "foo" {
  depends_on = [
    datadog_monitor.foo,
  ]
  monitor_tags_filter = [
    "test_datasource_monitor:%s",
  ]
}
`, testAccMonitorConfig(uniq), uniq)
}
