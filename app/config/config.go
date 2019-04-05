package config

// AppName of the application
const AppName = "dliver-project-skeleton"

// AppVersion Version of the application
var AppVersion string

type Config struct {
	CentrifugoHost     string `mapstructure:"centrifugo_host"`
	CentrifugoGrpcPort int    `mapstructure:"centrifugo_grpc_port" default:"10000"`

	ElasticAddress string `mapstructure:"elastic_address"`

	GebUsername string `mapstructure:"geb_username"`
	GebPassword string `mapstructure:"geb_password"`
	GebHost     string `mapstructure:"geb_host"`
	GebPort     int    `mapstructure:"geb_port" default:"5672"`

	LogLevel string `mapstructure:"log_level" default:"info"`

	Port int `mapstructure:"port" default:"80"`

	RedisStoreHost            string `mapstructure:"redis_store_host"`
	RedisStorePort            int    `mapstructure:"redis_store_port" default:"6379"`
	RedisStoreDatabase        int    `mapstructure:"redis_store_database"`
	RedisStorePoolMaxIdle     int    `mapstructure:"redis_store_pool_max_idle" default:"10"`
	RedisStorePoolIdleTimeout string `mapstructure:"redis_store_pool_idle_timeout" default:"240s"`

	TracerSamplerType                 string `mapstructure:"tracer_sampler_type" default:"remote"`
	TracerSamplerParam                string `mapstructure:"tracer_sampler_param" default:"1.0"`
	TracerSamplerSamplingServerScheme string `mapstructure:"tracer_sampler_sampling_server_scheme" default:"http"`
	TracerSamplerSamplingServerHost   string `mapstructure:"tracer_sampler_sampling_server_host"`
	TracerSamplerSamplingServerPort   string `mapstructure:"tracer_sampler_sampling_server_port" default:"5778"`
	TracerReporterLocalAgentHost      string `mapstructure:"tracer_reporter_local_agent_host"`
	TracerReporterLocalAgentPort      int    `mapstructure:"tracer_reporter_local_agent_port" default:"6831"`

	YafudsHost string `mapstructure:"yafuds_host"`
	YafudsPort string `mapstructure:"yafuds_port" default:"7890"`
}
