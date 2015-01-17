package client

import (
	"errors"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	oapi "github.com/openshift/origin/pkg/oauth/api"
	"github.com/openshift/origin/pkg/oauth/registry/test"
)

func TestCreateValidationError(t *testing.T) {
	registry := test.ClientRegistry{}
	storage := REST{
		registry: &registry,
	}
	client := &oapi.Client{
	// ObjectMeta: api.ObjectMeta{Name: "authTokenName"}, // Missing required field
	}

	ctx := api.NewContext()
	_, err := storage.Create(ctx, client)
	if err == nil {
		t.Errorf("Expected validation error")
	}
}

func TestCreateStorageError(t *testing.T) {
	registry := test.ClientRegistry{
		Err: errors.New("Sample Error"),
	}
	storage := REST{
		registry: &registry,
	}
	client := &oapi.Client{
		ObjectMeta: api.ObjectMeta{Name: "clientName"},
	}

	ctx := api.NewContext()
	channel, err := storage.Create(ctx, client)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	select {
	case r := <-channel:
		switch r := r.Object.(type) {
		case *api.Status:
			if r.Message == registry.Err.Error() {
				// expected case
			} else {
				t.Errorf("Got back unexpected error: %#v", r)
			}
		default:
			t.Errorf("Got back non-status result: %v", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Error("Unexpected timeout from async channel")
	}
}

func TestCreateValid(t *testing.T) {
	registry := test.ClientRegistry{}
	storage := REST{
		registry: &registry,
	}
	client := &oapi.Client{
		ObjectMeta: api.ObjectMeta{Name: "clientName"},
	}

	ctx := api.NewContext()
	channel, err := storage.Create(ctx, client)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	select {
	case r := <-channel:
		switch r := r.Object.(type) {
		case *api.Status:
			t.Errorf("Got back unexpected status: %#v", r)
		case *oapi.Client:
			// expected case
		default:
			t.Errorf("Got unexpected type: %#v", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Error("Unexpected timeout from async channel")
	}
}

func TestGetError(t *testing.T) {
	registry := test.ClientRegistry{
		Err: errors.New("Sample Error"),
	}
	storage := REST{
		registry: &registry,
	}
	ctx := api.NewContext()
	_, err := storage.Get(ctx, "name")
	if err == nil {
		t.Errorf("expected error")
		return
	}
	if err != registry.Err {
		t.Errorf("got unexpected error: %v", err)
		return
	}
}

func TestGetValid(t *testing.T) {
	registry := test.ClientRegistry{
		Client: &oapi.Client{
			ObjectMeta: api.ObjectMeta{Name: "clientName"},
		},
	}
	storage := REST{
		registry: &registry,
	}
	ctx := api.NewContext()
	client, err := storage.Get(ctx, "name")
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
		return
	}
	if client != registry.Client {
		t.Errorf("got unexpected client: %v", client)
		return
	}
}

func TestListError(t *testing.T) {
	registry := test.ClientRegistry{
		Err: errors.New("Sample Error"),
	}
	storage := REST{
		registry: &registry,
	}
	ctx := api.NewContext()
	_, err := storage.List(ctx, labels.Everything(), labels.Everything())
	if err == nil {
		t.Errorf("expected error")
		return
	}
	if err != registry.Err {
		t.Errorf("got unexpected error: %v", err)
		return
	}
}

func TestListEmpty(t *testing.T) {
	registry := test.ClientRegistry{
		Clients: &oapi.ClientList{},
	}
	storage := REST{
		registry: &registry,
	}
	ctx := api.NewContext()
	clients, err := storage.List(ctx, labels.Everything(), labels.Everything())
	if err != registry.Err {
		t.Errorf("got unexpected error: %v", err)
		return
	}
	switch clients := clients.(type) {
	case *oapi.ClientList:
		if len(clients.Items) != 0 {
			t.Errorf("expected empty list, got %#v", clients)
		}
	default:
		t.Errorf("expected clientList, got: %v", clients)
		return
	}
}

func TestList(t *testing.T) {
	registry := test.ClientRegistry{
		Clients: &oapi.ClientList{
			Items: []oapi.Client{
				{},
				{},
			},
		},
	}
	storage := REST{
		registry: &registry,
	}
	ctx := api.NewContext()
	clients, err := storage.List(ctx, labels.Everything(), labels.Everything())
	if err != registry.Err {
		t.Errorf("got unexpected error: %v", err)
		return
	}
	switch clients := clients.(type) {
	case *oapi.ClientList:
		if len(clients.Items) != 2 {
			t.Errorf("expected list with 2 items, got %#v", clients)
		}
	default:
		t.Errorf("expected clientList, got: %v", clients)
		return
	}
}

func TestUpdateNotSupported(t *testing.T) {
	registry := test.ClientRegistry{
		Err: errors.New("Storage Error"),
	}
	storage := REST{
		registry: &registry,
	}
	client := &oapi.Client{
		ObjectMeta: api.ObjectMeta{Name: "clientName"},
	}

	ctx := api.NewContext()
	_, err := storage.Update(ctx, client)
	if err == nil {
		t.Errorf("expected unsupported error, but update succeeded")
		return
	}
	if err == registry.Err {
		t.Errorf("expected unsupported error, but registry was called")
		return
	}
}

func TestDeleteError(t *testing.T) {
	registry := test.ClientRegistry{
		Err: errors.New("Sample Error"),
	}
	storage := REST{
		registry: &registry,
	}

	ctx := api.NewContext()
	channel, err := storage.Delete(ctx, "foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	select {
	case r := <-channel:
		switch r := r.Object.(type) {
		case *api.Status:
			if r.Message == registry.Err.Error() {
				// expected case
			} else {
				t.Errorf("Got back unexpected error: %#v", r)
			}
		default:
			t.Errorf("Got back non-status result: %v", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Error("Unexpected timeout from async channel")
	}
}

func TestDeleteValid(t *testing.T) {
	registry := test.ClientRegistry{}
	storage := REST{
		registry: &registry,
	}

	ctx := api.NewContext()
	channel, err := storage.Delete(ctx, "foo")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	select {
	case r := <-channel:
		switch r := r.Object.(type) {
		case *api.Status:
			if r.Status != "Success" {
				t.Errorf("Got back non-success status: %#v", r)
			}
		default:
			t.Errorf("Got back non-status result: %v", r)
		}
	case <-time.After(time.Millisecond * 100):
		t.Error("Unexpected timeout from async channel")
	}

	if registry.DeletedClientName != "foo" {
		t.Error("Unexpected client deleted: %s", registry.DeletedClientName)
	}
}
