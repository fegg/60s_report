package calendar

import (
	"context"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
)

func (d *Dao) GetProjectTagData(where map[string]interface{}) (items []*model.ProjectTag, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	err = d.orm.Where(whereSql, params...).Find(&items).Error
	return
}

// 获取有效的标签数据
func (d *Dao) ListValidProjectTags(ctx context.Context) (tags []model.ProjectTag) {
	tags = make([]model.ProjectTag, 0)

	_ = d.orm.Model(&model.ProjectTag{}).Where("state = ?", 1).Find(&tags).Error

	return
}
