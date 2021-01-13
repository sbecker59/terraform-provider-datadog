package datadog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	datadogV1 "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDatadogMonitors() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve information about an existing monitor for use in other resources.",
		Read:        dataSourceDatadogMonitorsRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			// Computed values
			"monitors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Description: "Name of the monitor",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"message": {
							Description: "Message included with notifications for this monitor",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"escalation_message": {
							Description: "Message included with a re-notification for this monitor.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"query": {
							Description: "Query of the monitor.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "Type of the monitor.",
							Type:        schema.TypeString,
							Computed:    true,
						},

						// Options
						"thresholds": {
							Description: "Alert thresholds of the monitor.",
							Type:        schema.TypeMap,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ok": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"warning": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"critical": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"unknown": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"warning_recovery": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"critical_recovery": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"threshold_windows": {
							Description: "Mapping containing `recovery_window` and `trigger_window` values, e.g. `last_15m`. This is only used by anomaly monitors.",
							Type:        schema.TypeMap,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recovery_window": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"trigger_window": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"notify_no_data": {
							Description: "Whether or not this monitor notifies when data stops reporting.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"new_host_delay": {
							Description: "Time (in seconds) allowing a host to boot and applications to fully start before starting the evaluation of monitor results.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"evaluation_delay": {
							Description: "Time (in seconds) for which evaluation is delayed. This is only used by metric monitors.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"no_data_timeframe": {
							Description: "The number of minutes before the monitor notifies when data stops reporting.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"renotify_interval": {
							Description: "The number of minutes after the last notification before the monitor re-notifies on the current status.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"notify_audit": {
							Description: "Whether or not tagged users are notified on changes to the monitor.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"timeout_h": {
							Description: "Number of hours of the monitor not reporting data before it automatically resolves from a triggered state.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"require_full_window": {
							Description: "Whether or not the monitor needs a full window of data before it is evaluated.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"locked": {
							Description: "Whether or not changes to the monitor are restricted to the creator or admins.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"include_tags": {
							Description: "Whether or not notifications from the monitor automatically inserts its triggering tags into the title.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"tags": {
							Description: "List of tags associated with the monitor.",
							// we use TypeSet to represent tags, paradoxically to be able to maintain them ordered;
							// we order them explicitly in the read/create/update methods of this resource and using
							// TypeSet makes Terraform ignore differences in order when creating a plan
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enable_logs_sample": {
							Description: "Whether or not a list of log values which triggered the alert is included. This is only used by log monitors.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDatadogMonitorsRead(d *schema.ResourceData, meta interface{}) error {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV1 := providerConf.DatadogClientV1
	authV1 := providerConf.AuthV1

	req := datadogClientV1.MonitorsApi.ListMonitors(authV1)
	res, _, err := req.Execute()
	if err != nil {
		return translateClientError(err, "error querying monitors")
	}
	monitors := []map[string]interface{}{}
	d.SetId(strconv.FormatInt(m.GetId(), 10))

	for _, r := range res {

		monitor := map[string]interface{}{}
		thresholds := make(map[string]string)

		if v, ok := r.Options.Thresholds.GetOkOk(); ok {
			thresholds["ok"] = fmt.Sprintf("%v", *v)
		}
		if v, ok := r.Options.Thresholds.GetWarningOk(); ok {
			thresholds["warning"] = fmt.Sprintf("%v", *v)
		}
		if v, ok := r.Options.Thresholds.GetCriticalOk(); ok {
			thresholds["critical"] = fmt.Sprintf("%v", *v)
		}
		if v, ok := r.Options.Thresholds.GetUnknownOk(); ok {
			thresholds["unknown"] = fmt.Sprintf("%v", *v)
		}
		if v, ok := r.Options.Thresholds.GetWarningRecoveryOk(); ok {
			thresholds["warning_recovery"] = fmt.Sprintf("%v", *v)
		}
		if v, ok := r.Options.Thresholds.GetCriticalRecoveryOk(); ok {
			thresholds["critical_recovery"] = fmt.Sprintf("%v", *v)
		}

		thresholdWindows := make(map[string]string)
		for k, v := range map[string]string{
			"recovery_window": m.Options.ThresholdWindows.GetRecoveryWindow(),
			"trigger_window":  m.Options.ThresholdWindows.GetTriggerWindow(),
		} {
			if v != "" {
				thresholdWindows[k] = v
			}
		}

		var tags []string
		for _, s := range m.GetTags() {
			tags = append(tags, s)
		}
		sort.Strings(tags)
		monitor["name"] = r.GetName()
		monitor["message"] = r.GetMessage()
		monitor["query"] = r.GetQuery()
		monitor["type"] = r.GetType()

		monitor["thresholds"] = thresholds
		monitor["threshold_windows"] = thresholdWindows

		monitor["new_host_delay"] = r.Options.GetNewHostDelay()
		monitor["evaluation_delay"] = r.Options.GetEvaluationDelay()
		monitor["notify_no_data"] = r.Options.GetNotifyNoData()
		monitor["no_data_timeframe"] = r.Options.NoDataTimeframe
		monitor["renotify_interval"] = r.Options.GetRenotifyInterval()
		monitor["notify_audit"] = r.Options.GetNotifyAudit()
		monitor["timeout_h"] = r.Options.GetTimeoutH()
		monitor["escalation_message"] = r.Options.GetEscalationMessage()
		monitor["include_tags"] = r.Options.GetIncludeTags()
		monitor["tags"] = tags
		monitor["require_full_window"] = r.Options.GetRequireFullWindow() // TODO Is this one of those options that we neeed to check?
		monitor["locked"] = r.Options.GetLocked()

		if m.GetType() == datadogV1.MONITORTYPE_LOG_ALERT {
			monitor["enable_logs_sample"] = r.Options.GetEnableLogsSample()
		}

		monitors = append(monitors, monitor)

	}

	if f, fOk := d.GetOkExists("filter"); fOk {
		resources = ApplyFilters(f.(*schema.Set), resources, dataSourceDatadogMonitors().Schema["monitors"].Elem.(*schema.Resource).Schema)
	}

	return nil
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
