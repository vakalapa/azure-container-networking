// Copyright 2018 Microsoft. All rights reserved.
// MIT License
package npm

import (
	"os"
	"reflect"
	"testing"

	"github.com/Azure/azure-container-networking/npm/iptm"
	"github.com/Azure/azure-container-networking/npm/metrics"

	"github.com/Azure/azure-container-networking/npm/ipsm"
	"github.com/Azure/azure-container-networking/npm/util"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewNs(t *testing.T) {
	if _, err := newNs("test"); err != nil {
		t.Errorf("TestnewNs failed @ newNs")
	}
}

func TestAllNsList(t *testing.T) {
	npMgr := &NetworkPolicyManager{}

	ipsMgr := ipsm.NewIpsetManager()
	if err := ipsMgr.Save(util.IpsetTestConfigFile); err != nil {
		t.Errorf("TestAllNsList failed @ ipsMgr.Save")
	}

	defer func() {
		if err := ipsMgr.Restore(util.IpsetTestConfigFile); err != nil {
			t.Errorf("TestAllNsList failed @ ipsMgr.Restore")
		}
	}()

	if err := npMgr.InitAllNsList(); err != nil {
		t.Errorf("TestAllNsList failed @ InitAllNsList")
	}

	if err := npMgr.UninitAllNsList(); err != nil {
		t.Errorf("TestAllNsList failed @ UninitAllNsList")
	}
}

func TestAddNamespace(t *testing.T) {
	npMgr := &NetworkPolicyManager{
		NsMap:            make(map[string]*Namespace),
		PodMap:           make(map[string]*NpmPod),
		RawNpMap:         make(map[string]*networkingv1.NetworkPolicy),
		ProcessedNpMap:   make(map[string]*networkingv1.NetworkPolicy),
		TelemetryEnabled: false,
	}

	allNs, err := newNs(util.KubeAllNamespacesFlag)
	if err != nil {
		panic(err.Error)
	}
	npMgr.NsMap[util.KubeAllNamespacesFlag] = allNs

	ipsMgr := ipsm.NewIpsetManager()
	if err := ipsMgr.Save(util.IpsetTestConfigFile); err != nil {
		t.Errorf("TestAddNamespace failed @ ipsMgr.Save")
	}

	defer func() {
		if err := ipsMgr.Restore(util.IpsetTestConfigFile); err != nil {
			t.Errorf("TestAddNamespace failed @ ipsMgr.Restore")
		}
	}()

	nsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
			Labels: map[string]string{
				"app": "test-namespace",
			},
		},
	}

	npMgr.Lock()
	if err := npMgr.AddNamespace(nsObj); err != nil {
		t.Errorf("TestAddNamespace @ npMgr.AddNamespace")
	}
	npMgr.Unlock()
}

func TestUpdateNamespace(t *testing.T) {
	npMgr := &NetworkPolicyManager{
		NsMap:            make(map[string]*Namespace),
		PodMap:           make(map[string]*NpmPod),
		RawNpMap:         make(map[string]*networkingv1.NetworkPolicy),
		ProcessedNpMap:   make(map[string]*networkingv1.NetworkPolicy),
		TelemetryEnabled: false,
	}

	allNs, err := newNs(util.KubeAllNamespacesFlag)
	if err != nil {
		panic(err.Error)
	}
	npMgr.NsMap[util.KubeAllNamespacesFlag] = allNs

	ipsMgr := ipsm.NewIpsetManager()
	if err := ipsMgr.Save(util.IpsetTestConfigFile); err != nil {
		t.Errorf("TestUpdateNamespace failed @ ipsMgr.Save")
	}

	defer func() {
		if err := ipsMgr.Restore(util.IpsetTestConfigFile); err != nil {
			t.Errorf("TestUpdateNamespace failed @ ipsMgr.Restore")
		}
	}()

	oldNsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "old-test-namespace",
			Labels: map[string]string{
				"app": "old-test-namespace",
			},
		},
	}

	newNsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "new-test-namespace",
			Labels: map[string]string{
				"app": "new-test-namespace",
			},
		},
	}

	npMgr.Lock()
	if err := npMgr.AddNamespace(oldNsObj); err != nil {
		t.Errorf("TestUpdateNamespace failed @ npMgr.AddNamespace")
	}

	if err := npMgr.UpdateNamespace(oldNsObj, newNsObj); err != nil {
		t.Errorf("TestUpdateNamespace failed @ npMgr.UpdateNamespace")
	}
	npMgr.Unlock()
}

func TestAddNamespaceLabel(t *testing.T) {
	npMgr := &NetworkPolicyManager{
		NsMap:            make(map[string]*Namespace),
		PodMap:           make(map[string]*NpmPod),
		RawNpMap:         make(map[string]*networkingv1.NetworkPolicy),
		ProcessedNpMap:   make(map[string]*networkingv1.NetworkPolicy),
		TelemetryEnabled: false,
	}

	allNs, err := newNs(util.KubeAllNamespacesFlag)
	if err != nil {
		t.Fatal(err.Error())
	}
	npMgr.NsMap[util.KubeAllNamespacesFlag] = allNs

	ipsMgr := ipsm.NewIpsetManager()
	if err := ipsMgr.Save(util.IpsetTestConfigFile); err != nil {
		t.Errorf("TestAddNamespaceLabel failed @ ipsMgr.Save")
	}

	defer func() {
		if err := ipsMgr.Restore(util.IpsetTestConfigFile); err != nil {
			t.Errorf("TestAddNamespaceLabel failed @ ipsMgr.Restore")
		}
	}()

	oldNsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "old-test-namespace",
			Labels: map[string]string{
				"app": "old-test-namespace",
			},
			ResourceVersion: "0",
		},
	}

	newNsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "old-test-namespace",
			Labels: map[string]string{
				"app":    "old-test-namespace",
				"update": "true",
			},

			ResourceVersion: "1",
		},
	}

	npMgr.Lock()
	if err := npMgr.AddNamespace(oldNsObj); err != nil {
		t.Errorf("TestAddNamespaceLabel failed @ npMgr.AddNamespace")
	}

	if err := npMgr.UpdateNamespace(oldNsObj, newNsObj); err != nil {
		t.Errorf("TestAddNamespaceLabel failed @ npMgr.UpdateNamespace")
	}

	if !reflect.DeepEqual(npMgr.NsMap["ns-"+newNsObj.Name].LabelsMap, newNsObj.ObjectMeta.Labels) {
		t.Errorf("TestAddNamespaceLabel failed @ npMgr.nsMap labelMap check")
	}

	npMgr.Unlock()
}

