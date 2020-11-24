package config

// AppName of the application
//%: {{ `
const AppName = "project-skeleton" //%: ` | replace "project-skeleton" .ProjectName | trim}}

// AppVersion Version of the application
var AppVersion string

type Config struct {
	Port int `mapstructure:"port" default:"80"`
	//%: {{- if .PublicRest }}
	PublicPort int `mapstructure:"public_port" default:"8080"`
	//%: {{- end }}
	DebugPProf bool `mapstructure:"debug_pprof" default:"false"`

	//%: {{ if .Elastic }}
	ElasticSearchScheme string `mapstructure:"elastic_search_scheme" default:"http"`
	ElasticSearchHost   string `mapstructure:"elastic_search_host" validate:"required"`
	ElasticSearchPort   int    `mapstructure:"elastic_search_port" default:"9200"`
	//%: {{ end }}

	//%: {{ if .RedisCache }}
	RedisCacheHost            string `mapstructure:"redis_cache_host" validate:"required"`
	RedisCachePort            int    `mapstructure:"redis_cache_port" default:"6379"`
	RedisCacheDatabase        int    `mapstructure:"redis_cache_database" validate:"required"`
	RedisCachePoolMaxIdle     int    `mapstructure:"redis_cache_pool_max_idle" default:"10"`
	RedisCachePoolIdleTimeout string `mapstructure:"redis_cache_pool_idle_timeout" default:"240s"`
	//%: {{ end }}

	//%: {{ if .RedisStore }}
	RedisStoreHost            string `mapstructure:"redis_store_host" validate:"required"`
	RedisStorePort            int    `mapstructure:"redis_store_port" default:"6379"`
	RedisStoreDatabase        int    `mapstructure:"redis_store_database" validate:"required"`
	RedisStorePoolMaxIdle     int    `mapstructure:"redis_store_pool_max_idle" default:"10"`
	RedisStorePoolIdleTimeout string `mapstructure:"redis_store_pool_idle_timeout" default:"240s"`
	//%: {{ end }}

	//%: {{ if .RedisNotice }}
	RedisNoticeHost            string `mapstructure:"redis_notice_host" validate:"required"`
	RedisNoticePort            int    `mapstructure:"redis_notice_port" default:"6379"`
	RedisNoticeDatabase        int    `mapstructure:"redis_notice_database" validate:"required"`
	RedisNoticePoolMaxIdle     int    `mapstructure:"redis_notice_pool_max_idle" default:"10"`
	RedisNoticePoolIdleTimeout string `mapstructure:"redis_notice_pool_idle_timeout" default:"240s"`
	//%: {{ end }}

	TracerSamplerType                 string `mapstructure:"tracer_sampler_type" default:"remote"`
	TracerSamplerParam                string `mapstructure:"tracer_sampler_param" default:"1.0"`
	TracerSamplerSamplingServerScheme string `mapstructure:"tracer_sampler_sampling_server_scheme" default:"http"`
	TracerSamplerSamplingServerHost   string `mapstructure:"tracer_sampler_sampling_server_host" validate:"required"`
	TracerSamplerSamplingServerPort   string `mapstructure:"tracer_sampler_sampling_server_port" default:"5778"`
	TracerReporterLocalAgentHost      string `mapstructure:"tracer_reporter_local_agent_host" validate:"required"`
	TracerReporterLocalAgentPort      int    `mapstructure:"tracer_reporter_local_agent_port" default:"6831"`

	//%: {{ if .ConfigFile }}
	ConfigFileContent map[string]Content `mapstructure:"content" validate:"required"`
	//%: {{ end }}
}

//%: {{ if .ConfigFile }}
type Content struct {
	Price float64 `mapstructure:"price"`
} //%: {{ end }}

//%: {{ if .Bootstrap }}
type BootstrapConfig struct {
} //%: {{ end }}
