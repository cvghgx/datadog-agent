module github.com/datadog-agent/DataDog/comp/otelcol/collector-contrib/impl

go 1.21.0

replace (
	github.com/DataDog/datadog-agent/comp/otelcol/collector-contrib/def => ../def
	github.com/DataDog/datadog-agent/pkg/obfuscate => ../../../../pkg/obfuscate
	github.com/DataDog/datadog-agent/pkg/proto => ../../../../pkg/proto
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state => ../../../../pkg/remoteconfig/state
	github.com/DataDog/datadog-agent/pkg/trace => ../../../../pkg/trace
	github.com/DataDog/datadog-agent/pkg/util/cgroups => ../../../../pkg/util/cgroups
	github.com/DataDog/datadog-agent/pkg/util/log => ../../../../pkg/util/log
	github.com/DataDog/datadog-agent/pkg/util/pointer => ../../../../pkg/util/pointer
	github.com/DataDog/datadog-agent/pkg/util/scrubber => ../../../../pkg/util/scrubber
)

require (
	github.com/DataDog/datadog-agent/comp/otelcol/collector-contrib/def v0.55.0-rc.3
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/countconnector v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/exceptionsconnector v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/grafanacloudconnector v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/routingconnector v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/servicegraphconnector v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alertmanagerexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/alibabacloudlogserviceexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awscloudwatchlogsexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsemfexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awskinesisexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awss3exporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/awsxrayexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuredataexplorerexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuremonitorexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/cassandraexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/coralogixexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datasetexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/elasticsearchexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudpubsubexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlemanagedprometheusexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/honeycombmarkerexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/influxdbexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/instanaexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kineticaexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/loadbalancingexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logicmonitorexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/mezmoexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opensearchexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/otelarrowexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/pulsarexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/rabbitmqexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sapmexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sentryexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/skywalkingexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sumologicexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/syslogexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/tencentcloudlogserviceexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/ackextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/asapauthextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/awsproxy v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/basicauthextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/bearertokenauthextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/headerssetterextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/healthcheckextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/httpforwarderextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/jaegerremotesampling v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oauth2clientauthextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/dockerobserver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecsobserver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/ecstaskobserver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/hostobserver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer/k8sobserver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oidcauthextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/opampextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/pprofextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/sigv4authextension v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/dbstorage v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/filestorage v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbyattrsprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/groupbytraceprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricsgenerationprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/metricstransformprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/probabilisticsamplerprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/redactionprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/remotetapprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/routingprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/spanprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/sumologicprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/tailsamplingprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/activedirectorydsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/aerospikereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/apachesparkreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscloudwatchreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsecscontainermetricsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsfirehosereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureeventhubreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azuremonitorreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/bigipreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/chronyreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudflarereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/cloudfoundryreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/collectdreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/couchdbreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/expvarreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filestatsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flinkmetricsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudspannerreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/haproxyreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/httpcheckreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/iisreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/influxdbreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/journaldreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sobjectsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkareceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/lokireceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/memcachedreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbatlasreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mongodbreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/mysqlreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/namedpipereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nginxreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/nsxtreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/oracledbreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/podmanreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/postgresqlreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/pulsarreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefbreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/riakreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sapmreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/signalfxreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/simpleprometheusreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snmpreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/snowflakereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/solacereceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlqueryreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sqlserverreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/sshcheckreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/tcplogreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/udplogreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vcenterreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/wavefrontreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/webhookeventreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowseventlogreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/windowsperfcountersreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.100.0
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zookeeperreceiver v0.100.0
	go.opentelemetry.io/collector/connector v0.100.0
	go.opentelemetry.io/collector/connector/forwardconnector v0.100.0
	go.opentelemetry.io/collector/exporter v0.100.0
	go.opentelemetry.io/collector/exporter/debugexporter v0.100.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.100.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.100.0
	go.opentelemetry.io/collector/extension v0.100.0
	go.opentelemetry.io/collector/extension/ballastextension v0.100.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.100.0
	go.opentelemetry.io/collector/otelcol v0.100.0
	go.opentelemetry.io/collector/processor v0.100.0
	go.opentelemetry.io/collector/processor/batchprocessor v0.100.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.100.0
	go.opentelemetry.io/collector/receiver v0.100.0
	go.opentelemetry.io/collector/receiver/nopreceiver v0.100.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.100.0
	go.uber.org/multierr v1.11.0
)

