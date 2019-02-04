// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package config

import (
// "bytes"
// "fmt"
// "os"
// "strings"
// "testing"

// "github.com/stretchr/testify/assert"
// "github.com/stretchr/testify/require"

// "github.com/mattermost/mattermost-server/model"
// "github.com/mattermost/mattermost-server/utils"
)

// func TestConfig(t *testing.T) {
// 	utils.TranslationsPreInit()
// 	_, _, _, err := LoadConfig("config.json")
// 	require.Nil(t, err)
// }

// func TestValidateLocales(t *testing.T) {
// 	utils.TranslationsPreInit()
// 	cfg, _, _, err := LoadConfig("config.json")
// 	require.Nil(t, err)

// 	*cfg.LocalizationSettings.DefaultServerLocale = "en"
// 	*cfg.LocalizationSettings.DefaultClientLocale = "en"
// 	*cfg.LocalizationSettings.AvailableLocales = ""

// 	// t.Logf("*cfg.LocalizationSettings.DefaultClientLocale: %+v", *cfg.LocalizationSettings.DefaultClientLocale)
// 	if err := ValidateLocales(cfg); err != nil {
// 		t.Fatal("Should have not returned an error")
// 	}

// 	// validate DefaultServerLocale
// 	*cfg.LocalizationSettings.DefaultServerLocale = "junk"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if *cfg.LocalizationSettings.DefaultServerLocale != "en" {
// 			t.Fatal("DefaultServerLocale should have assigned to en as a default value")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating DefaultServerLocale")
// 	}

// 	*cfg.LocalizationSettings.DefaultServerLocale = ""
// 	if err := ValidateLocales(cfg); err != nil {
// 		if *cfg.LocalizationSettings.DefaultServerLocale != "en" {
// 			t.Fatal("DefaultServerLocale should have assigned to en as a default value")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating DefaultServerLocale")
// 	}

// 	*cfg.LocalizationSettings.AvailableLocales = "en"
// 	*cfg.LocalizationSettings.DefaultServerLocale = "de"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if strings.Contains(*cfg.LocalizationSettings.AvailableLocales, *cfg.LocalizationSettings.DefaultServerLocale) {
// 			t.Fatal("DefaultServerLocale should not be added to AvailableLocales")
// 		}
// 		t.Fatal("Should have not returned an error validating DefaultServerLocale")
// 	}

// 	// validate DefaultClientLocale
// 	*cfg.LocalizationSettings.AvailableLocales = ""
// 	*cfg.LocalizationSettings.DefaultClientLocale = "junk"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if *cfg.LocalizationSettings.DefaultClientLocale != "en" {
// 			t.Fatal("DefaultClientLocale should have assigned to en as a default value")
// 		}
// 	} else {

// 		t.Fatal("Should have returned an error validating DefaultClientLocale")
// 	}

// 	*cfg.LocalizationSettings.DefaultClientLocale = ""
// 	if err := ValidateLocales(cfg); err != nil {
// 		if *cfg.LocalizationSettings.DefaultClientLocale != "en" {
// 			t.Fatal("DefaultClientLocale should have assigned to en as a default value")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating DefaultClientLocale")
// 	}

// 	*cfg.LocalizationSettings.AvailableLocales = "en"
// 	*cfg.LocalizationSettings.DefaultClientLocale = "de"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if !strings.Contains(*cfg.LocalizationSettings.AvailableLocales, *cfg.LocalizationSettings.DefaultClientLocale) {
// 			t.Fatal("DefaultClientLocale should have added to AvailableLocales")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating DefaultClientLocale")
// 	}

// 	// validate AvailableLocales
// 	*cfg.LocalizationSettings.DefaultServerLocale = "en"
// 	*cfg.LocalizationSettings.DefaultClientLocale = "en"
// 	*cfg.LocalizationSettings.AvailableLocales = "junk"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if *cfg.LocalizationSettings.AvailableLocales != "" {
// 			t.Fatal("AvailableLocales should have assigned to empty string as a default value")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating AvailableLocales")
// 	}

// 	*cfg.LocalizationSettings.AvailableLocales = "en,de,junk"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if *cfg.LocalizationSettings.AvailableLocales != "" {
// 			t.Fatal("AvailableLocales should have assigned to empty string as a default value")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating AvailableLocales")
// 	}

// 	*cfg.LocalizationSettings.DefaultServerLocale = "fr"
// 	*cfg.LocalizationSettings.DefaultClientLocale = "de"
// 	*cfg.LocalizationSettings.AvailableLocales = "en"
// 	if err := ValidateLocales(cfg); err != nil {
// 		if strings.Contains(*cfg.LocalizationSettings.AvailableLocales, *cfg.LocalizationSettings.DefaultServerLocale) {
// 			t.Fatal("DefaultServerLocale should not be added to AvailableLocales")
// 		}
// 		if !strings.Contains(*cfg.LocalizationSettings.AvailableLocales, *cfg.LocalizationSettings.DefaultClientLocale) {
// 			t.Fatal("DefaultClientLocale should have added to AvailableLocales")
// 		}
// 	} else {
// 		t.Fatal("Should have returned an error validating AvailableLocales")
// 	}
// }
