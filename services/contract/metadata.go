package contract

import (
	"errors"

	"blockwatch.cc/tzindex/micheline"
)

const (
	metadataThumbnailUri           = "thumbnailUri"
	metadataThumbnailUriJsonNaming = "thumbnail_uri"

	MetaDataEntrypoint = "tokenmetadata"
)

func ParseMetadata(data []byte) (fields map[string]interface{}, err error) {

	prim := &micheline.Prim{}

	err = prim.UnmarshalJSON(data)
	if err != nil {
		return fields, err
	}

	if len(prim.Args) != 2 {
		return fields, errors.New("wrong args len")
	}

	fields = map[string]interface{}{}
	var name string
	for i := range prim.Args[1].Args {

		if len(prim.Args[1].Args[i].Args) != 2 {
			return fields, errors.New("wrong elem args len")
		}

		name = prim.Args[1].Args[i].Args[0].String

		if name == metadataThumbnailUri {
			name = metadataThumbnailUriJsonNaming
		}

		fields[name] = string(prim.Args[1].Args[i].Args[1].Bytes)
	}

	return fields, nil
}
