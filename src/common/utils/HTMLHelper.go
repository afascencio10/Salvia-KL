package utils

import (
	"encoding/json"
	"html"
)

type HTMLInput struct {
	Locale           map[string]map[string]string
	InputId          string `json:"inputId"`
	ErrorSpanId      string `json:"errorSpanId"`
	InputType        string `json:"inputType"`
	InputClass       string `json:"class"`
	InputLabelClass  string `json:"labelClass"`
	InputErrorClass  string `json:"errorClass"`
	InputLabel       string `json:"inputLabel"`
	InputPlaceHolder string `json:"inputPlaceHolder"`
	EntityName       string `json:"entityName"`
	EntityPath       string `json:"entityPath"`
	EntityAttr       string `json:"entityAttr"`
	Mandatory        bool   `json:"mandatory"`
	InputPattern     string `json:"pattern"`
	InputMaxLength   string `json:"maxLength"`
}

func (i *HTMLInput) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLInput {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLInput = HTMLInput{}
	json.Unmarshal([]byte(jsonStr), &tmp)

	tmp.Locale = locale
	tmp.InputLabel = tmp.Translate(lang, tmp.InputLabel)
	tmp.InputPlaceHolder = tmp.Translate(lang, tmp.InputPlaceHolder)

	return tmp
}

func (i *HTMLInput) Translate(lang string, code string) string {
	return translate(lang, code, i.Locale)
}

type HTMLTextArea struct {
	Locale              map[string]map[string]string
	TextAreaId          string `json:"textAreaId"`
	ErrorSpanId         string `json:"errorSpanId"`
	TextAreaRows        string `json:"rows"`
	TextAreaClass       string `json:"class"`
	TextAreaLabelClass  string `json:"labelClass"`
	TextAreaErrorClass  string `json:"errorClass"`
	TextAreaLabel       string `json:"textAreaLabel"`
	TextAreaPlaceHolder string `json:"textAreaPlaceHolder"`
	EntityName          string `json:"entityName"`
	EntityPath          string `json:"entityPath"`
	EntityAttr          string `json:"entityAttr"`
	Mandatory           bool   `json:"mandatory"`
}

func (i *HTMLTextArea) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLTextArea {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLTextArea = HTMLTextArea{}
	json.Unmarshal([]byte(jsonStr), &tmp)

	tmp.Locale = locale
	tmp.TextAreaLabel = tmp.Translate(lang, tmp.TextAreaLabel)
	tmp.TextAreaPlaceHolder = tmp.Translate(lang, tmp.TextAreaPlaceHolder)

	return tmp
}

func (i *HTMLTextArea) Translate(lang string, code string) string {
	return translate(lang, code, i.Locale)
}

type HTMLSelect struct {
	Locale            map[string]map[string]string
	Multiple          string `json:"multiple"`
	SelectId          string `json:"selectId"`
	ErrorSpanId       string `json:"errorSpanId"`
	ContainerClass    string `json:"containerClass"`
	SelectClass       string `json:"class"`
	SelectErrorClass  string `json:"errorClass"`
	SelectLabel       string `json:"selectLabel"`
	SelectPlaceHolder string `json:"selectPlaceHolder"`
	EntityName        string `json:"entityName"`
	EntityAttr        string `json:"entityAttr"`
	EntityPath        string `json:"entityPath"`
	EnumsAttr         string `json:"enumsAttr"`
	EntityEnumLabel   string `json:"entityEnumLabel"`
	EntityEnumId      string `json:"entityEnumId"`
	EmptyLabel        string `json:"emptyLabel"`
	EmptyValue        string `json:"emptyValue"`
	Mandatory         bool   `json:"mandatory"`
	OnChanged         string `json:"onChanged"`
	ReturnObject      bool   `json:"returnObject"`
}

func (i *HTMLSelect) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLSelect {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLSelect = HTMLSelect{}
	tmp.ReturnObject = false

	json.Unmarshal([]byte(jsonStr), &tmp)
	tmp.Locale = locale
	tmp.SelectLabel = tmp.Translate(lang, tmp.SelectLabel)
	tmp.SelectPlaceHolder = tmp.Translate(lang, tmp.SelectPlaceHolder)
	tmp.EmptyLabel = tmp.Translate(lang, tmp.EmptyLabel)
	tmp.EmptyValue = tmp.Translate(lang, tmp.EmptyValue)

	return tmp
}

func (i *HTMLSelect) Translate(lang string, code string) string {
	return translate(lang, code, i.Locale)
}

type HTMLDatalist struct {
	Locale              map[string]map[string]string
	Multiple            bool   `json:"multiple"`
	DatalistId          string `json:"selectId"`
	ErrorSpanId         string `json:"errorSpanId"`
	DatalistClass       string `json:"class"`
	DatalistErrorClass  string `json:"errorClass"`
	DatalistLabel       string `json:"datalistLabel"`
	DatalistPlaceHolder string `json:"datalistPlaceHolder"`
	EntityName          string `json:"entityName"`
	EntityAttr          string `json:"entityAttr"`
	EntityPath          string `json:"entityPath"`
	EnumsAttr           string `json:"enumsAttr"`
	EntityEnumLabel     string `json:"entityEnumLabel"`
	EntityEnumId        string `json:"entityEnumId"`
	Mandatory           bool   `json:"mandatory"`
}

