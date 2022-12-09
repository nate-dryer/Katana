package output

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"github.com/projectdiscovery/fileutil"
	"github.com/projectdiscovery/sliceutil"
	stringsutil "github.com/projectdiscovery/utils/strings"
	"gopkg.in/yaml.v2"
)

// CustomFieldsMap is the global custom field data instance
// it is used for parsing the header and body of request
var CustomFieldsMap = make(map[string]CustomFieldConfig)

// CustomFieldConfig contains suggestions for field filling
type CustomFieldConfig struct {
	Name         string           `yaml:"name,omitempty"`
	Type         string           `yaml:"type,omitempty"`
	Group        int              `yaml:"group,omitempty"`
	Regex        []string         `yaml:"regex,omitempty"`
	CompileRegex []*regexp.Regexp `yaml:"-"`
}

var DefaultFieldConfigData = []CustomFieldConfig{
	{
		Name:  "email",
		Type:  "regex",
		Regex: []string{`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`},
	},
}

func (c *CustomFieldConfig) SetCompiledRegexp(r *regexp.Regexp) {
	c.CompileRegex = append(c.CompileRegex, r)
}

func (c *CustomFieldConfig) GetName() string {
	return c.Name
}

func parseCustomFieldName(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "could not read field config")
	}
	defer file.Close()

	var data []CustomFieldConfig
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return errors.Wrap(err, "could not decode field config")
	}
	passedCustomFieldMap := make(map[string]CustomFieldConfig)
	for _, item := range data {
		if !regexp.MustCompile(`^[A-Za-z0-9_-]+$`).MatchString(item.Name) {
			return fmt.Errorf("wrong custom field name %s", item.Name)
		}
		// check custom field name is pre-defined or not
		if sliceutil.Contains(FieldNames, item.Name) {
			return fmt.Errorf("could not register custom field. \"%s\" already pre-defined field", item.Name)
		}
		// check custom field name should be unqiue
		if _, ok := passedCustomFieldMap[item.Name]; ok {
			return fmt.Errorf("could not register custom field. \"%s\" custom field already exists", item.Name)
		}
		passedCustomFieldMap[item.Name] = item
	}
	return nil
}

func loadCustomFields(filePath string, fields string) error {
	var err error

	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "could not read field config")
	}
	defer file.Close()

	var data []CustomFieldConfig
	// read the field config file
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return errors.Wrap(err, "could not decode field config")
	}
	allCustomFields := make(map[string]CustomFieldConfig)
	for _, item := range data {
		for _, rg := range item.Regex {
			regex, err := regexp.Compile(rg)
			if err != nil {
				return errors.Wrap(err, "could not parse regex in field config")
			}
			item.SetCompiledRegexp(regex)
		}
		allCustomFields[item.Name] = item
	}
	// Set the passed custom field value globally
	for _, f := range stringsutil.SplitAny(fields, ",") {
		if val, ok := allCustomFields[f]; ok {
			CustomFieldsMap[f] = val
		}
	}
	return nil
}

func initCustomFieldConfigFile() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "could not get home directory")
	}
	defaultConfig := filepath.Join(homedir, ".config", "katana", "field-config.yaml")

	if fileutil.FileExists(defaultConfig) {
		return defaultConfig, nil
	}
	if err := os.MkdirAll(filepath.Dir(defaultConfig), 0775); err != nil {
		return "", err
	}
	customFieldConfig, err := os.Create(defaultConfig)
	if err != nil {
		return "", errors.Wrap(err, "could not get home directory")
	}
	defer customFieldConfig.Close()

	err = yaml.NewEncoder(customFieldConfig).Encode(DefaultFieldConfigData)
	if err != nil {
		return "", err
	}
	return defaultConfig, nil
}
