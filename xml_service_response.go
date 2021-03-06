package cas

import (
	"encoding/xml"
	"errors"
	"strings"
	"time"
)

type xmlServiceResponse struct {
	XMLName xml.Name `xml:"http://www.yale.edu/tp/cas serviceResponse"`

	Failure *xmlAuthenticationFailure
	Success *xmlAuthenticationSuccess
}

type xmlAuthenticationFailure struct {
	XMLName xml.Name `xml:"authenticationFailure"`
	Code    string   `xml:"code,attr"`
	Message string   `xml:",innerxml"`
}

type xmlAuthenticationSuccess struct {
	XMLName             xml.Name           `xml:"authenticationSuccess"`
	User                string             `xml:"user"`
	ProxyGrantingTicket string             `xml:"proxyGrantingTicket,omitempty"`
	Proxies             *xmlProxies        `xml:"proxies"`
	Attributes          *xmlAttributes     `xml:"attributes"`
	ExtraAttributes     []*xmlAnyAttribute `xml:",any"`
}

type xmlProxies struct {
	XMLName xml.Name `xml:"proxies"`
	Proxies []string `xml:"proxy"`
}

func (p *xmlProxies) AddProxy(proxy string) {
	p.Proxies = append(p.Proxies, proxy)
}

type xmlAttributes struct {
	XMLName                                xml.Name   `xml:"attributes"`
	AuthenticationDate                     *fixedTime `xml:"authenticationDate"`
	LongTermAuthenticationRequestTokenUsed bool       `xml:"longTermAuthenticationRequestTokenUsed"`
	IsFromNewLogin                         bool       `xml:"isFromNewLogin"`
	MemberOf                               []string   `xml:"memberOf"`
	UserAttributes                         *xmlUserAttributes
	ExtraAttributes                        []*xmlAnyAttribute `xml:",any"`
}

type fixedTime struct {
	time.Time
}

func (t *fixedTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	dataToken, err := d.Token()
	if err != nil {
		return err
	}

	charData, ok := dataToken.(xml.CharData)
	if !ok {
		return errors.New("Expected chardata")
	}

	timeStr := strings.SplitN(string(charData), "[", 2)[0]
	timeD, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}

	*t = fixedTime{timeD}

	_, err = d.Token()
	return err
}

type xmlUserAttributes struct {
	XMLName       xml.Name             `xml:"userAttributes"`
	Attributes    []*xmlNamedAttribute `xml:"attribute"`
	AnyAttributes []*xmlAnyAttribute   `xml:",any"`
}

type xmlNamedAttribute struct {
	XMLName xml.Name `xml:"attribute"`
	Name    string   `xml:"name,attr,omitempty"`
	Value   string   `xml:",innerxml"`
}

type xmlAnyAttribute struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (xsr *xmlServiceResponse) marshalXML(indent int) ([]byte, error) {
	if indent == 0 {
		return xml.Marshal(xsr)
	}

	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += " "
	}

	return xml.MarshalIndent(xsr, "", indentStr)
}

func failureServiceResponse(code, message string) *xmlServiceResponse {
	return &xmlServiceResponse{
		Failure: &xmlAuthenticationFailure{
			Code:    code,
			Message: message,
		},
	}
}

func successServiceResponse(username, pgt string) *xmlServiceResponse {
	return &xmlServiceResponse{
		Success: &xmlAuthenticationSuccess{
			User:                username,
			ProxyGrantingTicket: pgt,
		},
	}
}
