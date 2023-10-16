// Code generated by counterfeiter. DO NOT EDIT.
package k8sfakes

import (
	"context"
	"sync"

	"github.com/maliksalman/spring-boot-scanner/k8s"
	v1a "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v1b "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FakeAppsProvider struct {
	GetServiceAccountStub        func(context.Context, string, string) (*v1.ServiceAccount, error)
	getServiceAccountMutex       sync.RWMutex
	getServiceAccountArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}
	getServiceAccountReturns struct {
		result1 *v1.ServiceAccount
		result2 error
	}
	getServiceAccountReturnsOnCall map[int]struct {
		result1 *v1.ServiceAccount
		result2 error
	}
	ListDeploymentsStub        func(context.Context, string, v1b.ListOptions) (*v1a.DeploymentList, error)
	listDeploymentsMutex       sync.RWMutex
	listDeploymentsArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 v1b.ListOptions
	}
	listDeploymentsReturns struct {
		result1 *v1a.DeploymentList
		result2 error
	}
	listDeploymentsReturnsOnCall map[int]struct {
		result1 *v1a.DeploymentList
		result2 error
	}
	ListNamespacesStub        func(context.Context, v1b.ListOptions) (*v1.NamespaceList, error)
	listNamespacesMutex       sync.RWMutex
	listNamespacesArgsForCall []struct {
		arg1 context.Context
		arg2 v1b.ListOptions
	}
	listNamespacesReturns struct {
		result1 *v1.NamespaceList
		result2 error
	}
	listNamespacesReturnsOnCall map[int]struct {
		result1 *v1.NamespaceList
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAppsProvider) GetServiceAccount(arg1 context.Context, arg2 string, arg3 string) (*v1.ServiceAccount, error) {
	fake.getServiceAccountMutex.Lock()
	ret, specificReturn := fake.getServiceAccountReturnsOnCall[len(fake.getServiceAccountArgsForCall)]
	fake.getServiceAccountArgsForCall = append(fake.getServiceAccountArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetServiceAccountStub
	fakeReturns := fake.getServiceAccountReturns
	fake.recordInvocation("GetServiceAccount", []interface{}{arg1, arg2, arg3})
	fake.getServiceAccountMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAppsProvider) GetServiceAccountCallCount() int {
	fake.getServiceAccountMutex.RLock()
	defer fake.getServiceAccountMutex.RUnlock()
	return len(fake.getServiceAccountArgsForCall)
}

func (fake *FakeAppsProvider) GetServiceAccountCalls(stub func(context.Context, string, string) (*v1.ServiceAccount, error)) {
	fake.getServiceAccountMutex.Lock()
	defer fake.getServiceAccountMutex.Unlock()
	fake.GetServiceAccountStub = stub
}

func (fake *FakeAppsProvider) GetServiceAccountArgsForCall(i int) (context.Context, string, string) {
	fake.getServiceAccountMutex.RLock()
	defer fake.getServiceAccountMutex.RUnlock()
	argsForCall := fake.getServiceAccountArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeAppsProvider) GetServiceAccountReturns(result1 *v1.ServiceAccount, result2 error) {
	fake.getServiceAccountMutex.Lock()
	defer fake.getServiceAccountMutex.Unlock()
	fake.GetServiceAccountStub = nil
	fake.getServiceAccountReturns = struct {
		result1 *v1.ServiceAccount
		result2 error
	}{result1, result2}
}

func (fake *FakeAppsProvider) GetServiceAccountReturnsOnCall(i int, result1 *v1.ServiceAccount, result2 error) {
	fake.getServiceAccountMutex.Lock()
	defer fake.getServiceAccountMutex.Unlock()
	fake.GetServiceAccountStub = nil
	if fake.getServiceAccountReturnsOnCall == nil {
		fake.getServiceAccountReturnsOnCall = make(map[int]struct {
			result1 *v1.ServiceAccount
			result2 error
		})
	}
	fake.getServiceAccountReturnsOnCall[i] = struct {
		result1 *v1.ServiceAccount
		result2 error
	}{result1, result2}
}

func (fake *FakeAppsProvider) ListDeployments(arg1 context.Context, arg2 string, arg3 v1b.ListOptions) (*v1a.DeploymentList, error) {
	fake.listDeploymentsMutex.Lock()
	ret, specificReturn := fake.listDeploymentsReturnsOnCall[len(fake.listDeploymentsArgsForCall)]
	fake.listDeploymentsArgsForCall = append(fake.listDeploymentsArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 v1b.ListOptions
	}{arg1, arg2, arg3})
	stub := fake.ListDeploymentsStub
	fakeReturns := fake.listDeploymentsReturns
	fake.recordInvocation("ListDeployments", []interface{}{arg1, arg2, arg3})
	fake.listDeploymentsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAppsProvider) ListDeploymentsCallCount() int {
	fake.listDeploymentsMutex.RLock()
	defer fake.listDeploymentsMutex.RUnlock()
	return len(fake.listDeploymentsArgsForCall)
}

