package printerUsecase

import (
	"encoding/hex"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"golang.org/x/text/encoding/charmap"
)

func (p *printerUsecase) AddOperationToCheck(op models.Operation) {
	//data, cmdLen := p.createCommandData(consts.OperationV2)
	//byes :=
	//p.client.Send(bytes)

	buf, cmdLen := p.createCommandBuffer(consts.OperationV2, p.password)

	// Запись типа операции
	buf.WriteByte(op.Type)

	// Запись количества товара
	// Количество записывается в миллиграммах
	amount, err := p.intToBytesWithLen(op.Amount*consts.Milligram, 6)
	if err != nil {
		p.logger.Error(err)
		return
	}
	p.logger.Debug("amount:\n", hex.Dump(amount))
	buf.Write(amount)

	// запись цены товара
	// цена записывается в копейках
	price, err := p.intToBytesWithLen(op.Price, 5) // одна копейка
	if err != nil {
		p.logger.Error(err)
		return
	}

	buf.Write(price)

	// запись суммы товара
	// Сумма записывается в копейках
	summ, err := p.intToBytesWithLen(op.Sum, 5) // две копейки
	if err != nil {
		p.logger.Error(err)
		return
	}
	buf.Write(summ)

	// Запись налогов на товар
	// Налог записывается в копейках
	//tax, err := intToBytesWithLen(0, 5)
	//if err != nil {
	//	p.logger.Fatal(err)
	//}
	buf.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff}) // если нет налога надо отправлять 0xff*6
	//buf.Write(tax)

	// Запись налоговой ставки
	buf.WriteByte(consts.VAT0)
	// Запись номера отдела
	buf.WriteByte(1)

	// Запись признака способа рассчета
	buf.WriteByte(consts.FullPayment)

	// Запись признака предмета рассчета
	buf.WriteByte(op.Subject)

	// Запись название товара 0 - 128 байт строка
	// кодировка win1251
	str, err := charmap.Windows1251.NewEncoder().String(op.Name)
	if err != nil {
		p.logger.Error(err)
		return
	}
	// создаем массив с длинной 128 байт
	rStrBytes := make([]byte, 128)
	copy(rStrBytes, str)
	buf.Write(rStrBytes[:128]) // записываем только первые 128 байт

	p.logger.Debug("длинна сообщения в байтах: ", buf.Len())
	p.logger.Debug("\n", hex.Dump(buf.Bytes()))

	p.logger.Debug("cmdlen", cmdLen)
	rFrame, err := p.client.Send(buf.Bytes(), cmdLen)

	if err != nil {
		p.logger.Error(err)
	}

	if err := models.CheckOnPrinterError(rFrame.ERR); err != nil {
		p.logger.Error(err)
	}

	p.logger.Debug("frame in: \n", hex.Dump(rFrame.Bytes()))

}
