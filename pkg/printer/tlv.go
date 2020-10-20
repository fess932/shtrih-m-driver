package printer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
)

func (p *Printer) TLVWriteCashierINN(INN string) error {
	cmdBinary, cmdLen := p.createCommandData(consts.SendTLVToOp)
	buf := bytes.NewBuffer(cmdBinary)
	p.logger.Debug(buf)

	if utf8.RuneCountInString(INN) > int(consts.CashierINN.Length) {
		return errors.New("cлишком длинное инн")
	}

	cpStr, err := charmap.CodePage866.NewEncoder().String(INN)
	if err != nil {
		return err
	}

	p.logger.Debug(cpStr)

	// tlv структура
	tlv, err := newTLV(consts.CashierINN.Code, consts.CashierINN.Length, []byte(cpStr))
	if err != nil {
		return err
	}
	buf.Write(tlv.Tag)
	buf.Write(tlv.Len)
	buf.Write(tlv.Value)

	p.logger.Debug("Код, длинна, значение:", tlv.Tag, tlv.Len, tlv.Value)
	p.logger.Debug("Команда с тлв структурой\n", hex.Dump(buf.Bytes()))

	rFrame, err := p.send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Fatal(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.bytes()))

	//
	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		return err
	}

	return nil
}

type TLV struct {
	Tag   []byte
	Len   []byte
	Value []byte
}

func newTLV(Tag, Len uint16, value []byte) (TLV, error) {
	tlv := TLV{
		Tag:   make([]byte, 2),
		Len:   make([]byte, 2),
		Value: make([]byte, Len),
	}
	binary.LittleEndian.PutUint16(tlv.Tag, Tag) // код тега
	binary.LittleEndian.PutUint16(tlv.Len, Len) // длинна тега

	copy(tlv.Value, value) // значение тега

	if len(tlv.Value) != int(Len) {
		return TLV{}, errors.New("длинна не совпадает со значением")
	}

	return tlv, nil
}

//ff 4d 1e 00 00 00 b3 04  0c 00 32 36 33 32 30 39 37 34 35 33 35 37
//type Writer struct {
//	buf bufio.Writer
//}

// TEXT в TLV кодировка : CP866

// tlv types:
// string cp866
// uint16
//
//const (
//	StringTLV = iota
//	Uint16TLV
//	Uint32TLV
//)
//
//type TLVData []byte
//
//type TLVStruct struct {
//	Tag  string
//	Len  string
//	Data struct {
//		Type  int
//		Value []byte
//	}
//}
//
//// тип   длинна  значение (длинна сейчас 16)1
////11 04 | 10 00 | 39 32 38 31 30 30 30 31 30 30 30 30 37 34 34 32
//
//func New(tag uint16, dataLen uint16) TLVData {
//	buf := bytes.NewBuffer([]byte{})
//	binary.Write(buf, binary.LittleEndian, tag)
//	binary.Write(buf, binary.LittleEndian, dataLen)
//
//	return buf.Bytes()
//}
//
//func WriteUint16() {
//
//}
//
//func WriteString() {
//
//}
