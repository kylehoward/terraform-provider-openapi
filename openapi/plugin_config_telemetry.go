package openapi

// TelemetryProvider holds the behaviour expected to be implemented for the Telemetry Providers supported. At the moment
// only Graphite is supported.
type TelemetryProvider interface {
	// Validate performs a check to confirm that the telemetry configuration is valid
	Validate() error
	// IncOpenAPIPluginVersionTotalRunsCounter is the method responsible for submitting to the corresponding telemetry platform the counter increase for the OpenAPI plugin Version used
	IncOpenAPIPluginVersionTotalRunsCounter(openAPIPluginVersion string) error
	// IncServiceProviderTotalRunsCounter is the method responsible for submitting to the corresponding telemetry platform the counter increase for the service provider used
	IncServiceProviderTotalRunsCounter(providerName string) error
}
