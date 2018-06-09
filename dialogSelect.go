package slack

// SelectDataSource types of select datasource
type SelectDataSource string

const (
	// StaticDataSource menu with static Options/OptionGroups
	StaticDataSource SelectDataSource = "static"
	// ExternalDataSource dynamic datasource
	ExternalDataSource SelectDataSource = "external"
	// ConversationsDataSource provides a list of conversations
	ConversationsDataSource SelectDataSource = "conversations"
	// ChannelsDataSource provides a list of channels
	ChannelsDataSource SelectDataSource = "channels"
	// UsersDataSource provides a list of users
	UsersDataSource SelectDataSource = "users"
)

// baseSelect a menu select for dialogs
type baseSelect struct {
	DialogInput
	DataSource SelectDataSource `json:"data_source"`
}

// StaticSelectDialogInput can support all type except Dynamic menu
type StaticSelectDialogInput struct {
	baseSelect
	Value        string         `json:"value"` //This option is invalid in external, where you must use selected_options
	Options      []SelectOption `json:"options,omitempty"`
	OptionGroups []OptionGroup  `json:"option_groups,omitempty"`
}

// SelectOption is an option for the user to select from the menu
type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// OptionGroup is a collection of options for creating a segmented table
type OptionGroup struct {
	Label   string         `json:"label"`
	Options []SelectOption `json:"options"`
}

// ExternalSelectInputElement is a special case of `SelectInputElement``
type ExternalSelectInputElement struct {
	baseSelect
	SelectedOptions []SelectOption `json:"selected_options"` //This option is invalid in external, where you must use selected_options
}

// NewStaticSelectDialogInput constructor for a `static` datasource menu input
func NewStaticSelectDialogInput(name, label string, options []SelectOption) *StaticSelectDialogInput {
	return &StaticSelectDialogInput{
		baseSelect: baseSelect{
			DialogInput: DialogInput{
				Type:     InputTypeSelect,
				Name:     name,
				Label:    label,
				Optional: true,
			},
			DataSource: StaticDataSource,
		},
		Options: options,
	}
}

// NewGroupedSelectDialoginput a grouped options select input for Dialogs
func NewGroupedSelectDialoginput(name, label string, groups map[string]map[string]string) *StaticSelectDialogInput {
	optionGroups := []OptionGroup{}
	for groupName, options := range groups {
		optionGroups = append(optionGroups, OptionGroup{
			Label:   groupName,
			Options: optionsFromMap(options),
		})
	}
	return &StaticSelectDialogInput{
		baseSelect: baseSelect{
			DialogInput: DialogInput{
				Type:  InputTypeSelect,
				Name:  name,
				Label: label,
			},
			DataSource: StaticDataSource,
		},
		OptionGroups: optionGroups,
	}
}

func optionsFromArray(options []string) []SelectOption {
	selectOptions := make([]SelectOption, len(options))
	for idx, value := range options {
		selectOptions[idx] = SelectOption{
			Label: value,
			Value: value,
		}
	}
	return selectOptions
}

func optionsFromMap(options map[string]string) []SelectOption {
	selectOptions := make([]SelectOption, len(options))
	idx := 0
	var option SelectOption
	for key, value := range options {
		option = SelectOption{
			Label: key,
			Value: value,
		}
		selectOptions[idx] = option
		idx++
	}
	return selectOptions
}

// NewConversationsSelect returns a `Conversations` select
func NewConversationsSelect(name, label string) *StaticSelectDialogInput {
	return newPresetSelect(name, label, ConversationsDataSource)
}

// NewChannelsSelect returns a `Channels` select
func NewChannelsSelect(name, label string) *StaticSelectDialogInput {
	return newPresetSelect(name, label, ChannelsDataSource)
}

// NewUsersSelect returns a `Users` select
func NewUsersSelect(name, label string) *StaticSelectDialogInput {
	return newPresetSelect(name, label, UsersDataSource)
}

func newPresetSelect(name, label string, dataSourceType SelectDataSource) *StaticSelectDialogInput {
	return &StaticSelectDialogInput{
		baseSelect: baseSelect{
			DialogInput: DialogInput{
				Type:  InputTypeSelect,
				Label: label,
				Name:  name,
			},
			DataSource: dataSourceType,
		},
	}

}
