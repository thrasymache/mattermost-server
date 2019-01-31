package config

import (
	"github.com/mattermost/mattermost-server/model"
)

func desanitize(actualCfg, newCfg *model.Config) {
	if newCfg.LdapSettings.BindPassword != nil && *newCfg.LdapSettings.BindPassword == model.FAKE_SETTING {
		*newCfg.LdapSettings.BindPassword = *actualCfg.LdapSettings.BindPassword
	}

	if *newCfg.FileSettings.PublicLinkSalt == model.FAKE_SETTING {
		*newCfg.FileSettings.PublicLinkSalt = *actualCfg.FileSettings.PublicLinkSalt
	}
	if *newCfg.FileSettings.AmazonS3SecretAccessKey == model.FAKE_SETTING {
		newCfg.FileSettings.AmazonS3SecretAccessKey = actualCfg.FileSettings.AmazonS3SecretAccessKey
	}

	if *newCfg.EmailSettings.InviteSalt == model.FAKE_SETTING {
		newCfg.EmailSettings.InviteSalt = actualCfg.EmailSettings.InviteSalt
	}
	if *newCfg.EmailSettings.SMTPPassword == model.FAKE_SETTING {
		newCfg.EmailSettings.SMTPPassword = actualCfg.EmailSettings.SMTPPassword
	}

	if *newCfg.GitLabSettings.Secret == model.FAKE_SETTING {
		newCfg.GitLabSettings.Secret = actualCfg.GitLabSettings.Secret
	}

	if *newCfg.SqlSettings.DataSource == model.FAKE_SETTING {
		*newCfg.SqlSettings.DataSource = *actualCfg.SqlSettings.DataSource
	}
	if *newCfg.SqlSettings.AtRestEncryptKey == model.FAKE_SETTING {
		newCfg.SqlSettings.AtRestEncryptKey = actualCfg.SqlSettings.AtRestEncryptKey
	}

	if *newCfg.ElasticsearchSettings.Password == model.FAKE_SETTING {
		*newCfg.ElasticsearchSettings.Password = *actualCfg.ElasticsearchSettings.Password
	}

	for i := range newCfg.SqlSettings.DataSourceReplicas {
		newCfg.SqlSettings.DataSourceReplicas[i] = actualCfg.SqlSettings.DataSourceReplicas[i]
	}

	for i := range newCfg.SqlSettings.DataSourceSearchReplicas {
		newCfg.SqlSettings.DataSourceSearchReplicas[i] = actualCfg.SqlSettings.DataSourceSearchReplicas[i]
	}
}
