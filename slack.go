package main

// slack slackNotification
type SlackNotification struct {
    Channel     string       `json:"channel,omitempty"`
    Username    string       `json:"username,omitempty"`
    IconEmoji   string       `json:"icon_emoji,omitempty"`
    IconURL     string       `json:"icon_url,omitempty"`
    LinkNames   bool         `json:"link_names,omitempty"`
    Attachments []Attachment `json:"attachments"`
}

// attachment is used to display a richly-formatted message block.
type Attachment struct {
    Title      string               `json:"title,omitempty"`
    TitleLink  string               `json:"title_link,omitempty"`
    Pretext    string               `json:"pretext,omitempty"`
    Text       string               `json:"text"`
    Fallback   string               `json:"fallback"`
    CallbackID string               `json:"callback_id"`
    Fields     []SlackField         `json:"fields,omitempty"`
    Actions    []SlackAction        `json:"actions,omitempty"`
    ImageURL   string               `json:"image_url,omitempty"`
    ThumbURL   string               `json:"thumb_url,omitempty"`
    Footer     string               `json:"footer"`
    Color      string               `json:"color,omitempty"`
    MrkdwnIn   []string             `json:"mrkdwn_in,omitempty"`
}

// SlackField configures a single Slack field that is sent with each slackNotification.
// Each field must contain a title, value, and optionally, a boolean value to indicate if the field
// is short enough to be displayed next to other fields designated as short.
// See https://api.slack.com/docs/message-attachments#fields for more information.
type SlackField struct {
    Title string `yaml:"title,omitempty" json:"title,omitempty"`
    Value string `yaml:"value,omitempty" json:"value,omitempty"`
    Short *bool  `yaml:"short,omitempty" json:"short,omitempty"`
}

// SlackAction configures a single Slack action that is sent with each slackNotification.
// See https://api.slack.com/docs/message-attachments#action_fields and https://api.slack.com/docs/message-buttons
// for more information.
type SlackAction struct {
    Type         string                  `yaml:"type,omitempty"  json:"type,omitempty"`
    Text         string                  `yaml:"text,omitempty"  json:"text,omitempty"`
    URL          string                  `yaml:"url,omitempty"   json:"url,omitempty"`
    Style        string                  `yaml:"style,omitempty" json:"style,omitempty"`
    Name         string                  `yaml:"name,omitempty"  json:"name,omitempty"`
    Value        string                  `yaml:"value,omitempty"  json:"value,omitempty"`
    ConfirmField *SlackConfirmationField `yaml:"confirm,omitempty"  json:"confirm,omitempty"`
}

// SlackConfirmationField protect users from destructive actions or particularly distinguished decisions
// by asking them to confirm their button click one more time.
// See https://api.slack.com/docs/interactive-message-field-guide#confirmation_fields for more information.
type SlackConfirmationField struct {
    Text        string `yaml:"text,omitempty"  json:"text,omitempty"`
    Title       string `yaml:"title,omitempty"  json:"title,omitempty"`
    OkText      string `yaml:"ok_text,omitempty"  json:"ok_text,omitempty"`
    DismissText string `yaml:"dismiss_text,omitempty"  json:"dismiss_text,omitempty"`
}
