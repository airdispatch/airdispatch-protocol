package common

import (
	"airdispat.ch/airdispatch"
	"code.google.com/p/goprotobuf/proto"
	"errors"
)

type ADMail struct {
	payload map[string]*ADComponent

	byteload       []byte
	encryptionType string
	encrypted      bool

	FromAddress *ADAddress
	ToAddress   *ADAddress
	Timestamp   uint64
}

func (a *ADMail) HasDataType(typeName string) bool {
	_, ok := a.payload[typeName]
	return ok
}

func (a *ADMail) GetADComponentForType(typeName string) (*ADComponent, error) {
	v, ok := a.payload[typeName]
	if !ok {
		return nil, errors.New("ADMail doesn't contain that Type")
	}

	return v, nil
}

func CreateADMailFromADMessage(message *ADMessage, key *ADKey) (*ADMail, error) {
	if message.MessageType != MAIL_MESSAGE {
		return nil, errors.New("Cannot translate an ADMessage with incorrect message type to ADMail.")
	}

	output := &ADMail{}

	theMessage := &airdispatch.Mail{}
	err := proto.Unmarshal(message.Payload, theMessage)
	if err != nil {
		return nil, ADUnmarshallingError
	}

	output.FromAddress = message.FromAddress

	// TODO: Verify Addresses Match

	output.ToAddress = CreateADAddress(theMessage.GetToAddress())
	if output.ToAddress == nil {
		return nil, errors.New("Couldn't resolve To Address")
	}

	output.Timestamp = theMessage.GetTimestamp()
	output.byteload = theMessage.GetData()
	output.encryptionType = theMessage.GetEncryption()
	output.encrypted = !(theMessage.GetEncryption() == ADEncryptionNone)

	if output.encrypted && key == nil {
		return output, nil
	}

	output.DecryptPayload(key)
	output.Unmarshal()

	return output, nil
}

func (a *ADMail) DecryptPayload(key *ADKey) bool {
	if !a.encrypted {
		return false
	}

	var err error
	a.byteload, err = key.DecryptPayload(a.byteload)
	if err != nil {
		return false
	}

	a.encrypted = false
	return true
}

func (a *ADMail) EncryptPayload(address *ADAddress, key *ADKey, trackerList *ADTrackerList) bool {
	if a.encrypted {
		return false
	}

	encryptionKey, err := address.GetEncryptionKey(key, trackerList)
	if err != nil {
		return false
	}

	a.byteload, err = EncryptPayload(a.byteload, encryptionKey)
	if err != nil {
		return false
	}

	a.encrypted = true
	return true
}

func (a *ADMail) Unmarshal() error {
	if a.encrypted {
		return ADDecryptionError
	}

	theData := &airdispatch.MailData{}
	err := proto.Unmarshal(a.byteload, theData)
	if err != nil {
		return ADUnmarshallingError
	}

	componentMap := make(map[string]*ADComponent)
	for _, v := range theData.GetPayload() {
		if v.GetEncryption() != "" {
			return ADUnmarshallingError
		}

		c := CreateADComponent(v.GetTypeName(), v.GetPayload())
		componentMap[c.DataTypeValue()] = c
	}

	a.payload = componentMap

	return nil
}

func (a *ADMail) Marshal(address *ADAddress, key *ADKey, trackerList *ADTrackerList) (*ADMessage, error) {
	innerComponents := make([]*airdispatch.MailData_DataType, len(a.payload))
	incrementer := 0
	for _, v := range a.payload {
		innerComponents[incrementer] = v.ToPrimative()
		incrementer++
	}

	dataComponents := &airdispatch.MailData{
		Payload: innerComponents,
	}

	var err error
	a.byteload, err = proto.Marshal(dataComponents)
	if err != nil {
		return nil, err
	}

	if a.encryptionType != ADEncryptionNone {
		a.EncryptPayload(address, key, trackerList)
		a.encrypted = true
	}

	newMessage := &ADMessage{}
	newMessage.FromAddress = a.FromAddress
	newMessage.MessageType = MAIL_MESSAGE
	newMessage.Payload = a.byteload

	return newMessage, nil
}

func (a *ADMail) HashContents() []byte {
	if a.byteload != nil {
		return HashSHA(a.byteload)
	}
	return nil
}

// A simple message to output an Airdispatch Message to String
func (a *ADMail) PrintMessage() string {
	output := ""
	output += ("---- Message from " + a.FromAddress.ToString() + " ----\n")
	output += ("Encryption Type: " + a.encryptionType + "\n")

	for _, value := range a.payload {
		output += ("### " + value.DataTypeValue() + "\n")
		output += (value.StringValue() + "\n")
	}

	output += ("---- END ----")

	return output
}

func CreateADMail(fromAddress *ADAddress, toAddress *ADAddress, timestamp uint64, payload []*ADComponent) *ADMail {
	output := &ADMail{}

	output.FromAddress = fromAddress
	output.ToAddress = toAddress
	output.Timestamp = timestamp

	componentMap := make(map[string]*ADComponent)
	for _, v := range payload {
		componentMap[v.DataTypeValue()] = v
	}

	output.payload = componentMap

	return output
}

// ----------
// AD COMPONENT
// ----------e

type ADComponent struct {
	data_type      string
	data_component []byte
}

func (a *ADComponent) StringValue() string {
	return string(a.data_component)
}

func (a *ADComponent) ByteValue() []byte {
	return a.data_component
}

func (a *ADComponent) DataTypeValue() string {
	return a.data_type
}

func CreateADComponent(name string, data []byte) *ADComponent {
	return &ADComponent{name, data}
}

func (a *ADComponent) ToPrimative() *airdispatch.MailData_DataType {
	newDataType := &airdispatch.MailData_DataType{}

	newDataType.Payload = a.data_component
	newDataType.TypeName = &a.data_type

	return newDataType
}