// Code generated by counterfeiter. DO NOT EDIT.
package k8sfakes

import (
	"context"
	"sync"

	"github.com/maliksalman/spring-boot-scanner/k8s"
	v1 "k8s.io/api/core/v1"
)

type FakeSecretsProvider struct {
	GetSecretStub        func(context.Context, string, string) (*v1.Secret, error)
	getSecretMutex       sync.RWMutex
	getSecretArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}
	getSecretReturns struct {
		result1 *v1.Secret
		result2 error
	}
	getSecretReturnsOnCall map[int]struct {
		result1 *v1.Secret
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSecretsProvider) GetSecret(arg1 context.Context, arg2 string, arg3 string) (*v1.Secret, error) {
	fake.getSecretMutex.Lock()
	ret, specificReturn := fake.getSecretReturnsOnCall[len(fake.getSecretArgsForCall)]
	fake.getSecretArgsForCall = append(fake.getSecretArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetSecretStub
	fakeReturns := fake.getSecretReturns
	fake.recordInvocation("GetSecret", []interface{}{arg1, arg2, arg3})
	fake.getSecretMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSecretsProvider) GetSecretCallCount() int {
	fake.getSecretMutex.RLock()
	defer fake.getSecretMutex.RUnlock()
	return len(fake.getSecretArgsForCall)
}

func (fake *FakeSecretsProvider) GetSecretCalls(stub func(context.Context, string, string) (*v1.Secret, error)) {
	fake.getSecretMutex.Lock()
	defer fake.getSecretMutex.Unlock()
	fake.GetSecretStub = stub
}

func (fake *FakeSecretsProvider) GetSecretArgsForCall(i int) (context.Context, string, string) {
	fake.getSecretMutex.RLock()
	defer fake.getSecretMutex.RUnlock()
	argsForCall := fake.getSecretArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeSecretsProvider) GetSecretReturns(result1 *v1.Secret, result2 error) {
	fake.getSecretMutex.Lock()
	defer fake.getSecretMutex.Unlock()
	fake.GetSecretStub = nil
	fake.getSecretReturns = struct {
		result1 *v1.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeSecretsProvider) GetSecretReturnsOnCall(i int, result1 *v1.Secret, result2 error) {
	fake.getSecretMutex.Lock()
	defer fake.getSecretMutex.Unlock()
	fake.GetSecretStub = nil
	if fake.getSecretReturnsOnCall == nil {
		fake.getSecretReturnsOnCall = make(map[int]struct {
			result1 *v1.Secret
			result2 error
		})
	}
	fake.getSecretReturnsOnCall[i] = struct {
		result1 *v1.Secret
		result2 error
	}{result1, result2}
}

func (fake *FakeSecretsProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getSecretMutex.RLock()
	defer fake.getSecretMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSecretsProvider) recordInvocation(key string, args []interface{}) {
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

var _ k8s.SecretsProvider = new(FakeSecretsProvider)