require (
	bitbucket.org/atlassian/go-asap/v2 v2.8.0 // indirect
	cloud.google.com/go v0.112.2 // indirect
	cloud.google.com/go/auth v0.3.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.7 // indirect
	cloud.google.com/go/logging v1.9.0 // indirect
	cloud.google.com/go/longrunning v0.5.6 // indirect
	cloud.google.com/go/monitoring v1.18.1 // indirect
	cloud.google.com/go/pubsub v1.38.0 // indirect
	cloud.google.com/go/spanner v1.61.0 // indirect
	cloud.google.com/go/trace v1.10.6 // indirect
	code.cloudfoundry.org/clock v0.0.0-20180518195852-02e53af36e6c // indirect
	code.cloudfoundry.org/go-diodes v0.0.0-20211115184647-b584dd5df32c // indirect
	code.cloudfoundry.org/go-loggregator v7.4.0+incompatible // indirect
	code.cloudfoundry.org/rfc5424 v0.0.0-20201103192249-000122071b78 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/AthenZ/athenz v1.10.39 // indirect
	github.com/Azure/azure-amqp-common-go/v4 v4.2.0 // indirect
	github.com/Azure/azure-event-hubs-go/v3 v3.6.2 // indirect
	github.com/Azure/azure-kusto-go v0.15.2 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.5.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.5.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor v0.11.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4 v4.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.3.2 // indirect
	github.com/Azure/azure-storage-queue-go v0.0.0-20230531184854-c06a8eff66fe // indirect
	github.com/Azure/go-amqp v1.0.5 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.29 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.23 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/ClickHouse/ch-go v0.61.5 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.23.0 // indirect
	github.com/Code-Hex/go-generics-cache v1.3.1 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/GehirnInc/crypt v0.0.0-20200316065508-bb7000b8a962 // indirect
	github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.23.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector v0.47.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector/googlemanagedprometheus v0.47.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.23.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.47.0 // indirect
	github.com/IBM/sarama v1.43.2 // indirect
	github.com/JohnCGriffin/overflow v0.0.0-20211019200055-46fa312c352c // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ReneKroon/ttlcache/v2 v2.11.0 // indirect
	github.com/SAP/go-hdb v1.8.15 // indirect
	github.com/SermoDigital/jose v0.9.2-0.20180104203859-803625baeddc // indirect
	github.com/Showmax/go-fqdn v1.0.0 // indirect
	github.com/aerospike/aerospike-client-go/v6 v6.13.0 // indirect
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	github.com/alecthomas/units v0.0.0-20231202071711-9a357b53e9c9 // indirect
	github.com/aliyun/aliyun-log-go-sdk v0.1.72 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/apache/arrow/go/v14 v14.0.2 // indirect
	github.com/apache/arrow/go/v15 v15.0.0 // indirect
	github.com/apache/pulsar-client-go v0.8.1 // indirect
	github.com/apache/pulsar-client-go/oauth2 v0.0.0-20220120090717-25e59572242e // indirect
	github.com/apache/thrift v0.20.0 // indirect
	github.com/ardielle/ardielle-go v1.5.2 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aws/aws-sdk-go v1.53.2 // indirect
	github.com/aws/aws-sdk-go-v2 v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.27.13 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.13 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.27.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.53.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.29.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.24.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.7 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudfoundry-incubator/uaago v0.0.0-20190307164349-8136b7bbe76e // indirect
	github.com/cncf/xds/go v0.0.0-20240318125728-8a4994d93e50 // indirect
	github.com/coreos/go-oidc/v3 v3.10.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/devigned/tab v0.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/digitalocean/godo v1.109.0 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/docker/docker v25.0.5+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/eapache/go-resiliency v1.6.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/go-elasticsearch/v7 v7.17.10 // indirect
	github.com/elastic/go-structform v0.0.10 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/envoyproxy/go-control-plane v0.12.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/expr-lang/expr v1.16.7 // indirect
	github.com/facebook/time v0.0.0-20240109160331-d1456d1a6bac // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-openapi/swag v0.22.9 // indirect
	github.com/go-resty/resty/v2 v2.12.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gocql/gocql v1.6.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.4 // indirect
	github.com/gophercloud/gophercloud v1.8.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/gosnmp/gosnmp v1.37.0 // indirect
	github.com/grafana/loki/pkg/push v0.0.0-20240514112848-a1b1eeb09583 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/grobie/gomemcache v0.0.0-20230213081705-239240bbc445 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hamba/avro/v2 v2.20.1 // indirect
	github.com/hashicorp/consul/api v1.28.2 // indirect
	github.com/hashicorp/cronexpr v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20240306004928-3e7191ccb702 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hetznercloud/hcloud-go/v2 v2.6.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/influxdata/go-syslog/v3 v3.0.1-0.20230911200830-875f5bc594a4 // indirect
	github.com/influxdata/influxdb-observability/common v0.5.8 // indirect
	github.com/influxdata/influxdb-observability/influx2otel v0.5.8 // indirect
	github.com/influxdata/influxdb-observability/otel2influx v0.5.8 // indirect
	github.com/influxdata/line-protocol/v2 v2.2.1 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.11 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jaegertracing/jaeger v1.57.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/kineticadb/kinetica-api-go v0.0.5 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20181214104525-299bdde78165 // indirect
	github.com/leoluk/perflib_exporter v0.2.1 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/linkedin/goavro/v2 v2.9.8 // indirect
	github.com/linode/linodego v1.33.0 // indirect
	github.com/logicmonitor/lm-data-sdk-go v1.3.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-ieproxy v0.0.11 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/microsoft/ApplicationInsights-Go v0.4.4 // indirect
	github.com/microsoft/go-mssqldb v1.7.1 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mongodb-forks/digest v1.1.0 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/mostynb/go-grpc-compression v1.2.2 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nginxinc/nginx-prometheus-exporter v0.11.0 // indirect
	github.com/oklog/ulid/v2 v2.1.0 // indirect
	github.com/open-telemetry/opamp-go v0.14.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/extension/sumologicextension v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/awsutil v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/cwlogs v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/collectd v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/docker v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kafka v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/kubelet v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sqlquery v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpersignal v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/azure v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/loki v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalking v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.100.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/winperfcounters v0.100.0 // indirect
	github.com/open-telemetry/otel-arrow v0.23.0 // indirect
	github.com/open-telemetry/otel-arrow/collector v0.23.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opensearch-project/opensearch-go/v2 v2.3.0 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openshift/client-go v0.0.0-20210521082421-73d9475a9142 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/ovh/go-ovh v1.4.3 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.13.6 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.15.0 // indirect
	github.com/prometheus/prometheus v0.51.2-0.20240405174432-b4a973753c6e // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/redis/go-redis/v9 v9.5.1 // indirect
	github.com/relvacode/iso8601 v1.4.0 // indirect
	github.com/rs/cors v1.11.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/samber/lo v1.39.0 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.25 // indirect
	github.com/scalyr/dataset-go v0.18.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.4 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/sapm-proto v0.14.0 // indirect
	github.com/sijms/go-ora/v2 v2.8.14 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/snowflakedb/gosnowflake v1.9.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/cobra v1.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.18.2 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.920 // indirect
	github.com/tg123/go-htpasswd v1.2.2 // indirect
	github.com/tidwall/gjson v1.14.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tidwall/wal v1.1.7 // indirect
	github.com/tilinna/clock v1.1.0 // indirect
	github.com/tinylib/msgp v1.1.9 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	github.com/vmware/go-vmware-nsxt v0.0.0-20230223012718-d31b8a1ca05e // indirect
	github.com/vmware/govmomi v0.36.3 // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yuin/gopher-lua v0.0.0-20220504180219-658193537a64 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	github.com/ztrue/tracerr v0.4.0 // indirect
	go.etcd.io/bbolt v1.3.9 // indirect
	go.mongodb.org/atlas v0.36.0 // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector v0.100.0 // indirect
	go.opentelemetry.io/collector/component v0.100.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.100.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.8.0 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.100.0 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.100.0 // indirect
	go.opentelemetry.io/collector/config/confignet v0.100.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.7.0 // indirect
	go.opentelemetry.io/collector/config/configretry v0.100.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.100.0 // indirect
	go.opentelemetry.io/collector/config/configtls v0.100.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/converter/expandconverter v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/envprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/fileprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpsprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v0.100.0 // indirect
	go.opentelemetry.io/collector/consumer v0.100.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.100.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.8.0 // indirect
	go.opentelemetry.io/collector/filter v0.100.0 // indirect
	go.opentelemetry.io/collector/pdata v1.8.0 // indirect
	go.opentelemetry.io/collector/semconv v0.100.0 // indirect
	go.opentelemetry.io/collector/service v0.100.0 // indirect
	go.opentelemetry.io/contrib/config v0.6.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.51.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.51.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.26.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.51.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/bridge/opencensus v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.48.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/oauth2 v0.20.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	gonum.org/v1/gonum v0.15.0 // indirect
	google.golang.org/api v0.178.0 // indirect
	google.golang.org/genproto v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240429193739-8cf5692501f6 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240513163218-0867130af1f8 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.29.3 // indirect
	k8s.io/apimachinery v0.29.3 // indirect
	k8s.io/client-go v0.29.3 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/kube-openapi v0.0.0-20231010175941-2dd684a91f00 // indirect
	k8s.io/kubelet v0.29.3 // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
	sigs.k8s.io/controller-runtime v0.17.3 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
	skywalking.apache.org/repo/goapi v0.0.0-20240104145220-ba7202308dd4 // indirect
)