func TestDeleteandUpdateNamespaceLabel(t *testing.T) {
	npMgr := &NetworkPolicyManager{
		NsMap:            make(map[string]*Namespace),
		PodMap:           make(map[string]*NpmPod),
		RawNpMap:         make(map[string]*networkingv1.NetworkPolicy),
		ProcessedNpMap:   make(map[string]*networkingv1.NetworkPolicy),
		TelemetryEnabled: false,
	}

	allNs, err := newNs(util.KubeAllNamespacesFlag)
	if err != nil {
		t.Fatal(err.Error())
	}
	npMgr.NsMap[util.KubeAllNamespacesFlag] = allNs

	ipsMgr := ipsm.NewIpsetManager()
	if err := ipsMgr.Save(util.IpsetTestConfigFile); err != nil {
		t.Errorf("TestDeleteandUpdateNamespaceLabel failed @ ipsMgr.Save")
	}

	defer func() {
		if err := ipsMgr.Restore(util.IpsetTestConfigFile); err != nil {
			t.Errorf("TestDeleteandUpdateNamespaceLabel failed @ ipsMgr.Restore")
		}
	}()

	oldNsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "old-test-namespace",
			Labels: map[string]string{
				"app":    "old-test-namespace",
				"update": "true",
				"group":  "test",
			},
			ResourceVersion: "0",
		},
	}

	newNsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "old-test-namespace",
			Labels: map[string]string{
				"app":    "old-test-namespace",
				"update": "false",
			},
			ResourceVersion: "1",
		},
	}

	npMgr.Lock()
	if err := npMgr.AddNamespace(oldNsObj); err != nil {
		t.Errorf("TestDeleteandUpdateNamespaceLabel failed @ npMgr.AddNamespace")
	}

	if err := npMgr.UpdateNamespace(oldNsObj, newNsObj); err != nil {
		t.Errorf("TestDeleteandUpdateNamespaceLabel failed @ npMgr.UpdateNamespace")
	}

	if !reflect.DeepEqual(npMgr.NsMap["ns-"+newNsObj.Name].LabelsMap, newNsObj.ObjectMeta.Labels) {
		t.Errorf("TestDeleteandUpdateNamespaceLabel failed @ npMgr.nsMap labelMap check")
	}
	npMgr.Unlock()
}

func TestDeleteNamespace(t *testing.T) {
	npMgr := &NetworkPolicyManager{
		NsMap:            make(map[string]*Namespace),
		PodMap:           make(map[string]*NpmPod),
		RawNpMap:         make(map[string]*networkingv1.NetworkPolicy),
		ProcessedNpMap:   make(map[string]*networkingv1.NetworkPolicy),
		TelemetryEnabled: false,
	}

	allNs, err := newNs(util.KubeAllNamespacesFlag)
	if err != nil {
		panic(err.Error)
	}
	npMgr.NsMap[util.KubeAllNamespacesFlag] = allNs

	ipsMgr := ipsm.NewIpsetManager()
	if err := ipsMgr.Save(util.IpsetTestConfigFile); err != nil {
		t.Errorf("TestDeleteNamespace failed @ ipsMgr.Save")
	}

	defer func() {
		if err := ipsMgr.Restore(util.IpsetTestConfigFile); err != nil {
			t.Errorf("TestDeleteNamespace failed @ ipsMgr.Restore")
		}
	}()

	nsObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-namespace",
			Labels: map[string]string{
				"app": "test-namespace",
			},
		},
	}

	npMgr.Lock()
	if err := npMgr.AddNamespace(nsObj); err != nil {
		t.Errorf("TestDeleteNamespace @ npMgr.AddNamespace")
	}

	if err := npMgr.DeleteNamespace(nsObj); err != nil {
		t.Errorf("TestDeleteNamespace @ npMgr.DeleteNamespace")
	}

	if _, exists := npMgr.NsMap["ns-"+nsObj.Name]; exists {
		t.Errorf("TestDeleteNamespace failed @ npMgr.nsMap check")
	}
	npMgr.Unlock()
}

func TestMain(m *testing.M) {
	metrics.InitializeAll()
	iptMgr := iptm.NewIptablesManager()
	iptMgr.Save(util.IptablesConfigFile)

	ipsMgr := ipsm.NewIpsetManager()
	ipsMgr.Save(util.IpsetConfigFile)

	exitCode := m.Run()

	iptMgr.Restore(util.IptablesConfigFile)
	ipsMgr.Restore(util.IpsetConfigFile)

	os.Exit(exitCode)
}
