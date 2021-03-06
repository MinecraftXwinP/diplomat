package data

type SimpleConfiguration struct {
	Preprocessors []Preprocessor `navigate:"preprocessors"`
	Outputs       []Output `navigate:"outputs"`
}

func (s SimpleConfiguration) GetPreprocessors() []Preprocessor {
	return s.Preprocessors
}

func (s SimpleConfiguration) GetOutputs() []Output {
	return s.Outputs
}

type SimpleOutput struct {
	Selectors []Selector `navigate:"selectors"`
	Templates []Template `navigate:"templates"`
}

func (s SimpleOutput) GetSelectors() []Selector {
	return s.Selectors
}

func (s SimpleOutput) GetTemplates() []Template {
	return s.Templates
}

type SimpleTemplate struct {
	Type    string `navigate:"type"`
	Options TemplateOption `navigate:"options"`
}

func (s SimpleTemplate) GetType() string {
	return s.Type
}

func (s SimpleTemplate) GetOptions() TemplateOption {
	return s.Options
}

type SimpleTemplateOption map[string]interface{}

func (s SimpleTemplateOption) GetMapElement() map[string]interface{} {
	return s
}

func (s SimpleTemplateOption) GetFilename() string {
	return s.GetMapElement()["filename"].(string)
}

type SimplePreprocessor struct {
	Type    string `navigate:"type"`
	Options interface{} `navigate:"options"`
}

func (s SimplePreprocessor) GetType() string {
	return s.Type
}

func (s SimplePreprocessor) GetOptions() interface{} {
	return s.Options
}