func (i *HTMLDatalist) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLDatalist {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLDatalist = HTMLDatalist{}
	json.Unmarshal([]byte(jsonStr), &tmp)
	tmp.Locale = locale
	tmp.DatalistLabel = tmp.Translate(lang, tmp.DatalistLabel)
	tmp.DatalistPlaceHolder = tmp.Translate(lang, tmp.DatalistPlaceHolder)

	return tmp
}

func (i *HTMLDatalist) Translate(lang string, code string) string {
	return translate(lang, code, i.Locale)
}

type HTMLButton struct {
	Locale            map[string]map[string]string
	ButtonId          string `json:"buttonId"`
	ButtonLabel       string `json:"buttonLabel"`
	ButtonClass       string `json:"class"`
	ButtonPlaceHolder string `json:"buttonPlaceHolder"`
	OnClickHandler    string `json:"onClickHandler"`
}

func (b *HTMLButton) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLButton {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLButton = HTMLButton{}
	json.Unmarshal([]byte(jsonStr), &tmp)
	tmp.Locale = locale
	tmp.ButtonLabel = tmp.Translate(lang, tmp.ButtonLabel)
	tmp.ButtonPlaceHolder = tmp.Translate(lang, tmp.ButtonPlaceHolder)

	return tmp
}

func (b *HTMLButton) Translate(lang string, code string) string {
	return translate(lang, code, b.Locale)
}

type HTMLDateTime struct {
	Locale              map[string]map[string]string
	DateTimeId          string `json:"dateTimeId"`
	ErrorSpanId         string `json:"errorSpanId"`
	DateTimeClass       string `json:"class"`
	DateTimeErrorClass  string `json:"errorClass"`
	DateTimeLabel       string `json:"dateTimeLabel"`
	DateTimePlaceHolder string `json:"dateTimePlaceHolder"`
	DateTimeFormat      string `json:"dateTimeFormat"`
	EntityName          string `json:"entityName"`
	EntityAttr          string `json:"entityAttr"`
	EntityPath          string `json:"entityPath"`
	EnumsAttr           string `json:"enumsAttr"`
	EntityEnumLabel     string `json:"entityEnumLabel"`
	EntityEnumId        string `json:"entityEnumId"`
	Mandatory           bool   `json:"mandatory"`
}

func (dt *HTMLDateTime) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLDateTime {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLDateTime = HTMLDateTime{}
	json.Unmarshal([]byte(jsonStr), &tmp)
	tmp.Locale = locale
	tmp.DateTimeLabel = tmp.Translate(lang, tmp.DateTimeLabel)
	tmp.DateTimePlaceHolder = tmp.Translate(lang, tmp.DateTimePlaceHolder)

	return tmp
}

func (b *HTMLDateTime) Translate(lang string, code string) string {
	return translate(lang, code, b.Locale)
}

type HTMLText struct {
	Locale     map[string]map[string]string
	Text       string `json:"text"`
	HTMLDecode bool   `json:"htmlDecode"`
}

func (b *HTMLText) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLText {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLText = HTMLText{}
	json.Unmarshal([]byte(jsonStr), &tmp)
	tmp.Locale = locale
	tmp.Text = tmp.Translate(lang, tmp.Text)
	if tmp.HTMLDecode {
		tmp.Text = html.UnescapeString(tmp.Text)
	}

	return tmp
}

func (b *HTMLText) Translate(lang string, code string) string {
	return translate(lang, code, b.Locale)
}

type HTMLStarRating struct {
	Locale          map[string]map[string]string
	RatingId        string `json:"ratingId"`
	ErrorClass      string `json:"errorClass"`
	RatingClass     string `json:"ratingClass"`
	RatingMaxValue  int    `json:"ratingMaxValue"`
	StarClass       string `json:"starClass"`
	EntityName      string `json:"entityName"`
	EntityAttr      string `json:"entityAttr"`
	EntityPath      string `json:"entityPath"`
	FilledStarClass string `json:"filledStarClass"`
	EmptyStarClass  string `json:"emptyStarClass"`
	RatingLabel     string `json:"ratingLabel"`
	Mandatory       bool   `json:"mandatory"`
	LabelClass      string `json:"labelClass"`
	StarGroupClass  string `json:"starGroupClass"`
}

func (i *HTMLStarRating) New(jsonStr string, lang string, locale map[string]map[string]string) HTMLStarRating {
	//Aquí viene el string con todos los valores de los campos.
	var tmp HTMLStarRating = HTMLStarRating{}
	json.Unmarshal([]byte(jsonStr), &tmp)

	tmp.Locale = locale
	tmp.RatingLabel = tmp.Translate(lang, tmp.RatingLabel)

	return tmp
}

func (i *HTMLStarRating) Translate(lang string, code string) string {
	return translate(lang, code, i.Locale)
}

func translate(lang string, code string, locale map[string]map[string]string) string {
	if langVal, found := locale[lang]; found {
		if labelVal, found := langVal[code]; found {
			return labelVal
		}
	}
	return ""
}
