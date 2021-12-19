package calendar

import (
	"context"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"med-common/library/code"

	"github.com/go-kratos/kratos/pkg/log"
)

func (d *Dao) FetchOneDept(ctx context.Context, filter map[string]interface{}) (dept *model.Dept, err error) {
	dept = &model.Dept{}

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.Dept{}).Where(whereSql, params...)
	err = db.Order("dept_code desc").First(dept).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchOneDept] 查询失败，err：%+v,", err)
	}
	return
}

func (d *Dao) FetchAllDept(ctx context.Context, filter map[string]interface{}, order string) (resp []*model.Dept, err error) {
	resp = make([]*model.Dept, 0)

	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.Dept{}).Where(whereSql, params...)
	err = db.Order(order).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchAllDept] 查询失败，err：%+v,err")
	}
	return
}

func (d *Dao) UpdateDept(ctx context.Context, data *model.Dept, filter map[string]interface{}) (err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.Dept{}).Where(whereSql, params...)
	err = db.Updates(data).Error
	if err != nil {
		log.Errorc(ctx, "[dao][[UpdateDept]修改部门，err:%v", err)
	}
	return
}

// 创建
func (d *Dao) CreateDept(ctx context.Context, project *model.Dept) (err error) {
	db := d.orm.Model(&model.Dept{})
	err = db.Create(project).Error
	if err != nil {
		log.Errorc(ctx, "[dao][CreateDept] 创建日常事务，err：%+v,", err)
	}
	return
}

func (d *Dao) FetchAllDeptMap() (resp map[string]*model.Dept, err error) {
	resp = make(map[string]*model.Dept, 0)

	var list []*model.Dept
	err = d.orm.Model(&model.Dept{}).Find(&list).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		return
	}
	for _, dept := range list {
		resp[dept.DeptCode] = dept
	}
	return
}

func (d *Dao) FetchAllDeptMapViaDeptId() (resp map[string]*model.Dept, err error) {
	resp = make(map[string]*model.Dept, 0)

	var list []*model.Dept
	err = d.orm.Model(&model.Dept{}).Find(&list).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		return
	}
	for _, dept := range list {
		resp[dept.DeptId] = dept
	}
	return
}

// 处理部门标签
func (d *Dao) PostDeptTags(ctx context.Context, deptTag *model.DeptTag, relationDeptTag []*model.RelationDeptDeptTag) (deptTagId int64, err error) {
	new_db := d.orm.Begin()
	if err := new_db.Error; err != nil {
		log.Errorc(ctx, "[dao][PostDeptTags] 启动事务失败，err：%+v,", err)
		return 0, code.MedErrDbTransactionErr
	}

	defer func() {
		if err != nil {
			new_db.Rollback()
			log.Errorc(ctx, "[dao][dept] PostDeptTags 修改数据失败，err：%+v,", err)
			return
		}
	}()

	var dbDeptTag *model.DeptTag
	if err = d.orm.Model(&model.DeptTag{}).Where("name = ?", deptTag.Name).Where("state = ?", 1).
		First(&dbDeptTag).Error; err != nil {
		if err = dbutil.IgnoreNoRecordErr(err); err != nil {
			return 0, code.MedErrDBFailed
		}
	}
	// 名字已经存在
	if dbDeptTag != nil && dbDeptTag.Id > 0 {
		if dbDeptTag.Id == 0 {	// 新增
			return 0, code.MedErrDeptTagRepeat
		} else {	// 修改
			if dbDeptTag.Id != deptTag.Id {
				return 0, code.MedErrDeptTagRepeat
			}
		}
	}
	if deptTag.Id == 0 {
		if err = new_db.Model(&model.DeptTag{}).Create(&deptTag).Error; err != nil {
			log.Errorc(ctx, "[dao|dept] PostDeptTags 创建部门标签失败", err)
			return
		}
		for _, v := range relationDeptTag {
			v.DeptTagId = deptTag.Id
		}

		if err = new_db.Model(&model.RelationDeptDeptTag{}).Create(&relationDeptTag).Error; err != nil {
			log.Errorc(ctx, "[dao|dept] PostDeptTags 创建部门标签关联关系失败", err)
			return
		}
	} else {
		if err = new_db.Model(&model.DeptTag{}).Where("id = ?", deptTag.Id).Updates(&deptTag).Error; err != nil {
			log.Errorc(ctx, "[dao|dept] PostDeptTags 更新部门标签失败", err)
			return
		}
		// 删除老的
		if err = new_db.Model(&model.RelationDeptDeptTag{}).Where("dept_tag_id = ?", deptTag.Id).Where("state = ?", 1).
			Updates(map[string]interface{}{"state": -1}).Error; err != nil {
			log.Errorc(ctx, "[dao|dept] PostDeptTags 更行部门标签失败", err)
			return
		}
		for _, v := range relationDeptTag {
			v.DeptTagId = deptTag.Id
		}
		if err = new_db.Model(&model.RelationDeptDeptTag{}).Create(&relationDeptTag).Error; err != nil {
			log.Errorc(ctx, "[dao|dept] PostDeptTags 创建部门标签关联关系失败2", err)
			return
		}
	}
	if err := new_db.Commit().Error; err != nil {
		new_db.Rollback()
		log.Errorc(ctx, "[dao][PostDeptTags] 提交事务失败，err：%+v,", err)
		return 0, code.MedErrDbTransactionErr
	}
	deptTagId = deptTag.Id
	return deptTagId, nil
}

