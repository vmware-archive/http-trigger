package controller

import (
	"testing"
	"time"

	httpTriggerApi "github.com/kubeless/http-trigger/pkg/apis/kubeless/v1beta1"
	httpTriggerFake "github.com/kubeless/http-trigger/pkg/client/clientset/versioned/fake"
	kubelessApi "github.com/kubeless/kubeless/pkg/apis/kubeless/v1beta1"
	"github.com/sirupsen/logrus"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestHTTPFunctionAddedUpdated(t *testing.T) {
	myNsFoo := metav1.ObjectMeta{
		Namespace: "myns",
		Name:      "foo",
	}

	f := kubelessApi.Function{
		ObjectMeta: myNsFoo,
	}

	httpTrigger := httpTriggerApi.HTTPTrigger{
		ObjectMeta: myNsFoo,
	}

	triggerClientset := httpTriggerFake.NewSimpleClientset(&httpTrigger)

	ingress := extensionsv1beta1.Ingress{
		ObjectMeta: myNsFoo,
	}
	clientset := fake.NewSimpleClientset(&ingress)

	controller := HTTPTriggerController{
		clientset:  clientset,
		httpclient: triggerClientset,
		logger:     logrus.WithField("controller", "http-trigger-controller"),
	}

	// no-op for when the function is not deleted
	err := controller.functionAddedDeletedUpdated(&f, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	list, err := controller.httpclient.KubelessV1beta1().HTTPTriggers("myns").List(metav1.ListOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].ObjectMeta.Name != "foo" {
		t.Errorf("Missing trigger in list: %v", list.Items)
	}
}

func TestHTTPFunctionDeleted(t *testing.T) {
	myNsFoo := metav1.ObjectMeta{
		Namespace: "myns",
		Name:      "foo",
	}

	f := kubelessApi.Function{
		ObjectMeta: myNsFoo,
	}

	httpTrigger := httpTriggerApi.HTTPTrigger{
		ObjectMeta: myNsFoo,
		Spec: httpTriggerApi.HTTPTriggerSpec{
			FunctionName: myNsFoo.Name,
		},
	}

	triggerClientset := httpTriggerFake.NewSimpleClientset(&httpTrigger)

	ingress := extensionsv1beta1.Ingress{
		ObjectMeta: myNsFoo,
	}
	clientset := fake.NewSimpleClientset(&ingress)

	controller := HTTPTriggerController{
		clientset:  clientset,
		httpclient: triggerClientset,
		logger:     logrus.WithField("controller", "http-trigger-controller"),
	}

	err := controller.functionAddedDeletedUpdated(&f, true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	list, err := controller.httpclient.KubelessV1beta1().HTTPTriggers("myns").List(metav1.ListOptions{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(list.Items) != 0 {
		t.Errorf("Trigger should be deleted from list: %v", list.Items)
	}
}

func TestHTTPTriggerObjChanged(t *testing.T) {
	type testObj struct {
		old             *httpTriggerApi.HTTPTrigger
		new             *httpTriggerApi.HTTPTrigger
		expectedChanged bool
	}
	t1 := metav1.Time{
		Time: time.Now(),
	}
	t2 := metav1.Time{
		Time: time.Now(),
	}
	testObjs := []testObj{
		{
			old:             &httpTriggerApi.HTTPTrigger{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			new:             &httpTriggerApi.HTTPTrigger{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
			expectedChanged: false,
		},
		{
			old:             &httpTriggerApi.HTTPTrigger{ObjectMeta: metav1.ObjectMeta{DeletionTimestamp: &t1}},
			new:             &httpTriggerApi.HTTPTrigger{ObjectMeta: metav1.ObjectMeta{DeletionTimestamp: &t2}},
			expectedChanged: true,
		},
		{
			old:             &httpTriggerApi.HTTPTrigger{ObjectMeta: metav1.ObjectMeta{ResourceVersion: "1"}},
			new:             &httpTriggerApi.HTTPTrigger{ObjectMeta: metav1.ObjectMeta{ResourceVersion: "2"}},
			expectedChanged: true,
		},
		{
			old:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{HostName: "a"}},
			new:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{HostName: "a"}},
			expectedChanged: false,
		},
		{
			old:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{HostName: "a"}},
			new:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{HostName: "b"}},
			expectedChanged: true,
		},
		{
			old:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{TLSAcme: true}},
			new:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{TLSAcme: false}},
			expectedChanged: true,
		},
		{
			old:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{Path: "a"}},
			new:             &httpTriggerApi.HTTPTrigger{Spec: httpTriggerApi.HTTPTriggerSpec{Path: "b"}},
			expectedChanged: true,
		},
	}
	for _, to := range testObjs {
		changed := httpTriggerObjChanged(to.old, to.new)
		if changed != to.expectedChanged {
			t.Errorf("%v != %v expected to be %v", to.old, to.new, to.expectedChanged)
		}
	}
}
