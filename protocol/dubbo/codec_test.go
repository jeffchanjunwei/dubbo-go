/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dubbo

import (
	"bytes"
	"testing"
	"time"
)

import (
	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go/protocol/dubbo/hessian2"
)

func TestDubboPackageMarshalAndUnmarshal(t *testing.T) {
	pkg := &DubboPackage{}
	pkg.Body = []interface{}{"a"}
	pkg.Header.Type = hessian2.PackageHeartbeat
	pkg.Header.SerialID = byte(S_Dubbo)
	pkg.Header.ID = 10086

	// heartbeat
	data, err := pkg.Marshal()
	assert.NoError(t, err)

	pkgres := &DubboPackage{}
	pkgres.Body = []interface{}{}
	err = pkgres.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, hessian2.PackageHeartbeat|hessian2.PackageRequest|hessian2.PackageRequest_TwoWay, pkgres.Header.Type)
	assert.Equal(t, byte(S_Dubbo), pkgres.Header.SerialID)
	assert.Equal(t, int64(10086), pkgres.Header.ID)
	assert.Equal(t, 0, len(pkgres.Body.([]interface{})))

	// request
	pkg.Header.Type = hessian2.PackageRequest
	pkg.Service.Interface = "Service"
	pkg.Service.Path = "path"
	pkg.Service.Version = "2.6"
	pkg.Service.Method = "Method"
	pkg.Service.Timeout = time.Second
	data, err = pkg.Marshal()
	assert.NoError(t, err)

	pkgres = &DubboPackage{}
	pkgres.Body = make([]interface{}, 7)
	err = pkgres.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, hessian2.PackageRequest, pkgres.Header.Type)
	assert.Equal(t, byte(S_Dubbo), pkgres.Header.SerialID)
	assert.Equal(t, int64(10086), pkgres.Header.ID)
	assert.Equal(t, "2.0.2", pkgres.Body.([]interface{})[0])
	assert.Equal(t, "path", pkgres.Body.([]interface{})[1])
	assert.Equal(t, "2.6", pkgres.Body.([]interface{})[2])
	assert.Equal(t, "Method", pkgres.Body.([]interface{})[3])
	assert.Equal(t, "Ljava/lang/String;", pkgres.Body.([]interface{})[4])
	assert.Equal(t, []interface{}{"a"}, pkgres.Body.([]interface{})[5])
	assert.Equal(t, map[string]interface{}{"dubbo": "2.0.2", "interface": "Service", "path": "path", "timeout": "1000", "version": "2.6"}, pkgres.Body.([]interface{})[6])
}

func TestIssue380(t *testing.T) {
	pkg := &DubboPackage{}
	buf := bytes.NewBuffer([]byte("hello"))
	err := pkg.Unmarshal(buf)
	assert.True(t, perrors.Cause(err) == hessian2.ErrHeaderNotEnough)
}
