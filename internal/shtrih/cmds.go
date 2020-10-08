package shtrih

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"net"
)

func (p *Printer) createCommandData(command uint16) (data []byte, cmdLen int) {
	dataBuffer := bytes.NewBuffer([]byte{})

	cb := make([]byte, 2)
	binary.BigEndian.PutUint16(cb, command)
	cb = bytes.TrimPrefix(cb, []byte{0})

	dataBuffer.Write(cb)

	passwordBinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(passwordBinary, p.password)
	dataBuffer.Write(passwordBinary) // write password

	return dataBuffer.Bytes(), len(cb)
}

func (p *Printer) sendCommand(command uint16) ([]byte, error) {
	cmdBinary, cmdLen := p.createCommandData(command)
	frame := p.client.createFrame(cmdBinary)

	con, _ := net.Dial("tcp", p.client.host)
	if err := p.client.sendFrame(frame, con); err != nil {
		return nil, err
	}

	rFrame, err := p.client.receiveFrame(con, byte(cmdLen))
	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		return nil, err
	}

	return rFrame.DATA, nil
}

func (p *Printer) WriteTable(tableNumber, rowNumber byte, fieldNumber uint16, fieldValue []byte) {
	data, cmdLen := p.createCommandData(WriteTable)

	buf := bytes.NewBuffer(data)

	buf.WriteByte(tableNumber)
	buf.WriteByte(rowNumber)

	cb := make([]byte, 2)
	binary.BigEndian.PutUint16(cb, fieldNumber)
	cb = bytes.TrimPrefix(cb, []byte{0})
	buf.Write(cb)
	buf.Write(fieldValue)
	frame := p.client.createFrame(buf.Bytes())

	p.logger.Debug("write table frame:")
	p.logger.Debug("\n", hex.Dump(frame))

	con, _ := net.Dial("tcp", p.client.host)
	if err := p.client.sendFrame(frame, con); err != nil {
		p.logger.Fatal(err)
	}

	rFrame, err := p.client.receiveFrame(con, byte(cmdLen))
	p.logger.Debug("recived table frame:")
	p.logger.Debug("\n", hex.Dump(rFrame.bytes()))
	if err != nil {
		p.logger.Fatal(err)
	}

	if err := checkOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Fatal(err)
	}

	//stream.write((v >>> 0) & 0xFF);
	//stream.write((v >>> 8) & 0xFF);
	//out.writeShort(rowNumber);
	//out.writeByte(fieldNumber);
	//out.writeBytes(fieldValue);

	//_, err := p.sendCommand(WriteTable)
	//if err != nil {
	//	p.logger.Fatal(err)
	//	return
	//}

	//params := make([]string, 4)
	//params[0] = strconv.Itoa(tableNumber)
	//params[1] = strconv.Itoa(rowNumber)
	//params[2] = strconv.Itoa(fieldNumber)
	//params[3] = fieldValue
	//
	////public void execute(int[] data, Object object) throws Exception {
	////	DIOUtils.checkDataMinLength(data, 3);
	////	DIOUtils.checkObjectMinLength((String[]) object, 1);
	////
	////	int tableNumber = data[0];
	////	int rowNumber = data[1];
	////	int fieldNumber = data[2];
	////	String fieldValue = ((String[]) (object))[0];
	////fieldValue = service.decodeText(fieldValue)
	//p.writeTable(tableNumber, rowNumber, fieldNumber, fieldValue)
	////service.printer.check();
	//}

}
