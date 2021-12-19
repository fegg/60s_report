package calendar

import (
	"context"
	"github.com/go-kratos/kratos/pkg/log"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
)

// 批量创建
func (d *Dao) BatchCreateUserDailyReport(ctx context.Context, data []*model.UserDailyReport) (err error) {
	db := d.orm.Model(&model.UserDailyReport{})
	err = db.Create(&data).Error
	if err != nil {
		log.Errorc(ctx, "[dao][BatchCreateUserDailyReport] 批量创建统计日报，err：%+v,", err)
	}
	return
}

// 修改
func (d *Dao) UpdateUserDailyReport(ctx context.Context, data *model.UserDailyReport, filter map[string]interface{}) (err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.UserDailyReport{}).Where(whereSql, params...)
	err = db.Updates(data).Error
	if err != nil {
		log.Errorc(ctx, "[dao][[UpdateUserDailyReport]修改统计填报失败，err:%v", err)
	}
	return
}
