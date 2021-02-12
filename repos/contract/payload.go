package contract

import (
	"errors"
	"tezosign/models"

	"gorm.io/gorm"
)

const PayloadsTable = "requests"

func (r *Repository) GetPayloadByHash(id string) (payload models.Request, isFound bool, err error) {
	err = r.db.
		Table(PayloadsTable).
		Model(models.Request{}).
		Where("req_hash = ?", id).
		First(&payload).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return payload, false, nil
		}
		return payload, false, err
	}

	return payload, true, nil
}

func (r *Repository) SavePayload(request models.Request) (err error) {
	err = r.db.
		Table(PayloadsTable).
		Model(models.Request{}).
		Create(&request).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPayloadsReportByContractID(id uint64) (requests []models.RequestReport, err error) {
	err = r.db.Select("requests.*, signatures").
		Model(models.Request{}).
		Table(PayloadsTable).
		Joins("LEFT JOIN request_json_signatures as rjs on rjs.req_id = requests.req_id").
		Where("ctr_id = ?", id).
		//TODO order by time
		Order("req_id desc").
		Find(&requests).Error
	if err != nil {
		return requests, err
	}

	return requests, nil
}

func (r *Repository) UpdatePayload(request models.Request) (err error) {
	err = r.db.Model(&models.Request{ID: request.ID}).
		Updates(models.Request{
			Status:      request.Status,
			OperationID: request.OperationID,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPayloadByContractAndCounter(contractID uint64, counter int64) (request models.Request, isFound bool, err error) {
	err = r.db.
		Table(PayloadsTable).
		Model(models.Request{}).
		Where("ctr_id = ? and req_counter = ?", contractID, counter).
		First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return request, false, nil
		}
		return request, false, err
	}

	return request, true, nil
}
