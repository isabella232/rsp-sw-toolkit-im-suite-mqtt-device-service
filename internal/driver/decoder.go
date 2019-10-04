/*
 * INTEL CONFIDENTIAL
 * Copyright (2019) Intel Corporation.
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
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.impcloud.net/RSP-Inventory-Suite/tagcode/bittag"
	"github.impcloud.net/RSP-Inventory-Suite/tagcode/epc"
	"strings"
)

type TagDecoder func(tagData []byte) (URI string, err error)

type NamedDecoder struct {
	Name string
	TagDecoder
}

type DecoderRing struct {
	Decoders []NamedDecoder
}

func (dr *DecoderRing) AddBitTagDecoder(authority, date string, widths []int) error {
	btd, err := bittag.NewDecoder(authority, date, widths)
	if err != nil {
		return err
	}

	decoder := func(tagData []byte) (string, error) {
		bitTag, err := btd.Decode(tagData)
		if err != nil {
			return "", err
		}
		return bitTag.URI(), nil
	}

	dr.Decoders = append(dr.Decoders, NamedDecoder{Name: btd.Prefix(), TagDecoder: decoder})
	return nil
}

func (dr *DecoderRing) AddSGTINDecoder(strict bool) {
	decoder := func(tagData []byte) (URI string, err error) {
		var s epc.SGTIN
		s, err = epc.DecodeSGTIN(tagData)
		if strict && err == nil {
			err = s.ValidateRanges()
		}
		if err == nil {
			URI = s.URI()
		}
		return
	}

	dr.Decoders = append(dr.Decoders, NamedDecoder{Name: "SGTIN", TagDecoder: decoder})
}

func (dr *DecoderRing) TagDataToURI(tagData string) (string, error) {
	tagDataBytes, err := hex.DecodeString(tagData)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode tag hex data")
	}

	var decodingErrors []string
	for _, decoder := range dr.Decoders {
		tagUri, err := decoder.TagDecoder(tagDataBytes)
		if err == nil {
			return tagUri, nil
		}
		decodingErrors = append(decodingErrors, fmt.Sprintf("%s: %v",
			decoder.Name, err))
	}

	return "", errors.Errorf("no decoder successfully decoded the tag "+
		"data: individual decoder errors are as follows:\n%s",
		strings.Join(decodingErrors, "\n"))
}
