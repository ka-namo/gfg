package config

// ENVConfig is an ENV configuration.
type ENVConfig struct {
	NotifyEmail bool `envconfig:"NOTIFY_SMS"`
	NotifySMS   bool `envconfig:"NOTIFY_EMAIL"`
}
