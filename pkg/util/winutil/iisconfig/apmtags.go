// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.
//go:build windows

package iisconfig

/*
the file datadog.json can be located anywhere; it is path-relative to a .net application
give the path name, read the json and return it as a map of string/string
*/

import (
	"encoding/json"
	"encoding/xml"
	"os"
)

type APMTags struct {
	DDService string
	DDEnv     string
	DDVersion string
}

func ReadDatadogJson(datadogJsonPath string) (APMTags, error) {
	var datadogJson map[string]string
	var apmtags APMTags

	file, err := os.Open(datadogJsonPath)
	if err != nil {
		return apmtags, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&datadogJson)
	if err != nil {
		return apmtags, err
	}
	apmtags.DDService = datadogJson["DD_SERVICE"]
	apmtags.DDEnv = datadogJson["DD_ENV"]
	apmtags.DDVersion = datadogJson["DD_VERSION"]
	return apmtags, nil
}

type iisAppSetting struct {
	Key   string `xml:"key,attr"`
	Value string `xml:"value,attr"`
}
type iisAppSettings struct {
	XMLName xml.Name        `xml:"appSettings"`
	Adds    []iisAppSetting `xml:"add"`
}

type appConfiguration struct {
	XMLName     xml.Name `xml:"configuration"`
	AppSettings iisAppSettings
}

func ReadDotNetConfig(cfgpath string) (APMTags, error) { //(APMTags, error) {
	var newcfg appConfiguration
	var apmtags APMTags
	f, err := os.ReadFile(cfgpath)
	if err != nil {
		return apmtags, err
	}
	err = xml.Unmarshal(f, &newcfg)
	if err != nil {
		return apmtags, err
	}
	for _, setting := range newcfg.AppSettings.Adds {
		switch setting.Key {
		case "DD_SERVICE":
			apmtags.DDService = setting.Value
		case "DD_ENV":
			apmtags.DDEnv = setting.Value
		case "DD_VERSION":
			apmtags.DDVersion = setting.Value
		}
	}
	return apmtags, nil
}
