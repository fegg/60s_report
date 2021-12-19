package calendar

import (
	"github.com/go-kratos/kratos/pkg/log"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
)

func (d *Dao) GetUserAutoFillData(where map[string]interface{}) (data []*model.UserAutoFill, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	err = d.orm.Model(&model.UserAutoFill{}).Where(whereSql, params...).Find(&data).Error
	if err != nil {
		log.Error("[dao][[GetUserProjectResultData]修改统计填报失败，err:%v", err)
	}
	return
}
