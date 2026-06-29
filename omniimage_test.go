package omniimage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/plexusone/omniimage/provider"
)

func TestNewClient(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"created": time.Now().Unix(),
			"data": []map[string]string{
				{"url": "https://example.com/image.png"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Test with mock OpenAI provider
	client, err := NewClient(ClientConfig{
		Providers: []ProviderConfig{
			{
				Provider: ProviderNameOpenAI,
				APIKey:   "test-key",
				BaseURL:  server.URL,
			},
		},
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	// Test generate
	resp, err := client.Generate(context.Background(), &provider.GenerateRequest{
		Model:  ModelGPTImage2,
		Prompt: "A test image",
	})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if len(resp.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(resp.Images))
	}

	if resp.Images[0].URL != "https://example.com/image.png" {
		t.Errorf("unexpected image URL: %s", resp.Images[0].URL)
	}
}

func TestNoProviders(t *testing.T) {
	_, err := NewClient(ClientConfig{})
	if err != ErrNoProviders {
		t.Errorf("expected ErrNoProviders, got %v", err)
	}
}

func TestCapabilities(t *testing.T) {
	caps := GetCapabilities(ProviderNameOpenAI)
	if caps == nil {
		t.Fatal("expected capabilities for OpenAI")
	}

	if !caps.Generate {
		t.Error("expected Generate capability")
	}
	if !caps.Edit {
		t.Error("expected Edit capability")
	}
	if !caps.Variations {
		t.Error("expected Variations capability")
	}
	if caps.Upscale {
		t.Error("did not expect Upscale capability for OpenAI")
	}

	caps = GetCapabilities(ProviderNameFal)
	if caps == nil {
		t.Fatal("expected capabilities for Fal")
	}

	if !caps.Generate {
		t.Error("expected Generate capability")
	}
	if caps.Edit {
		t.Error("did not expect Edit capability for Fal")
	}
	if !caps.Upscale {
		t.Error("expected Upscale capability for Fal")
	}
}

func TestModelInfo(t *testing.T) {
	info := GetModelInfo(ModelGPTImage2)
	if info == nil {
		t.Fatal("expected model info for GPT Image 2")
	}

	if info.Provider != ProviderNameOpenAI {
		t.Errorf("expected OpenAI provider, got %s", info.Provider)
	}
	if info.MaxImages != 10 {
		t.Errorf("expected MaxImages=10 for GPT Image 2, got %d", info.MaxImages)
	}
	if !info.SupportsHD {
		t.Error("expected GPT Image 2 to support HD")
	}

	info = GetModelInfo(ModelFluxPro)
	if info == nil {
		t.Fatal("expected model info for FLUX Pro")
	}

	if info.Provider != ProviderNameFal {
		t.Errorf("expected Fal provider, got %s", info.Provider)
	}
}

func TestProviderTypes(t *testing.T) {
	tests := []struct {
		size   provider.ImageSize
		width  int
		height int
	}{
		{provider.Size256x256, 256, 256},
		{provider.Size512x512, 512, 512},
		{provider.Size1024x1024, 1024, 1024},
		{provider.Size1024x1792, 1024, 1792},
		{provider.Size1792x1024, 1792, 1024},
	}

	for _, tt := range tests {
		t.Run(string(tt.size), func(t *testing.T) {
			// Just verify the constants are defined correctly
			if tt.size == "" {
				t.Error("size should not be empty")
			}
		})
	}
}
