package types

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"

	"blockwatch.cc/tzindex/micheline"
)

//TODO reformat consts
const (
	MichelineTypeInt    byte = 0b_0000_0000
	MichelineTypeBytes       = 0b_0010_0000
	MichelineTypeString      = 0b_0100_0000
	MichelineTypeArray       = 0b_0110_0000

	AnnotationTypeField    = 0b_0100_0000
	AnnotationTypeType     = 0b_1000_0000
	AnnotationTypeVariable = 0b_1100_0000
)

type TZKTPrim micheline.Prim

func (p *TZKTPrim) UnmarshalBinary(data []byte) (err error) {
	return p.DecodeBuffer(bytes.NewBuffer(data))
}

func (p TZKTPrim) MichelinePrim() *micheline.Prim {
	prim := micheline.Prim(p)
	return &prim
}

//TODO add more OpCodes
func PrimTypeFromTypeCode(o micheline.OpCode) micheline.PrimType {

	switch o {
	case micheline.D_PAIR, micheline.T_PAIR, micheline.T_OR, micheline.T_LAMBDA:
		return micheline.PrimBinary
	case micheline.K_STORAGE, micheline.K_PARAMETER, micheline.T_LIST, micheline.T_OPTION, micheline.D_LEFT, micheline.D_RIGHT, micheline.D_SOME,
		micheline.I_NIL, micheline.I_DROP, micheline.I_SOME:
		return micheline.PrimUnary
	case micheline.T_KEY, micheline.T_STRING, micheline.T_NAT, micheline.T_INT, micheline.T_MUTEZ,
		micheline.T_TIMESTAMP, micheline.T_SIGNATURE, micheline.T_UNIT, micheline.T_ADDRESS, micheline.T_KEY_HASH, micheline.T_OPERATION,
		micheline.D_NONE:
		return micheline.PrimNullary
	default:
		return micheline.PrimBytes
	}

}

//Implementation of TZKT decoder
//https://github.com/baking-bad/netezos/blob/master/Netezos/Encoding/Micheline/Micheline.cs#L72
func (p *TZKTPrim) DecodeBuffer(buf *bytes.Buffer) (err error) {
	tag, _ := buf.ReadByte()

	if tag >= 0x80 {

		bt, err := buf.ReadByte()
		if err != nil {
			return err
		}

		p.OpCode = micheline.OpCode(bt)

		p.Type = PrimTypeFromTypeCode(p.OpCode)

		//Init int field
		if p.OpCode == micheline.T_INT || p.OpCode == micheline.T_NAT {
			p.Int = big.NewInt(0)
		}

		var args = (tag & 0x70) >> 4

		if args > 0 {
			if args == 0x07 {
				args, err = Read7BitInt(buf)
				if err != nil {
					return err
				}
			}

			var arr []*micheline.Prim
			for ; args > 0; args-- {

				pr := TZKTPrim{}

				err = pr.DecodeBuffer(buf)
				if err != nil {
					return err
				}

				pRim := micheline.Prim(pr)

				arr = append(arr, &pRim)
			}

			p.Args = arr
		}

		var annotsLen = tag & 0x0F

		if annotsLen > 0 {
			if annotsLen == 0x0F {
				annotsLen, err = Read7BitInt(buf)
				if err != nil {
					return err
				}
			}

			var annots []string
			for ; annotsLen > 0; annotsLen-- {
				anno, err := ReadAnno(buf)
				if err != nil {
					return err
				}

				annots = append(annots, anno)
			}

			//Change to Anno type
			p.Type = p.Type + 1

			p.Anno = annots
		}

	} else {

		cnt := tag & 0x1F
		if cnt == 0x1F {
			cnt, err = Read7BitInt(buf)
			if err != nil {
				return err
			}
		}

		opCode := tag & 0xE0

		switch opCode {
		case MichelineTypeArray:
			var arr []*micheline.Prim

			p.Type = micheline.PrimSequence
			p.OpCode = micheline.T_LIST

			for ; cnt > 0; cnt-- {
				pr := TZKTPrim{}

				err = pr.DecodeBuffer(buf)
				if err != nil {
					return err
				}

				pRim := micheline.Prim(pr)

				arr = append(arr, &pRim)
			}

			p.Args = arr

		case MichelineTypeBytes:
			p.Type = micheline.PrimBytes
			p.OpCode = micheline.T_BYTES
			p.Bytes = buf.Next(int(cnt))

		case MichelineTypeInt:
			p.Type = micheline.PrimInt
			p.OpCode = micheline.T_INT
			p.Int = big.NewInt(0).SetBytes(buf.Next(int(cnt)))

		case MichelineTypeString:
			p.Type = micheline.PrimString
			p.OpCode = micheline.T_STRING
			p.String = string(buf.Next(int(cnt)))

		default:
			return errors.New("Wrong tag")
		}
	}

	return nil
}

func Read7BitInt(data *bytes.Buffer) (res uint8, err error) {

	var b byte

	for bits := 0; bits < 28; {
		b, err = data.ReadByte()
		if err != nil {
			return res, err
		}

		res |= (b & 0x7F) << bits
		bits += 7

		if b < 0x80 {
			return res, nil
		}

	}

	b, err = data.ReadByte()
	if err != nil {
		return res, err
	}

	if b > 0x0F {
		return res, errors.New("Int32 overflow")
	}

	res |= b << 28

	return res, nil
}

func ReadAnno(buf *bytes.Buffer) (anno string, err error) {
	tag, err := buf.ReadByte()
	if err != nil {
		return anno, err
	}

	var cnt = tag & 0x3F

	if cnt == 0x3F {
		cnt, err = Read7BitInt(buf)
		if err != nil {
			return anno, err
		}
	}

	var annoPrefix string

	//TzKt prefixes do not match with tzIndex
	switch tag & AnnotationTypeVariable {
	case AnnotationTypeField:
		annoPrefix = micheline.VarAnnoPrefix
	case AnnotationTypeType:
		annoPrefix = micheline.TypeAnnoPrefix
	case AnnotationTypeVariable:
		annoPrefix = micheline.FieldAnnoPrefix
	default:
		return anno, errors.New("invalid annotation tag")
	}

	return fmt.Sprint(annoPrefix, string(buf.Next(int(cnt)))), nil
}

func (p *TZKTPrim) Scan(value interface{}) (err error) {
	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	if len(data) == 0 {
		return nil
	}

	err = p.UnmarshalBinary(data)
	if err != nil {
		return fmt.Errorf("UnmarshalBinary: %s", err.Error())
	}

	return nil
}
func (p TZKTPrim) Value() (driver.Value, error) {

	return []byte{}, errors.New("Not implemented")
}
