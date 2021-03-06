package clientcmd

import (
	"net/http"
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"

	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/diff"
	"k8s.io/kubernetes/pkg/util/flowcontrol"
)

func TestAnonymousConfig(t *testing.T) {
	f := fuzz.New().NilChance(0.0).NumElements(1, 1)
	f.Funcs(
		func(r *runtime.Codec, f fuzz.Continue) {},
		func(r *http.RoundTripper, f fuzz.Continue) {},
		func(fn *func(http.RoundTripper) http.RoundTripper, f fuzz.Continue) {},
		func(r *restclient.AuthProviderConfigPersister, f fuzz.Continue) {},
		func(r *runtime.NegotiatedSerializer, f fuzz.Continue) {},
		func(r *flowcontrol.RateLimiter, f fuzz.Continue) {},
		func(r *api.AuthProviderConfig, f fuzz.Continue) {
			r.Config = map[string]string{}
		},
	)
	for i := 0; i < 20; i++ {
		original := &restclient.Config{}
		f.Fuzz(original)
		actual := AnonymousClientConfig(original)
		expected := *original

		// this is the list of known security related fields, add to this list if a new field
		// is added to restclient.Config, update AnonymousClientConfig to preserve the field otherwise.
		expected.Impersonate = ""
		expected.BearerToken = ""
		expected.Username = ""
		expected.Password = ""
		expected.AuthProvider = nil
		expected.AuthConfigPersister = nil
		expected.TLSClientConfig.CertData = nil
		expected.TLSClientConfig.CertFile = ""
		expected.TLSClientConfig.KeyData = nil
		expected.TLSClientConfig.KeyFile = ""

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("AnonymousClientConfig dropped unexpected fields, identify whether they are security related or not: %s", diff.ObjectGoPrintDiff(expected, actual))
		}
	}
}
