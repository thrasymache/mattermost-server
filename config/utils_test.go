package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/model"
)

func TestDesanitize(t *testing.T) {
	actual := &model.Config{}
	actual.SetDefaults()

	// These setting should be ignored
	actual.LdapSettings.Enable = bToP(false)
	actual.FileSettings.DriverName = sToP("s3")

	// These settings should be desanitized into target.
	actual.LdapSettings.BindPassword = sToP("bind_password")
	actual.FileSettings.PublicLinkSalt = sToP("public_link_salt")
	actual.FileSettings.AmazonS3SecretAccessKey = sToP("amazon_s3_secret_access_key")
	actual.EmailSettings.InviteSalt = sToP("invite_salt")
	actual.EmailSettings.SMTPPassword = sToP("smtp_password")
	actual.GitLabSettings.Secret = sToP("secret")
	actual.SqlSettings.DataSource = sToP("data_source")
	actual.SqlSettings.AtRestEncryptKey = sToP("at_rest_encrypt_key")
	actual.ElasticsearchSettings.Password = sToP("password")
	actual.SqlSettings.DataSourceReplicas = append(actual.SqlSettings.DataSourceReplicas, "replica0")
	actual.SqlSettings.DataSourceReplicas = append(actual.SqlSettings.DataSourceReplicas, "replica1")
	actual.SqlSettings.DataSourceSearchReplicas = append(actual.SqlSettings.DataSourceSearchReplicas, "search_replica0")
	actual.SqlSettings.DataSourceSearchReplicas = append(actual.SqlSettings.DataSourceSearchReplicas, "search_replica1")

	target := &model.Config{}
	target.SetDefaults()

	// These setting should be ignored
	target.LdapSettings.Enable = bToP(true)
	target.FileSettings.DriverName = sToP("file")

	// These settings should be updated from actual
	target.LdapSettings.BindPassword = sToP(model.FAKE_SETTING)
	target.FileSettings.PublicLinkSalt = sToP(model.FAKE_SETTING)
	target.FileSettings.AmazonS3SecretAccessKey = sToP(model.FAKE_SETTING)
	target.EmailSettings.InviteSalt = sToP(model.FAKE_SETTING)
	target.EmailSettings.SMTPPassword = sToP(model.FAKE_SETTING)
	target.GitLabSettings.Secret = sToP(model.FAKE_SETTING)
	target.SqlSettings.DataSource = sToP(model.FAKE_SETTING)
	target.SqlSettings.AtRestEncryptKey = sToP(model.FAKE_SETTING)
	target.ElasticsearchSettings.Password = sToP(model.FAKE_SETTING)
	target.SqlSettings.DataSourceReplicas = append(target.SqlSettings.DataSourceReplicas, "old_replica0")
	target.SqlSettings.DataSourceSearchReplicas = append(target.SqlSettings.DataSourceReplicas, "old_search_replica0")

	actual_clone := actual.Clone()
	desanitize(actual, target)
	assert.Equal(t, actual_clone, actual, "actual should not have been changed")

	// Verify the settings that should have been left untouched in target
	assert.True(t, *target.LdapSettings.Enable, "LdapSettings.Enable should not have changed")
	assert.Equal(t, "file", *target.FileSettings.DriverName, "FileSettings.DriverName should not have been changed")

	// Verify the settings that should have been desanitized into target
	assert.Equal(t, *actual.LdapSettings.BindPassword, *target.LdapSettings.BindPassword)
	assert.Equal(t, *actual.FileSettings.PublicLinkSalt, *target.FileSettings.PublicLinkSalt)
	assert.Equal(t, *actual.FileSettings.AmazonS3SecretAccessKey, *target.FileSettings.AmazonS3SecretAccessKey)
	assert.Equal(t, *actual.EmailSettings.InviteSalt, *target.EmailSettings.InviteSalt)
	assert.Equal(t, *actual.EmailSettings.SMTPPassword, *target.EmailSettings.SMTPPassword)
	assert.Equal(t, *actual.GitLabSettings.Secret, *target.GitLabSettings.Secret)
	assert.Equal(t, *actual.SqlSettings.DataSource, *target.SqlSettings.DataSource)
	assert.Equal(t, *actual.SqlSettings.AtRestEncryptKey, *target.SqlSettings.AtRestEncryptKey)
	assert.Equal(t, *actual.ElasticsearchSettings.Password, *target.ElasticsearchSettings.Password)
	assert.Equal(t, actual.SqlSettings.DataSourceReplicas, target.SqlSettings.DataSourceReplicas)
	assert.Equal(t, actual.SqlSettings.DataSourceSearchReplicas, target.SqlSettings.DataSourceSearchReplicas)
}

func sToP(s string) *string {
	return &s
}

func bToP(b bool) *bool {
	return &b
}

func iToP(i int) *int {
	return &i
}
