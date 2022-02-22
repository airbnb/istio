//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package echo

import (
	"regexp"
	"strings"

	"istio.io/istio/pkg/test/echo/proto"
)

var (
	requestIDFieldRegex      = regexp.MustCompile("(?i)" + string(RequestIDField) + "=(.*)")
	serviceVersionFieldRegex = regexp.MustCompile(string(ServiceVersionField) + "=(.*)")
	servicePortFieldRegex    = regexp.MustCompile(string(ServicePortField) + "=(.*)")
	statusCodeFieldRegex     = regexp.MustCompile(string(StatusCodeField) + "=(.*)")
	hostFieldRegex           = regexp.MustCompile(string(HostField) + "=(.*)")
	hostnameFieldRegex       = regexp.MustCompile(string(HostnameField) + "=(.*)")
	responseHeaderFieldRegex = regexp.MustCompile(string(ResponseHeader) + "=(.*)")
	URLFieldRegex            = regexp.MustCompile(string(URLField) + "=(.*)")
	ClusterFieldRegex        = regexp.MustCompile(string(ClusterField) + "=(.*)")
	IstioVersionFieldRegex   = regexp.MustCompile(string(IstioVersionField) + "=(.*)")
	IPFieldRegex             = regexp.MustCompile(string(IPField) + "=(.*)")
)

func ParseResponses(req *proto.ForwardEchoRequest, resp *proto.ForwardEchoResponse) Responses {
	responses := make([]Response, len(resp.Output))
	for i, output := range resp.Output {
		responses[i] = parseResponse(output)
		responses[i].RequestURL = req.Url
	}
	return responses
}

func parseResponse(output string) Response {
	out := Response{
		Body: output,
	}

	match := requestIDFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.ID = match[1]
	}

	match = serviceVersionFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.Version = match[1]
	}

	match = servicePortFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.Port = match[1]
	}

	match = statusCodeFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.Code = match[1]
	}

	match = hostFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.Host = match[1]
	}

	match = hostnameFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.Hostname = match[1]
	}

	match = URLFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.URL = match[1]
	}

	match = ClusterFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.Cluster = match[1]
	}

	match = IstioVersionFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.IstioVersion = match[1]
	}

	match = IPFieldRegex.FindStringSubmatch(output)
	if match != nil {
		out.IP = match[1]
	}

	out.RawResponse = map[string]string{}

	matches := responseHeaderFieldRegex.FindAllStringSubmatch(output, -1)
	for _, kv := range matches {
		sl := strings.SplitN(kv[1], ":", 2)
		if len(sl) != 2 {
			continue
		}
		out.RawResponse[sl[0]] = sl[1]
	}

	for _, l := range strings.Split(output, "\n") {
		prefixSplit := strings.Split(l, "body] ")
		if len(prefixSplit) != 2 {
			continue
		}
		kv := strings.SplitN(prefixSplit[1], "=", 2)
		if len(kv) != 2 {
			continue
		}
		out.RawResponse[kv[0]] = kv[1]
	}

	return out
}
