package config

import (
	"strings"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils"
)

// desanitize replaces fake settings with their actual values.
func desanitize(actual, target *model.Config) {
	if target.LdapSettings.BindPassword != nil && *target.LdapSettings.BindPassword == model.FAKE_SETTING {
		*target.LdapSettings.BindPassword = *actual.LdapSettings.BindPassword
	}

	if *target.FileSettings.PublicLinkSalt == model.FAKE_SETTING {
		*target.FileSettings.PublicLinkSalt = *actual.FileSettings.PublicLinkSalt
	}
	if *target.FileSettings.AmazonS3SecretAccessKey == model.FAKE_SETTING {
		target.FileSettings.AmazonS3SecretAccessKey = actual.FileSettings.AmazonS3SecretAccessKey
	}

	if *target.EmailSettings.InviteSalt == model.FAKE_SETTING {
		target.EmailSettings.InviteSalt = actual.EmailSettings.InviteSalt
	}
	if *target.EmailSettings.SMTPPassword == model.FAKE_SETTING {
		target.EmailSettings.SMTPPassword = actual.EmailSettings.SMTPPassword
	}

	if *target.GitLabSettings.Secret == model.FAKE_SETTING {
		target.GitLabSettings.Secret = actual.GitLabSettings.Secret
	}

	if *target.SqlSettings.DataSource == model.FAKE_SETTING {
		*target.SqlSettings.DataSource = *actual.SqlSettings.DataSource
	}
	if *target.SqlSettings.AtRestEncryptKey == model.FAKE_SETTING {
		target.SqlSettings.AtRestEncryptKey = actual.SqlSettings.AtRestEncryptKey
	}

	if *target.ElasticsearchSettings.Password == model.FAKE_SETTING {
		*target.ElasticsearchSettings.Password = *actual.ElasticsearchSettings.Password
	}

	target.SqlSettings.DataSourceReplicas = make([]string, len(actual.SqlSettings.DataSourceReplicas))
	for i := range target.SqlSettings.DataSourceReplicas {
		target.SqlSettings.DataSourceReplicas[i] = actual.SqlSettings.DataSourceReplicas[i]
	}

	target.SqlSettings.DataSourceSearchReplicas = make([]string, len(actual.SqlSettings.DataSourceReplicas))
	for i := range target.SqlSettings.DataSourceSearchReplicas {
		target.SqlSettings.DataSourceSearchReplicas[i] = actual.SqlSettings.DataSourceSearchReplicas[i]
	}
}

// fixConfig patches invalid or missing data in the configuration, returning true if changed.
func fixConfig(cfg *model.Config) bool {
	changed := false

	// Ensure SiteURL has no trailing slash.
	if strings.HasSuffix(*cfg.ServiceSettings.SiteURL, "/") {
		*cfg.ServiceSettings.SiteURL = strings.TrimRight(*cfg.ServiceSettings.SiteURL, "/")
		changed = true
	}

	// Ensure the directory for a local file store has a trailing slash.
	if *cfg.FileSettings.DriverName == model.IMAGE_DRIVER_LOCAL {
		if !strings.HasSuffix(*cfg.FileSettings.Directory, "/") {
			*cfg.FileSettings.Directory += "/"
			changed = true
		}
	}

	// Fix various invalid locale issues.
	locales := utils.GetSupportedLocales()
	if _, ok := locales[*cfg.LocalizationSettings.DefaultServerLocale]; !ok {
		*cfg.LocalizationSettings.DefaultServerLocale = model.DEFAULT_LOCALE
		changed = true
		mlog.Warn("DefaultServerLocale must be one of the supported locales. Setting DefaultServerLocale to en as default value.")
	}

	if _, ok := locales[*cfg.LocalizationSettings.DefaultClientLocale]; !ok {
		*cfg.LocalizationSettings.DefaultClientLocale = model.DEFAULT_LOCALE
		changed = true
		mlog.Warn("DefaultClientLocale must be one of the supported locales. Setting DefaultClientLocale to en as default value.")
	}

	if len(*cfg.LocalizationSettings.AvailableLocales) > 0 {
		isDefaultClientLocaleInAvailableLocales := false
		for _, word := range strings.Split(*cfg.LocalizationSettings.AvailableLocales, ",") {
			if _, ok := locales[word]; !ok {
				*cfg.LocalizationSettings.AvailableLocales = ""
				isDefaultClientLocaleInAvailableLocales = true
				changed = true
				mlog.Warn("AvailableLocales must include DefaultClientLocale. Setting AvailableLocales to all locales as default value.")
				break
			}

			if word == *cfg.LocalizationSettings.DefaultClientLocale {
				isDefaultClientLocaleInAvailableLocales = true
			}
		}

		availableLocales := *cfg.LocalizationSettings.AvailableLocales

		if !isDefaultClientLocaleInAvailableLocales {
			availableLocales += "," + *cfg.LocalizationSettings.DefaultClientLocale
			changed = true
			mlog.Warn("Adding DefaultClientLocale to AvailableLocales.")
		}

		*cfg.LocalizationSettings.AvailableLocales = strings.Join(utils.RemoveDuplicatesFromStringArray(strings.Split(availableLocales, ",")), ",")
	}

	return changed
}
