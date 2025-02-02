package exporters

import (
	"github.com/gophercloud/gophercloud/openstack/placement/v1/resourceproviders"
	"github.com/prometheus/client_golang/prometheus"
)

type PlacementExporter struct {
	BaseOpenStackExporter
}

var defaultPlacementMetrics = []Metric{
	{Name: "resource_total", Fn: ListPlacementResourceProviders, Labels: []string{"hostname", "resourcetype", "region_name"}},
	{Name: "resource_allocation_ratio", Labels: []string{"hostname", "resourcetype", "region_name"}},
	{Name: "resource_reserved", Labels: []string{"hostname", "resourcetype", "region_name"}},
	{Name: "resource_usage", Labels: []string{"hostname", "resourcetype", "region_name"}},
}

func NewPlacementExporter(config *ExporterConfig) (*PlacementExporter, error) {
	exporter := PlacementExporter{
		BaseOpenStackExporter{
			Name:           "placement",
			ExporterConfig: *config,
		},
	}
	for _, metric := range defaultPlacementMetrics {
		if exporter.isDeprecatedMetric(&metric) {
			continue
		}
		if !exporter.isSlowMetric(&metric) {
			exporter.AddMetric(metric.Name, metric.Fn, metric.Labels, metric.DeprecatedVersion, nil)
		}
	}
	return &exporter, nil
}

func ListPlacementResourceProviders(exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allResourceProviders []resourceproviders.ResourceProvider

	allPagesResourceProviders, err := resourceproviders.List(exporter.Client, resourceproviders.ListOpts{}).AllPages()
	if err != nil {
		return err
	}

	if allResourceProviders, err = resourceproviders.ExtractResourceProviders(allPagesResourceProviders); err != nil {
		return err
	}

	uuidToNameMap := map[string]string{}

	for _, resourceprovider := range allResourceProviders {
		uuidToNameMap[resourceprovider.UUID] = resourceprovider.Name

		inventoryResult, err := resourceproviders.GetInventories(exporter.Client, resourceprovider.UUID).Extract()
		if err != nil {
			return err
		}

		for k, v := range inventoryResult.Inventories {

			ch <- prometheus.MustNewConstMetric(exporter.Metrics["resource_total"].Metric,
				prometheus.GaugeValue, float64(v.Total), resourceprovider.Name, k,
				endpointOpts["placement"].Region)

			ch <- prometheus.MustNewConstMetric(exporter.Metrics["resource_allocation_ratio"].Metric,
				prometheus.GaugeValue, float64(v.AllocationRatio), resourceprovider.Name, k,
				endpointOpts["placement"].Region)

			ch <- prometheus.MustNewConstMetric(exporter.Metrics["resource_reserved"].Metric,
				prometheus.GaugeValue, float64(v.Reserved), resourceprovider.Name, k,
				endpointOpts["placement"].Region)
		}

		usagesResult, err := resourceproviders.GetUsages(exporter.Client, resourceprovider.UUID).Extract()
		if err != nil {
			return err
		}

		for k, v := range usagesResult.Usages {
			ch <- prometheus.MustNewConstMetric(exporter.Metrics["resource_usage"].Metric,
				prometheus.GaugeValue, float64(v), resourceprovider.Name, k,
				endpointOpts["placement"].Region)
		}

	}

	return nil

}
