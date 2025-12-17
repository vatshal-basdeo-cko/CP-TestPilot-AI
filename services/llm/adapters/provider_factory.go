package adapters

// ProviderFactory manages LLM providers
type ProviderFactory struct {
	providers       map[string]LLMProvider
	defaultProvider string
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory(openAIKey, anthropicKey, geminiKey, defaultProvider string) *ProviderFactory {
	factory := &ProviderFactory{
		providers:       make(map[string]LLMProvider),
		defaultProvider: defaultProvider,
	}

	// Register providers
	openai := NewOpenAIProvider(openAIKey)
	if openai.IsAvailable() {
		factory.providers["openai"] = openai
	}

	anthropic := NewAnthropicProvider(anthropicKey)
	if anthropic.IsAvailable() {
		factory.providers["anthropic"] = anthropic
	}

	gemini := NewGeminiProvider(geminiKey)
	if gemini.IsAvailable() {
		factory.providers["gemini"] = gemini
	}

	return factory
}

// GetProvider returns the specified provider or the default
func (f *ProviderFactory) GetProvider(name string) LLMProvider {
	if name != "" {
		if provider, ok := f.providers[name]; ok {
			return provider
		}
	}

	// Try default provider
	if provider, ok := f.providers[f.defaultProvider]; ok {
		return provider
	}

	// Return first available provider
	for _, provider := range f.providers {
		return provider
	}

	return nil
}

// GetDefaultProvider returns the default provider
func (f *ProviderFactory) GetDefaultProvider() LLMProvider {
	return f.GetProvider(f.defaultProvider)
}

// ListAvailableProviders returns names of available providers
func (f *ProviderFactory) ListAvailableProviders() []string {
	var names []string
	for name := range f.providers {
		names = append(names, name)
	}
	return names
}

// GetDefaultProviderName returns the name of the default provider
func (f *ProviderFactory) GetDefaultProviderName() string {
	if provider := f.GetDefaultProvider(); provider != nil {
		return provider.Name()
	}
	return "none"
}

