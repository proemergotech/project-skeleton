package config

// AppName of the application
const AppName = "dliver-project-skeleton"

// AppVersion Version of the application
var AppVersion string

type Config struct {
	Port       int  `mapstructure:"port" default:"80"`
	PublicPort int  `mapstructure:"public_port" default:"8080"`
	DebugPProf bool `mapstructure:"debug_pprof" default:"false"`

	CentrifugeScheme string `mapstructure:"centrifuge_scheme" default:"http"`
	CentrifugeHost   string `mapstructure:"centrifuge_host" validate:"required"`
	CentrifugePort   string `mapstructure:"centrifuge_port" default:"80"`

	ElasticSearchScheme string `mapstructure:"elastic_search_scheme" default:"http"`
	ElasticSearchHost   string `mapstructure:"elastic_search_host" validate:"required"`
	ElasticSearchPort   int    `mapstructure:"elastic_search_port" default:"9200"`

	GebUsername string `mapstructure:"geb_username" validate:"required"`
	GebPassword string `mapstructure:"geb_password" validate:"required"`
	GebHost     string `mapstructure:"geb_host" validate:"required"`
	GebPort     int    `mapstructure:"geb_port" default:"5672"`

	RedisStoreHost            string `mapstructure:"redis_store_host" validate:"required"`
	RedisStorePort            int    `mapstructure:"redis_store_port" default:"6379"`
	RedisStoreDatabase        int    `mapstructure:"redis_store_database" validate:"required"`
	RedisStorePoolMaxIdle     int    `mapstructure:"redis_store_pool_max_idle" default:"10"`
	RedisStorePoolIdleTimeout string `mapstructure:"redis_store_pool_idle_timeout" default:"240s"`

	TracerSamplerType                 string `mapstructure:"tracer_sampler_type" default:"remote"`
	TracerSamplerParam                string `mapstructure:"tracer_sampler_param" default:"1.0"`
	TracerSamplerSamplingServerScheme string `mapstructure:"tracer_sampler_sampling_server_scheme" default:"http"`
	TracerSamplerSamplingServerHost   string `mapstructure:"tracer_sampler_sampling_server_host" validate:"required"`
	TracerSamplerSamplingServerPort   string `mapstructure:"tracer_sampler_sampling_server_port" default:"5778"`
	TracerReporterLocalAgentHost      string `mapstructure:"tracer_reporter_local_agent_host" validate:"required"`
	TracerReporterLocalAgentPort      int    `mapstructure:"tracer_reporter_local_agent_port" default:"6831"`

	YafudsHost string `mapstructure:"yafuds_host" validate:"required"`
	YafudsPort string `mapstructure:"yafuds_port" default:"7890"`

	SiteConfigServiceScheme string `mapstructure:"site_config_service_scheme" default:"http"`
	SiteConfigServiceHost   string `mapstructure:"site_config_service_host" validate:"required"`
	SiteConfigServicePort   int    `mapstructure:"site_config_service_port" default:"80"`

	ConfigFileContent map[string]Content `mapstructure:"content" validate:"required"`
}

type Content struct {
	Price float64 `mapstructure:"price"`
}

type BootstrapConfig struct {
	YafudsHost string `mapstructure:"yafuds_host" validate:"required"`
	YafudsPort string `mapstructure:"yafuds_port" default:"7890"`
}