func (fake *FakeAppsProvider) ListDeploymentsCalls(stub func(context.Context, string, v1b.ListOptions) (*v1a.DeploymentList, error)) {
	fake.listDeploymentsMutex.Lock()
	defer fake.listDeploymentsMutex.Unlock()
	fake.ListDeploymentsStub = stub
}

func (fake *FakeAppsProvider) ListDeploymentsArgsForCall(i int) (context.Context, string, v1b.ListOptions) {
	fake.listDeploymentsMutex.RLock()
	defer fake.listDeploymentsMutex.RUnlock()
	argsForCall := fake.listDeploymentsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeAppsProvider) ListDeploymentsReturns(result1 *v1a.DeploymentList, result2 error) {
	fake.listDeploymentsMutex.Lock()
	defer fake.listDeploymentsMutex.Unlock()
	fake.ListDeploymentsStub = nil
	fake.listDeploymentsReturns = struct {
		result1 *v1a.DeploymentList
		result2 error
	}{result1, result2}
}

func (fake *FakeAppsProvider) ListDeploymentsReturnsOnCall(i int, result1 *v1a.DeploymentList, result2 error) {
	fake.listDeploymentsMutex.Lock()
	defer fake.listDeploymentsMutex.Unlock()
	fake.ListDeploymentsStub = nil
	if fake.listDeploymentsReturnsOnCall == nil {
		fake.listDeploymentsReturnsOnCall = make(map[int]struct {
			result1 *v1a.DeploymentList
			result2 error
		})
	}
	fake.listDeploymentsReturnsOnCall[i] = struct {
		result1 *v1a.DeploymentList
		result2 error
	}{result1, result2}
}

func (fake *FakeAppsProvider) ListNamespaces(arg1 context.Context, arg2 v1b.ListOptions) (*v1.NamespaceList, error) {
	fake.listNamespacesMutex.Lock()
	ret, specificReturn := fake.listNamespacesReturnsOnCall[len(fake.listNamespacesArgsForCall)]
	fake.listNamespacesArgsForCall = append(fake.listNamespacesArgsForCall, struct {
		arg1 context.Context
		arg2 v1b.ListOptions
	}{arg1, arg2})
	stub := fake.ListNamespacesStub
	fakeReturns := fake.listNamespacesReturns
	fake.recordInvocation("ListNamespaces", []interface{}{arg1, arg2})
	fake.listNamespacesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAppsProvider) ListNamespacesCallCount() int {
	fake.listNamespacesMutex.RLock()
	defer fake.listNamespacesMutex.RUnlock()
	return len(fake.listNamespacesArgsForCall)
}

func (fake *FakeAppsProvider) ListNamespacesCalls(stub func(context.Context, v1b.ListOptions) (*v1.NamespaceList, error)) {
	fake.listNamespacesMutex.Lock()
	defer fake.listNamespacesMutex.Unlock()
	fake.ListNamespacesStub = stub
}

func (fake *FakeAppsProvider) ListNamespacesArgsForCall(i int) (context.Context, v1b.ListOptions) {
	fake.listNamespacesMutex.RLock()
	defer fake.listNamespacesMutex.RUnlock()
	argsForCall := fake.listNamespacesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAppsProvider) ListNamespacesReturns(result1 *v1.NamespaceList, result2 error) {
	fake.listNamespacesMutex.Lock()
	defer fake.listNamespacesMutex.Unlock()
	fake.ListNamespacesStub = nil
	fake.listNamespacesReturns = struct {
		result1 *v1.NamespaceList
		result2 error
	}{result1, result2}
}

func (fake *FakeAppsProvider) ListNamespacesReturnsOnCall(i int, result1 *v1.NamespaceList, result2 error) {
	fake.listNamespacesMutex.Lock()
	defer fake.listNamespacesMutex.Unlock()
	fake.ListNamespacesStub = nil
	if fake.listNamespacesReturnsOnCall == nil {
		fake.listNamespacesReturnsOnCall = make(map[int]struct {
			result1 *v1.NamespaceList
			result2 error
		})
	}
	fake.listNamespacesReturnsOnCall[i] = struct {
		result1 *v1.NamespaceList
		result2 error
	}{result1, result2}
}

func (fake *FakeAppsProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getServiceAccountMutex.RLock()
	defer fake.getServiceAccountMutex.RUnlock()
	fake.listDeploymentsMutex.RLock()
	defer fake.listDeploymentsMutex.RUnlock()
	fake.listNamespacesMutex.RLock()
	defer fake.listNamespacesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAppsProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ k8s.AppsProvider = new(FakeAppsProvider)