// 删除部门标签
func (d *Dao) DeleteDeptTags(ctx context.Context, ids []int64) (err error) {
	if len(ids) == 0 {
		return
	}

	if err = d.orm.Model(&model.DeptTag{}).Where("id in (?)", ids).Updates(map[string]interface{}{"state": -1}).Error; err != nil {
		return code.MedErrDBFailed
	}
	if err = d.orm.Model(&model.RelationDeptDeptTag{}).Where("dept_tag_id in (?)", ids).Updates(map[string]interface{}{"state": -1}).Error; err != nil {
		return code.MedErrDBFailed
	}
	return
}

func (d *Dao) FetchAllDeptTagMap() (resp map[int64]*model.DeptTag, err error) {
	resp = make(map[int64]*model.DeptTag, 0)

	var list []*model.DeptTag
	err = d.orm.Model(&model.DeptTag{}).Find(&list).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		return
	}
	for _, deptTag := range list {
		resp[deptTag.Id] = deptTag
	}
	return
}

func (d *Dao) GetRelationDeptDeptTagPaginate(ctx context.Context, filter map[string]interface{}, page, limit int64) (resp []*model.RelationDeptDeptTag, err error) {
	resp = make([]*model.RelationDeptDeptTag, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.RelationDeptDeptTag{}).Where(whereSql, params...)
	err = db.Order("id desc").Limit(int(limit)).Offset((int(page)-1)*int(limit)).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetRelationDeptDeptTagPaginate] 查询失败，err：%+v,err", err)
	}
	return
}

func (d *Dao) GetRelationDeptDeptTagIdsPaginate(ctx context.Context, filter map[string]interface{}, page, limit int64) (resp []int64, err error) {
	resp = make([]int64, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.RelationDeptDeptTag{}).Where(whereSql, params...)
	err = db.Select("dept_tag_id").Order("id desc").Group("dept_tag_id").Limit(int(limit)).Offset((int(page)-1)*int(limit)).Pluck("dept_tag_id", &resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetRelationDeptDeptTagPaginate] 查询失败，err：%+v,err", err)
	}
	return
}

func (d *Dao) FetchRelationDeptDeptTag(ctx context.Context, deptTagIds []int64) (ret map[int64][]*model.RelationDeptDeptTag, err error){
	ret = make(map[int64][]*model.RelationDeptDeptTag, 0)
	var list []*model.RelationDeptDeptTag
	err = d.orm.Model(&model.RelationDeptDeptTag{}).Where("dept_tag_id in (?)", deptTagIds).Where("state = ?", 1).Find(&list).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchRelationDeptDeptTag] 查询失败，err：%+v,err", err)
		return
	}
	for _, v := range list {
		if _, ok := ret[v.DeptTagId]; !ok {
			ret[v.DeptTagId] = make([]*model.RelationDeptDeptTag, 0)
		}
		ret[v.DeptTagId] = append(ret[v.DeptTagId], v)
	}
	return
}

func (d *Dao) GetDeptPaginate(ctx context.Context, filter map[string]interface{}, page, limit int) (resp []*model.Dept, err error) {
	resp = make([]*model.Dept, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.Dept{}).Where(whereSql, params...)
	err = db.Order("updated_at desc").Limit(limit).Offset((page-1)*limit).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][GetDeptPaginate] 查询失败，err：%+v,err", err)
	}
	return
}
