// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

/*
 * INTEL CONFIDENTIAL
 * Copyright (2017) Intel Corporation.
 *
 * The source code contained or described herein and all documents related to the source code ("Material")
 * are owned by Intel Corporation or its suppliers or licensors. Title to the Material remains with
 * Intel Corporation or its suppliers and licensors. The Material may contain trade secrets and proprietary
 * and confidential information of Intel Corporation and its suppliers and licensors, and is protected by
 * worldwide copyright and trade secret laws and treaty provisions. No part of the Material may be used,
 * copied, reproduced, modified, published, uploaded, posted, transmitted, distributed, or disclosed in
 * any way without Intel/'s prior express written permission.
 * No license under any patent, copyright, trade secret or other intellectual property right is granted
 * to or conferred upon you by disclosure or delivery of the Materials, either expressly, by implication,
 * inducement, estoppel or otherwise. Any license under such intellectual property rights must be express
 * and approved by Intel in writing.
 * Unless otherwise agreed by Intel in writing, you may not remove or alter this notice or any other
 * notice embedded in Materials by Intel or Intel's suppliers or licensors in any way.
 */

package driver

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Intel added Command field to the struct
type configuration struct {
	Incoming SubscribeInfo
	Command  SubscribeInfo
	Response SubscribeInfo
}

type SubscribeInfo struct {
	Protocol     string
	Host         string
	Port         int
	Username     string
	Password     string
	Qos          int
	KeepAlive    int
	MqttClientId string
	Topics       []string
	MetaDataPort int
}

// LoadConfigFromFile use to load toml configuration
func LoadConfigFromFile() (*configuration, error) {
	config := new(configuration)

	confDir := flag.Lookup("confdir").Value.(flag.Getter).Get().(string)
	if confDir == "" {
		confDir = flag.Lookup("c").Value.(flag.Getter).Get().(string)
	}

	if confDir == "" {
		confDir = "./res"
	}

	filePath := fmt.Sprintf("%v/configuration-driver.toml", confDir)

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("could not load configuration file (%s): %v", filePath, err.Error())
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return config, fmt.Errorf("unable to parse configuration file (%s): %v", filePath, err.Error())
	}
	return config, err
}
