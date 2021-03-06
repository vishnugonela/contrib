/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rewrite

import (
	"errors"
	"strconv"

	"k8s.io/kubernetes/pkg/apis/extensions"
)

const (
	rewriteTo  = "ingress.kubernetes.io/rewrite-target"
	addBaseURL = "ingress.kubernetes.io/add-base-url"
)

// Redirect returns authentication configuration for an Ingress rule
type Redirect struct {
	// Target URI where the traffic must be redirected
	Target string
	// AddBaseURL indicates if is required to add a base tag in the head
	// of the responses from the upstream servers
	AddBaseURL bool
}

type ingAnnotations map[string]string

func (a ingAnnotations) addBaseURL() bool {
	val, ok := a[addBaseURL]
	if ok {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return false
}

func (a ingAnnotations) rewriteTo() string {
	val, ok := a[rewriteTo]
	if ok {
		return val
	}
	return ""
}

// ParseAnnotations parses the annotations contained in the ingress
// rule used to rewrite the defined paths
func ParseAnnotations(ing *extensions.Ingress) (*Redirect, error) {
	if ing.GetAnnotations() == nil {
		return &Redirect{}, errors.New("no annotations present")
	}

	rt := ingAnnotations(ing.GetAnnotations()).rewriteTo()
	abu := ingAnnotations(ing.GetAnnotations()).addBaseURL()
	return &Redirect{
		Target:     rt,
		AddBaseURL: abu,
	}, nil
}
